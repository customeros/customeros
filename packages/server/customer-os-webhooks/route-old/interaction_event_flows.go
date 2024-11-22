package route

import (
	"bytes"
	"encoding/json"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func (h *EmailWebhookHandler) processEmailForFlows(ctx context.Context, tenant string, data *model.PostmarkEmailWebhookData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.processEmailForFlows")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	senderUser, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(
		ctx, tenant, data.FromFull.Email,
	)
	if err != nil {
		return err
	}

	if senderUser == nil {
		return nil
	}

	participants := h.extractParticipants(tenant, data)
	for _, email := range participants {
		if err := h.updateFlowStatus(ctx, tenant, email, data.Subject); err != nil {
			return err
		}
	}

	return nil
}

func (h *EmailWebhookHandler) updateFlowStatus(ctx context.Context, tenant, email, subject string) error {
	contacts, err := h.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactsWithEmail(
		ctx, tenant, email,
	)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		contactEntity := mapper.MapDbNodeToContactEntity(contact)
		if err := h.updateContactFlowStatus(ctx, contactEntity, subject); err != nil {
			return err
		}
	}

	return nil
}

func (h *EmailWebhookHandler) updateContactFlowStatus(ctx context.Context, contact *entity.Contact, subject string) error {
	flows, err := h.services.CommonServices.FlowService.FlowsGetListWithParticipant(
		ctx, []string{contact.Id}, commonModel.CONTACT,
	)
	if err != nil {
		return err
	}

	for _, flow := range *flows {
		participant, err := h.services.CommonServices.FlowService.FlowParticipantByEntity(
			ctx, flow.Id, contact.Id, commonModel.CONTACT,
		)
		if err != nil {
			return err
		}

		if participant == nil || h.isCompletedStatus(participant.Status) {
			continue
		}

		newStatus := h.determineNewStatus(subject)
		participant.Status = newStatus

		_, err = h.services.CommonServices.Neo4jRepositories.FlowParticipantWriteRepository.Merge(
			ctx, nil, participant,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
