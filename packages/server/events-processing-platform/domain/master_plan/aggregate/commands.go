package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
func (a *MasterPlanAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.CreateMasterPlanCommand:
		return a.createMasterPlan(ctx, c)
	case *command.UpdateMasterPlanCommand:
		return a.updateMasterPlan(ctx, c)
	case *command.CreateMasterPlanMilestoneCommand:
		return a.createMasterPlanMilestone(ctx, c)
	case *command.UpdateMasterPlanMilestoneCommand:
		return a.updateMasterPlanMilestone(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *MasterPlanAggregate) createMasterPlan(ctx context.Context, cmd *command.CreateMasterPlanCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MasterPlanAggregate.createMasterPlan")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())

	createEvent, err := event.NewMasterPlanCreateEvent(a, cmd.Name, cmd.SourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewMasterPlanCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.GetAppSource(),
	})

	return a.Apply(createEvent)
}

func (a *MasterPlanAggregate) updateMasterPlan(ctx context.Context, cmd *command.UpdateMasterPlanCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MasterPlanAggregate.updateMasterPlan")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	updateEvent, err := event.NewMasterPlanUpdateEvent(a, cmd.Name, cmd.Retired, updatedAtNotNil, cmd.FieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewMasterPlanUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.GetAppSource(),
	})

	return a.Apply(updateEvent)
}

func (a *MasterPlanAggregate) createMasterPlanMilestone(ctx context.Context, cmd *command.CreateMasterPlanMilestoneCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MasterPlanAggregate.createMasterPlanMilestone")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())

	createEvent, err := event.NewMasterPlanMilestoneCreateEvent(a, cmd.MilestoneId, cmd.Name, cmd.DurationHours, cmd.Order, cmd.Items, cmd.Optional, cmd.SourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewMasterPlanMilestoneCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.GetAppSource(),
	})

	return a.Apply(createEvent)
}

func (a *MasterPlanAggregate) updateMasterPlanMilestone(ctx context.Context, cmd *command.UpdateMasterPlanMilestoneCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MasterPlanAggregate.updateMasterPlanMilestone")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	updateEvent, err := event.NewMasterPlanMilestoneUpdateEvent(a, cmd.MilestoneId, cmd.Name, cmd.DurationHours, cmd.Order,
		cmd.Items, cmd.FieldsMask, cmd.Optional, cmd.Retired, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewMasterPlanMilestoneUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.GetAppSource(),
	})

	return a.Apply(updateEvent)
}
