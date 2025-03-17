package mailstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const TEST_MAILBOX_DOMAIN = "testcustomeros.com"

type mailstackService struct {
	cfg      *config.CommonConfig
	events   *events.EventsService
	postgres *postgres_repository.Repositories
	neo4j    *neo4j_repository.Repositories
	opensrs  interfaces.OpenSrsService
	email    interfaces.EmailService
}

func NewMailstackService(cfg *config.CommonConfig, events *events.EventsService, postgres *postgres_repository.Repositories, neo4j *neo4j_repository.Repositories, opensrs interfaces.OpenSrsService, email interfaces.EmailService) interfaces.MailstackService {
	return &mailstackService{
		cfg:      cfg,
		events:   events,
		postgres: postgres,
		opensrs:  opensrs,
		neo4j:    neo4j,
		email:    email,
	}
}

func (s *mailstackService) CheckDomainAvailability(ctx context.Context, tenant, domain string) (int, string, bool, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.CheckDomainAvailability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("request.domain", domain)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/check-availability/"+domain, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", false, false, err
	}

	// Add required headers
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)
	req.Header.Set("Content-Type", "application/json")

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", false, false, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", false, false, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, false, false, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		IsAvailable bool `json:"isAvailable"`
		IsPremium   bool `json:"isPremium"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", false, false, err
	}

	return http.StatusOK, "", response.IsAvailable, response.IsPremium, nil
}

func (s *mailstackService) GetPaymentIntent(ctx context.Context, domains []string, usernames []string, amount int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.GetPaymentIntent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("request.domains", domains)
	span.LogKV("request.usernames", usernames)
	span.LogFields(tracingLog.Int64("request.amount", amount))

	tenant := common.GetTenantFromContext(ctx)
	email := common.GetUserEmailFromContext(ctx)

	if s.cfg.External.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return "", err
	}

	// validate domains before creating payment intent
	for _, domain := range domains {
		statusCode, errMsg, available, _, err := s.CheckDomainAvailability(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error checking domain availability"))
			return "", err
		}
		if statusCode != http.StatusOK {
			err = errors.New(errMsg)
			tracing.TraceErr(span, err)
			return "", err
		}
		if !available {
			err = errors.New("Domain not available for purchase")
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	// create stripe payment intent
	stipePaymentDescription := fmt.Sprintf("Mailstack purchase: %d domains (%s) with usernames (%s)", len(domains), strings.Join(domains, ", "), strings.Join(usernames, ", "))

	params := &stripe.PaymentIntentParams{
		Amount:       stripe.Int64(amount), // Amount in cents (e.g., 2000 for $20.00)
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
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
	stripe.Key = s.cfg.External.StripeConfig.ApiKey
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Fatalf("Failed to create payment intent: %v", err)
	}

	return pi.ClientSecret, nil
}

func (s *mailstackService) RegisterBuyDomainsWithMailboxes(ctx context.Context, test bool, paymentIntentId string, domains []string, usernames []string, redirectWebsite string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RegisterBuyDomainsWithMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogKV("domains", domains)
	span.LogKV("usernames", usernames)
	span.LogKV("redirectWebsite", redirectWebsite)

	tenant := common.GetTenantFromContext(ctx)

	if s.cfg.External.StripeConfig.ApiKey == "" {
		err := errors.New("Stripe API key not set")
		tracing.TraceErr(opentracing.SpanFromContext(ctx), err)
		return err
	}

	// call stripe and check if payment is successful
	stripe.Key = s.cfg.External.StripeConfig.ApiKey
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

	err = s.postgres.Db.Transaction(func(tx *gorm.DB) error {
		mailstackBuyRequestId, err = s.postgres.MailstackBuyRequestRepository.Store(ctx, tx, &postgres_entity.MailstackBuyRequest{
			Domains:         strings.Join(domains, ","),
			Usernames:       strings.Join(usernames, ","),
			Status:          postgres_entity.MailstackBuyRequestStatusPending,
			PaymentIntentId: paymentIntentId,
		})
		if err != nil {
			return err
		}

		for _, domain := range domains {
			err := s.postgres.MailstackBuyRequestRepository.StoreDomain(ctx, tx, &postgres_entity.MailstackBuyRequestDomain{
				MailstackBuyRequestId: mailstackBuyRequestId,
				Domain:                domain,
				RedirectWebsite:       redirectWebsite,
				Status:                postgres_entity.MailstackBuyRequestDomainStatusPendingProvisioning,
			})
			if err != nil {
				return err
			}

			for _, username := range usernames {
				result, err := s.RegisterMailbox(ctx, tenant, domain, interfaces.CreateMailboxRequest{
					IgnoreDomainOwnership: true,
					Domain:                domain,
					Username:              username,
					Password:              utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false),
					WebmailEnabled:        true,
					ForwardingTo:          []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))},
					LinkedUserEmail:       common.GetUserEmailFromContext(ctx),
				})
				if err != nil {
					return err
				}
				if result.StatusCode != http.StatusCreated {
					return errors.New(result.ErrorMsg)
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
		err = s.events.Publisher.PublishFanoutEvent(ctx, mailstackBuyRequestId, model.MAILSTACK_BUY_REQUEST, dto.MailstackProvisionBuyRequest{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *mailstackService) RegisterMailbox(ctx context.Context, tenant string, domain string, request interfaces.CreateMailboxRequest) (*interfaces.RegisterMailboxResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RegisterMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("request.domain", domain)
	tracing.LogObjectAsJson(span, "request.request", request)

	// Get user ID from linked email
	var userId string
	if request.LinkedUserEmail != "" {
		userDbNode, err := s.neo4j.UserReadRepository.GetFirstUserByEmail(ctx, tenant, request.LinkedUserEmail)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error finding linked user"))
			return nil, err
		}
		if userDbNode != nil {
			userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
			userId = userEntity.Id
		}
	}

	// Create request body for Mailstack API
	reqBody := struct {
		Username              string   `json:"username"`
		Password              string   `json:"password"`
		Domain                string   `json:"domain"`
		ForwardingTo          []string `json:"forwardingTo"`
		WebmailEnabled        bool     `json:"webmailEnabled"`
		UserId                string   `json:"userId"`
		IgnoreDomainOwnership bool     `json:"ignoreDomainOwnership"`
	}{
		Username:              request.Username,
		Password:              request.Password,
		Domain:                domain,
		ForwardingTo:          request.ForwardingTo,
		WebmailEnabled:        request.WebmailEnabled,
		UserId:                userId,
		IgnoreDomainOwnership: request.IgnoreDomainOwnership,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marshal request body")
	}

	// Create request to Mailstack API
	mailstackReq, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/mailboxes", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create request to Mailstack API")
	}

	// Add required headers
	mailstackReq.Header.Set("Content-Type", "application/json")
	mailstackReq.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	mailstackReq.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(mailstackReq.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport and timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(mailstackReq)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to Mailstack API")
	}
	defer resp.Body.Close()

	response := &interfaces.RegisterMailboxResponse{
		StatusCode: resp.StatusCode,
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		response.ErrorMsg = errorResponse.Error
		tracing.TraceErr(span, errors.New(errorResponse.Error))
		return response, nil
	}

	// Parse response
	var mailboxRecord interfaces.MailboxRecord
	if err := json.NewDecoder(resp.Body).Decode(&mailboxRecord); err != nil {
		return nil, errors.Wrap(err, "Unable to parse Mailstack API response")
	}

	// Create email node and link with user if identified
	emailFields := interfaces.EmailFields{
		Email: mailboxRecord.Email,
	}
	var linkWith *common_srv.LinkWith
	if userId != "" {
		linkWith = &common_srv.LinkWith{
			Type: model.USER,
			Id:   userId,
		}
	}
	_, err = s.email.Merge(ctx, nil, tenant, emailFields, linkWith)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating email node"))
		return nil, err
	}

	// Publish fanout event to provision mailbox
	err = s.events.Publisher.PublishFanoutEvent(ctx, mailboxRecord.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error publishing mailbox provision event"))
		return nil, err
	}

	response.Mailbox = &mailboxRecord
	return response, nil
}

func (s *mailstackService) ConfigureMailbox(ctx context.Context, tenant string, mailboxId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mailstackService.ConfigureMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, mailboxId)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/mailboxes/"+mailboxId+"/configure", nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))
		return errors.New(errorResponse.Error)
	}

	return nil
}

func (s *mailstackService) RegisterNewDomain(ctx context.Context, tenant, domain, website string) (int, string, *interfaces.DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RegisterNewDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request body
	reqBody := struct {
		Domain  string `json:"domain"`
		Website string `json:"website"`
	}{
		Domain:  domain,
		Website: website,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to marshal request body"))
		return http.StatusInternalServerError, "Unable to marshal request body", nil, err
	}

	// Create request to Mailstack API
	mailstackReq, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains", bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	mailstackReq.Header.Set("Content-Type", "application/json")
	mailstackReq.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	mailstackReq.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(mailstackReq.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(mailstackReq)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Domain interfaces.DomainRecord `json:"domain"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", &response.Domain, nil
}

