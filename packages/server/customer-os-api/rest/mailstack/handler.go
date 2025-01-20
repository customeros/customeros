package mailstack

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type MailstackHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewMailstackHandler(services *cosapi_services.Services, responseHandler *response.Response) *MailstackHandler {
	return &MailstackHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}
