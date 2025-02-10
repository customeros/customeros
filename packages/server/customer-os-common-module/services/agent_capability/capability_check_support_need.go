package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CheckSupportNeedCapability struct{}

type CheckSupportNeedInput struct {
	OrganizationID string `json:"organizationId"`
}

type CheckSupportNeedConfig struct {
	TagName ConfigSingleValue `json:"tagName"`
}

func NewCheckSupporNeedCapability() *CheckSupportNeedCapability {
	return &CheckSupportNeedCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[CheckSupportNeedInput, NoOutput, CheckSupportNeedConfig] = (*CheckSupportNeedCapability)(nil)
)

func (c *CheckSupportNeedCapability) Type() enum.AgentCapability {
	return enum.CapabilityCheckSupportNeed
}

func (c *CheckSupportNeedCapability) Name() string {
	return "Check support need"
}

func (c *CheckSupportNeedCapability) NewInput() CheckSupportNeedInput {
	return CheckSupportNeedInput{}
}

func (c *CheckSupportNeedCapability) NewConfig() CheckSupportNeedConfig {
	return CheckSupportNeedConfig{}
}

func (c *CheckSupportNeedCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *CheckSupportNeedCapability) ValidateConfig(config CheckSupportNeedConfig) error {
	return nil
}

func (c *CheckSupportNeedCapability) ValidateInput(input CheckSupportNeedInput) error {
	if input.OrganizationID == "" {
		return errors.New("OrganizationID cannot be empty")
	}
	return nil
}

func (c *CheckSupportNeedCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[CheckSupportNeedInput, CheckSupportNeedConfig]) (bool, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CheckSupportNeedCapability.Execute")
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

	if executionContainer.InputData.OrganizationID != "" {
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}
