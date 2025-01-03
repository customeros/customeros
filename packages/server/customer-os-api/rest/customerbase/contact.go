// @openapi 3.0.0
package customerbase

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
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
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// @Summary Create a new contact
// @Description Creates a contact from either an email address or linkedin_url
// @Tags CustomerBASE API
// @Accept json
// @Produce json
// @Param contact body ContactRecord false "Contact information"
// @Success 200 {object} SingleContactResponse "Successfully created single contact"
// @Failure 400 {object} rest.BaseResponse "Invalid request data"
// @Failure 401 {object} rest.BaseResponse "Unauthorized"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/contacts [post]
// @Security ApiKeyAuth
func CreateContact(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		handleJSONRequest(c, s)
	}
}

func handleJSONRequest(c *gin.Context, s *service.Services) {
	ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.handleJSONRequest", c.Request.Header)
	defer span.Finish()
	tracing.TagComponentRest(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	var contactRecord ContactRecord
	if err := c.BindJSON(&contactRecord); err == nil && (strings.TrimSpace(contactRecord.Email) != "" || strings.TrimSpace(contactRecord.LinkedInURL) != "") {
		err, errValue := validateContactRecord(&contactRecord)
		if err != nil {
			errMessage := fmt.Sprintf("%s | %s", errValue, err)
			span.LogFields(log.String("result.error", errMessage))
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage(errMessage))
			return
		}
		contactRecord.ContactId = processContact(c, s, contactRecord)
		c.JSON(http.StatusOK, SingleContactResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Contact:      contactRecord,
		})
		return
	}
}

func validateContactRecord(record *ContactRecord) (error, string) {
	var errValue string
	email := strings.TrimSpace(record.Email)
	linkedInUrl := strings.TrimSpace(record.LinkedInURL)
	if email == "" && linkedInUrl == "" {
		return errors.New("must provide either email or LinkedIn URL"), errValue
	}

	if email != "" {
		errValue = email
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

	if linkedInUrl != "" && !isValidLinkedinContactUrl(linkedInUrl) {
		errValue = linkedInUrl
		return errors.New("invalid LinkedIn URL format"), errValue
	}

	return nil, errValue
}

func processContact(ctx context.Context, s *service.Services, record ContactRecord) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Customerbase.processContact")
	defer span.Finish()
	tracing.TagComponentRest(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("email", record.Email), log.String("linkedin", record.LinkedInURL))

	linkedInUrl := strings.TrimSpace(record.LinkedInURL)
	email := strings.TrimSpace(record.Email)

	if linkedInUrl == "" && email == "" {
		span.LogFields(log.String("result", "No email or LinkedIn URL provided"))
		return ""
	}

	createdContactId := ""
	var err error
	if linkedInUrl != "" {
		createdContactId, err = s.CommonServices.ContactService.CreateContactByLinkedIn(ctx, nil, linkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save contact"))
			return ""
		}
	}

	if email != "" {
		if createdContactId == "" {
			createdContactId, err = s.CommonServices.ContactService.CreateContactByEmail(ctx, nil, email)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to save contact"))
				return ""
			}
		} else {
			associateEmailWithContact(ctx, s, email, createdContactId)
		}
	}
	return createdContactId
}

func associateEmailWithContact(ctx context.Context, s *service.Services, email, contactId string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Customerbase.associateEmailWithContact")
	defer span.Finish()
	tracing.TagComponentRest(span)

	_, err := s.CommonServices.EmailService.Merge(ctx, nil, common.GetTenantFromContext(ctx),
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
