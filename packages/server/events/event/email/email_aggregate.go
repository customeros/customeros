package email

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	emailevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	EmailAggregateType eventstore.AggregateType = "email"
)

type EmailAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Email *Email
}

func NewEmailAggregateWithTenantAndID(tenant, id string) *EmailAggregate {
	emailAggregate := EmailAggregate{}
	emailAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(EmailAggregateType, tenant, id)
	emailAggregate.SetWhen(emailAggregate.When)
	emailAggregate.Email = &Email{}
	emailAggregate.Tenant = tenant

	return &emailAggregate
}

func (a *EmailAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *emailpb.EmailValidationGrpcRequest:
		return nil, a.emailValidatedV2(ctx, r)
	case *emailpb.UpsertEmailRequest:
		return nil, a.UpsertEmail(ctx, r)
	case *emailpb.DeleteEmailRequest:
		return nil, a.DeleteEmail(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *EmailAggregate) emailValidatedV2(ctx context.Context, request *emailpb.EmailValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.emailValidatedV2")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	event, err := emailevent.NewEmailValidatedEventV2(a,
		request.Tenant,
		request.RawEmail,
		request.Email,
		request.Domain,
		request.Username,
		request.IsValidSyntax,
		request.IsRisky,
		request.IsFirewalled,
		request.Provider,
		request.Firewall,
		request.Deliverable,
		request.IsCatchAll,
		request.IsMailboxFull,
		request.IsRoleAccount,
		request.IsSystemGenerated,
		request.IsFreeAccount,
		request.SmtpSuccess,
		request.ResponseCode,
		request.ErrorCode,
		request.Description,
		request.IsPrimaryDomain,
		request.PrimaryDomain,
		request.AlternateEmail,
		request.RetryValidation,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailValidatedEventV2")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(event)
}

func (a *EmailAggregate) UpsertEmail(ctx context.Context, request *emailpb.UpsertEmailRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.UpsertEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	source := request.SourceFields.Source
	if source == "" {
		source = constants.SourceOpenline
	}

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	upsertEvent, err := emailevent.NewEmailUpsertEvent(a, request.Tenant, request.RawEmail, source, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailUpsertEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&upsertEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(upsertEvent)
}

func (a *EmailAggregate) DeleteEmail(ctx context.Context, request *emailpb.DeleteEmailRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.UpsertEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	deleteEvent, err := emailevent.NewEmailDeleteEvent(a, request.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailDeleteEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(deleteEvent)
}

func (a *EmailAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case emailevent.EmailValidatedV2:
		return a.OnEmailValidated(event)
	case emailevent.EmailValidationFailedV1,
		emailevent.EmailValidatedV1,
		emailevent.EmailCreateV1,
		emailevent.EmailUpdateV1,
		emailevent.EmailUpsertV1,
		emailevent.EmailDeleteV1:
		return nil
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *EmailAggregate) OnEmailValidated(event eventstore.Event) error {
	var eventData emailevent.EmailValidatedEventV2
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Email.Email = eventData.Email
	return nil
}