func (s *mailstackService) ConfigureDomain(ctx context.Context, tenant, domain, website string) (int, string, *interfaces.DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.ConfigureDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request body
	reqBody := struct {
		Domain  string `json:"domain"`
		Website string `json:"website"`
	}{
		Domain:  domain,
		Website: website,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to marshal request body"))
		return http.StatusInternalServerError, "Unable to marshal request body", nil, err
	}

	// Create request to Mailstack API
	mailstackReq, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/configure", bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	mailstackReq.Header.Set("Content-Type", "application/json")
	mailstackReq.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	mailstackReq.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(mailstackReq.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(mailstackReq)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Domain interfaces.DomainRecord `json:"domain"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", &response.Domain, nil
}

func (s *mailstackService) GetDomains(ctx context.Context, tenant string) (int, string, []interfaces.DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.GetDomains")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains", nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Domains []interfaces.DomainRecord `json:"domains"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", response.Domains, nil
}

func (s *mailstackService) RecommendDomain(ctx context.Context, tenant, baseName string) (int, string, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.RecommendDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/recommendations?baseName="+baseName, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Recommendations []string `json:"recommendations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", response.Recommendations, nil
}

func (s *mailstackService) PurchaseDomain(ctx context.Context, tenant, domain string) (int, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.PurchaseDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request body
	reqBody := struct {
		Domain string `json:"domain"`
	}{
		Domain: domain,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to marshal request body"))
		return http.StatusInternalServerError, "Unable to marshal request body", err
	}

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/purchase", bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, errors.New(errorResponse.Error)
	}

	return http.StatusOK, "", nil
}

