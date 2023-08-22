package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphOrganizationEventHandler struct {
	Repositories         *repository.Repositories
	organizationCommands *cmdhnd.OrganizationCommands
	log                  logger.Logger
}

func (e *GraphOrganizationEventHandler) OnOrganizationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnOrganizationCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.OrganizationRepository.CreateOrganization(ctx, organizationId, eventData)

	return err
}

func (e *GraphOrganizationEventHandler) OnOrganizationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnOrganizationUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	if eventData.IgnoreEmptyFields {
		return e.Repositories.OrganizationRepository.UpdateOrganizationIgnoreEmptyInputParams(ctx, organizationId, eventData)
	} else {
		return e.Repositories.OrganizationRepository.UpdateOrganization(ctx, organizationId, eventData)
	}
}

func (e *GraphOrganizationEventHandler) OnPhoneNumberLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnPhoneNumberLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (e *GraphOrganizationEventHandler) OnEmailLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnEmailLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (e *GraphOrganizationEventHandler) OnDomainLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnDomainLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.OrganizationRepository.LinkWithDomain(ctx, eventData.Tenant, organizationId, eventData.Domain)

	return err
}

func (e *GraphOrganizationEventHandler) OnSocialAddedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnSocialAddedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.SocialRepository.CreateSocialFor(ctx, eventData.Tenant, organizationId, "Organization", eventData)

	return err
}

func (h *GraphOrganizationEventHandler) OnRenewalLikelihoodUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnRenewalLikelihoodUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateRenewalLikelihoodEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	orgDbNode, err := h.Repositories.OrganizationRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "GetOrganization")
	}
	organizationEntity := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)

	err = h.Repositories.OrganizationRepository.UpdateRenewalLikelihood(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if organizationEntity.RenewalLikelihood.RenewalLikelihood != eventData.GetRenewalLikelihoodAsStringForGraphDb() {
		err := h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
	}

	return err
}

func (h *GraphOrganizationEventHandler) OnRenewalForecastUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnRenewalForecastUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateRenewalForecastEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.Repositories.OrganizationRepository.UpdateRenewalForecast(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if eventData.UpdatedBy != "" && eventData.Amount == nil {
		err := h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
	}

	return err
}

func (h *GraphOrganizationEventHandler) OnBillingDetailsUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphOrganizationEventHandler.OnRenewalForecastUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateBillingDetailsEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.Repositories.OrganizationRepository.UpdateBillingDetails(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if eventData.UpdatedBy != "" {
		err = h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
		err = h.organizationCommands.RequestNextCycleDateCommand.Handle(ctx, cmd.NewRequestNextCycleDateCommand(eventData.Tenant, organizationId))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestNextCycleDateCommand failed: %v", err.Error())
		}
	}

	return err
}
