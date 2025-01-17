package events

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
)

type EventsService struct {
	Publisher  *RabbitMQPublisher
	Subscriber *RabbitMQSubscriber
	Handlers   *HandlerRegistry
}

func NewEventsService(rabbitmqURL string, log logger.Logger) (*EventsService, error) {
	publisher, err := NewRabbitMQPublisher(rabbitmqURL, log)
	if err != nil {
		return nil, err
	}

	subscriber, err := NewRabbitMQSubscriber(rabbitmqURL, log)
	if err != nil {
		return nil, err
	}

	return &EventsService{
		Publisher:  publisher,
		Subscriber: subscriber,
		Handlers:   NewHandlerRegistry(),
	}, nil
}

func (s *EventsService) RegisterHandler(eventType interface{}, handler interfaces.EventHandler) {
	s.Handlers.handlers[handler.EventType] = handler
}

func (s *EventsService) GetHandler(eventTypeName string) (interfaces.EventHandler, bool) {
	handler, exists := s.Handlers.handlers[eventTypeName]
	return handler, exists
}

func (s *EventsService) Close() error {
	var errs []error

	if s.Publisher != nil {
		if err := s.Publisher.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if s.Subscriber != nil {
		if err := s.Subscriber.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing events service: %v", errs)
	}

	return nil
}
