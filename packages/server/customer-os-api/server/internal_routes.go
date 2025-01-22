package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/private"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

const InternalPath = "/internal/v1"

func registerInternalRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/askAI", InternalPath),
		handler:   h.AskAI.AskAI(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/integrations", InternalPath),
		handler:   h.PrivateIntegrations.GetIntegrations(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/integrations", InternalPath),
		handler:   h.PrivateIntegrations.CreateIntegration(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/settings/integrations/:identifier", InternalPath),
		handler:   h.PrivateIntegrations.DeleteIntegrations(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/mailboxes", InternalPath),
		handler:   private.GetMailboxes(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/personal_integrations", InternalPath),
		handler:   private.CreatePersonalIntegrations(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/personal_integrations", InternalPath),
		handler:   private.GetPersonalIntegrations(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/personal_integrations/:integrationName", InternalPath),
		handler:   private.GetPersonalIntegrationByName(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/tenant/settings/organizationStage/:id", InternalPath),
		handler:   private.CreateOrganizationStage(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/tenant/settings/apiKey", InternalPath),
		handler:   private.GetAPIKey(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/user/settings/oauth/:tenant", InternalPath),
		handler:   private.GetOAuthSettings(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/user/settings/slack", InternalPath),
		handler:   private.GetSlackSettings(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/revoke", InternalPath),
		handler:   private.Revoke(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/quickbooks/requestAccess", InternalPath),
		handler:   private.RequestAccessQuickbooks(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/quickbooks/oauth/callback", InternalPath),
		handler:   private.CallbackQuickbooks(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/slack/requestAccess", InternalPath),
		handler:   private.RequestAccessSlack(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/slack/oauth/callback", InternalPath),
		handler:   private.CallbackSlack(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/slack/revoke", InternalPath),
		handler:   private.RevokeSlack(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/mail/send", InternalPath),
		handler:   h.Mail.SendEmail(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/mail/:customerOSInternalIdentifier/track", InternalPath),
		handler:   h.Mail.TrackEmail(),
		routeType: RouteInternal,
		services:  s,
	})
}
