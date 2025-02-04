package agent_capability

import (
	"context"
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
		executor, ok := GetTypedExecutor[AnalyzeWebSessionInput, AnalyzeWebSessionOutput, NoConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityAnalyzeWebSessionIntent)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityApplyTag:
		executor, ok := GetTypedExecutor[ApplyTagInput, ApplyTagOutput, ApplyTagConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityApplyTag)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityCreateAndEnrichCompany:
		executor, ok := GetTypedExecutor[CreateOrganizationInput, CreateOrganizationOutput, NoConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateAndEnrichCompany)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityCreateMarkdownTimelineEvent:
		executor, ok := GetTypedExecutor[CreateMarkdownTimelineEventInput, CreateMarkdownTimelineEventOutput, NoConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateMarkdownTimelineEvent)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityEvaluateCompanyICPFit:
		executor, ok := GetTypedExecutor[ICPQualificationInput, ICPQualificationOutput, ICPQualificationConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityEvaluateCompanyICPFit)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityIdentifyWebVisitor:
		executor, ok := GetTypedExecutor[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorOutput, IdentifyWebsiteVisitorConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIdentifyWebVisitor)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilitySendSlackNotification:
		executor, ok := GetTypedExecutor[SendSlackNotificationInput, SendSlackNotificationOutput, SendSlackNotificationConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendSlackNotification)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilitySendWebVisitorSlackNotification:
		executor, ok := GetTypedExecutor[SendWebVisitorSlackNotificationInput, SendWebVisitorSlackNotificationOutput, SendWebVisitorSlackNotificationConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, capability, params, executor)

	case enum.CapabilityGenerateInvoice:
		executor, ok := GetTypedExecutor[GenerateInvoiceInput, GenerateInvoiceOutput, GenerateInvoiceConfig](executors, capability.Type)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, capability, params, executor)

	default:
		err := fmt.Errorf("capability not configured")
		span.LogKV("capability", capability.Type)
		tracing.TraceErr(span, err)
		return nil, err
	}
}

func (f *agentCapabilityExecutionService) handleGetTypedExecutorError(ctx context.Context, capability enum.AgentCapability) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "agentCapabilityExecutionService.handleGetTypedExecutorError")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := fmt.Errorf("failed to get typed executor for %s", capability.String())
	tracing.TraceErr(span, err)
	return err
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

	input := cap.NewInput()
	err := utils.MapToStruct(allParams, &input)
	if err != nil {
		return nil, err
	}

	config := cap.NewConfig()
	err = capability.GetConfig(&config)
	if err != nil {
		return nil, err
	}

	output, err := cap.Execute(ctx, input, config)
	if err != nil {
		return nil, err
	}

	return utils.StructToMap(output)
}
