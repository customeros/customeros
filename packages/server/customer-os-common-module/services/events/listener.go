package events

import (
	"context"
	"errors"
	"reflect"

	"github.com/opentracing/opentracing-go"

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

	message, ok := input.(*dto.Event)
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

	return message, nil
}

func GetEventType[T any]() string {
	var t T
	eventType := reflect.TypeOf(t)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	return eventType.Name()
}
