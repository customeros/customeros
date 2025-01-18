package base

import (
	"github.com/customeros/customeros/packages/server/events-processing-platform/config"
	"github.com/customeros/customeros/packages/server/events-processing-platform/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform/repository"
	"github.com/customeros/customeros/packages/server/events/eventstore"
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
