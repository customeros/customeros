package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type IngestEmailInput struct {
	EntityId        string           `json:"entityId"`
	EntityType      model.EntityType `json:"entityType"`
	OrganizationIds []string         `json:"organizationIds"`
}

type IngestEmailCapability struct {
	postgres    *postgres_repository.Repositories
	events      *events.EventsService
	mailService interfaces.MailService
	opensearch  interfaces.OpensearchService
}

func NewIngestEmailCapability(postgresRepositories *postgres_repository.Repositories, events *events.EventsService, mailService interfaces.MailService, opensearch interfaces.OpensearchService) *IngestEmailCapability {
	return &IngestEmailCapability{
		postgres:    postgresRepositories,
		events:      events,
		mailService: mailService,
		opensearch:  opensearch,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[IngestEmailInput, NoOutput, postgres_entity.NoConfig] = (*IngestEmailCapability)(nil)
)

func (c *IngestEmailCapability) Type() enum.AgentCapability {
	return enum.CapabilityIngestEmail
}

func (c *IngestEmailCapability) Name() string {
	return "Ingest email"
}

func (c *IngestEmailCapability) NewInput() IngestEmailInput {
	return IngestEmailInput{}
}

func (c *IngestEmailCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *IngestEmailCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *IngestEmailCapability) DefaultActive() bool {
	return true
}

func (c *IngestEmailCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *IngestEmailCapability) ValidateInput(input IngestEmailInput) error {
	if input.EntityId == "" {
		return errors.New("entityId is required")
	}
	if input.EntityType == "" {
		return errors.New("entityType is required")
	}
	if len(input.OrganizationIds) == 0 {
		return errors.New("organizationIds is required")
	}
	return nil
}

func (c *IngestEmailCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[IngestEmailInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	ingestEmailMessage, err := c.postgres.IngestEmailMessageRepository.GetEmail(ctx, executionContainer.InputData.EntityId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if ingestEmailMessage == nil {
		tracing.TraceErr(span, errors.Wrap(err, "email is nil"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	emailMessageData, err := c.mailService.LoadIngestEmailMessage(ctx, ingestEmailMessage)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to load email"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	for _, organizationId := range executionContainer.InputData.OrganizationIds {

		osTo := []map[string]string{}
		for _, to := range emailMessageData.Participants.To {
			osTo = append(osTo, map[string]string{
				"email":     to.Email,
				"firstName": to.FirstName,
				"lastName":  to.LastName,
			})
		}

		osCc := []map[string]string{}
		for _, cc := range emailMessageData.Participants.Cc {
			osCc = append(osCc, map[string]string{
				"email":     cc.Email,
				"firstName": cc.FirstName,
				"lastName":  cc.LastName,
			})
		}

		osBcc := []map[string]string{}
		for _, bcc := range emailMessageData.Participants.Bcc {
			osBcc = append(osBcc, map[string]string{
				"email":     bcc.Email,
				"firstName": bcc.FirstName,
				"lastName":  bcc.LastName,
			})
		}

		osData := map[string]interface{}{
			"id":             ingestEmailMessage.Id,
			"organizationId": organizationId,
			"from": map[string]string{
				"email":     emailMessageData.Participants.From.Email,
				"firstName": emailMessageData.Participants.From.FirstName,
				"lastName":  emailMessageData.Participants.From.LastName,
			},
			"to":         osTo,
			"cc":         osCc,
			"bcc":        osBcc,
			"subject":    emailMessageData.Content.Subject,
			"sentDate":   emailMessageData.Content.SentDate,
			"html":       emailMessageData.Content.Html,
			"text":       emailMessageData.Content.Text,
			"providerId": emailMessageData.Identifiers.ProviderMessageId,
			"threadId":   emailMessageData.Identifiers.EmailThreadId,
		}

		indexName := "email-" + ingestEmailMessage.SentAt.Format("2006-01")
		err := c.opensearch.UpsertDocument(ctx, indexName, &ingestEmailMessage.Id, osData)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to index document"))
		}
	}

	err = c.postgres.AgentExecutionRepository.GoalAchieved(ctx, executionContainer.AgentExecutionID, true, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish ingest email event"))
	}

	return enum.CapabilityExecutionCompleted, NoOutput{}, nil
}
