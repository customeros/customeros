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
	ctx context.Context, executionContainer interfaces.ExecutionContainer,
) (map[string]any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityExecutionService.Execute")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch executionContainer.Capability.Type {
	case enum.CapabilityAnalyzeWebSessionIntent:
		executor, ok := GetTypedExecutor[AnalyzeWebSessionInput, AnalyzeWebSessionOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityAnalyzeWebSessionIntent)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityApplyTag:
		executor, ok := GetTypedExecutor[ApplyTagInput, NoOutput, ApplyTagConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityApplyTag)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityCreateAndEnrichCompany:
		executor, ok := GetTypedExecutor[CreateOrganizationInput, CreateOrganizationOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateAndEnrichCompany)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityCreateMarkdownTimelineEvent:
		executor, ok := GetTypedExecutor[CreateMarkdownTimelineEventInput, CreateMarkdownTimelineEventOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateMarkdownTimelineEvent)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityEvaluateCompanyICPFit:
		executor, ok := GetTypedExecutor[EvaluateICPFitInput, EvaluateICPFitOutput, EvaluateICPFitConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityEvaluateCompanyICPFit)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityGatherCompanyIntelligence:
		executor, ok := GetTypedExecutor[GatherCompanyIntilligenceInput, GatherCompanyIntelligenceOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityEvaluateCompanyICPFit)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityGenerateInvoice:
		executor, ok := GetTypedExecutor[GenerateInvoiceInput, GenerateInvoiceOutput, GenerateInvoiceConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityIdentifyWebVisitor:
		executor, ok := GetTypedExecutor[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorOutput, IdentifyWebsiteVisitorConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIdentifyWebVisitor)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendSlackNotification:
		executor, ok := GetTypedExecutor[SendSlackNotificationInput, NoOutput, SendSlackNotificationConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendWebVisitorSlackNotification:
		executor, ok := GetTypedExecutor[SendWebVisitorSlackNotificationInput, NoOutput, SendWebVisitorSlackNotificationConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityUpdateCompanyStatus:
		executor, ok := GetTypedExecutor[UpdateCompanyStatusInput, NoOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	default:
		err := fmt.Errorf("capability not configured")
		span.LogKV("capability", executionContainer.Capability.Type)
		tracing.TraceErr(span, err)
		return nil, err
	}
}

func (f *agentCapabilityExecutionService) handleGetTypedExecutorError(ctx context.Context, capability enum.AgentCapability) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityExecutionService.handleGetTypedExecutorError")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := fmt.Errorf("failed to get typed executor for %s", capability.String())
	tracing.TraceErr(span, err)
	return err
}

func executeCapability[I, O, C any](
	ctx context.Context,
	cap interfaces.AgentCapability[I, O, C],
	executionContainer interfaces.ExecutionContainer,
) (map[string]any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("AgentCapabilityExecutionService.executeCapability.%s", executionContainer.Capability.Type.String()))
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	input := cap.NewInput()
	err := utils.MapToStruct(executionContainer.ExecutionParams, &input)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	config := cap.NewConfig()
	err = executionContainer.Capability.GetConfig(&config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	typedExecutionContainer := interfaces.TypedExecutionContainer[I, C]{
		AgentExecutionID: executionContainer.AgentExecutionID,
		InputData:        input,
		ConfigData:       config,
	}

	ok, output, err := cap.Execute(ctx, typedExecutionContainer)
	if !ok {
		// todo -- what do we do when input validation does not pass?
		// for when we implement retry logic
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return utils.StructToMap(output)
}
