package integrations

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type IntegrationHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewIntegrationHandler(services *cosapi_services.Services, responseHandler *response.Response) *IntegrationHandler {
	return &IntegrationHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}
