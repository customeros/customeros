package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
func (a *OpportunityAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.CreateOpportunityCommand:
		return a.createOpportunity(ctx, c)
	case *command.CreateRenewalOpportunityCommand:
		return a.createRenewalOpportunity(ctx, c)
	case *command.UpdateRenewalOpportunityNextCycleDateCommand:
		return a.updateRenewalOpportunityNextCycleDate(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *OpportunityAggregate) createOpportunity(ctx context.Context, cmd *command.CreateOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.createOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := event.NewOpportunityCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *OpportunityAggregate) createRenewalOpportunity(ctx context.Context, cmd *command.CreateRenewalOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.createRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createRenewalEvent, err := event.NewOpportunityCreateRenewalEvent(a, cmd.ContractId, cmd.Source, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCreateRenewalEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createRenewalEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createRenewalEvent)
}

func (a *OpportunityAggregate) updateRenewalOpportunityNextCycleDate(ctx context.Context, cmd *command.UpdateRenewalOpportunityNextCycleDateCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateRenewalOpportunityNextCycleDate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	// if opportunity is not renewal or status is closed, return error
	if a.Opportunity.InternalType != model.OpportunityInternalTypeStringRenewal {
		err := errors.New(constants.Validate + ": Opportunity is not renewal")
		tracing.TraceErr(span, err)
		return err
	} else if a.Opportunity.InternalStage != model.OpportunityInternalStageStringOpen {
		err := errors.New(constants.Validate + ": Opportunity is closed")
		tracing.TraceErr(span, err)
		return err
	}

	updateRenewalNextCycleDateEvent, err := event.NewOpportunityUpdateNextCycleDateEvent(a, updatedAtNotNil, cmd.RenewedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateRenewalNextCycleDateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateRenewalNextCycleDateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(updateRenewalNextCycleDateEvent)
}
