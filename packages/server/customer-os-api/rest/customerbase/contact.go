package customerbase

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// @Summary Create a new contact
// @Description Creates a contact from either JSON or CSV upload
// @Tags CustomerBASE API
// @Accept json,multipart/form-data
// @Produce json
// @Param file formData file false "CSV file with contact data"
// @Param contact body ContactRecord false "Contact information"
// @Success 200 {object} ContactResult
// @Success 201 {object} ContactsResponse
// @Failure 400 {object} BaseResponse
// @Failure 401 {object} BaseResponse
// @Failure 500 {object} BaseResponse
// @Router /customerbase/v1/contacts [post]
// @Security ApiKeyAuth
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
	if len(results) == 0 {
		sendError(c, http.StatusBadRequest, "No valid contacts found in file")
		return
	}

	c.JSON(http.StatusCreated, ContactsResponse{
		Status:   "success",
		Contacts: results,
	})
}

func handleJSONRequest(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string) {
	// Try single contact
	var singleContact ContactRecord
	if err := c.BindJSON(&singleContact); err == nil && (singleContact.Email != "" || singleContact.LinkedInURL != "") {
		result := validateAndProcessContact(c, ctx, span, services, tenant, singleContact)
		c.JSON(http.StatusOK, result)
		return
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

	c.JSON(http.StatusCreated, ContactsResponse{
		Status:   "success",
		Contacts: results,
	})
}

func validateAndProcessContact(c *gin.Context, ctx context.Context, span opentracing.Span, services *service.Services, tenant string, record ContactRecord) ContactResult {
	if err := validateContact(&record); err != nil {
		return ContactResult{
			BaseResponse: BaseResponse{
				Status:  "error",
				Message: err.Error(),
			},
			Email:       record.Email,
			LinkedInURL: record.LinkedInURL,
		}
	}

	contactId := processContact(c, ctx, span, services, tenant, record)
	if contactId == "" {
		return ContactResult{
			BaseResponse: BaseResponse{
				Status:  "error",
				Message: "Failed to process contact",
			},
			Email:       record.Email,
			LinkedInURL: record.LinkedInURL,
		}
	}

	return ContactResult{
		BaseResponse: BaseResponse{
			Status: "success",
		},
		ContactId:   contactId,
		Email:       record.Email,
		LinkedInURL: record.LinkedInURL,
	}
}

func validateContact(record *ContactRecord) error {
	if record.Email == "" && record.LinkedInURL == "" {
		return errors.New("must provide either email or LinkedIn URL")
	}

	if record.Email != "" {
		emailSyntax := mailvalidate.ValidateEmailSyntax(record.Email)
		switch {
		case !emailSyntax.IsValid:
			return errors.New("invalid email format")
		case emailSyntax.IsRoleAccount:
			return errors.New("email is a role account")
		case emailSyntax.IsSystemGenerated:
			return errors.New("email is system generated")
		default:
			record.Email = emailSyntax.CleanEmail
		}
	}

	if record.LinkedInURL != "" && !isValidLinkedinUrl(record.LinkedInURL) {
		return errors.New("invalid LinkedIn URL format")
	}

	return nil
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
			results = append(results, ContactResult{
				BaseResponse: BaseResponse{
					Status:  "error",
					Message: "Failed to read record",
				},
			})
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
		contactId = createNewContact(ctx, span, services, record.LinkedInURL)
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

func createNewContact(ctx context.Context, span opentracing.Span, services *service.Services, linkedInURL string) string {
	contactId, err := services.CommonServices.ContactService.CreateContactByLinkedIn(ctx, nil, linkedInURL)
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
