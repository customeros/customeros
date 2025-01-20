package enrich

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type EnrichHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewEnrichHandler(services *cosapi_services.Services, responseHandler *response.Response) *EnrichHandler {
	return &EnrichHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}
