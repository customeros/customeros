package aggregate

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	PhoneNumberAggregateType eventstore.AggregateType = "phone_number"
)

type PhoneNumberAggregate struct {
	*aggregate.CommonTenantIdAggregate
	PhoneNumber *models.PhoneNumber
}

func NewPhoneNumberAggregateWithTenantAndID(tenant, id string) *PhoneNumberAggregate {
	phoneNumberAggregate := PhoneNumberAggregate{}
	phoneNumberAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(PhoneNumberAggregateType, tenant, id)
	phoneNumberAggregate.SetWhen(phoneNumberAggregate.When)
	phoneNumberAggregate.PhoneNumber = &models.PhoneNumber{}
	phoneNumberAggregate.Tenant = tenant

	return &phoneNumberAggregate
}

func (a *PhoneNumberAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *phonenumberpb.PassPhoneNumberValidationGrpcRequest:
		return nil, a.phoneNumberValidated(ctx, r)
	case *phonenumberpb.FailPhoneNumberValidationGrpcRequest:
		return nil, a.phoneNumberValidationFailed(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *PhoneNumberAggregate) phoneNumberValidated(ctx context.Context, request *phonenumberpb.PassPhoneNumberValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberTempAggregate.requestPhoneNumberValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	if request.PhoneNumber != a.PhoneNumber.RawPhoneNumber {
		span.LogFields(log.String("result", fmt.Sprintf("phone number does not match. validated %s, current %s", request.PhoneNumber, a.PhoneNumber.RawPhoneNumber)))
		return nil
	}

	event, err := events.NewPhoneNumberValidatedEvent(a, request.Tenant, request.PhoneNumber, request.E164, request.CountryCodeA2)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberValidatedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) phoneNumberValidationFailed(ctx context.Context, request *phonenumberpb.FailPhoneNumberValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberTempAggregate.requestPhoneNumberValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	if request.PhoneNumber != "" && request.PhoneNumber != a.PhoneNumber.RawPhoneNumber {
		span.LogFields(log.String("result", fmt.Sprintf("phone number does not match. validated %s, current %s", request.PhoneNumber, a.PhoneNumber.RawPhoneNumber)))
		return nil
	}

	event, err := events.NewPhoneNumberFailedValidationEvent(a, request.Tenant, request.PhoneNumber, request.CountryCodeA2, request.ErrorMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case events.PhoneNumberCreateV1:
		return a.onPhoneNumberCreate(event)
	case events.PhoneNumberUpdateV1:
		return a.onPhoneNumberUpdate(event)
	case events.PhoneNumberValidationSkippedV1:
		return a.OnPhoneNumberSkippedValidation(event)
	case events.PhoneNumberValidationFailedV1:
		return a.OnPhoneNumberFailedValidation(event)
	case events.PhoneNumberValidatedV1:
		return a.OnPhoneNumberValidated(event)
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *PhoneNumberAggregate) onPhoneNumberCreate(event eventstore.Event) error {
	var eventData events.PhoneNumberCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.RawPhoneNumber = eventData.RawPhoneNumber
	a.PhoneNumber.CreatedAt = eventData.CreatedAt
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	if eventData.SourceFields.Available() {
		a.PhoneNumber.Source = eventData.SourceFields
	} else {
		a.PhoneNumber.Source.Source = eventData.Source
		a.PhoneNumber.Source.SourceOfTruth = eventData.SourceOfTruth
		a.PhoneNumber.Source.AppSource = eventData.AppSource
	}
	return nil
}

func (a *PhoneNumberAggregate) onPhoneNumberUpdate(event eventstore.Event) error {
	var eventData events.PhoneNumberUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.PhoneNumber.Source.SourceOfTruth = eventData.Source
	}
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	a.PhoneNumber.RawPhoneNumber = eventData.RawPhoneNumber
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberSkippedValidation(event eventstore.Event) error {
	var eventData events.PhoneNumberSkippedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.SkipReason = eventData.Reason
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberFailedValidation(event eventstore.Event) error {
	var eventData events.PhoneNumberFailedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.ValidationError = eventData.ValidationError
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberValidated(event eventstore.Event) error {
	var eventData events.PhoneNumberValidatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.ValidationError = ""
	a.PhoneNumber.PhoneNumberValidation.SkipReason = ""
	a.PhoneNumber.E164 = eventData.E164
	a.PhoneNumber.CountryCodeA2 = eventData.CountryCodeA2
	a.PhoneNumber.UpdatedAt = eventData.ValidatedAt
	return nil
}
