package aggregate

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
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
	*aggregate.CommonTenantIdAggregate
	Email *models.Email
}

func NewEmailAggregateWithTenantAndID(tenant, id string) *EmailAggregate {
	emailAggregate := EmailAggregate{}
	emailAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(EmailAggregateType, tenant, id)
	emailAggregate.SetWhen(emailAggregate.When)
	emailAggregate.Email = &models.Email{}
	emailAggregate.Tenant = tenant

	return &emailAggregate
}

func (a *EmailAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *emailpb.PassEmailValidationGrpcRequest:
		return nil, a.emailValidated(ctx, r)
	case *emailpb.FailEmailValidationGrpcRequest:
		return nil, a.emailValidationFailed(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *EmailAggregate) emailValidated(ctx context.Context, request *emailpb.PassEmailValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailTempAggregate.requestEmailValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	if request.RawEmail != a.Email.RawEmail {
		span.LogFields(log.String("result", fmt.Sprintf("email does not match. validated %s, current %s", request.RawEmail, a.Email.RawEmail)))
		return nil
	}

	event, err := events.NewEmailValidatedEvent(a, request.Tenant, request.RawEmail, request.IsReachable, request.ErrorMessage,
		request.Domain, request.Username, request.Email, request.AcceptsMail, request.CanConnectSmtp, request.HasFullInbox, request.IsCatchAll,
		request.IsDisabled, request.IsValidSyntax, request.IsDeliverable, request.IsDisposable, request.IsRoleAccount)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailValidatedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(event)
}

func (a *EmailAggregate) emailValidationFailed(ctx context.Context, request *emailpb.FailEmailValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailTempAggregate.requestEmailValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	if request.RawEmail != "" && request.RawEmail != a.Email.RawEmail {
		span.LogFields(log.String("result", fmt.Sprintf("email does not match. validated %s, current %s", request.RawEmail, a.Email.RawEmail)))
		return nil
	}

	event, err := events.NewEmailFailedValidationEvent(a, request.Tenant, request.ErrorMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(event)
}

func (a *EmailAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case events.EmailCreateV1:
		return a.onEmailCreate(event)
	case events.EmailUpdateV1:
		return a.onEmailUpdated(event)
	case events.EmailValidationFailedV1:
		return a.OnEmailFailedValidation(event)
	case events.EmailValidatedV1:
		return a.OnEmailValidated(event)
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
	var eventData events.EmailCreateEvent
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
	var eventData events.EmailUpdateEvent
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

func (a *EmailAggregate) OnEmailFailedValidation(event eventstore.Event) error {
	var eventData events.EmailFailedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.EmailValidation.ValidationError = eventData.ValidationError
	return nil
}

func (a *EmailAggregate) OnEmailValidated(event eventstore.Event) error {
	var eventData events.EmailValidatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.Email = eventData.EmailAddress
	a.Email.EmailValidation.IsReachable = eventData.IsReachable
	a.Email.EmailValidation.ValidationError = eventData.ValidationError
	a.Email.EmailValidation.AcceptsMail = eventData.AcceptsMail
	a.Email.EmailValidation.CanConnectSmtp = eventData.CanConnectSmtp
	a.Email.EmailValidation.HasFullInbox = eventData.HasFullInbox
	a.Email.EmailValidation.IsCatchAll = eventData.IsCatchAll
	a.Email.EmailValidation.IsDeliverable = eventData.IsDeliverable
	a.Email.EmailValidation.IsDisabled = eventData.IsDisabled
	a.Email.EmailValidation.IsDisposable = eventData.IsDisposable
	a.Email.EmailValidation.IsRoleAccount = eventData.IsRoleAccount
	a.Email.EmailValidation.Domain = eventData.Domain
	a.Email.EmailValidation.IsValidSyntax = eventData.IsValidSyntax
	a.Email.EmailValidation.Username = eventData.Username
	return nil
}
