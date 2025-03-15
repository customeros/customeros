package events_listeners

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	common_model "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type MailstackProvisionBuyRequestListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewMailstackProvisionBuyRequestListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &MailstackProvisionBuyRequestListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.MailstackProvisionBuyRequest](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *MailstackProvisionBuyRequestListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *MailstackProvisionBuyRequestListener) handle(ctx context.Context, entityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	mailstackBuyRequest, err := l.dependencies.PostgresRepositories.MailstackBuyRequestRepository.GetById(ctx, entityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if mailstackBuyRequest == nil {
		err = errors.New("mailstack buy request not found")
		tracing.TraceErr(span, err)
		return err
	}

	span.LogKV("mailstackBuyRequest.Status", mailstackBuyRequest.Status)

	if mailstackBuyRequest.Status != postgres_entity.MailstackBuyRequestStatusPending {
		return nil
	}

	err = l.processDomains(ctx, mailstackBuyRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = l.processMailboxes(ctx, mailstackBuyRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// reload latest state for domains
	domains, err := l.dependencies.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, domain := range domains {
		if domain.Status != postgres_entity.MailstackBuyRequestDomainStatusCompleted {
			continue
		}

		mailboxes, err := l.dependencies.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		for _, mailbox := range mailboxes {
			if mailbox.Status != postgres_entity.MailboxStatusPendingProvisioning {
				continue
			}

			err := l.dependencies.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, mailbox.ID, common_model.MAILBOX, dto.MailstackProvisionMailbox{})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	// mark buy request as completed or failed
	if mailstackBuyRequest.Status == postgres_entity.MailstackBuyRequestStatusPending {

		completed := true

		for _, domain := range domains {
			if domain.Status != postgres_entity.MailstackBuyRequestDomainStatusCompleted {
				completed = false
			}
		}

		if completed {
			mailstackBuyRequest.Status = postgres_entity.MailstackBuyRequestStatusCompleted
		}

		_, err = l.dependencies.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, nil, mailstackBuyRequest)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return nil
}

func (l *MailstackProvisionBuyRequestListener) configureDomainInMailstack(ctx context.Context, domain string, website string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.configureDomainInMailstack")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("missing tenant in context")
		tracing.TraceErr(span, err)
		return err
	}

	statusCode, errorMsg, _, err := l.dependencies.CommonServices.MailstackService.ConfigureDomain(ctx, tenant, domain, website)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if statusCode != http.StatusOK {
		err := errors.New(errorMsg)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (l *MailstackProvisionBuyRequestListener) purchaseDomainInMailstack(ctx context.Context, tenant string, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.purchaseDomainInMailstack")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// Create request body
	reqBody := struct {
		Domain string `json:"domain"`
	}{
		Domain: domain,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to marshal request body"))
		return err
	}

	// Create request to Mailstack API
	req, err := http.NewRequestWithContext(ctx, "POST", l.dependencies.CommonConfig.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/purchase", bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create request to Mailstack API"))
		return err
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CUSTOMER-OS-API-KEY", l.dependencies.CommonConfig.Internal.MailstackApiConfig.ApiKey)
	req.Header.Set("tenant", tenant)

	// Forward Jaeger trace context
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(tracingLog.Error(err))
	}

	// Create HTTP client with default transport
	client := &http.Client{}

	// Make request to Mailstack API
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to connect to Mailstack API"))
		return err
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
		return errors.New(errorResponse.Error)
	}

	return nil
}

func (l *MailstackProvisionBuyRequestListener) processDomains(ctx context.Context, mailstackBuyRequest *postgres_entity.MailstackBuyRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.processDomains")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	domains, err := l.dependencies.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Semaphore to limit to 5 goroutines

	for _, domain := range domains {
		wg.Add(1)

		// Create a new scope to avoid data races on the `domain` variable
		go func(domain *postgres_entity.MailstackBuyRequestDomain) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Step 1 - Purchase domain in Mailstack
			if domain.Status == postgres_entity.MailstackBuyRequestDomainStatusPendingProvisioning {
				err := l.purchaseDomainInMailstack(ctx, domain.Tenant, domain.Domain)
				if err != nil {
					tracing.TraceErr(span, err)
					mailstackBuyRequest.Status = postgres_entity.MailstackBuyRequestStatusFailed
					domain.Status = postgres_entity.MailstackBuyRequestDomainStatusFailed
				} else {
					domain.Status = postgres_entity.MailstackBuyRequestDomainStatusPendingConfiguration
				}

				err = l.dependencies.PostgresRepositories.CommonRepository.UpdateProperty(ctx, domain.Tenant, postgres_entity.MailstackBuyRequestDomain{}, domain.ID, "Status", domain.Status)
				if err != nil {
					tracing.TraceErr(span, err)
					return // Exit the goroutine on error
				}
			}

			// Step 2 - Configure domain in Mailstack
			if domain.Status == postgres_entity.MailstackBuyRequestDomainStatusPendingConfiguration {
				err := l.configureDomainInMailstack(ctx, domain.Domain, domain.RedirectWebsite)
				if err != nil {
					tracing.TraceErr(span, err)
				} else {
					domain.Status = postgres_entity.MailstackBuyRequestDomainStatusCompleted
				}

				err = l.dependencies.PostgresRepositories.CommonRepository.UpdateProperty(ctx, domain.Tenant, postgres_entity.MailstackBuyRequestDomain{}, domain.ID, "Status", domain.Status)
				if err != nil {
					tracing.TraceErr(span, err)
					return // Exit the goroutine on error
				}
			}
		}(domain)
	}

	wg.Wait() // Wait for all goroutines to finish
	return nil
}

func (l *MailstackProvisionBuyRequestListener) processMailboxes(ctx context.Context, mailstackBuyRequest *postgres_entity.MailstackBuyRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionBuyRequestListener.processMailboxes")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	domains, err := l.dependencies.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // Semaphore to limit to 4 goroutines

	for _, domain := range domains {
		if domain.Status != postgres_entity.MailstackBuyRequestDomainStatusCompleted {
			continue
		}

		wg.Add(1)

		// Create a new scope to avoid data races on the `domain` variable
		go func(domain *postgres_entity.MailstackBuyRequestDomain) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			mailboxes, err := l.dependencies.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain.Domain)
			if err != nil {
				tracing.TraceErr(span, err)
				return // Exit the goroutine on error
			}

			for _, mailbox := range mailboxes {
				if mailbox.Status != postgres_entity.MailboxStatusPendingProvisioning {
					continue
				}

				err := l.dependencies.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, mailbox.ID, common_model.MAILBOX, dto.MailstackProvisionMailbox{})
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			}
		}(domain)
	}

	wg.Wait() // Wait for all goroutines to finish
	return nil
}
