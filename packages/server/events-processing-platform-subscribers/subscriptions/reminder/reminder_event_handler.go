package reminder

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"golang.org/x/net/context"
)

type ReminderEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	cfg          config.Config
}

func NewReminderEventHandler(log logger.Logger, repositories *repository.Repositories, cfg config.Config) *ReminderEventHandler {
	return &ReminderEventHandler{
		log:          log,
		repositories: repositories,
		cfg:          cfg,
	}
}

func (h *ReminderEventHandler) onReminderCreateV1(ctx context.Context, evt eventstore.Event) error {
	return nil
}

func (h *ReminderEventHandler) onReminderUpdateV1(ctx context.Context, evt eventstore.Event) error {
	return nil
}
