package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	local_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *UserAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	switch c := cmd.(type) {
	case *command.UpsertUserCommand:
		if c.IsCreateCommand {
			return a.createUser(ctx, c)
		} else {
			return a.updateUser(ctx, c)
		}
	case *command.AddPlayerInfoCommand:
		return a.addPlayerInfo(ctx, c)
	case *command.LinkEmailCommand:
		return a.linkEmail(ctx, c)
	case *command.LinkPhoneNumberCommand:
		return a.linkPhoneNumber(ctx, c)
	case *command.AddRoleCommand:
		return a.addRole(ctx, c)
	case *command.RemoveRoleCommand:
		return a.removeRole(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *UserAggregate) createUser(ctx context.Context, cmd *command.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.createUser")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := events.NewUserCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&createEvent, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(createEvent)
}

func (a *UserAggregate) updateUser(ctx context.Context, cmd *command.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.updateUser")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	event, err := events.NewUserUpdateEvent(a, cmd.DataFields, cmd.Source.Source, updatedAtNotNil, cmd.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *UserAggregate) addPlayerInfo(ctx context.Context, cmd *command.AddPlayerInfoCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.addPlayerInfo")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	timestampNotNil := utils.IfNotNilTimeWithDefault(cmd.Timestamp, utils.Now())
	cmd.Source.SetDefaultValues()

	event, err := events.NewUserAddPlayerInfoEvent(a, models.PlayerInfo{
		Provider:   cmd.Provider,
		AuthId:     cmd.AuthId,
		IdentityId: cmd.IdentityId,
	}, cmd.Source, timestampNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserAddPlayerInfoEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *UserAggregate) LinkJobRole(ctx context.Context, tenant, jobRoleId, loggedInUserId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkJobRole")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkJobRoleEvent(a, tenant, jobRoleId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkJobRoleEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, loggedInUserId)

	return a.Apply(event)
}

func (a *UserAggregate) linkPhoneNumber(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.linkPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkPhoneNumberEvent(a, cmd.Tenant, cmd.PhoneNumberId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.User.PhoneNumbers {
			if k != cmd.PhoneNumberId && v.Primary {
				if err = a.SetPhoneNumberNonPrimary(ctx, cmd.Tenant, k, cmd.LoggedInUserId); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *UserAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId, loggedInUserId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())

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

		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, loggedInUserId)
		return a.Apply(event)
	}
	return nil
}

func (a *UserAggregate) linkEmail(ctx context.Context, cmd *command.LinkEmailCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.linkEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckIfEventIsRedundant(cmd.AppSource, cmd.LoggedInUserId) {
		if a.User.HasEmail(cmd.EmailId, cmd.Label, cmd.Primary) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkEmailEvent(a, cmd.Tenant, cmd.EmailId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkEmailEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: cmd.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.User.Emails {
			if k != cmd.EmailId && v.Primary {
				if err = a.SetEmailNonPrimary(ctx, cmd.Tenant, k, cmd.LoggedInUserId); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *UserAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId, loggedInUserId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())

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

		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, loggedInUserId)
		return a.Apply(event)
	}
	return nil
}

func (a *UserAggregate) addRole(ctx context.Context, cmd *command.AddRoleCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.addRole")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	event, err := events.NewUserAddRoleEvent(a, cmd.Role, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserAddRoleEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *UserAggregate) removeRole(ctx context.Context, cmd *command.RemoveRoleCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.removeRole")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	event, err := events.NewUserRemoveRoleEvent(a, cmd.Role, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserRemoveRoleEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}
