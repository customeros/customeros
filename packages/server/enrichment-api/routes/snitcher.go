package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
)

func snitcherData(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "snitcherData")
		defer span.Finish()

		ipAddress, exists := c.GetQuery("ipAddress")
		if !exists || ipAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ipAddress is required"})
			return
		}

		data, err := s.IPIdentityService.IPIdentity(ctx, ipAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get snitcher data"})
			return
		}

		c.JSON(http.StatusOK, model.SnitcherDataResponse{
			Status:       "success",
			CompanyFound: data.CompanyFound(),
			Data:         data.Company,
		})
		return
	}
}
