package email

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
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
		request.IsFreeAccount,
		request.SmtpSuccess,
		request.ResponseCode,
		request.ErrorCode,
		request.Description,
		request.IsPrimaryDomain,
		request.PrimaryDomain,
		request.AlternateEmail,
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

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	upsertEvent, err := emailevent.NewEmailUpsertEvent(a, request.Tenant, request.RawEmail, sourceFields, createdAtNotNil)
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

func (a *EmailAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case emailevent.EmailCreateV1:
		return a.onEmailCreate(event)
	case emailevent.EmailUpdateV1:
		return a.onEmailUpdated(event)
	case emailevent.EmailValidatedV2:
		return a.OnEmailValidated(event)
	case emailevent.EmailValidationFailedV1,
		emailevent.EmailValidatedV1,
		emailevent.EmailUpsertV1:
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

func (a *EmailAggregate) onEmailCreate(event eventstore.Event) error {
	var eventData emailevent.EmailCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.RawEmail = eventData.RawEmail
	if eventData.SourceFields.Available() {
		a.Email.Source = eventData.SourceFields
	} else {
		a.Email.Source.Source = eventData.Source
		a.Email.Source.SourceOfTruth = eventData.SourceOfTruth
		a.Email.Source.AppSource = eventData.AppSource
	}
	a.Email.CreatedAt = eventData.CreatedAt
	a.Email.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *EmailAggregate) onEmailUpdated(event eventstore.Event) error {
	var eventData emailevent.EmailUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Email.Source.SourceOfTruth = eventData.Source
	}
	a.Email.UpdatedAt = eventData.UpdatedAt
	a.Email.RawEmail = eventData.RawEmail
	return nil
}

func (a *EmailAggregate) OnEmailValidated(event eventstore.Event) error {
	var eventData emailevent.EmailValidatedEventV2
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Email.Email = eventData.Email
	return nil
}
