package listeners

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"strings"
	"sync"
)

func Handle_MailstackProvisionBuyRequest(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_MailstackProvisionBuyRequest")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)

	mailstackBuyRequest, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetById(ctx, message.Event.EntityId)
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

	if mailstackBuyRequest.Status != entity.MailstackBuyRequestStatusPending {
		return nil
	}

	err = processDomains(ctx, services, mailstackBuyRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = processMailboxes(ctx, services, mailstackBuyRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//reload latest state for domains
	domains, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, domain := range domains {
		if domain.Status != entity.MailstackBuyRequestDomainStatusCompleted {
			continue
		}

		mailboxes, err := services.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain.Domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		for _, mailbox := range mailboxes {
			if mailbox.Status != entity.MailboxStatusPendingProvisioning {
				continue
			}

			err := services.RabbitMQService.PublishEvent(ctx, mailbox.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	// mark buy request as completed or failed
	if mailstackBuyRequest.Status == entity.MailstackBuyRequestStatusPending {

		completed := true

		for _, domain := range domains {
			if domain.Status != entity.MailstackBuyRequestDomainStatusCompleted {
				completed = false
			}
		}

		if completed {
			mailstackBuyRequest.Status = entity.MailstackBuyRequestStatusCompleted
		}

		_, err = services.PostgresRepositories.MailstackBuyRequestRepository.Store(ctx, nil, mailstackBuyRequest)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return nil
}

func processDomains(ctx context.Context, services *service.Services, mailstackBuyRequest *entity.MailstackBuyRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_MailstackProvisionBuyRequest.processDomains")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	domains, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Semaphore to limit to 5 goroutines

	for _, domain := range domains {
		wg.Add(1)

		// Create a new scope to avoid data races on the `domain` variable
		go func(domain *entity.MailstackBuyRequestDomain) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Step 1 - Purchase domain in Namecheap
			if domain.Status == entity.MailstackBuyRequestDomainStatusPendingProvisioning {
				err := services.NamecheapService.PurchaseDomain(ctx, domain.Tenant, domain.Domain)
				if err != nil {
					tracing.TraceErr(span, err)
					mailstackBuyRequest.Status = entity.MailstackBuyRequestStatusFailed
					domain.Status = entity.MailstackBuyRequestDomainStatusFailed
				} else {
					domain.Status = entity.MailstackBuyRequestDomainStatusPendingConfiguration
				}

				err = services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, domain.Tenant, entity.MailstackBuyRequestDomain{}, domain.ID, "Status", domain.Status)
				if err != nil {
					tracing.TraceErr(span, err)
					return // Exit the goroutine on error
				}
			}

			// Step 2 - Configure domain in Mailstack
			if domain.Status == entity.MailstackBuyRequestDomainStatusPendingConfiguration {
				err := services.MailstackService.ConfigureMailstackDomain(ctx, domain.Domain, domain.RedirectWebsite)
				if err != nil {
					tracing.TraceErr(span, err)
				} else {
					domain.Status = entity.MailstackBuyRequestDomainStatusCompleted
				}

				err = services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, domain.Tenant, entity.MailstackBuyRequestDomain{}, domain.ID, "Status", domain.Status)
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

func processMailboxes(ctx context.Context, services *service.Services, mailstackBuyRequest *entity.MailstackBuyRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_MailstackProvisionBuyRequest.processMailboxes")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	domains, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // Semaphore to limit to 4 goroutines

	for _, domain := range domains {
		if domain.Status != entity.MailstackBuyRequestDomainStatusCompleted {
			continue
		}

		wg.Add(1)

		// Create a new scope to avoid data races on the `domain` variable
		go func(domain *entity.MailstackBuyRequestDomain) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			mailboxes, err := services.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain.Domain)
			if err != nil {
				tracing.TraceErr(span, err)
				return // Exit the goroutine on error
			}

			for _, mailbox := range mailboxes {
				if mailbox.Status != entity.MailboxStatusPendingProvisioning {
					continue
				}

				err := services.RabbitMQService.PublishEvent(ctx, mailbox.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
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

func Handle_MailstackProvisionMailbox(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_MailstackProvisionMailbox")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)

	mailbox, err := services.PostgresRepositories.TenantSettingsMailboxRepository.GetById(ctx, message.Event.EntityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = services.OpenSrsService.SetupMailbox(ctx, mailbox.Tenant, mailbox.MailboxUsername, mailbox.MailboxPassword, strings.Split(mailbox.ForwardingTo, ","), mailbox.WebmailEnabled)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	mailbox.Status = entity.MailboxStatusProvisioned

	err = services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, mailbox.Tenant, entity.TenantSettingsMailbox{}, mailbox.ID, "Status", mailbox.Status)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}
