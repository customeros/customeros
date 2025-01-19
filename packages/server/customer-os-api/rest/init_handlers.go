package handlers

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/public"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type RestHandlers struct {
	WebsiteTrackerEvents *public.WebsiteTrackerEventsHandler
}

func InitRestHandlers(services *cosapi_services.Services) *RestHandlers {
	handlers := RestHandlers{
		WebsiteTrackerEvents: public.NewWebsiteTrackerEventsHandler(services),
	}

	return &handlers
}
