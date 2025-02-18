package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	mailService interfaces.MailService
}

func NewClassifyEmailCapability(postgresRepositories *postgres_repository.Repositories, mailService interfaces.MailService) *ClassifyEmailCapability {
	return &ClassifyEmailCapability{
		postgres:    postgresRepositories,
		mailService: mailService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ClassifyEmailInput, ClassifyEmailOutput, postgres_entity.NoConfig] = (*ClassifyEmailCapability)(nil)
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

func (c *ClassifyEmailCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *ClassifyEmailCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ClassifyEmailCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *ClassifyEmailCapability) ValidateInput(input ClassifyEmailInput) error {
	if input.RawEmailId == "" {
		return errors.New("raw email id is empty")
	}
	return nil
}

func (c *ClassifyEmailCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ClassifyEmailInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, ClassifyEmailOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ClassifyEmailCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	ingestEmailMessage, err := c.postgres.IngestEmailMessageRepository.GetEmail(ctx, executionContainer.InputData.RawEmailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	emailMessageData, err := c.mailService.LoadIngestEmailMessage(ctx, ingestEmailMessage)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to load email"))
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	if emailMessageData.Identifiers.ProviderMessageId == "" {
		err := errors.New("provider message id is empty")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	if len(emailMessageData.Participants.AllEmails) == 0 {
		err := errors.New("no email participants")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, ClassifyEmailOutput{}, err
	}

	check := c.mailService.ProcessEmailCheck(ctx, ingestEmailMessage.Tenant, &emailMessageData)
	if !check.ProcessEmail {
		// TODO: trigger goal achieved = false

		// set all bounced emails to undeliverable
		//for _, e := range check.BouncedEmails {
		//	err := s.neo4j.EmailWriteRepository.SetDeliverableByEmailForAllTenants(ctx, e, "false")
		//	if err != nil {
		//		tracing.TraceErr(span, errors.Wrap(err, "failed to set deliverable by email for all tenants"))
		//	}
		//}
		return enum.CapabilityExecutionStop, ClassifyEmailOutput{}, nil
	}

	return enum.CapabilityExecutionCompleted, ClassifyEmailOutput{
		EntityId:   ingestEmailMessage.Id,
		EntityType: model.INGEST_EMAIL_MESSAGE,
	}, nil
}
