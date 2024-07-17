package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"net/http"
)

func addTrackingRoutes(rg *gin.RouterGroup, services *service.Services) {
	rg.POST("", func(ginContext *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "Tracking.track", ginContext.Request.Header)
		defer span.Finish()

		origin := ginContext.GetHeader("Origin")
		referer := ginContext.GetHeader("Referer")
		userAgent := ginContext.GetHeader("User-Agent")

		span.LogFields(tracingLog.String("origin", origin))
		span.LogFields(tracingLog.String("referer", referer))
		span.LogFields(tracingLog.String("userAgent", userAgent))

		if origin == "" || referer == "" || userAgent == "" {
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		tenant, err := services.CommonServices.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(ctx, origin)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil || *tenant == "" {
			span.LogFields(tracingLog.String("result.info", "tenant not found for origin"))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		tracking := entity.Tracking{}
		if err := ginContext.BindJSON(&tracking); err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		tracking.Tenant = *tenant
		tracking.State = entity.TrackingIdentificationStateNew

		_, err = services.CommonServices.PostgresRepositories.TrackingRepository.Store(ctx, tracking)
		if err != nil {
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to store tracking: %v", err.Error()),
			})
			return
		}

		ginContext.JSON(http.StatusOK, gin.H{})
	})

}
