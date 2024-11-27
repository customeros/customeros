package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"gorm.io/gorm"
	"log"
	"strings"
)

type mailstackService struct {
	cfg      *config.GlobalConfig
	services *Services
}

type MailstackService interface {
	GetPaymentIntent(ctx context.Context, domains []string, usernames []string, amount int64) (string, error)                                         // stripe client secret
	RegisterBuyDomainsWithMailboxes(ctx context.Context, test bool, paymentIntentId string, domains []string, usernames []string, amount int64) error // id, stripe client secret
}

func NewMailstackService(cfg *config.GlobalConfig, services *Services) MailstackService {
	return &mailstackService{
		cfg:      cfg,
		services: services,
	}
}

func (s *mailstackService) GetPaymentIntent(ctx context.Context, domains []string, usernames []string, amount int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.GetPaymentIntent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("request.domains", domains)
	span.LogKV("request.usernames", usernames)

	tenant := common.GetTenantFromContext(ctx)
	email := common.GetUserEmailFromContext(ctx)

	if s.cfg.ExternalServices.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return "", err
	}

	//create stripe payment intent
	stipePaymentDescription := fmt.Sprintf("Mailstack purchase: %d domains (%s) with usernames (%s)", len(domains), strings.Join(domains, ", "), strings.Join(usernames, ", "))

	params := &stripe.PaymentIntentParams{
		Amount:       stripe.Int64(amount),                      // Amount in cents (e.g., $20.00)
		Currency:     stripe.String(string(stripe.CurrencyUSD)), // Currency (e.g., USD)
		Description:  stripe.String(stipePaymentDescription),
		ReceiptEmail: stripe.String(email),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: utils.BoolPtr(true),
		},
		Metadata: map[string]string{
			"tenant":    tenant,
			"username":  email,
			"domains":   strings.Join(domains, ","),
			"usernames": strings.Join(usernames, ","),
		},
	}

	// Create a PaymentIntent
	stripe.Key = s.cfg.ExternalServices.StripeConfig.ApiKey
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Fatalf("Failed to create payment intent: %v", err)
	}

	return pi.ClientSecret, nil
}

func (s *mailstackService) RegisterBuyDomainsWithMailboxes(ctx context.Context, test bool, paymentIntentId string, domains []string, usernames []string, amount int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RegisterBuyDomainsWithMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("request.domains", domains)
	span.LogKV("request.usernames", usernames)

	if s.cfg.ExternalServices.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return err
	}

	//call stripe and check if payment is successful
	stripe.Key = s.cfg.ExternalServices.StripeConfig.ApiKey
	stripePaymentIntent, err := paymentintent.Get(paymentIntentId, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if stripePaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		err := errors.New("Payment not successful")
		tracing.TraceErr(span, err)
		return err
	}

	mailstackBuyRequestId := ""

	err = s.services.PostgresRepositories.Db.Transaction(func(tx *gorm.DB) error {

		mailstackBuyRequestId, err = s.services.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, tx, &entity.MailstackBuyRequest{
			Domains:         strings.Join(domains, ","),
			Usernames:       strings.Join(usernames, ","),
			Status:          entity.MailstackBuyRequestStatusPending,
			PaymentIntentId: paymentIntentId,
		})
		if err != nil {
			return err
		}

		for _, domain := range domains {
			err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreDomain(ctx, tx, &entity.MailstackBuyRequestDomain{
				MailstackBuyRequestId: mailstackBuyRequestId,
				Domain:                domain,
				Status:                entity.MailstackBuyRequestDomainStatusPending,
			})
			if err != nil {
				return err
			}

			for _, username := range usernames {
				err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreMailbox(ctx, tx, &entity.MailstackBuyRequestMailbox{
					MailstackBuyRequestId: mailstackBuyRequestId,
					Domain:                domain,
					Username:              username,
					Mailbox:               username + "@" + domain,
					Status:                entity.MailstackBuyRequestDomainMailboxPending,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if !test {
		err = s.services.RabbitMQService.PublishEvent(ctx, mailstackBuyRequestId, model.MAILSTACK_BUY_REQUEST, dto.MailstackBuyRequest{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
