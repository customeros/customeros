package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
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
	default:
		return errors.New("invalid command type")
	}
}

func (a *UserAggregate) createUser(ctx context.Context, cmd *command.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.createUser")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)

	createEvent, err := events.NewUserCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&createEvent, &span, a.Tenant, cmd.UserID)

	return a.Apply(createEvent)
}

func (a *UserAggregate) updateUser(ctx context.Context, cmd *command.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.updateUser")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	sourceOfTruth := cmd.Source.SourceOfTruth
	if sourceOfTruth == "" {
		sourceOfTruth = a.User.Source.SourceOfTruth
	}

	// do not change data if user was modified by openline
	if sourceOfTruth != a.User.Source.SourceOfTruth && a.User.Source.SourceOfTruth == constants.SourceOpenline {
		sourceOfTruth = a.User.Source.SourceOfTruth
		if a.User.Name != "" {
			cmd.DataFields.Name = a.User.Name
		}
		if a.User.FirstName != "" {
			cmd.DataFields.Name = a.User.FirstName
		}
		if a.User.LastName != "" {
			cmd.DataFields.Name = a.User.LastName
		}
		if a.User.Timezone != "" {
			cmd.DataFields.Name = a.User.Timezone
		}
		if a.User.ProfilePhotoUrl != "" {
			cmd.DataFields.Name = a.User.ProfilePhotoUrl
		}
	}

	event, err := events.NewUserUpdateEvent(a, cmd.DataFields, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *UserAggregate) addPlayerInfo(ctx context.Context, cmd *command.AddPlayerInfoCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.addPlayerInfo")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	timestampNotNil := utils.IfNotNilTimeWithDefault(cmd.Timestamp, utils.Now())

	event, err := events.NewUserAddPlayerInfoEvent(a, models.PlayerInfo{
		Provider:   cmd.Provider,
		AuthId:     cmd.AuthId,
		IdentityId: cmd.IdentityId,
	}, cmd.Source, timestampNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserAddPlayerInfoEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

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
