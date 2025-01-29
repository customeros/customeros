package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type ApplyTagCapability struct {
	tagService interfaces.TagService
}

type ApplyTagInput struct {
	OrganizationID string `json:"organizationId"`
}

type ApplyTagOutput struct {
	CapabilityOutput
}

type ApplyTagConfig struct {
	TagName TagConfig `json:"tagName"`
}

type TagConfig struct {
	Value string `json:"value"`
	Error string `json:"error"`
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

func (c *ApplyTagCapability) GetInput() any {
	return &ApplyTagInput{}
}

func (c *ApplyTagCapability) GetConfig() any {
	return &ApplyTagConfig{}
}

func (c *ApplyTagCapability) GetOutput() any {
	return &ApplyTagOutput{}
}

func NewApplyTagCapability(tagService interfaces.TagService) *ApplyTagCapability {
	return &ApplyTagCapability{
		tagService: tagService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[ApplyTagInput, ApplyTagOutput, ApplyTagConfig] = (*ApplyTagCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                  = (*ApplyTagCapability)(nil)
)

func (c *ApplyTagCapability) Execute(ctx context.Context, data ApplyTagInput, config ApplyTagConfig) (ApplyTagOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := ApplyTagOutput{}

	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}
	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}

	result.ExecutionValidated = true

	var entityType model.EntityType
	var entityID string
	if data.OrganizationID != "" {
		entityType = model.ORGANIZATION
		entityID = data.OrganizationID
	}

	err := c.applyTag(ctx, entityType, entityID, config.TagName.Value)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

func (c *ApplyTagCapability) applyTag(ctx context.Context, entityType model.EntityType, entityId string, tagName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.applyTag")
	defer span.Finish()

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

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *ApplyTagCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*ApplyTagInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected ApplyTagInput")
	}

	typedConfig, ok := config.(*ApplyTagConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected ApplyTagConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
