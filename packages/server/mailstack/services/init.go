package services

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services/email_filter"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
)

type Services struct {
	Cache              *caches.Cache
	EventsService      *events.EventsService
	EmailFilterService interfaces.EmailFilterService
	IMAPService        interfaces.IMAPService
}

func InitServices(rabbitmqURL string, log logger.Logger, repos *repository.Repositories) (*Services, error) {
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
		Cache:              caches.NewCommonCache(),
		EventsService:      events,
		EmailFilterService: email_filter.NewEmailFilterService(),
		IMAPService:        imap.NewIMAPService(repos),
	}

	return &services, nil
}
