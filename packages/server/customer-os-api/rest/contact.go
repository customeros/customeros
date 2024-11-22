package rest

import (
	"context"
	"encoding/csv"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type ContactRecord struct {
	Email       string `json:"email" csv:"email"`
	LinkedInURL string `json:"linkedinUrl" csv:"linkedin_url"`
}

func CreateContact(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := validateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		contentType := c.GetHeader("Content-Type")

		switch {
		case strings.HasPrefix(contentType, "multipart/form-data"):
			createContactsFromCsvUpload(c, ctx, span, services, tenant)
		case strings.HasPrefix(contentType, "application/json"):
			createContactFromJson(c, ctx, span, services, tenant)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported Content-Type"})
		}
	}
}

func createContactsFromCsvUpload(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracing.TagComponentRest(span)

		file, err := validateAndOpenFile(c)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		if err := validateHeaders(c, reader); err != nil {
			return
		}

		processRecords(c, ctx, span, reader, services, tenant)
	}
}

func createContactFromJson(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracing.TagComponentRest(span)

		var record ContactRecord
		if err := c.BindJSON(&record); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		processContact(c, ctx, span, services, tenant, record)
	}
}

func validateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "API key invalid or expired"})
		return ""
	}
	return tenant
}

func validateAndOpenFile(c *gin.Context) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to parse file"})
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid file type"})
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func validateHeaders(c *gin.Context, reader *csv.Reader) error {
	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to read file"})
		return err
	}

	if headers[0] != "email" && headers[1] != "linkedin_url" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid headers, must be 'email' and 'linkedin_url'"})
		return errors.New("invalid headers")
	}
	return nil
}

func processRecords(c *gin.Context, ctx context.Context, span opentracing.Span, reader *csv.Reader, services *service.Services, tenant string) {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to read file"})
			return
		}

		contactRecord := ContactRecord{
			Email:       record[0],
			LinkedInURL: record[1],
		}

		if contactRecord.Email == "" && contactRecord.LinkedInURL == "" {
			continue
		}

		processContact(c, ctx, span, services, tenant, contactRecord)
	}
}

func processContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string, record ContactRecord) {
	emailEntity, contactId := findExistingContact(c, ctx, span, services, record.Email)
	if emailEntity == nil && contactId == "" {
		contactId = createNewContact(c, ctx, span, services, record.LinkedInURL)
	}

	if emailEntity == nil && record.Email != "" {
		associateEmail(c, ctx, span, services, tenant, record.Email, contactId)
	}
}

func findExistingContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, email string) (*neo4jentity.EmailEntity, string) {
	if email == "" {
		return nil, ""
	}

	emailEntity, err := services.EmailService.GetByEmailAddress(ctx, email)
	if err != nil {
		span.LogFields(log.String("result", "Failed to get email entity"))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get email entity"})
		return nil, ""
	}

	if emailEntity == nil {
		return nil, ""
	}

	contacts, err := services.ContactService.GetContactsForEmails(ctx, []string{emailEntity.Id})
	if err != nil {
		span.LogFields(log.String("result", "Failed to get contacts for email"))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to get contacts for email"})
		return nil, ""
	}

	if contacts != nil && len(*contacts) > 0 {
		return emailEntity, (*contacts)[0].Id
	}

	return emailEntity, ""
}

func createNewContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, linkedInURL string) string {
	contactId, err := services.CommonServices.ContactService.Save(ctx, nil, neo4jrepo.ContactFields{}, linkedInURL, neo4jmodel.ExternalSystem{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to save contact"))
		return ""
	}
	return contactId
}

func associateEmail(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant, email, contactId string) {
	_, err := services.CommonServices.EmailService.Merge(ctx, tenant,
		commonservice.EmailFields{
			Email:     email,
			Source:    neo4jentity.DataSourceOpenline,
			AppSource: constants.AppSourceCustomerOsApiRest,
		}, &commonservice.LinkWith{
			Type: commonModel.CONTACT,
			Id:   contactId,
		})
	if err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.String("result", "Failed to upsert email"))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upsert email"})
		return
	}
}
