package listeners

import (
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

func OnIntentEvent(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnIntentEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message, err := validateEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	intentData, err := validateIntentEvent(ctx, message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if common.GetTenantFromContext(ctx) == "" {
		ctx = common.SetTenantInContext(ctx, intentData.Tenant)
	}
	return distributeIntentEventToSubscribers(ctx, dependencies, intentData)
}

func distributeIntentEventToSubscribers(ctx context.Context, dep *model.DependencyContainer, intentData *dto.IntentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.distributeIntentEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var agentTypes []enum.AgentType
	switch intentData.IntentType {
	case enum.IntentSupportRequired:
		agentTypes = append(agentTypes, enum.AgentSupport)
	default:
		err := errors.New("IntentType not supported")
		tracing.TraceErr(span, err)
		return err
	}

	if len(agentTypes) == 0 {
		err := errors.New("No agent types configured for intent")
		tracing.TraceErr(span, err)
		return err
	}

	activeAgents := lookupActiveAgents(ctx, dep, agentTypes)
	var errs error
	for _, agent := range activeAgents {
		initialParams, err := utils.StructToMap(intentData)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = dep.CommonServices.AgentRunnerService.Run(ctx, agent, intentData.IntentType.String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func validateIntentEvent(ctx context.Context, message *dto.Event) (*dto.IntentEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.validateIntentEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	messageData, ok := message.Event.Data.(*dto.IntentEvent)
	if !ok {
		err := errors.New("cannot cast type to dto.IntentEvent")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return messageData, nil
}

func lookupActiveAgents(ctx context.Context, dep *model.DependencyContainer, agentTypes []enum.AgentType) []postgres_entity.Agents {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentTypes", fmt.Sprintf("%v", agentTypes)))

	agents, err := dep.PostgresRepositories.AgentsRepository.GetActiveAgentsByTypes(ctx, agentTypes)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
