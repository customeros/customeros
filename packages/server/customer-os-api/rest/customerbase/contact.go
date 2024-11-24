package customerbase

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
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

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		handleJSONRequest(httpContext)
	}
}

func handleJSONRequest(ctx rest.HTTPContext) {
	// Try single contact
	var contact ContactRecord
	if err := ctx.GinContext.BindJSON(&contact); err == nil && (contact.Email != "" || contact.LinkedInURL != "") {
		err, errValue := validateContact(&contact)
		if err != nil {
			errMessage := fmt.Sprintf("%s | %s", errValue, err)
			rest.SendError(ctx.GinContext, http.StatusBadRequest, rest.ErrBadRequest.WithMessage(errMessage))

		}
		ctx.GinContext.JSON(http.StatusOK, SingleContactResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Contact:      contact,
		})
		return
	}
}

func validateContact(record *ContactRecord) (error, string) {
	var errValue string
	if record.Email == "" && record.LinkedInURL == "" {
		return errors.New("must provide either email or LinkedIn URL"), errValue
	}

	if record.Email != "" {
		errValue = record.Email
		emailSyntax := mailvalidate.ValidateEmailSyntax(record.Email)
		switch {
		case !emailSyntax.IsValid:
			return errors.New("invalid email format"), errValue
		case emailSyntax.IsRoleAccount:
			return errors.New("email is a role account"), errValue
		case emailSyntax.IsSystemGenerated:
			return errors.New("email is system generated"), errValue
		default:
			record.Email = emailSyntax.CleanEmail
		}
	}

	if record.LinkedInURL != "" && !isValidLinkedinContactUrl(record.LinkedInURL) {
		errValue = record.LinkedInURL
		return errors.New("invalid LinkedIn URL format"), errValue
	}

	return nil, errValue
}

func processContact(ctx rest.HTTPContext, record ContactRecord) string {
	emailEntity, contactId := findExistingContact(ctx, record.Email)
	if emailEntity == nil && contactId == "" {
		contactId = createNewContact(ctx, record.LinkedInURL)
	}

	if emailEntity == nil && record.Email != "" {
		associateEmail(ctx, record.Email, contactId)
	}
	return contactId
}

func findExistingContact(ctx rest.HTTPContext, email string) (*neo4jentity.EmailEntity, string) {
	if email == "" {
		return nil, ""
	}
	emailEntity, err := ctx.Services.EmailService.GetByEmailAddress(*ctx.ServiceContext, email)
	if err != nil {
		ctx.Span.LogFields(log.String("result", "Failed to get email entity"))
		return nil, ""
	}

	if emailEntity == nil {
		return nil, ""
	}

	contacts, err := ctx.Services.ContactService.GetContactsForEmails(*ctx.ServiceContext, []string{emailEntity.Id})
	if err != nil {
		ctx.Span.LogFields(log.String("result", "Failed to get contacts for email"))
		return nil, ""
	}

	if contacts != nil && len(*contacts) > 0 {
		return emailEntity, (*contacts)[0].Id
	}

	return emailEntity, ""
}

func createNewContact(ctx rest.HTTPContext, linkedInURL string) string {
	contactId, err := ctx.Services.CommonServices.ContactService.CreateContactByLinkedIn(*ctx.ServiceContext, nil, linkedInURL)
	if err != nil {
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to save contact"))
		return ""
	}
	return contactId
}

func associateEmail(ctx rest.HTTPContext, email, contactId string) {
	_, err := ctx.Services.CommonServices.EmailService.Merge(*ctx.ServiceContext, ctx.Tenant,
		commonservice.EmailFields{
			Email:     email,
			Source:    neo4jentity.DataSourceOpenline,
			AppSource: constants.AppSourceCustomerOsApiRest,
		}, &commonservice.LinkWith{
			Type: commonModel.CONTACT,
			Id:   contactId,
		})
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		ctx.Span.LogFields(log.String("result", "Failed to upsert email"))
	}
}

func isValidLinkedinContactUrl(s string) bool {
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
