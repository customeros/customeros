// @openapi 3.0.0
package billing

import (
	"net/http"
	"sort"
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type BillingHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewBillingHandler(services *cosapi_services.Services, responseHandler *response.Response) *BillingHandler {
	return &BillingHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

// @Summary Get organization's invoices
// @Description Retrieves all non-dry-run invoices for a specific organization, sorted by due date descending
// @Tags Billing API
// @Accept json
// @Produce json
// @Param id path string true "Organization ID or COS ID" example(org_123)
// @Success 200 {object} InvoicesResponse "List of invoices retrieved successfully"
// @Success 206 {object} InvoicesResponse "Invoices retrieved with some errors fetching public URLs"
// @Failure 400 {object} rest.BaseResponse "Invalid organization ID format"
// @Failure 401 {object} rest.BaseResponse "Unauthorized - Invalid or missing API key"
// @Failure 404 {object} rest.BaseResponse "Organization not found"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /billing/v1/organizations/{id}/invoices [get]
// @Security ApiKeyAuth
func (h *BillingHandler) GetInvoicesForOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetInvoicesForOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)

		// Extract organization ID from the path
		orgID := c.Param("id")
		if orgID == "" {
			message := "Invalid organization ID"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Check organization exists
		organizationDbNode, err := h.services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, tenant, orgID)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Organization does not exist"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}
		if organizationDbNode == nil {
			message := "Organization does not exist"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		invoiceEntities, err := h.services.CommonServices.InvoiceService.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationEntity.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		response := InvoicesResponse{
			Invoices: make([]InvoiceRecord, 0, len(*invoiceEntities)), // Pre-allocate slice
		}

		workerCount := 10
		recordsChan := make(chan *InvoiceRecord, len(*invoiceEntities))
		errChan := make(chan error, len(*invoiceEntities))
		var wg sync.WaitGroup
		sem := make(chan struct{}, workerCount)

		for _, invoiceEntity := range *invoiceEntities {
			record := InvoiceRecord{
				ID:            invoiceEntity.Id,
				Number:        invoiceEntity.Number,
				DueDate:       invoiceEntity.DueDate,
				InvoiceStatus: invoiceEntity.Status.String(),
				Amount:        invoiceEntity.TotalAmount,
				Currency:      invoiceEntity.Currency.String(),
			}

			if (invoiceEntity.Status == neo4jenum.InvoiceStatusDue || invoiceEntity.Status == neo4jenum.InvoiceStatusOverdue) &&
				(invoiceEntity.PaymentDetails.PaymentLink != "") {
				record.PaymentLink = h.services.Cfg.Common.Internal.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
			}

			wg.Add(1)
			go func(invoiceEntity neo4jentity.InvoiceEntity, record InvoiceRecord) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire a spot
				defer func() { <-sem }() // Release spot in defer

				publicUrl, err := h.services.CommonServices.FileService.GetFilePublicUrl(ctx, invoiceEntity.RepositoryFileId)
				if err != nil {
					errChan <- errors.Wrap(err, "failed to get invoice public url")
					return
				}
				record.PublicUrl = publicUrl
				recordsChan <- &record
			}(invoiceEntity, record)
		}

		go func() {
			wg.Wait()
			close(recordsChan)
			close(errChan)
		}()

		// Collect results
		for record := range recordsChan {
			response.Invoices = append(response.Invoices, *record)
		}

		// Check for errors
		if len(errChan) > 0 {
			err := <-errChan
			tracing.TraceErr(span, err)
		}

		// Sort invoices by due date descending
		sort.Slice(response.Invoices, func(i, j int) bool {
			return response.Invoices[i].DueDate.After(response.Invoices[j].DueDate)
		})

		h.responseHandler.HandleSuccess(c, response)
	}
}
