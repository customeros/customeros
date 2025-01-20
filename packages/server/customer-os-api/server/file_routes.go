package server

import (
	"context"
	"fmt"

	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/gin-gonic/gin"
)

const FilePath = "/files/v1"

func registerFileRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/files", FilePath),
		handler:   h.Files.UploadFile(FilePath),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id", FilePath),
		handler:   h.Files.GetFileByID(FilePath),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/download", FilePath),
		handler:   h.Files.DownloadFile(),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/base64", FilePath),
		handler:   h.Files.GetBase64(),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/public-url", FilePath),
		handler:   h.Files.GetPublicURL(),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/jwt", FilePath),
		handler:   h.Files.GetJWT(),
		routeType: RouteFiles,
		services:  s,
	})
}
