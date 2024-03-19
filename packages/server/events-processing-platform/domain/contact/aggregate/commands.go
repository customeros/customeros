package aggregate

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *ContactAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpsertContactCommand:
		if c.IsCreateCommand {
			return a.createContact(ctx, c)
		} else {
			return a.updateContact(ctx, c)
		}
	case *command.LinkEmailCommand:
		return a.linkEmail(ctx, c)
	case *command.LinkPhoneNumberCommand:
		return a.linkPhoneNumber(ctx, c)
	case *command.LinkLocationCommand:
		return a.linkLocation(ctx, c)
	case *command.LinkOrganizationCommand:
		return a.linkOrganization(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *ContactAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contactpb.ContactAddSocialGrpcRequest:
		return a.addSocial(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContactAggregate) createContact(ctx context.Context, cmd *command.UpsertContactCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.createContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := event.NewContactCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactCreateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *ContactAggregate) updateContact(ctx context.Context, cmd *command.UpsertContactCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.createContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckForNoChanges(cmd.Source.AppSource, cmd.LoggedInUserId) {
		if a.Contact.SameData(cmd.DataFields, cmd.ExternalSystem) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	updateEvent, err := event.NewContactUpdateEvent(a, cmd.Source.Source, cmd.DataFields, cmd.ExternalSystem, updatedAtNotNil, cmd.FieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactUpdateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *ContactAggregate) linkEmail(ctx context.Context, cmd *command.LinkEmailCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.linkEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckForNoChanges(cmd.AppSource, cmd.LoggedInUserId) {
		if a.Contact.HasEmail(cmd.EmailId, cmd.Label, cmd.Primary) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := event.NewContactLinkEmailEvent(a, cmd.EmailId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkEmailEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.Contact.Emails {
			if k != cmd.EmailId && v.Primary {
				if err = a.SetEmailNonPrimary(ctx, k, cmd.LoggedInUserId, cmd.AppSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *ContactAggregate) SetEmailNonPrimary(ctx context.Context, emailId, loggedInUserId, appSource string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("emailId", emailId), log.String("loggedInUserId", loggedInUserId))

	updatedAtNotNil := utils.Now()

	email, ok := a.Contact.Emails[emailId]
	if !ok {
		return nil
	}

	if email.Primary {
		event, err := event.NewContactLinkEmailEvent(a, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewContactLinkEmailEvent")
		}

		aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
			Tenant: a.Tenant,
			UserId: loggedInUserId,
			App:    appSource,
		})
		return a.Apply(event)
	}
	return nil
}

func (a *ContactAggregate) linkPhoneNumber(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.linkPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckForNoChanges(cmd.AppSource, cmd.LoggedInUserId) {
		if a.Contact.HasPhoneNumber(cmd.PhoneNumberId, cmd.Label, cmd.Primary) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := event.NewContactLinkPhoneNumberEvent(a, cmd.PhoneNumberId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkPhoneNumberEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.Contact.PhoneNumbers {
			if k != cmd.PhoneNumberId && v.Primary {
				if err = a.SetPhoneNumberNonPrimary(ctx, k, cmd.LoggedInUserId, cmd.AppSource); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *ContactAggregate) SetPhoneNumberNonPrimary(ctx context.Context, phoneNumberId, loggedInUserId, appSource string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("phoneNumberId", phoneNumberId), log.String("loggedInUserId", loggedInUserId))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Contact.PhoneNumbers[phoneNumberId]
	if !ok {
		return nil
	}

	if phoneNumber.Primary {
		event, err := event.NewContactLinkPhoneNumberEvent(a, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewContactLinkPhoneNumberEvent")
		}

		aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
			Tenant: a.Tenant,
			UserId: loggedInUserId,
			App:    appSource,
		})
		return a.Apply(event)
	}
	return nil
}

func (a *ContactAggregate) linkLocation(ctx context.Context, cmd *command.LinkLocationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.linkLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckForNoChanges(cmd.AppSource, cmd.LoggedInUserId) {
		if a.Contact.HasLocation(cmd.LocationId) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	updatedAtNotNil := utils.Now()

	event, err := event.NewContactLinkLocationEvent(a, cmd.LocationId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkLocationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *ContactAggregate) linkOrganization(ctx context.Context, cmd *command.LinkOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.linkOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	if aggregate.AllowCheckForNoChanges(cmd.Source.AppSource, cmd.LoggedInUserId) {
		if a.Contact.HasJobRoleInOrganization(cmd.OrganizationId, cmd.JobRoleFields, cmd.Source) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	event, err := event.NewContactLinkWithOrganizationEvent(a, cmd.OrganizationId, cmd.JobRoleFields.JobTitle, cmd.JobRoleFields.Description,
		cmd.JobRoleFields.Primary, cmd.Source, createdAtNotNil, updatedAtNotNil, cmd.JobRoleFields.StartedAt, cmd.JobRoleFields.EndedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkWithOrganizationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(event)
}

func (a *ContactAggregate) addSocial(ctx context.Context, r *contactpb.ContactAddSocialGrpcRequest) (any, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.addSocial")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(r.SourceFields)

	if aggregate.AllowCheckForNoChanges(sourceFields.AppSource, r.LoggedInUserId) {
		if a.Contact.HasSocialUrl(r.Url) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return "", nil
		}
	}

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.CreatedAt), utils.Now())

	socialId := utils.StringFirstNonEmpty(a.Contact.GetSocialIdForUrl(r.Url), uuid.New().String())

	addSocialEvent, err := event.NewAddSocialEvent(a, r, socialId, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewAddSocialEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addSocialEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.SourceFields.AppSource,
	})

	return socialId, a.Apply(addSocialEvent)
}
