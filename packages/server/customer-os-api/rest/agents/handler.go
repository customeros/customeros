package agents

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type AgentHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewAgentHandler(services *cosapi_services.Services, responseHandler *response.Response) *AgentHandler {
	return &AgentHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}
