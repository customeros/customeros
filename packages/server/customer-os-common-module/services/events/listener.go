package events

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

// EventListener interface defines what all listeners must implement
type EventListener interface {
	Handle(ctx context.Context, baseEvent any) error
	GetEventType() string
	GetQueueName() string
}

// BaseEventListener provides common functionality for all listeners
type BaseEventListener struct {
	logger    logger.Logger
	eventType string
	queueName string
}

// NewBaseEventListener creates a new base event listener
func NewBaseEventListener(logger logger.Logger, eventType, queueName string) BaseEventListener {
	return BaseEventListener{
		logger:    logger,
		eventType: eventType,
		queueName: queueName,
	}
}

func (b BaseEventListener) GetEventType() string {
	return b.eventType
}

func (b BaseEventListener) GetQueueName() string {
	return b.queueName
}

func (b BaseEventListener) ValidateBaseEvent(ctx context.Context, input any) (*dto.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Events.ValidateEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set on context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	message, ok := input.(dto.Event)
	if !ok {
		err := errors.New("unable to cast to event type")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Event.EntityId == "" {
		err := errors.New("entity id is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Event.Tenant == "" {
		err := errors.New("tenant is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Event.EventType == "" {
		err := errors.New("event type is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &message, nil
}

func DecodeEventData[T any](ctx context.Context, event *dto.Event) (T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listener.DecodeEventData")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var decoded T
	if err := mapstructure.Decode(event.Event.Data, &decoded); err != nil {
		err = errors.Wrap(err, "failed to decode event data")
		tracing.LogObjectAsJson(span, "event", event)
		tracing.TraceErr(span, err)
		return decoded, err
	}

	return decoded, nil
}

func GetEventType[T any]() string {
	var t T
	eventType := reflect.TypeOf(t)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	return eventType.Name()
}
