package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	s "github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"log"
	"net/http"
)

func AddQueryRoutes(rg *gin.RouterGroup, cosService s.CustomerOSService, redisService s.RedisService) {
	rg.POST("/query", func(ctx *gin.Context) {
		isActive, tenant := redisService.GetKeyInfo(ctx, "tenantKey", ctx.Request.Header.Get("x-openline-tenant-key"))
		if isActive {
			var request model.ForwardQuery
			if err := ctx.BindJSON(&request); err != nil {
				log.Printf("unable to parse json: %v", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}
			response, err := cosService.ForwardQuery(tenant, &request.Query)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to forward query: %v", err)})
				return
			} else {
				ctx.Data(http.StatusOK, "application/json", response)
				return
			}
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("tenant key is not active")})
			return
		}
	})
}
