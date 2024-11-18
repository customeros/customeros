package rest

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
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

func CreateContactsFromCsvUpload(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContactsFromCsvUpload", c.Request.Header)
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
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid headers, must be 'email' and 'linkedin_url'"})
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

			inputEmail := record[0]
			inputSocialUrl := record[1]

			if inputEmail == "" && inputSocialUrl == "" {
				continue
			}

			// if email provided, check if email exists in db
			var emailEntity *neo4jentity.EmailEntity
			if inputEmail != "" {
				emailEntity, err = services.EmailService.GetByEmailAddress(ctx, inputEmail)
				if err != nil {
					span.LogFields(tracingLog.String("result", "Failed to get email entity"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get email entity"})
					return
				}
			}

			emailId := ""
			contactId := ""

			// if email exists, get contact id associated with email
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

			// if contact not exists, create new contact
			if contactId == "" {
				contactId, err = services.CommonServices.ContactService.SaveContact(ctx, nil, neo4jrepo.ContactFields{}, inputSocialUrl, neo4jmodel.ExternalSystem{})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to save contact"))
					continue
				}
			}

			// associate email with created contact if email address provided
			if emailId == "" && inputEmail != "" {
				_, err := services.CommonServices.EmailService.Merge(ctx, tenant,
					commonservice.EmailFields{
						Email:     inputEmail,
						Source:    neo4jentity.DataSourceOpenline,
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
