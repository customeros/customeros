package rest

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

func CreateContact(services *service.Services, grpcClients *grpc_client.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContact", c.Request.Header)
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

		// Parse the uploaded CSV file
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to parse file"})
			return
		}
		defer file.Close()

		// Validate file type
		if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid file type"})
			return
		}

		// Parse the CSV file
		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to read file"})
			return
		}

		if headers[0] != "email" && headers[1] != "linkedin_url" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid headers"})
			return
		}

		// Read and validate each email
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to read file"})
				return
			}

			contactEmail := record[0]
			contactSocialUrl := record[1]

			emailEntity, err := services.EmailService.GetByEmailAddress(ctx, contactEmail)
			if err != nil {
				span.LogFields(tracingLog.String("result", "Failed to get email entity"))
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get email entity"})
				return
			}

			emailId := ""
			contactId := ""

			if emailEntity != nil {
				emailId = emailEntity.Id
				contactsWithEmail, err := services.ContactService.GetContactsForEmails(ctx, []string{emailEntity.Id})
				if err != nil {
					span.LogFields(tracingLog.String("result", "Failed to get contacts for email"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get contacts for email"})
					return
				}

				if contactsWithEmail != nil && len(*contactsWithEmail) > 0 {
					contactId = (*contactsWithEmail)[0].Id
				}
			}

			if contactId == "" {
				contactId, err = services.CommonServices.ContactService.CreateContact(ctx, tenant, neo4jrepo.ContactFields{}, contactSocialUrl, neo4jmodel.ExternalSystem{})
				if err != nil {
					span.LogFields(tracingLog.String("result", "Failed to upsert contact"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upsert contact"})
					return
				}
			}

			if emailId == "" {
				_, err := services.CommonServices.EmailService.Merge(ctx, tenant,
					commonservice.EmailFields{
						Email:     contactEmail,
						Source:    neo4jentity.DataSourceOpenline.String(),
						AppSource: constants.AppSourceCustomerOsApiRest,
					}, &commonservice.LinkWith{
						Type: commonModel.CONTACT,
						Id:   contactId,
					})
				if err != nil {
					tracing.TraceErr(span, err)
					span.LogFields(tracingLog.String("result", "Failed to upsert email"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upsert email"})
					return
				}
			}
		}
	}
}
