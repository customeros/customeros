package public

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/customerbase"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type BrowserExtensionHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewBrowserExtensionHandler(services *cosapi_services.Services, responseHandler *response.Response) *BrowserExtensionHandler {
	return &BrowserExtensionHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

func (h *BrowserExtensionHandler) CreateContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BrowserExtensionHandler.CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		h.handleCreateContactJSONRequest(c)
	}
}

func (h *BrowserExtensionHandler) handleCreateContactJSONRequest(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "BrowserExtensionHandler.handleCreateContactJSONRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var contactRecord customerbase.ContactRecord
	if err := c.BindJSON(&contactRecord); err == nil && (strings.TrimSpace(contactRecord.Email) != "" || strings.TrimSpace(contactRecord.LinkedInURL) != "") {
		err, errValue := h.validateContactRecord(&contactRecord)
		if err != nil {
			errMessage := fmt.Sprintf("%s | %s", errValue, err)
			span.LogFields(log.String("result.error", errMessage))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &errMessage)
			return
		}
		contactRecord.ContactId = h.processCreateContact(c.Request.Context(), contactRecord)
		resp := customerbase.SingleContactResponse{
			Contact: contactRecord,
		}
		h.responseHandler.HandleSuccess(c, resp)
		return
	}
}

func (h *BrowserExtensionHandler) processCreateContact(ctx context.Context, record customerbase.ContactRecord) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrowserExtensionHandler.processCreateContact")
	defer span.Finish()
	tracing.TagComponentRest(span)
	tenant := common.GetTenantFromContext(ctx)

	tracing.TagTenant(span, tenant)

	span.LogFields(log.String("email", record.Email), log.String("linkedin", record.LinkedInURL))

	linkedInUrl := strings.TrimSpace(record.LinkedInURL)

	if linkedInUrl == "" {
		span.LogFields(log.String("result", "No email or LinkedIn URL provided"))
		return ""
	}

	createdContactId := ""
	var err error
	if linkedInUrl != "" {
		createdContactId, err = h.services.CommonServices.ContactService.CreateContactByLinkedIn(ctx, nil, linkedInUrl)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save contact"))
			return ""
		}

		h.services.CommonServices.Events.Publisher.PublishNotification(ctx, tenant, createdContactId, model.CONTACT, utils.NewEventCompletedDetails().WithCreate())
	}

	return createdContactId
}

func (h *BrowserExtensionHandler) validateContactRecord(record *customerbase.ContactRecord) (error, string) {
	var errValue string
	linkedInUrl := strings.TrimSpace(record.LinkedInURL)
	if linkedInUrl == "" {
		return errors.New("must provide either email or LinkedIn URL"), errValue
	}

	if linkedInUrl != "" && !isValidLinkedinContactUrl(linkedInUrl) {
		errValue = linkedInUrl
		return errors.New("invalid LinkedIn URL format"), errValue
	}

	return nil, errValue
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

func (h *BrowserExtensionHandler) GetContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BrowserExtensionHandler.GetContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// get linked in query param
		linkedInUrl := c.Query("linkedin")
		if linkedInUrl == "" {
			h.responseHandler.HandleError(c, http.StatusBadRequest, utils.StringPtr("linkedin param is required"))
			return
		}

		contactFound, contactId, err := h.services.CommonServices.ContactService.CheckContactExistsWithLinkedIn(c.Request.Context(), linkedInUrl, "", "")
		if err != nil {
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		if !contactFound {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		var contactRecord customerbase.ContactRecord
		contactRecord.ContactId = contactId
		resp := customerbase.SingleContactResponse{
			Contact: contactRecord,
		}
		h.responseHandler.HandleSuccess(c, resp)
		return
	}
}

func (h *BrowserExtensionHandler) TouchContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BrowserExtensionHandler.GetContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// get id path param
		contactId := c.Param("id")
		if contactId == "" {
			h.responseHandler.HandleError(c, http.StatusBadRequest, utils.StringPtr("id param is missing"))
			return
		}

		err := h.services.CommonServices.ContactService.TouchContact(c.Request.Context(), nil, contactId)
		if err != nil {
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, nil)
		return
	}
}
