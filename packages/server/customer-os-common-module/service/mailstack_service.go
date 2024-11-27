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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"gorm.io/gorm"
	"log"
	"strings"
)

type mailstackService struct {
	cfg      *config.GlobalConfig
	services *Services
}

type MailstackService interface {
	RegisterBuyDomainsWithMailboxes(ctx context.Context, domains []string, usernames []string, amount int64) (string, string, error) // id, stripe client secret
	MarkBuyRequestAsPaid(ctx context.Context, id string) error
}

func NewMailstackService(cfg *config.GlobalConfig, services *Services) MailstackService {
	return &mailstackService{
		cfg:      cfg,
		services: services,
	}
}

func (s *mailstackService) RegisterBuyDomainsWithMailboxes(ctx context.Context, domains []string, usernames []string, amount int64) (string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RegisterBuyDomainsWithMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("request.domains", domains)
	span.LogKV("request.usernames", usernames)

	tenant := common.GetTenantFromContext(ctx)
	email := common.GetUserEmailFromContext(ctx)

	if s.cfg.ExternalServices.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return "", "", err
	}

	id, err := s.services.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, nil, &entity.MailstackBuyRequest{
		Domains:   strings.Join(domains, ","),
		Usernames: strings.Join(usernames, ","),
		Status:    entity.MailstackBuyRequestStatusAwaitingPayment,
	})
	if err != nil {
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return "", "", nil
	}

	//create stripe payment intent
	stipePaymentDescription := fmt.Sprintf("Mailstack purchase: %d domains (%s) with usernames (%s)", len(domains), strings.Join(domains, ", "), strings.Join(usernames, ", "))

	params := &stripe.PaymentIntentParams{
		Amount:       stripe.Int64(amount),                      // Amount in cents (e.g., $20.00)
		Currency:     stripe.String(string(stripe.CurrencyUSD)), // Currency (e.g., USD)
		Description:  stripe.String(stipePaymentDescription),
		ReceiptEmail: stripe.String(email),
		Params: stripe.Params{
			Metadata: map[string]string{
				"mailstack_buy_request_id": id,
				"tenant":                   tenant,
				"username":                 email,
				"domains":                  strings.Join(domains, ","),
				"usernames":                strings.Join(usernames, ","),
			},
		},
	}

	// Create a PaymentIntent
	stripe.Key = s.cfg.ExternalServices.StripeConfig.ApiKey
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Fatalf("Failed to create payment intent: %v", err)
	}

	err = s.services.PostgresRepositories.Db.Transaction(func(tx *gorm.DB) error {

		mailstackBuyRequest, err := s.services.PostgresRepositories.MailstackBuyRequestRepository.GetById(ctx, id)
		if err != nil {
			tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
			return nil
		}

		mailstackBuyRequest.PaymentIntentId = pi.ID
		mailstackBuyRequest.PaymentIntentClientSecret = pi.ClientSecret

		_, err = s.services.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, tx, mailstackBuyRequest)
		if err != nil {
			tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
			return nil
		}

		for _, domain := range strings.Split(mailstackBuyRequest.Domains, ",") {
			err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreDomain(ctx, tx, &entity.MailstackBuyRequestDomain{
				MailstackBuyRequestId: id,
				Domain:                domain,
				Status:                entity.MailstackBuyRequestDomainStatusAwaitingPayment,
			})
			if err != nil {
				return err
			}

			for _, username := range strings.Split(mailstackBuyRequest.Usernames, ",") {
				err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreMailbox(ctx, tx, &entity.MailstackBuyRequestMailbox{
					MailstackBuyRequestId: id,
					Domain:                domain,
					Username:              username,
					Mailbox:               username + "@" + domain,
					Status:                entity.MailstackBuyRequestDomainMailboxAwaitingPayment,
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
		return "", "", err
	}

	return id, pi.ClientSecret, nil
}

func (s *mailstackService) MarkBuyRequestAsPaid(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.MarkBuyRequestAsPaid")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("id", id)

	if s.cfg.ExternalServices.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return err
	}

	mailstackBuyRequest, err := s.services.PostgresRepositories.MailstackBuyRequestRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return err
	}

	if mailstackBuyRequest == nil {
		err := errors.New("Mailstack buy request not found")
		tracing.TraceErr(span, err)
		return err
	}

	//call stripe and check if payment is successful
	//stripe.Key = s.cfg.ExternalServices.StripeConfig.ApiKey
	//stripePaymentIntent, err := paymentintent.Get(mailstackBuyRequest.PaymentIntentId, nil)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return err
	//}
	//
	//if stripePaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
	//	err := errors.New("Payment not successful")
	//	tracing.TraceErr(span, err)
	//	return err
	//}

	mailstackBuyRequest.Status = entity.MailstackBuyRequestStatusPending

	err = s.services.PostgresRepositories.Db.Transaction(func(tx *gorm.DB) error {
		_, err = s.services.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, tx, mailstackBuyRequest)
		if err != nil {
			return err
		}

		domains, err := s.services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, id)
		if err != nil {
			return err
		}

		for _, domain := range domains {

			domain.Status = entity.MailstackBuyRequestDomainStatusPending

			err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreDomain(ctx, tx, domain)
			if err != nil {
				return err
			}
		}

		mailboxes, err := s.services.PostgresRepositories.MailstackBuyRequestRepository.GetMailboxes(ctx, id)
		if err != nil {
			return err
		}

		for _, mailbox := range mailboxes {
			mailbox.Status = entity.MailstackBuyRequestDomainMailboxPending

			err := s.services.PostgresRepositories.MailstackBuyRequestRepository.StoreMailbox(ctx, tx, mailbox)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, mailstackBuyRequest.ID, model.MAILSTACK_BUY_REQUEST, dto.MailstackBuyRequest{})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
