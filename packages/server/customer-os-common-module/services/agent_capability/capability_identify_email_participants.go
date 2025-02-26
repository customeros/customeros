package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type IdentifyEmailParticipantsInput struct {
	EntityId   string           `json:"entityId"`
	EntityType model.EntityType `json:"entityType"`
}

type IdentifyEmailParticipantsOutput struct {
	OrganizationIds []string `json:"organizationIds"`
}

type IdentifyEmailParticipantsCapability struct {
	postgres    *postgres_repository.Repositories
	mailService interfaces.MailService
}

func NewIdentifyEmailParticipantsCapability(postgresRepositories *postgres_repository.Repositories, mailService interfaces.MailService) *IdentifyEmailParticipantsCapability {
	return &IdentifyEmailParticipantsCapability{
		postgres:    postgresRepositories,
		mailService: mailService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[IdentifyEmailParticipantsInput, IdentifyEmailParticipantsOutput, postgres_entity.NoConfig] = (*IdentifyEmailParticipantsCapability)(nil)
)

func (c *IdentifyEmailParticipantsCapability) Type() enum.AgentCapability {
	return enum.CapabilityIdentifyEmailParticipants
}

func (c *IdentifyEmailParticipantsCapability) Name() string {
	return "Identify participants"
}

func (c *IdentifyEmailParticipantsCapability) NewInput() IdentifyEmailParticipantsInput {
	return IdentifyEmailParticipantsInput{}
}

func (c *IdentifyEmailParticipantsCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *IdentifyEmailParticipantsCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *IdentifyEmailParticipantsCapability) DefaultActive() bool {
	return true
}

func (c *IdentifyEmailParticipantsCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *IdentifyEmailParticipantsCapability) ValidateInput(input IdentifyEmailParticipantsInput) error {
	if input.EntityId == "" {
		return errors.New("entityId is required")
	}
	if input.EntityType == "" {
		return errors.New("entityType is required")
	}
	return nil
}

func (c *IdentifyEmailParticipantsCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[IdentifyEmailParticipantsInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, IdentifyEmailParticipantsOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyEmailParticipantsCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
	}

	ingestEmailMessage, err := c.postgres.IngestEmailMessageRepository.GetEmail(ctx, executionContainer.InputData.EntityId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email"))
		return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
	}

	emailMessageData, err := c.mailService.LoadIngestEmailMessage(ctx, ingestEmailMessage)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to load email"))
		return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
	}

	orgIds := make([]string, 0)

	//process FROM email
	fromOrgId, err := c.mailService.GetOrganizationIdForEmail(ctx, nil, ingestEmailMessage.Tenant, emailMessageData.Participants.From.Email, ingestEmailMessage.Provider)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email id for email"))
		return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
	}
	if fromOrgId != "" {
		orgIds = append(orgIds, fromOrgId)
	}

	//process FROM email
	for _, to := range emailMessageData.Participants.To {
		toOrgId, err := c.mailService.GetEmailIdForEmail(ctx, nil, ingestEmailMessage.Tenant, to.Email, ingestEmailMessage.Provider)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email id for email"))
			return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
		}
		if toOrgId != "" {
			orgIds = utils.AddToListIfNotExists(orgIds, toOrgId)
		}
	}

	//process CC email
	for _, cc := range emailMessageData.Participants.Cc {
		ccOrgId, err := c.mailService.GetEmailIdForEmail(ctx, nil, ingestEmailMessage.Tenant, cc.Email, ingestEmailMessage.Provider)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email id for email"))
			return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
		}
		if ccOrgId != "" {
			orgIds = utils.AddToListIfNotExists(orgIds, ccOrgId)
		}
	}

	//process BCC email
	for _, bcc := range emailMessageData.Participants.Bcc {
		bccOrgId, err := c.mailService.GetEmailIdForEmail(ctx, nil, ingestEmailMessage.Tenant, bcc.Email, ingestEmailMessage.Provider)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email id for email"))
			return enum.CapabilityExecutionError, IdentifyEmailParticipantsOutput{}, err
		}
		if bccOrgId != "" {
			orgIds = utils.AddToListIfNotExists(orgIds, bccOrgId)
		}
	}

	return enum.CapabilityExecutionCompleted, IdentifyEmailParticipantsOutput{
		OrganizationIds: orgIds,
	}, nil
}
