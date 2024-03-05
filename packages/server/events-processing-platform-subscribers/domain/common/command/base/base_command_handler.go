package base

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type BaseCommandHandler struct {
	Log          logger.Logger
	Cfg          *config.Config
	Es           eventstore.AggregateStore
	Repositories *repository.Repositories
}

func NewBaseCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *BaseCommandHandler {
	return &BaseCommandHandler{Log: log, Cfg: cfg, Es: es}
}
