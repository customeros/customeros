package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentCapabilityExecutionService struct{}

func NewAgentCapabilityExecutionService() interfaces.AgentCapabilityExecutionService {
	return &agentCapabilityExecutionService{}
}

func (f *agentCapabilityExecutionService) Execute(
	ctx context.Context, executionContainer interfaces.ExecutionContainer,
) (enum.CapabilityExecutionStatus, map[string]any, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AgentCapabilityExecutionService.Execute")
	defer spans.Finish()

	spans.LogKV("capabilityType", executionContainer.Capability.Type.String())

	switch executionContainer.Capability.Type {

	case enum.CapabilityAddMeetingNotesToCompany:
		executor, ok := GetTypedExecutor[AddMeetingNotesToCompanyInput, AddMeetingNotesToCompanyOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityAddMeetingNotesToCompany)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityAnalyzeWebSessionIntent:
		executor, ok := GetTypedExecutor[AnalyzeWebSessionInput, AnalyzeWebSessionOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityAnalyzeWebSessionIntent)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityApplyTagToCompany:
		executor, ok := GetTypedExecutor[ApplyTagToCompanyInput, NoOutput, ApplyTagToCompanyConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityApplyTagToCompany)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityCreateAndEnrichCompany:
		executor, ok := GetTypedExecutor[CreateOrganizationInput, CreateOrganizationOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateAndEnrichCompany)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityCreateAndEnrichContact:
		executor, ok := GetTypedExecutor[CreateContactInput, CreateContactOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateAndEnrichContact)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityCreateMarkdownTimelineEvent:
		executor, ok := GetTypedExecutor[CreateMarkdownTimelineEventInput, CreateMarkdownTimelineEventOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityCreateMarkdownTimelineEvent)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityEvaluateCompanyICPFit:
		executor, ok := GetTypedExecutor[EvaluateICPFitInput, EvaluateICPFitOutput, EvaluateICPFitConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityEvaluateCompanyICPFit)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityExtractMeetingHighlights:
		executor, ok := GetTypedExecutor[ExtractMeetingHighlightsInput, ExtractMeetingHighlightsOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityExtractMeetingHighlights)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityExtractSupportSignalsFromMeeting:
		executor, ok := GetTypedExecutor[ExtractSupportSignalsFromMeetingInput, ExtractSupportSignalsFromMeetingOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityExtractSupportSignalsFromMeeting)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityGatherCompanyIntelligence:
		executor, ok := GetTypedExecutor[GatherCompanyIntilligenceInput, GatherCompanyIntelligenceOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityEvaluateCompanyICPFit)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityGenerateInvoice:
		executor, ok := GetTypedExecutor[GenerateInvoiceInput, GenerateInvoiceOutput, GenerateInvoiceConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityProcessAutopayment:
		executor, ok := GetTypedExecutor[ProcessAutopaymentInput, ProcessAutopaymentOutput, ProcessAutopaymentConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityProcessAutopayment)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendInvoiceViaEmail:
		executor, ok := GetTypedExecutor[SendInvoiceViaEmailInput, SendInvoiceViaEmailOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendInvoiceViaEmail)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendPaidNotification:
		executor, ok := GetTypedExecutor[SendPaidNotificationInput, SendPaidNotificationOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendPaidNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendInvoiceVoidedNotification:
		executor, ok := GetTypedExecutor[SendInvoiceVoidedNotificationInput, SendInvoiceVoidedNotificationOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendInvoiceVoidedNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendPastDueNotification:
		executor, ok := GetTypedExecutor[SendPastDueNotificationInput, SendPastDueNotificationOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendPastDueNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySyncInvoiceToAccounting:
		executor, ok := GetTypedExecutor[SyncInvoiceToAccountingInput, SyncInvoiceToAccountingOutput, SyncInvoiceToAccountingConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySyncInvoiceToAccounting)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityIdentifyMeetingParticipants:
		executor, ok := GetTypedExecutor[IdentifyMeetingParticipantsInput, IdentifyMeetingParticipantsOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIdentifyMeetingParticipants)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityIdentifyWebVisitor:
		executor, ok := GetTypedExecutor[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIdentifyWebVisitor)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendSlackNotification:
		executor, ok := GetTypedExecutor[SendSlackNotificationInput, NoOutput, SendSlackNotificationConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySendWebVisitorSlackNotification:
		executor, ok := GetTypedExecutor[SendWebVisitorSlackNotificationInput, NoOutput, SendWebVisitorSlackNotificationConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityUpdateCompanyStatus:
		executor, ok := GetTypedExecutor[UpdateCompanyStatusInput, NoOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySendWebVisitorSlackNotification)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityValidateEmailAddressDeliverability:
		executor, ok := GetTypedExecutor[ValidateEmailDeliverabilityInput, ValidateEmailDeliverabilityOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityValidateEmailAddressDeliverability)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityClassifyEmail:
		executor, ok := GetTypedExecutor[ClassifyEmailInput, ClassifyEmailOutput, ClassifyEmailConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityClassifyEmail)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityIdentifyEmailParticipants:
		executor, ok := GetTypedExecutor[IdentifyEmailParticipantsInput, IdentifyEmailParticipantsOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIdentifyEmailParticipants)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySummarizeMessage:
		executor, ok := GetTypedExecutor[SummarizeMessageInput, NoOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySummarizeMessage)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilitySummarizeThread:
		executor, ok := GetTypedExecutor[SummarizeThreadInput, NoOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilitySummarizeThread)
		}
		return executeCapability(ctx, executor, executionContainer)

	case enum.CapabilityIngestEmail:
		executor, ok := GetTypedExecutor[IngestEmailInput, NoOutput, postgres_entity.NoConfig](
			executionContainer.UntypedExecutors,
			executionContainer.Capability.Type,
		)
		if !ok {
			return enum.CapabilityExecutionError, nil, f.handleGetTypedExecutorError(ctx, enum.CapabilityIngestEmail)
		}
		return executeCapability(ctx, executor, executionContainer)

	default:
		err := fmt.Errorf("capability %s not configured", executionContainer.Capability.Type.String())
		spans.LogKV("capability", executionContainer.Capability.Type)
		spans.TraceError(err)
		return enum.CapabilityExecutionError, nil, err
	}
}

func (f *agentCapabilityExecutionService) handleGetTypedExecutorError(ctx context.Context, capability enum.AgentCapability) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AgentCapabilityExecutionService.handleGetTypedExecutorError")
	defer spans.Finish()

	err := fmt.Errorf("failed to get typed executor for %s", capability.String())
	spans.TraceError(err)
	return err
}

func executeCapability[I, O, C any](
	ctx context.Context,
	cap interfaces.AgentCapability[I, O, C],
	executionContainer interfaces.ExecutionContainer,
) (enum.CapabilityExecutionStatus, map[string]any, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AgentCapabilityExecutionService.executeCapability")
	defer spans.Finish()

	spans.LogKV("capabilityType", executionContainer.Capability.Type.String())

	input := cap.NewInput()
	err := utils.MapToStruct(executionContainer.ExecutionParams, &input)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, nil, err
	}

	config := cap.NewConfig()
	err = executionContainer.Capability.GetConfig(&config)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, nil, err
	}

	typedExecutionContainer := interfaces.TypedExecutionContainer[I, C]{
		AgentExecutionID: executionContainer.AgentExecutionID,
		InputData:        input,
		ConfigData:       config,
	}

	capabilityExecutionStatus, output, err := cap.Execute(ctx, typedExecutionContainer)
	switch capabilityExecutionStatus {
	case enum.CapabilityExecutionCompleted:
		results, err := utils.StructToMap(output)
		if err != nil {
			return enum.CapabilityExecutionError, nil, err
		}
		return enum.CapabilityExecutionCompleted, results, nil

	case enum.CapabilityExecutionError:
		spans.TraceError(err)
		return enum.CapabilityExecutionError, nil, err

	case enum.CapabilityExecutionStop:
		return enum.CapabilityExecutionStop, nil, nil

	case enum.CapabilityExecutionRetry:
		return enum.CapabilityExecutionRetry, nil, err

	case enum.CapabilityExecutionPending:
		// todo  -- need to implement poll until done
		return enum.CapabilityExecutionPending, nil, nil

	default:
		err = fmt.Errorf("unexpected capability execution status {%s}", capabilityExecutionStatus.String())
		spans.TraceError(err)
		return enum.CapabilityExecutionError, nil, err
	}
}
