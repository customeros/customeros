package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"strings"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type ClassifyEmailInput struct {
	RawEmailId string `json:"rawEmailId"`
}

type ClassifyEmailOutput struct {
	EntityId   string           `json:"entityId"`
	EntityType model.EntityType `json:"entityType"`
}

type ClassifyEmailCapability struct {
	postgres    *postgres_repository.Repositories
	events      *events.EventsService
	mailService interfaces.MailService
}

type ClassifyEmailConfig struct {
	ImportEmails ConfigSingleValue `json:"importEmails"`
}

func NewClassifyEmailCapability(postgresRepositories *postgres_repository.Repositories, events *events.EventsService, mailService interfaces.MailService) *ClassifyEmailCapability {
	return &ClassifyEmailCapability{
		postgres:    postgresRepositories,
		events:      events,
		mailService: mailService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ClassifyEmailInput, ClassifyEmailOutput, ClassifyEmailConfig] = (*ClassifyEmailCapability)(nil)
)

func (c *ClassifyEmailCapability) Type() enum.AgentCapability {
	return enum.CapabilityClassifyEmail
}

func (c *ClassifyEmailCapability) Name() string {
	return "Classify email"
}

func (c *ClassifyEmailCapability) NewInput() ClassifyEmailInput {
	return ClassifyEmailInput{}
}

func (c *ClassifyEmailCapability) NewConfig() ClassifyEmailConfig {
	return ClassifyEmailConfig{}
}

func (c *ClassifyEmailCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.ImportEmails.Value = "AUTOMATICALLY"
	return &config
}

func (c *ClassifyEmailCapability) DefaultActive() bool {
	return true
}

func (c *ClassifyEmailCapability) ValidateConfig(config ClassifyEmailConfig) error {
	if config.ImportEmails.Value != "AUTOMATICALLY" && config.ImportEmails.Value != "MANUALLY" {
		return errors.New("import emails is invalid")
	}
	return nil
}

func (c *ClassifyEmailCapability) ValidateInput(input ClassifyEmailInput) error {
	if input.RawEmailId == "" {
		return errors.New("raw email id is empty")
	}
	return nil
}

func (c *ClassifyEmailCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ClassifyEmailInput, ClassifyEmailConfig]) (enum.CapabilityExecutionStatus, ClassifyEmailOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ClassifyEmailCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	ingestEmailMessage, err := c.postgres.IngestEmailMessageRepository.GetEmail(ctx, executionContainer.InputData.RawEmailId)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get email"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	emailMessageData, err := c.mailService.LoadIngestEmailMessage(ctx, ingestEmailMessage)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to load email"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	if emailMessageData.Identifiers.ProviderMessageId == "" {
		err := errors.New("provider message id is empty")
		spans.TraceError(err)
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	if len(emailMessageData.Participants.AllEmails) == 0 {
		err := errors.New("no email participants")
		spans.TraceError(err)
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	check := c.mailService.ProcessEmailCheck(ctx, ingestEmailMessage.Tenant, &emailMessageData)
	if !check.ProcessEmail {
		// TODO: trigger goal achieved = false

		// set all bounced emails to undeliverable
		//for _, e := range check.BouncedEmails {
		//	err := s.neo4j.EmailWriteRepository.SetDeliverableByEmailForAllTenants(ctx, e, "false")
		//	if err != nil {
		//		spans.TraceError( errors.Wrap(err, "failed to set deliverable by email for all tenants"))
		//	}
		//}

		err = c.postgres.AgentExecutionRepository.GoalAchieved(ctx, executionContainer.AgentExecutionID, false, nil)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to set goal achieved"))
		}

		return enum.CapabilityExecutionStop, ClassifyEmailOutput{}, nil
	}

	if strings.Contains(strings.ToLower(emailMessageData.Content.Subject), "out of office") ||
		strings.Contains(strings.ToLower(emailMessageData.Content.Subject), "confidential") ||
		strings.Contains(strings.ToLower(emailMessageData.Content.Html), "you are receiving this email because you are subscribed to calendar notifications") {
		return enum.CapabilityExecutionStop, ClassifyEmailOutput{}, nil
	}

	return enum.CapabilityExecutionCompleted, ClassifyEmailOutput{
		EntityId:   ingestEmailMessage.Id,
		EntityType: model.INGEST_EMAIL_MESSAGE,
	}, nil
}
