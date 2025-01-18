package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/flows"
	integrations "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/flows_integrations"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/public"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

func registerPublicRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/invoice/:invoiceId/pay",
		handler:   public.RedirectToPayInvoice(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/invoice/:invoiceId/paymentLink",
		handler:   public.GetInvoicePaymentLink(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/dmarc", WebhooksPath),
		handler:   integrations.PostmarkDMARCMonitor(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/:tenantId/i/:integrationId", FlowsPath),
		handler:   flows.HandleWebhook(s, FlowsPath),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/v1/l",
		handler:   public.TrackLinkRequest(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/v1/s",
		handler:   public.TrackOpenRequest(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/v1/u",
		handler:   public.TrackUnsubscribeRequest(s),
		routeType: RoutePublic,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      "/reveal",
		handler:   public.RevealWebsiteEvents(s),
		routeType: RoutePublic,
	})
}