func (s *mailstackService) AddDNSRecord(ctx context.Context, tenant, domain string, record interfaces.DNSRecord) (int, string, *interfaces.DNSRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.AddDNSRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request body
	jsonBody, err := json.Marshal(record)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to marshal request body"))
		return http.StatusInternalServerError, "Unable to marshal request body", nil, err
	}

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/domains/%s/dns", s.cfg.Internal.MailstackApiConfig.ApiUrl, domain), bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Record interfaces.DNSRecord `json:"dnsRecord"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", &response.Record, nil
}

func (s *mailstackService) DeleteDNSRecord(ctx context.Context, tenant, domain, dnsId string) (int, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.DeleteDNSRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v1/domains/%s/dns/%s", s.cfg.Internal.MailstackApiConfig.ApiUrl, domain, dnsId), nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, errors.New(errorResponse.Error)
	}

	return http.StatusOK, "", nil
}

func (s *mailstackService) GetDNSRecords(ctx context.Context, tenant, domain string) (int, string, []interfaces.DNSRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.GetDNSRecords")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/domains/%s/dns", s.cfg.Internal.MailstackApiConfig.ApiUrl, domain), nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Records []interfaces.DNSRecord `json:"dnsRecords"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", response.Records, nil
}

func (s *mailstackService) ProcessDMARCMonitoringReport(ctx context.Context, emailData []byte) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.ProcessDMARCMonitoringReport")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Internal.MailstackApiConfig.ApiUrl+"/v1/dmarc", bytes.NewBuffer(emailData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusAccepted {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))
		return errors.New(errorResponse.Error)
	}

	return nil
}

func (s *mailstackService) GetMailboxes(ctx context.Context, tenant, domain string) (int, string, []interfaces.MailboxRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackService.GetMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Create request to Mailstack API
	url := s.cfg.Internal.MailstackApiConfig.ApiUrl + "/v1/mailboxes"
	if domain != "" {
		url += "?domain=" + domain
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to create request to Mailstack API"))
		return http.StatusInternalServerError, "Unable to create request to Mailstack API", nil, err
	}

	// Add required headers
	req.Header.Set("X-CUSTOMER-OS-API-KEY", s.cfg.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport and timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to connect to Mailstack API"))
		return http.StatusInternalServerError, "Unable to connect to Mailstack API", nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read error response body
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			errorResponse.Error = "Unknown error occurred"
		}
		tracing.TraceErr(span, errors.New(errorResponse.Error))

		// For 500 errors, use a generic message
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
			return http.StatusInternalServerError, "Internal server error", nil, errors.New(errorResponse.Error)
		}

		// For other errors, propagate the message from Mailstack
		return resp.StatusCode, errorResponse.Error, nil, errors.New(errorResponse.Error)
	}

	// Parse response
	var response struct {
		Mailboxes []interfaces.MailboxRecord `json:"mailboxes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse Mailstack API response"))
		return http.StatusInternalServerError, "Unable to parse Mailstack API response", nil, err
	}

	return http.StatusOK, "", response.Mailboxes, nil
}
