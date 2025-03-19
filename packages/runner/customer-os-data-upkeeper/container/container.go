package container

import (
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	agent_producers "github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_event_producers"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
)

type Container struct {
	Cfg            *config.Config
	Log            logger.Logger
	Repositories   *repository.Repositories
	CommonServices *commonservice.CommonServices
	AgentProducers *agent_producers.AgentProducers
}
