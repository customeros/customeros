package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
)

func (h *EmailWebhookHandler) processEmailFlow(ctx context.Context, span opentracing.Span, tenant string, data *model.PostmarkEmailWebhookData) (string, error) {
	participants := h.extractParticipants(tenant, data)
	username, err := h.resolveMailbox(ctx, span, tenant, participants)
	if err != nil || username == "" {
		return "", errors.New("invalid mailbox")
	}

	messageId, err := h.getMessageId(data)
	if err != nil {
		return "", err
	}

	if err := h.processEmailForFlows(ctx, tenant, data); err != nil {
		return "", err
	}

	if err := h.processMailstackReply(ctx, tenant, data); err != nil {
		return "", err
	}

	if err := h.storeEmail(ctx, tenant, username, messageId, data); err != nil {
		return "", err
	}

	return username, nil
}

func (h *EmailWebhookHandler) extractParticipants(tenant string, data *model.PostmarkEmailWebhookData) []string {
	participants := make([]string, 0)
	participants = append(participants, data.FromFull.Email)

	for _, to := range data.ToFull {
		participants = append(participants, to.Email)
	}

	for _, cc := range data.CcFull {
		participants = append(participants, cc.Email)
	}

	tenantEmail := "bcc@" + strings.ToLower(tenant) + ".customeros.ai"
	for _, bcc := range data.BccFull {
		if bcc.Email != "" && bcc.Email != tenantEmail {
			participants = append(participants, bcc.Email)
		}
	}

	return participants
}

func (h *EmailWebhookHandler) resolveMailbox(ctx context.Context, span opentracing.Span, tenant string, participants []string) (string, error) {
	for _, email := range participants {
		user, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
		if err != nil {
			continue
		}
		if user != nil {
			return email, nil
		}
	}
	return "", errors.New("no valid mailbox found")
}

func (h *EmailWebhookHandler) storeEmail(ctx context.Context, tenant, username, messageId string, data *model.PostmarkEmailWebhookData) error {
	exists, err := h.services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(
		ctx, "mailstack", tenant, username, messageId,
	)
	if err != nil {
		return err
	}

	if !exists {
		emailData, err := h.mapToEmailRawData(tenant, data)
		if err != nil {
			return err
		}

		jsonContent, err := h.marshalEmailData(emailData)
		if err != nil {
			return err
		}

		return h.services.CommonServices.PostgresRepositories.RawEmailRepository.Store(
			ctx, "mailstack", tenant, username, emailData.ProviderMessageId,
			messageId, string(jsonContent), emailData.Sent, entity.REAL_TIME,
		)
	}

	return nil
}

func (h *EmailWebhookHandler) marshalEmailData(data entity.EmailRawData) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	return buffer.Bytes(), encoder.Encode(data)
}
