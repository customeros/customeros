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
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type ApplyTagCapability struct {
	tagService interfaces.TagService
}

type ApplyTagInput struct {
	OrganizationID string `json:"organizationId"`
}

type ApplyTagConfig struct {
	TagName ConfigSingleValue `json:"tagName"`
}

func NewApplyTagCapability(tagService interfaces.TagService) *ApplyTagCapability {
	return &ApplyTagCapability{
		tagService: tagService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ApplyTagInput, NoOutput, ApplyTagConfig] = (*ApplyTagCapability)(nil)
)

func (c *ApplyTagCapability) Type() enum.AgentCapability {
	return enum.CapabilityApplyTag
}

func (c *ApplyTagCapability) Name() string {
	return "Apply a tag"
}

func (c *ApplyTagCapability) NewInput() ApplyTagInput {
	return ApplyTagInput{}
}

func (c *ApplyTagCapability) NewConfig() ApplyTagConfig {
	return ApplyTagConfig{}
}

func (c *ApplyTagCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.TagName.Value = "Support"
	return &config
}

func (c *ApplyTagCapability) ValidateConfig(config ApplyTagConfig) error {
	if config.TagName.Value == "" {
		return errors.New("Tag not configured")
	}
	return nil
}

func (c *ApplyTagCapability) ValidateInput(input ApplyTagInput) error {
	if input.OrganizationID == "" {
		return errors.New("OrganizationID cannot be empty")
	}
	return nil
}

func (c *ApplyTagCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ApplyTagInput, ApplyTagConfig]) (bool, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := NoOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
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
		return true, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}

func (c *ApplyTagCapability) applyTag(ctx context.Context, entityType model.EntityType, entityId string, tagName string) error {
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
