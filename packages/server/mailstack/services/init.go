package services

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
)

type Services struct {
	Cache         *caches.Cache
	EventsService *events.EventsService
	IMAPService   interfaces.IMAPService
}

func InitServices(rabbitmqURL string, log logger.Logger) (*Services, error) {
	// events
	publisherConfig := &events.PublisherConfig{
		MessageTTL:          events.DefaultMessageTTL,
		MaxRetries:          events.DefaultMaxRetries,
		PublishTimeout:      events.DefaultPublishTimeout,
		ReconnectBackoff:    events.DefaultReconnectBackoff,
		MaxReconnectBackoff: events.DefaultMaxReconnectBackoff,
	}

	subscriberConfig := &events.SubscriberConfig{
		MaxRetries:          events.DefaultMaxRetries,
		ReconnectBackoff:    events.DefaultReconnectBackoff,
		MaxReconnectBackoff: events.DefaultMaxReconnectBackoff,
	}

	events, err := events.NewEventsService(rabbitmqURL, log, publisherConfig, subscriberConfig)
	if err != nil {
		return nil, err
	}

	services := Services{
		Cache:         caches.NewCommonCache(),
		EventsService: events,
		IMAPService:   imap.NewIMAPService(),
	}

	return &services, nil
}
