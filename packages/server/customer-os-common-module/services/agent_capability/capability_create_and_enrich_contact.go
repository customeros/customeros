package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateContactCapability struct {
	contactService interfaces.ContactService
}

func NewCreateContactCapability(contactService interfaces.ContactService) *CreateContactCapability {
	return &CreateContactCapability{
		contactService: contactService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[CreateContactInput, CreateContactOutput, postgres_entity.NoConfig] = (*CreateContactCapability)(nil)
)

func (c *CreateContactCapability) Type() enum.AgentCapability {
	return enum.CapabilityCreateAndEnrichContact
}

func (c *CreateContactCapability) Name() string {
	return "Create and enrich contacts"
}

func (c *CreateContactCapability) NewInput() CreateContactInput {
	return CreateContactInput{}
}

func (c *CreateContactCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *CreateContactCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *CreateContactCapability) DefaultActive() bool {
	return true
}

func (c *CreateContactCapability) ValidateInput(input CreateContactInput) error {
	if len(input.ContactEmails) == 0 {
		return coserrors.ErrCapabilityContactMissing
	}
	return nil
}

func (c *CreateContactCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type CreateContactInput struct {
	ContactEmails []string `json:"contactEmails"`
}

type CreateContactOutput struct {
	ContactIDs []string `json:"contactIds"`
}

func (c *CreateContactCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[CreateContactInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, CreateContactOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateContactCapability.Execute")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := CreateContactOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	for _, email := range executionContainer.InputData.ContactEmails {
		contactExists, id, err := c.contactService.CheckContactExistsWithEmail(ctx, email)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		if contactExists {
			result.ContactIDs = append(result.ContactIDs, id)
			continue
		}
		id, err = c.contactService.CreateContactByEmail(ctx, nil, email)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		result.ContactIDs = append(result.ContactIDs, id)
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
