package events

import (
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type EventsService struct {
	Publisher  *RabbitMQPublisher
	Subscriber *RabbitMQSubscriber
}

func NewEventsService(rabbitmqURL string, log logger.Logger, publisherConfig *PublisherConfig, subscriberConfig *SubscriberConfig) (*EventsService, error) {
	publisher, err := NewRabbitMQPublisher(rabbitmqURL, log, publisherConfig)
	if err != nil {
		return nil, err
	}

	subscriber, err := NewRabbitMQSubscriber(rabbitmqURL, log, subscriberConfig)
	if err != nil {
		return nil, err
	}

	return &EventsService{
		Publisher:  publisher,
		Subscriber: subscriber,
	}, nil
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
