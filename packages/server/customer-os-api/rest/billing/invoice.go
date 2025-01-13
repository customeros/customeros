// @openapi 3.0.0
package billing

import (
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

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
func GetInvoicesForOrganization(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetInvoicesForOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := handlers.ValidateTenant(c, ctx, span)

		// Extract organization ID from the path
		orgID := c.Param("id")
		if orgID == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Invalid organization ID"))
			return
		}

		// Check organization exists
		organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, tenant, orgID)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("Organization does not exist"))
			return
		}
		if organizationDbNode == nil {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("Organization does not exist"))
			return
		}
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		invoiceEntities, err := services.CommonServices.InvoiceService.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationEntity.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		response := InvoicesResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Invoices:     make([]InvoiceRecord, 0, len(*invoiceEntities)), // Pre-allocate slice
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
				record.PaymentLink = services.Cfg.CommonServices.InternalServices.CustomerOsApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
			}

			wg.Add(1)
			go func(invoiceEntity neo4jentity.InvoiceEntity, record InvoiceRecord) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire a spot
				defer func() { <-sem }() // Release spot in defer

				publicUrl, err := services.CommonServices.FileService.GetFilePublicUrl(ctx, invoiceEntity.RepositoryFileId)
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

		c.JSON(http.StatusOK, response)
	}
}
