package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/private"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

const InternalPath = "/internal/v1"

func registerInternalRoutes(ctx context.Context, r *gin.Engine, s *service.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/askAI", InternalPath),
		handler:   private.AskAI(s),
		routeType: RouteInternal,
		services:  s,
	})
}
