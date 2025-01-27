package agent_capability

import (
	"context"
	"errors"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type ApplyTagCapability struct {
	tagService interfaces.TagService
}

func (c *ApplyTagCapability) ValidateConfig() error {
	//TODO implement me
	panic("implement me")
}

func (c *ApplyTagCapability) ValidateInput() error {
	//TODO implement me
	panic("implement me")
}

func (c *ApplyTagCapability) GetInput() any {
	return &ApplyTagInput{}
}

func (c *ApplyTagCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *ApplyTagCapability) GetOutput() any {
	return &ApplyTagResult{}
}

func NewApplyTagCapability(tagService interfaces.TagService) *ApplyTagCapability {
	return &ApplyTagCapability{
		tagService: tagService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[ApplyTagInput, ApplyTagResult, NoConfig] = (*ApplyTagCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                            = (*ApplyTagCapability)(nil)
)

type ApplyTagInput struct {
	EntityType model.EntityType `json:"entityType"`
	EntityID   string           `json:"entityId"`
	TagID      string           `json:"tagId"`
}

type ApplyTagResult struct {
	Success bool `json:"success"`
}

func (c *ApplyTagCapability) Execute(ctx context.Context, data ApplyTagInput, config NoConfig) (ApplyTagResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApplyTagCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return ApplyTagResult{Success: false}, err
	}

	_, err := c.tagService.AddTagToEntity(ctx, nil, tenant, data.EntityID, data.EntityType, data.TagID, "")
	if err != nil {
		tracing.TraceErr(span, err)
		return ApplyTagResult{Success: false}, err
	}

	return ApplyTagResult{Success: true}, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *ApplyTagCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*ApplyTagInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected ApplyTagInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected NoCOnfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
