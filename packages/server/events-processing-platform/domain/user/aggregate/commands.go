package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	local_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *UserAggregate) HandleCommand(ctx context.Context, command eventstore.Command) error {
	switch c := command.(type) {
	case *cmd.UpsertUserCommand:
		if c.IsCreateCommand {
			return a.createUser(ctx, c)
		} else {
			return a.updateUser(ctx, c)
		}
	default:
		return errors.New("invalid command type")
	}
}

func (a *UserAggregate) createUser(ctx context.Context, command *cmd.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.createUser")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(command.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, createdAtNotNil)

	createEvent, err := events.NewUserCreateEvent(a, command.DataFields, command.Source, command.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&createEvent, &span, a.Tenant, command.UserID)

	return a.Apply(createEvent)
}

func (a *UserAggregate) updateUser(ctx context.Context, command *cmd.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.updateUser")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, utils.Now())
	sourceOfTruth := command.Source.SourceOfTruth
	if sourceOfTruth == "" {
		sourceOfTruth = a.User.Source.SourceOfTruth
	}

	// do not change data if user was modified by openline
	if sourceOfTruth != a.User.Source.SourceOfTruth && a.User.Source.SourceOfTruth == constants.SourceOpenline {
		sourceOfTruth = a.User.Source.SourceOfTruth
		if a.User.Name != "" {
			command.DataFields.Name = a.User.Name
		}
		if a.User.FirstName != "" {
			command.DataFields.Name = a.User.FirstName
		}
		if a.User.LastName != "" {
			command.DataFields.Name = a.User.LastName
		}
		if a.User.Timezone != "" {
			command.DataFields.Name = a.User.Timezone
		}
		if a.User.ProfilePhotoUrl != "" {
			command.DataFields.Name = a.User.ProfilePhotoUrl
		}
	}

	event, err := events.NewUserUpdateEvent(a, command.DataFields, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *UserAggregate) LinkJobRole(ctx context.Context, tenant, jobRoleId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkJobRole")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkJobRoleEvent(a, tenant, jobRoleId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkJobRoleEvent")
	}

	// TODO add user id
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}

func (a *UserAggregate) LinkPhoneNumber(ctx context.Context, tenant, phoneNumberId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkPhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkPhoneNumberEvent(a, tenant, phoneNumberId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
	}

	// TODO add user id
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}

func (a *UserAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.User.PhoneNumbers[phoneNumberId]
	if !ok {
		return local_errors.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := events.NewUserLinkPhoneNumberEvent(a, tenant, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
		}

		// TODO add user id
		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")
		return a.Apply(event)
	}
	return nil
}

func (a *UserAggregate) LinkEmail(ctx context.Context, tenant, emailId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkEmailEvent(a, tenant, emailId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkEmailEvent")
	}

	// TODO add user id
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}

func (a *UserAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	email, ok := a.User.Emails[emailId]
	if !ok {
		return local_errors.ErrEmailNotFound
	}

	if email.Primary {
		event, err := events.NewUserLinkEmailEvent(a, tenant, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewUserLinkEmailEvent")
		}

		// TODO add user id
		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")
		return a.Apply(event)
	}
	return nil
}
