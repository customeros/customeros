package public

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/customerbase"
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

type ContactResponse struct {
	ContactID   string `json:"contactId"`
	LinkedinURL string `json:"linkedinUrl,omitempty"`
}

func (h *BrowserExtensionHandler) CreateContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "BrowserExtensionHandler.CreateContact")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		h.handleCreateContactJSONRequest(c)
	}
}

func (h *BrowserExtensionHandler) handleCreateContactJSONRequest(c *gin.Context) {
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "BrowserExtensionHandler.handleCreateContactJSONRequest")
	defer spans.Finish()

	var contactRecord customerbase.ContactRecord
	err := c.BindJSON(&contactRecord)
	if err != nil {
		spans.TraceError(err)
		errMessage := "cannot parse request"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &errMessage)
		return
	}

	if strings.TrimSpace(contactRecord.Email) == "" && strings.TrimSpace(contactRecord.LinkedInURL) == "" {
		errMessage := "email or linkedin_url must be provided"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &errMessage)
	}

	err, errValue := h.validateContactRecord(&contactRecord)
	if err != nil {
		errMessage := fmt.Sprintf("%s | %s", errValue, err)
		spans.LogKV("result.error", errMessage)
		h.responseHandler.HandleError(c, http.StatusBadRequest, &errMessage)
		return
	}

	resp := ContactResponse{}
	resp.ContactID, resp.LinkedinURL = h.processCreateContact(c.Request.Context(), contactRecord)
	h.responseHandler.HandleSuccess(c, resp)
}

func (h *BrowserExtensionHandler) processCreateContact(ctx context.Context, record customerbase.ContactRecord) (string, string) {
	spans, ctx := telemetry.StartRestSpan(ctx, "BrowserExtensionHandler.processCreateContact")
	defer spans.Finish()

	spans.LogKV("email", record.Email, "linkedin", record.LinkedInURL)

	linkedInUrl := strings.TrimSpace(record.LinkedInURL)

	if linkedInUrl == "" {
		spans.LogKV("result", "No email or LinkedIn URL provided")
		return "", ""
	}

	createdContactId := ""
	var err error
	createdContactId, linkedInUrl, err = h.services.CommonServices.ContactService.CreateContactByLinkedIn(ctx, nil, linkedInUrl)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to save contact"))
		return "", ""
	}

	h.services.CommonServices.Events.Publisher.PublishNotification(ctx, common.GetTenantFromContext(ctx), createdContactId, model.CONTACT, utils.NewEventCompletedDetails().WithCreate())

	return createdContactId, linkedInUrl
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

	patterns := []string{
		// Pattern for public profiles
		`^https?:\/\/(www\.)?linkedin\.com\/in\/[a-zA-Z0-9\-_.]{3,100}\/?$`,
		// Pattern for sales navigator profiles
		`^https?:\/\/(www\.)?linkedin\.com\/sales\/lead\/[A-Za-z0-9_-]{3,100}\/?$`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func (h *BrowserExtensionHandler) GetContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "BrowserExtensionHandler.GetContact")
		defer spans.Finish()

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

		var resp ContactResponse
		resp.ContactID = contactId
		h.responseHandler.HandleSuccess(c, resp)
		return
	}
}

func (h *BrowserExtensionHandler) TouchContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "BrowserExtensionHandler.TouchContact")
		defer spans.Finish()

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
