package agent_capability

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type ApplyTagCapability struct {
	tagService interfaces.TagService
}

func NewApplyTagCapability(tagService interfaces.TagService) *ApplyTagCapability {
	return &ApplyTagCapability{
		tagService: tagService,
	}
}

// Compile-time interface check
var _ interfaces.AgentCapabilityExecution[ApplyTagInput, ApplyTagResult] = (*ApplyTagCapability)(nil)

type ApplyTagInput struct {
	EntityType model.EntityType
	EntityID   string
	TagID      string
}

type ApplyTagResult struct {
	Success bool
}

func (c *ApplyTagCapability) Execute(ctx context.Context, data ApplyTagInput) (ApplyTagResult, error) {
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
