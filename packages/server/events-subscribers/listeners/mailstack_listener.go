package listeners

import (
	"context"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"strings"
)

func Handle_MailstackBuyRequest(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.Handle_MailstackBuyRequest")
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

	//buy domains
	domains, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, domain := range domains {
		if domain.Status != entity.MailstackBuyRequestDomainStatusPending {
			continue
		}

		err := services.NamecheapService.PurchaseDomain(ctx, domain.Tenant, domain.Domain)
		if err != nil {
			tracing.TraceErr(span, err)

			mailstackBuyRequest.Status = entity.MailstackBuyRequestStatusFailed
			domain.Status = entity.MailstackBuyRequestDomainStatusFailed
		} else {
			domain.Status = entity.MailstackBuyRequestDomainStatusCompleted
		}

		err = services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, domain.Tenant, entity.MailstackBuyRequestDomain{}, domain.ID, "Status", domain.Status)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	//buy mailboxes
	mailboxes, err := services.PostgresRepositories.MailstackBuyRequestRepository.GetMailboxes(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//reload latest state for domains
	domains, err = services.PostgresRepositories.MailstackBuyRequestRepository.GetDomains(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, mailbox := range mailboxes {
		if mailbox.Status != entity.MailstackBuyRequestDomainMailboxPending {
			continue
		}

		err := services.MailboxService.AddMailbox(ctx, service.AddMailboxRequest{
			Domain:         mailbox.Domain,
			Username:       mailbox.Username,
			Password:       utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false),
			WebmailEnabled: true,
			ForwardingTo:   []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(mailbox.Tenant))},
		})
		if err != nil && err != coserrors.ErrMailboxExists {
			tracing.TraceErr(span, err)

			mailbox.Status = entity.MailstackBuyRequestDomainMailboxFailed
		} else {
			mailbox.Status = entity.MailstackBuyRequestDomainMailboxCompleted
		}

		err = services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, mailbox.Tenant, entity.MailstackBuyRequestMailbox{}, mailbox.ID, "Status", mailbox.Status)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	//reload latest state and update request status
	mailboxes, err = services.PostgresRepositories.MailstackBuyRequestRepository.GetMailboxes(ctx, mailstackBuyRequest.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	hasFailures := false
	allMailboxesCompleted := true
	for _, mailbox := range mailboxes {
		if mailbox.Status == entity.MailstackBuyRequestDomainMailboxFailed {
			hasFailures = true
			break
		}
		if mailbox.Status != entity.MailstackBuyRequestDomainMailboxCompleted {
			allMailboxesCompleted = false
			break
		}
	}

	if hasFailures {
		err := services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, mailstackBuyRequest.Tenant, entity.MailstackBuyRequest{}, mailstackBuyRequest.ID, "Status", entity.MailstackBuyRequestStatusFailed)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if !hasFailures && allMailboxesCompleted {
		err := services.PostgresRepositories.CommonRepository.UpdateProperty(ctx, mailstackBuyRequest.Tenant, entity.MailstackBuyRequest{}, mailstackBuyRequest.ID, "Status", entity.MailstackBuyRequestStatusCompleted)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
