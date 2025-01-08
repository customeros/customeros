package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/flows"
	integrations "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/flows_integrations"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/public"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func register(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
}

func registerPublicRoutes(ctx context.Context, r *gin.Engine, s *service.Services) {
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
		handler:   public.RevealWebsiteVisitor(s),
		routeType: RoutePublic,
	})
}
