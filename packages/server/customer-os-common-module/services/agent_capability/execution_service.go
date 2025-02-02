package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentCapabilityExecutionService struct{}

func NewAgentCapabilityExecutionService() interfaces.AgentCapabilityExecutionService {
	return &agentCapabilityExecutionService{}
}

func (f *agentCapabilityExecutionService) Execute(
	ctx context.Context,
	capability postgres_entity.Capability,
	params map[string]any,
	executors map[enum.AgentCapability]interfaces.AgentCapabilityUntyped,
) (map[string]any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentCapabilityExecutionService.Execute")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch capability.Type {
	case enum.CapabilityAnalyzeWebSessionIntent:
		executor, ok := GetTypedExecutor[AnalyzeWebSessionInput, AnalyzeWebSessionOutput, NoConfig](
			executors,
			capability.Type,
		)
		if !ok {
			err := fmt.Errorf("failed to get typed executor for AnalyzeWebSession")
			tracing.TraceErr(span, err)
			return nil, err
		}
		return executeCapability(ctx, capability, params, executor)

	default:

		return nil, nil
	}
}

func executeCapability[I, O, C any](
	ctx context.Context,
	capability postgres_entity.Capability,
	allParams map[string]any,
	cap interfaces.AgentCapability[I, O, C],
) (map[string]any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentCapabilityExecutionService.executeCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	input := cap.GetInput()
	if err := utils.MapToStruct(allParams, &input); err != nil {
		return nil, err
	}

	config := cap.GetConfig()
	if capability.Config != "" {
		if err := json.Unmarshal([]byte(capability.Config), &config); err != nil {
			return nil, err
		}
	}

	output, err := cap.Execute(ctx, input, config)
	if err != nil {
		return nil, err
	}

	return utils.StructToMap(output)
}
