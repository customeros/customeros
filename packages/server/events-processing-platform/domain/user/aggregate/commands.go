package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	local_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *UserAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpsertUserCommand:
		if c.IsCreateCommand {
			return a.createUser(ctx, c)
		} else {
			return a.updateUser(ctx, c)
		}
	case *command.AddRoleCommand:
		return a.addRole(ctx, c)
	case *command.RemoveRoleCommand:
		return a.removeRole(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *UserAggregate) HandleRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *userpb.LinkPhoneNumberToUserGrpcRequest:
		return nil, a.linkPhoneNumber(ctx, r)
	case *userpb.LinkEmailToUserGrpcRequest:
		return nil, a.linkEmail(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
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

	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *UserAggregate) updateUser(ctx context.Context, cmd *command.UpsertUserCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.updateUser")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if eventstore.AllowCheckForNoChanges(cmd.Source.AppSource, cmd.LoggedInUserId) {
		if a.User.SameUserData(cmd.DataFields, cmd.ExternalSystem) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	event, err := events.NewUserUpdateEvent(a, cmd.DataFields, cmd.Source.Source, updatedAtNotNil, cmd.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserUpdateEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

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

	eventstore.EnrichEventWithMetadata(&event, &span, a.Tenant, loggedInUserId)

	return a.Apply(event)
}

func (a *UserAggregate) linkPhoneNumber(ctx context.Context, request *userpb.LinkPhoneNumberToUserGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.linkPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))

	if eventstore.AllowCheckForNoChanges(request.AppSource, request.LoggedInUserId) {
		if a.User.HasPhoneNumber(request.PhoneNumberId, request.Label) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkPhoneNumberEvent(a, request.Tenant, request.PhoneNumberId, request.Label, request.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if request.Primary {
		for k, v := range a.User.PhoneNumbers {
			if k != request.PhoneNumberId && v.Primary {
				if err = a.SetPhoneNumberNonPrimary(ctx, request.Tenant, k, request.LoggedInUserId); err != nil {
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

		eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
			Tenant: a.Tenant,
			UserId: loggedInUserId,
		})
		return a.Apply(event)
	}
	return nil
}

func (a *UserAggregate) linkEmail(ctx context.Context, request *userpb.LinkEmailToUserGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.linkEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))

	if eventstore.AllowCheckForNoChanges(request.AppSource, request.LoggedInUserId) {
		if a.User.HasEmail(request.EmailId) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkEmailEvent(a, request.Tenant, request.EmailId, request.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkEmailEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if request.Primary {
		for k, v := range a.User.Emails {
			if k != request.EmailId && v.Primary {
				if err = a.SetEmailNonPrimary(ctx, request.Tenant, k, request.LoggedInUserId); err != nil {
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
		event, err := events.NewUserLinkEmailEvent(a, tenant, emailId, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewUserLinkEmailEvent")
		}

		eventstore.EnrichEventWithMetadata(&event, &span, a.Tenant, loggedInUserId)
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
	eventstore.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

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
	eventstore.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}
