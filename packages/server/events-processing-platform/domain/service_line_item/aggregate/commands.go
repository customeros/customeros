package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// HandleCommand processes commands and applies the resulting events to the aggregate.
func (a *ServiceLineItemAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpdateServiceLineItemCommand:
		return a.updateServiceLineItem(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *ServiceLineItemAggregate) updateServiceLineItem(ctx context.Context, cmd *command.UpdateServiceLineItemCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.updateServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if a.ServiceLineItem.IsDeleted {
		err := errors.New(constants.Validate + ": cannot update a deleted service line item")
		tracing.TraceErr(span, err)
		return err
	}

	// fail if quantity or price is negative
	if cmd.DataFields.Quantity < 0 || cmd.DataFields.Price < 0 {
		err := errors.New(constants.FieldValidation + ": quantity and price must not be negative")
		tracing.TraceErr(span, err)
		return err
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	if (a.ServiceLineItem.Billed == model.OnceBilled.String() && cmd.DataFields.Billed != model.OnceBilled) ||
		(a.ServiceLineItem.Billed != model.OnceBilled.String() && cmd.DataFields.Billed == model.OnceBilled) {
		return errors.New(constants.FieldValidation + ": cannot change billed type from 'once' to a frequency-based option or vice versa")
	}

	// Adjust vat rate
	if cmd.DataFields.VatRate < 0 {
		cmd.DataFields.VatRate = 0
	}
	cmd.DataFields.VatRate = utils.TruncateFloat64(cmd.DataFields.VatRate, 2)

	// Prepare the data for the update event
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		a,
		cmd.DataFields,
		cmd.Source,
		updatedAtNotNil,
		cmd.StartedAt,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}
