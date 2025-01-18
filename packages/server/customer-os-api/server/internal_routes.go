package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/private"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

const InternalPath = "/internal/v1"

func registerInternalRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/askAI", InternalPath),
		handler:   private.AskAI(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/settings/intergrations", InternalPath),
		handler:   private.GetIntegrations(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/settings/intergrations", InternalPath),
		handler:   private.CreateIntegration(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/settings/intergrations/:identifier", InternalPath),
		handler:   private.DeleteIntegrations(s),
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
		path:      fmt.Sprintf("%s/tenant/settings/organizationStage/:id", InternalPath),
		handler:   private.CreateOrganizationStage(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/tenant/settings/apiKey", InternalPath),
		handler:   private.GetAPIKey(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/user/settings/oauth/:tenant", InternalPath),
		handler:   private.GetOAuthSettings(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/user/settings/slack", InternalPath),
		handler:   private.GetSlackSettings(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/mail/send", InternalPath),
		handler:   private.SendEmail(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/mail/:customerOSInternalIdentifier/track", InternalPath),
		handler:   private.TrackEmail(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/revoke", InternalPath),
		handler:   private.Revoke(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/slack/requestAccess", InternalPath),
		handler:   private.RequestAccessSlack(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/slack/oauth/callback", InternalPath),
		handler:   private.CallbackSlack(s),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/slack/revoke", InternalPath),
		handler:   private.RevokeSlack(s),
		routeType: RouteInternal,
		services:  s,
	})
}
