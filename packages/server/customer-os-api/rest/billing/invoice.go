package billing

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"sort"
	"sync"
)

// GetInvoicesForOrganization retrieves the list of invoices for a given organization
// @Summary Get organization invoices
// @Description Retrieves a list of invoices for the organization with the given ID
// @Tags Billing API
// @Accept  json
// @Produce  json
// @Param   id   path     string  true  "Organization ID or Organization COS ID"
// @Success 200  {array}  InvoiceResponse "List of invoices for the organization"
// @Failure 400  "Invalid organization ID"
// @Failure 401  "Unauthorized"
// @Failure 404  "Organization not found"
// @Failure 500  "Internal server error"
// @Router /billing/v1/organizations/{id}/invoices [get]
// @Security ApiKeyAuth
func GetInvoicesForOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetInvoicesForOrganization", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Extract organization ID from the path
		orgID := c.Param("id")
		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid organization ID"})
			span.LogFields(tracingLog.String("result", "Invalid organization ID"))
			return
		}

		// Check organization exists
		organizationDbNode, err := services.Repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId(ctx, tenant, orgID)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}
		if organizationDbNode == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Organization not found"})
			return
		}
		organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		invoiceEntities, err := services.CommonServices.InvoiceService.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationEntity.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal server error"})
			return
		}

		response := InvoicesResponse{
			Status: "success",
		}

		workerCount := 10
		invoicesChan := make(chan *InvoiceResponse, len(*invoiceEntities))
		errChan := make(chan error, len(*invoiceEntities))
		var wg sync.WaitGroup
		sem := make(chan struct{}, workerCount) // Semaphore to limit workers

		for _, invoiceEntity := range *invoiceEntities {
			invoiceResponse := InvoiceResponse{
				ID:            invoiceEntity.Id,
				Number:        invoiceEntity.Number,
				DueDate:       invoiceEntity.DueDate,
				InvoiceStatus: invoiceEntity.Status.String(),
				Amount:        invoiceEntity.TotalAmount,
				Currency:      invoiceEntity.Currency.String(),
			}
			if (invoiceEntity.Status == neo4jenum.InvoiceStatusDue || invoiceEntity.Status == neo4jenum.InvoiceStatusOverdue) &&
				(invoiceEntity.PaymentDetails.PaymentLink != "") {
				invoiceResponse.PaymentLink = services.Cfg.InternalServices.CustomerOsApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
			}

			wg.Add(1)
			go func(invoiceEntity neo4jentity.InvoiceEntity, invoiceResponse InvoiceResponse) {
				defer wg.Done()
				sem <- struct{}{} // Acquire a spot

				publicUrl, err := services.FileStoreApiService.GetFilePublicUrl(ctx, tenant, invoiceEntity.RepositoryFileId)
				if err != nil {
					errChan <- errors.New("failed to get invoice public url: " + err.Error())
					<-sem // Release a spot
					return
				}
				invoiceResponse.PublicUrl = publicUrl
				invoicesChan <- &invoiceResponse
				<-sem // Release a spot
			}(invoiceEntity, invoiceResponse)
		}

		go func() {
			wg.Wait()
			close(invoicesChan)
			close(errChan)
		}()

		// Collect results
		for invoiceResponse := range invoicesChan {
			response.Invoices = append(response.Invoices, *invoiceResponse)
		}

		// Check for errors
		if len(errChan) > 0 {
			tracing.TraceErr(span, <-errChan)
		}

		// sort returned invoices by due date descending
		sort.Slice(response.Invoices, func(i, j int) bool {
			// Compare the due dates, sorting by descending order (latest first)
			return (response.Invoices)[i].DueDate.After((response.Invoices)[j].DueDate)
		})

		c.JSON(http.StatusOK, response)
	}
}
