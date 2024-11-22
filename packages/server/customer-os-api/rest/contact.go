package rest

import (
	"context"
	"encoding/csv"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
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

type ContactResult struct {
	Status      string `json:"status"`              // success or error
	Message     string `json:"message,omitempty"`   // error message or success details
	ContactId   string `json:"contactId,omitempty"` // returned on successful contact creation/update
	Email       string `json:"email,omitempty"`
	LinkedInURL string `json:"linkedinUrl,omitempty"`
}

type ContactsResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	Contacts []ContactResult `json:"contacts,omitempty"`
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
			handleCSVUpload(c, ctx, span, services, tenant)
		case strings.HasPrefix(contentType, "application/json"):
			handleJSONRequest(c, ctx, span, services, tenant)
		default:
			sendError(c, http.StatusBadRequest, "Unsupported Content-Type")
		}
	}
}

func handleCSVUpload(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string) {
	file, err := validateAndOpenFile(c)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if err := validateFileHeaders(c, reader); err != nil {
		return
	}

	results := processCSVRecords(c, ctx, span, reader, services, tenant)
	c.JSON(http.StatusOK, ContactsResponse{
		Status:   "success",
		Contacts: results,
	})
}

func handleJSONRequest(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string) {
	// Try single contact first
	var singleContact ContactRecord
	if err := c.BindJSON(&singleContact); err == nil {
		if singleContact.Email != "" || singleContact.LinkedInURL != "" {
			result := validateAndProcessContact(c, ctx, span, services, tenant, singleContact)
			c.JSON(http.StatusOK, result)
			return
		}
	}

	// Try multiple contacts
	var request struct {
		Contacts []ContactRecord `json:"contacts"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		sendError(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(request.Contacts) == 0 {
		sendError(c, http.StatusBadRequest, "No contacts provided")
		return
	}

	results := make([]ContactResult, 0, len(request.Contacts))
	for _, record := range request.Contacts {
		result := validateAndProcessContact(c, ctx, span, services, tenant, record)
		results = append(results, result)
	}

	c.JSON(http.StatusOK, ContactsResponse{
		Status:   "success",
		Contacts: results,
	})
}

func validateAndProcessContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string, record ContactRecord) ContactResult {
	result := ContactResult{
		Status:      "success",
		Email:       record.Email,
		LinkedInURL: record.LinkedInURL,
	}

	if record.Email == "" && record.LinkedInURL == "" {
		return createError("Must provide either email or LinkedIn URL")
	}

	if record.Email != "" {
		emailSyntax := mailvalidate.ValidateEmailSyntax(record.Email)
		switch {
		case !emailSyntax.IsValid:
			return createError("Invalid email format")
		case emailSyntax.IsRoleAccount:
			return createError("Email is a role account")
		case emailSyntax.IsSystemGenerated:
			return createError("Email is system generated")
		default:
			record.Email = emailSyntax.CleanEmail
			result.Email = emailSyntax.CleanEmail
		}
	}

	if record.LinkedInURL != "" && !isValidLinkedinUrl(record.LinkedInURL) {
		return createError("Invalid LinkedIn URL format")
	}

	contactId := processContact(c, ctx, span, services, tenant, record)
	if contactId == "" {
		return createError("Failed to process contact")
	}

	result.ContactId = contactId
	return result
}

func validateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		sendError(c, http.StatusUnauthorized, "API key invalid or expired")
		return ""
	}
	return tenant
}

func validateAndOpenFile(c *gin.Context) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sendError(c, http.StatusBadRequest, "Failed to parse file")
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		sendError(c, http.StatusBadRequest, "Invalid file type")
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func validateFileHeaders(c *gin.Context, reader *csv.Reader) error {
	headers, err := reader.Read()
	if err != nil {
		sendError(c, http.StatusBadRequest, "Failed to read file")
		return err
	}

	hasEmail := false
	hasLinkedIn := false
	for _, header := range headers {
		if header == "email" {
			hasEmail = true
		}
		if header == "linkedin_url" {
			hasLinkedIn = true
		}
	}

	if !hasEmail || !hasLinkedIn {
		sendError(c, http.StatusBadRequest, "Missing required headers: email, linkedin_url")
		return errors.New("invalid headers")
	}
	return nil
}

func processCSVRecords(c *gin.Context, ctx context.Context, span opentracing.Span, reader *csv.Reader, services *service.Services, tenant string) []ContactResult {
	results := make([]ContactResult, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			results = append(results, createError("Failed to read record"))
			continue
		}

		contactRecord := ContactRecord{
			Email:       record[0],
			LinkedInURL: record[1],
		}

		if contactRecord.Email == "" && contactRecord.LinkedInURL == "" {
			continue
		}

		result := validateAndProcessContact(c, ctx, span, services, tenant, contactRecord)
		results = append(results, result)
	}

	return results
}

func processContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string, record ContactRecord) string {
	emailEntity, contactId := findExistingContact(c, ctx, span, services, record.Email)
	if emailEntity == nil && contactId == "" {
		contactId = createNewContact(c, ctx, span, services, record.LinkedInURL)
	}

	if emailEntity == nil && record.Email != "" {
		associateEmail(c, ctx, span, services, tenant, record.Email, contactId)
	}
	return contactId
}

func findExistingContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, email string) (*neo4jentity.EmailEntity, string) {
	if email == "" {
		return nil, ""
	}

	emailEntity, err := services.EmailService.GetByEmailAddress(ctx, email)
	if err != nil {
		span.LogFields(log.String("result", "Failed to get email entity"))
		return nil, ""
	}

	if emailEntity == nil {
		return nil, ""
	}

	contacts, err := services.ContactService.GetContactsForEmails(ctx, []string{emailEntity.Id})
	if err != nil {
		span.LogFields(log.String("result", "Failed to get contacts for email"))
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
	}
}

func isValidLinkedinUrl(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "http") {
		s = "https://" + s
	}

	pattern := `^https?:\/\/(www\.)?linkedin\.com\/in\/[a-zA-Z0-9\-_.]{3,100}\/?$`
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, ContactResult{
		Status:  "error",
		Message: message,
	})
}

func createError(message string) ContactResult {
	return ContactResult{
		Status:  "error",
		Message: message,
	}
}
