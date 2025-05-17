package events_listeners

import (
	"context"
	"net/http"
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	common_model "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
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
	spans, ctx := telemetry.StartListenerSpan(ctx, "MailstackProvisionBuyRequestListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
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

		statusCode, errMessage, mailboxes, err := l.dependencies.CommonServices.MailstackService.GetMailboxes(ctx, domain.Tenant, domain.Domain, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if statusCode != http.StatusOK {
			err = errors.New(errMessage)
			tracing.TraceErr(span, err)
			return err
		}

		for _, mailbox := range mailboxes {
			if mailbox.Provisioned {
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

	statusCode, errorMsg, err := l.dependencies.CommonServices.MailstackService.PurchaseDomain(ctx, tenant, domain)
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

			statusCode, errMessage, mailboxes, err := l.dependencies.CommonServices.MailstackService.GetMailboxes(ctx, domain.Tenant, domain.Domain, "")
			if err != nil {
				tracing.TraceErr(span, err)
				return // Exit the goroutine on error
			}

			if statusCode != http.StatusOK {
				err = errors.New(errMessage)
				tracing.TraceErr(span, err)
				return // Exit the goroutine on error
			}

			for _, mailbox := range mailboxes {
				if mailbox.Provisioned {
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
