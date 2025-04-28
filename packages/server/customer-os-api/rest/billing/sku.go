// @openapi 3.0.0
package billing

import (
	"net/http"
	"sort"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get skus
// @Description Retrieves all SKUs
// @Tags Billing API
// @Accept json
// @Produce json
// @Success 200 {object} SkusResponse "List of SKUs retrieved successfully"
// @Failure 401 {object} rest.BaseResponse "Unauthorized - Invalid or missing API key"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /billing/v1/skus [get]
// @Security ApiKeyAuth
func (h *BillingHandler) GetSkus() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "GetSkus")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)

		skuEntities, err := h.services.CommonServices.PostgresRepositories.SkuRepository.GetAll(ctx, tenant, utils.TruePtr())
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		response := SkusResponse{
			Skus: make([]SkuRecord, 0, len(skuEntities)), // Pre-allocate slice
		}
		for _, skuEntity := range skuEntities {
			response.Skus = append(response.Skus, SkuRecord{
				ID:      skuEntity.ID,
				Name:    skuEntity.Name,
				Price:   skuEntity.Price,
				SkuType: skuEntity.Type.String(),
			})
		}

		// Sort skus by name ascending, ignoring case
		sort.Slice(response.Skus, func(i, j int) bool {
			return strings.ToLower(response.Skus[i].Name) < strings.ToLower(response.Skus[j].Name)
		})

		h.responseHandler.HandleSuccess(c, response)
	}
}
