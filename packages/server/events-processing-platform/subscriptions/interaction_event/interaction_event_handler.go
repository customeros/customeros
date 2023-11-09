package interactionEvent

import (
	"context"
	"encoding/json"
	"fmt"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/ai"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type interactionEventHandler struct {
	repositories             *repository.Repositories
	interactionEventCommands *command_handler.CommandHandlers
	log                      logger.Logger
	cfg                      *config.Config
}

func (h *interactionEventHandler) GenerateSummaryForEmail(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.GenerateSummaryForEmail")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.InteractionEventRequestSummaryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	interactionEvent, err := h.repositories.InteractionEventRepository.GetInteractionEvent(ctx, eventData.Tenant, interactionEventId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting interaction event with id %s: %v", interactionEvent, err)
		return nil
	}

	interactionEventChannel := utils.GetStringPropOrEmpty(interactionEvent.Props, "channel")
	interactionEventContent := utils.GetStringPropOrEmpty(interactionEvent.Props, "content")

	if interactionEventChannel != "EMAIL" {
		tracing.TraceErr(span, errors.New("interaction event is not an email"))
		h.log.Warnf("Interaction event with id %s is not an email, skipping", interactionEventId)
		return nil
	}
	if interactionEventContent == "" {
		tracing.TraceErr(span, errors.New("interaction event content is empty"))
		h.log.Warnf("Interaction event with id %s has no content, skipping", interactionEventId)
		return nil
	}

	summaryPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.EmailSummaryPrompt, interactionEventContent)

	promptLog := commonEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatform,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_EmailSummary,
		Tenant:         &eventData.Tenant,
		NodeId:         &interactionEventId,
		NodeLabel:      utils.StringPtr(constants.NodeLabel_InteractionEvent),
		PromptTemplate: &h.cfg.Services.Anthropic.EmailSummaryPrompt,
		Prompt:         summaryPrompt,
	}
	promptStoreLogId, err := h.repositories.CommonRepositories.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId", promptStoreLogId))
	}

	aiResponse, err := ai.InvokeAnthropic(ctx, h.cfg, h.log, summaryPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err.Error())
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil
	} else {
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId, aiResponse)
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}
	summary := utils.ExtractAfterColon(aiResponse)

	err = h.interactionEventCommands.ReplaceSummary.Handle(ctx, command.NewReplaceSummaryCommand(eventData.Tenant, interactionEventId, summary, "text/plain", nil))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error replacing summary: %v", err)
		return err
	}

	return nil
}

func (h *interactionEventHandler) GenerateActionItemsForEmail(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.GenerateActionItemsForEmail")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.InteractionEventRequestSummaryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	interactionEvent, err := h.repositories.InteractionEventRepository.GetInteractionEvent(ctx, eventData.Tenant, interactionEventId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting interaction event with id %s: %v", interactionEvent, err)
		return nil
	}

	interactionEventChannel := utils.GetStringPropOrEmpty(interactionEvent.Props, "channel")
	interactionEventContent := utils.GetStringPropOrEmpty(interactionEvent.Props, "content")

	if interactionEventChannel != "EMAIL" {
		tracing.TraceErr(span, errors.New("interaction event is not an email"))
		h.log.Warnf("Interaction event with id %s is not an email, skipping", interactionEventId)
		return nil
	}
	if interactionEventContent == "" {
		tracing.TraceErr(span, errors.New("interaction event content is empty"))
		h.log.Warnf("Interaction event with id %s has no content, skipping", interactionEventId)
		return nil
	}

	actionItemsPrompt := fmt.Sprintf(h.cfg.Services.Anthropic.EmailActionsItemsPrompt, interactionEventContent)

	promptLog := commonEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatform,
		Provider:       constants.Anthropic,
		Model:          "claude-2",
		PromptType:     constants.PromptType_EmailActionItems,
		Tenant:         &eventData.Tenant,
		NodeId:         &interactionEventId,
		NodeLabel:      utils.StringPtr(constants.NodeLabel_InteractionEvent),
		PromptTemplate: &h.cfg.Services.Anthropic.EmailActionsItemsPrompt,
		Prompt:         actionItemsPrompt,
	}
	promptStoreLogId, err := h.repositories.CommonRepositories.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log: %v", err)
	} else {
		span.LogFields(log.String("promptStoreLogId", promptStoreLogId))
	}

	aiResponse, err := ai.InvokeAnthropic(ctx, h.cfg, h.log, actionItemsPrompt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking AI: %v", err.Error())
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil
	} else {
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId, aiResponse)
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}

	actionItems, err := extractActionItemsFromAiResponse(aiResponse)
	span.LogFields(log.Object("output - actionItems", actionItems))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error extracting action items from ai response: %v", err)
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil
	}
	if len(actionItems) == 0 {
		storeErr := h.repositories.CommonRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, storeErr)
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
	}

	err = h.interactionEventCommands.ReplaceActionItems.Handle(ctx, command.NewReplaceActionItemsCommand(eventData.Tenant, interactionEventId, actionItems, nil))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error replacing action items: %v", err)
		return err
	}

	return nil
}

func extractActionItemsFromAiResponse(str string) ([]string, error) {
	jsonStr, err := utils.ExtractJsonFromString(str)
	if err != nil {
		return []string{}, err
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(jsonStr), &data)

	items, ok := data["items"].([]interface{})
	if !ok {
		return []string{}, fmt.Errorf("invalid JSON format")
	}

	var actionItems []string
	for _, item := range items {
		actionItems = append(actionItems, item.(string))
	}

	return actionItems, nil
}
