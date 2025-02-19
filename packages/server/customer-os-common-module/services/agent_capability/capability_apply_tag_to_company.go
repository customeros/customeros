package agent_capability

import (
	"context"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type ApplyTagToCompanyCapability struct {
	events     *events.EventsService
	tagService interfaces.TagService
}

type ApplyTagToCompanyInput struct {
	OrganizationID string `json:"organizationId"`
}

type ApplyTagToCompanyConfig struct {
	TagName ConfigSingleValue `json:"tagName"`
}

func NewApplyTagToCompanyCapability(tagService interfaces.TagService, events *events.EventsService) *ApplyTagToCompanyCapability {
	return &ApplyTagToCompanyCapability{
		events:     events,
		tagService: tagService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ApplyTagToCompanyInput, NoOutput, ApplyTagToCompanyConfig] = (*ApplyTagToCompanyCapability)(nil)
)

func (c *ApplyTagToCompanyCapability) Type() enum.AgentCapability {
	return enum.CapabilityApplyTagToCompany
}

func (c *ApplyTagToCompanyCapability) Name() string {
	return "Apply a tag to company"
}

func (c *ApplyTagToCompanyCapability) NewInput() ApplyTagToCompanyInput {
	return ApplyTagToCompanyInput{}
}

func (c *ApplyTagToCompanyCapability) NewConfig() ApplyTagToCompanyConfig {
	return ApplyTagToCompanyConfig{}
}

func (c *ApplyTagToCompanyCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.TagName.Value = "Support"
	return &config
}

func (c *ApplyTagToCompanyCapability) ValidateConfig(config ApplyTagToCompanyConfig) error {
	if config.TagName.Value == "" {
		return errors.New("Tag not configured")
	}
	return nil
}

func (c *ApplyTagToCompanyCapability) ValidateInput(input ApplyTagToCompanyInput) error {
	if input.OrganizationID == "" {
		return errors.New("OrganizationID cannot be empty")
	}
	return nil
}

func (c *ApplyTagToCompanyCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ApplyTagToCompanyInput, ApplyTagToCompanyConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := NoOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	var entityType model.EntityType
	var entityID string
	if executionContainer.InputData.OrganizationID != "" {
		entityType = model.ORGANIZATION
		entityID = executionContainer.InputData.OrganizationID
	}

	err := c.applyTag(ctx, entityType, entityID, executionContainer.ConfigData.TagName.Value)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *ApplyTagToCompanyCapability) applyTag(ctx context.Context, entityType model.EntityType, entityId string, tagName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.applyTag")
	defer span.Finish()
	tracing.TagComponentService(span)

	// check if tag exists
	tagEntity, err := c.tagService.GetTagByEntityTypeAndName(ctx, entityType, tagName)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if tagEntity == nil {
		// create tag
		input := neo4jentity.TagEntity{
			Name:       tagName,
			EntityType: entityType,
		}
		tagEntity, err = c.tagService.Save(ctx, nil, &input)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	_, err = c.tagService.AddTagToEntity(ctx, nil, common.GetTenantFromContext(ctx), entityId, entityType, tagEntity.Id, "")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
