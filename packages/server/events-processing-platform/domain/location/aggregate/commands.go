package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *LocationAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	switch c := cmd.(type) {
	case *command.UpsertLocationCommand:
		if c.IsCreateCommand {
			return a.createLocation(ctx, c)
		} else {
			return a.updateLocation(ctx, c)
		}
	case *command.FailedLocationValidationCommand:
		return a.failLocationValidation(ctx, c)
	case *command.SkippedLocationValidationCommand:
		return a.skipLocationValidation(ctx, c)
	case *command.LocationValidatedCommand:
		return a.locationValidated(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *LocationAggregate) createLocation(ctx context.Context, cmd *command.UpsertLocationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.createLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	locationAddress := models.LocationAddress{}
	locationAddress.From(cmd.LocationAddressFields)

	event, err := events.NewLocationCreateEvent(a, cmd.Name, cmd.RawAddress, cmd.Source, createdAtNotNil, updatedAtNotNil, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *LocationAggregate) updateLocation(ctx context.Context, cmd *command.UpsertLocationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.updateLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	if cmd.Source.Source == "" {
		cmd.Source.Source = constants.SourceOpenline
	}

	locationAddress := models.LocationAddress{}
	locationAddress.From(cmd.LocationAddressFields)

	event, err := events.NewLocationUpdateEvent(a, cmd.Name, cmd.RawAddress, cmd.Source.Source, updatedAtNotNil, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *LocationAggregate) failLocationValidation(ctx context.Context, cmd *command.FailedLocationValidationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.failLocationValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	event, err := events.NewLocationFailedValidationEvent(a, cmd.RawAddress, cmd.Country, cmd.ValidationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *LocationAggregate) skipLocationValidation(ctx context.Context, cmd *command.SkippedLocationValidationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.SkipLocationValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	event, err := events.NewLocationSkippedValidationEvent(a, cmd.RawAddress, cmd.ValidationSkipReason)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationSkippedValidationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *LocationAggregate) locationValidated(ctx context.Context, cmd *command.LocationValidatedCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.locationValidated")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	locationAddress := models.LocationAddress{}
	locationAddress.From(cmd.LocationAddressFields)

	event, err := events.NewLocationValidatedEvent(a, cmd.RawAddress, cmd.CountryForValidation, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationValidatedEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}
