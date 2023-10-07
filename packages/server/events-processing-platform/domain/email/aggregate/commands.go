package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *EmailAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	switch c := cmd.(type) {
	case *command.UpsertEmailCommand:
		if c.IsCreateCommand {
			return a.createEmail(ctx, c)
		} else {
			return a.updateEmail(ctx, c)
		}
	case *command.EmailValidatedCommand:
		return a.emailValidated(ctx, c)
	case *command.FailedEmailValidationCommand:
		return a.failEmailValidation(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *EmailAggregate) createEmail(ctx context.Context, cmd *command.UpsertEmailCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.createEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	event, err := events.NewEmailCreateEvent(a, cmd.Tenant, cmd.RawEmail, cmd.Source, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *EmailAggregate) updateEmail(ctx context.Context, cmd *command.UpsertEmailCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.updateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	if cmd.Source.Source == "" {
		cmd.Source.Source = constants.SourceOpenline
	}

	event, err := events.NewEmailUpdateEvent(a, cmd.RawEmail, cmd.Tenant, cmd.Source.Source, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *EmailAggregate) failEmailValidation(ctx context.Context, cmd *command.FailedEmailValidationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.updateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewEmailFailedValidationEvent(a, cmd.Tenant, cmd.ValidationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *EmailAggregate) emailValidated(ctx context.Context, cmd *command.EmailValidatedCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.emailValidated")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewEmailValidatedEvent(a, cmd.Tenant, cmd.RawEmail, cmd.IsReachable, cmd.ValidationError,
		cmd.Domain, cmd.Username, cmd.EmailAddress, cmd.AcceptsMail, cmd.CanConnectSmtp, cmd.HasFullInbox, cmd.IsCatchAll,
		cmd.IsDeliverable, cmd.IsDisabled, cmd.IsValidSyntax)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailValidatedEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}
