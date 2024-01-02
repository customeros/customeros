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
		App:    cmd.SourceFields.AppSource,
	})

	return a.Apply(createEvent)
}
