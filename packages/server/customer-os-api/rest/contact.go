package rest

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
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
				ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
				contactGrpcResponse, err := utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
					return grpcClients.ContactClient.UpsertContact(ctx, &contactpb.UpsertContactGrpcRequest{
						Tenant:    tenant,
						SocialUrl: contactSocialUrl,
					})
				})
				if err != nil {
					span.LogFields(tracingLog.String("result", "Failed to upsert contact"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upsert contact"})
					return
				}

				contactId = contactGrpcResponse.Id
			}

			if emailId == "" {
				contact := commonModel.CONTACT.String()
				_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
					return grpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
						Tenant:       tenant,
						RawEmail:     contactEmail,
						LinkWithId:   &contactId,
						LinkWithType: &contact,
					})
				})

				if err != nil {
					span.LogFields(tracingLog.String("result", "Failed to upsert email"))
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upsert email"})
					return
				}
			}
		}
	}
}