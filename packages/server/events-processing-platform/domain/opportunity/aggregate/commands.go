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
	case *command.UpdateOpportunityCommand:
		return a.updateOpportunity(ctx, c)
	case *command.UpdateRenewalOpportunityCommand:
		return a.updateRenewalOpportunity(ctx, c)
	case *command.CloseWinOpportunityCommand:
		return a.closeWinOpportunity(ctx, c)
	case *command.CloseLooseOpportunityCommand:
		return a.closeLooseOpportunity(ctx, c)
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

	renewalLikelihood := cmd.RenewalLikelihood
	if string(renewalLikelihood) == "" {
		renewalLikelihood = model.RenewalLikelihoodStringHigh
	}

	createRenewalEvent, err := event.NewOpportunityCreateRenewalEvent(a, cmd.ContractId, string(renewalLikelihood), cmd.Source, createdAtNotNil, updatedAtNotNil)
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

func (a *OpportunityAggregate) updateOpportunity(ctx context.Context, cmd *command.UpdateOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	cmd.Source.SetDefaultValues()

	updateEvent, err := event.NewOpportunityUpdateEvent(a, cmd.DataFields, cmd.Source.Source, cmd.ExternalSystem, updatedAtNotNil, cmd.FieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *OpportunityAggregate) updateRenewalOpportunity(ctx context.Context, cmd *command.UpdateRenewalOpportunityCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.updateRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	cmd.Source.SetDefaultValues()

	renewalLikelihood := cmd.RenewalLikelihood
	if string(renewalLikelihood) == "" {
		renewalLikelihood = model.RenewalLikelihoodStringHigh
	}

	updateRenewalEvent, err := event.NewOpportunityUpdateRenewalEvent(a, string(renewalLikelihood), cmd.Comments, cmd.LoggedInUserId, cmd.Source.Source, cmd.Amount, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityUpdateRenewalEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateRenewalEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateRenewalEvent)
}

func (a *OpportunityAggregate) closeWinOpportunity(ctx context.Context, cmd *command.CloseWinOpportunityCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpportunityAggregate.closeWinOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	now := utils.Now()
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, now)
	closedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.ClosedAt, now)

	closeWinEvent, err := event.NewOpportunityCloseWinEvent(a, updatedAtNotNil, closedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCloseWinEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&closeWinEvent, span, aggregate.EventMetadata{
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

	now := utils.Now()
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, now)
	closedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.ClosedAt, now)

	closeLooseEvent, err := event.NewOpportunityCloseLooseEvent(a, updatedAtNotNil, closedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOpportunityCloseLooseEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&closeLooseEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(closeLooseEvent)
}
