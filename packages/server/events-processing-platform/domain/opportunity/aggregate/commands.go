package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
func (a *OpportunityAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpdateRenewalOpportunityNextCycleDateCommand:
		return a.updateRenewalOpportunityNextCycleDate(ctx, c)
	case *command.CloseWinOpportunityCommand:
		return a.closeWinOpportunity(ctx, c)
	case *command.CloseLooseOpportunityCommand:
		return a.closeLooseOpportunity(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *OpportunityAggregate) updateRenewalOpportunityNextCycleDate(ctx context.Context, cmd *command.UpdateRenewalOpportunityNextCycleDateCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateRenewalOpportunityNextCycleDate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	// skip if no changes on aggregate
	if a.Opportunity.RenewalDetails.RenewedAt != nil &&
		cmd.RenewedAt != nil &&
		a.Opportunity.RenewalDetails.RenewedAt.Equal(*cmd.RenewedAt) {
		return nil
	}

	// if opportunity is not renewal or status is closed, return error
	if a.Opportunity.InternalType != neo4jenum.OpportunityInternalTypeRenewal.String() {
		err := errors.New(events.Validate + ": Opportunity is not renewal")
		tracing.TraceErr(span, err)
		return err
	} else if a.Opportunity.InternalStage != neo4jenum.OpportunityInternalStageOpen.String() {
		err := errors.New(events.Validate + ": Opportunity is closed")
		tracing.TraceErr(span, err)
		return err
	}

	updateRenewalNextCycleDateEvent, err := event.NewOpportunityUpdateNextCycleDateEvent(a, updatedAtNotNil, cmd.RenewedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateRenewalNextCycleDateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateRenewalNextCycleDateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(updateRenewalNextCycleDateEvent)
}

func (a *OpportunityAggregate) closeWinOpportunity(ctx context.Context, cmd *command.CloseWinOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.closeWinOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// skip if opportunity is already closed won
	if a.Opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedWon.String() {
		span.LogFields(log.String("result", "Opportunity is already closed won"))
		return nil
	}

	now := utils.Now()
	closedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.ClosedAt, now)

	closeWinEvent, err := event.NewOpportunityCloseWinEvent(a, closedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCloseWinEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&closeWinEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(closeWinEvent)
}

func (a *OpportunityAggregate) closeLooseOpportunity(ctx context.Context, cmd *command.CloseLooseOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.closeLooseOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	// skip if opportunity is already closed lost
	if a.Opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedLost.String() {
		span.LogFields(log.String("result", "Opportunity is already closed lost"))
		return nil
	}

	now := utils.Now()
	closedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.ClosedAt, now)

	closeLooseEvent, err := event.NewOpportunityCloseLooseEvent(a, closedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCloseLooseEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&closeLooseEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(closeLooseEvent)
}
