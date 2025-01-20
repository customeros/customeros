package verify

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type VerifyHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewVerifyHandler(services *cosapi_services.Services, responseHandler *response.Response) *VerifyHandler {
	return &VerifyHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}
