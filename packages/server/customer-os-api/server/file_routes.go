package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/files"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

const FilePath = "/files/v1"

func registerFileRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/files", FilePath),
		handler:   files.UploadFile(s, FilePath),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id", FilePath),
		handler:   files.GetFileByID(s, FilePath),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/download", FilePath),
		handler:   files.DownloadFile(s),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/base64", FilePath),
		handler:   files.GetBase64(s),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/public-url", FilePath),
		handler:   files.GetPublicURL(s),
		routeType: RouteFiles,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/files/:id/jwt", FilePath),
		handler:   files.GetJWT(s),
		routeType: RouteFiles,
		services:  s,
	})
}
