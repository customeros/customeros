package private

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
)

type AskAIForEmailRequest struct {
	EmailFrom        string `json:"emailFrom"`
	FromEmailAddress string `json:"fromEmailAddress"`
	EmailTo          string `json:"emailTo"`
	ToEmailAddress   string `json:"toEmailAddress"`
	EmailBodyText    string `json:"emailBodyText"`
	EmailBodyHTML    string `json:"emailBodyHtml"`
}

type AskAIForEmailResponse struct {
	Answer data_fields.EmailResponse `json:"answer"`
}

// AskAI handles AI requests
func (h *AskAIHandler) AskAIForEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, "/askAIForEmail", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		// parse request
		var request AskAIForEmailRequest
		err := c.BindJSON(&request)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// validate request
		err = h.validateAskAIForEmailRequest(&request)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Invalid request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		systemPrompt := `I will provide you with the raw body of an email that may or may not contain an email signature.  Your job is to process this email into structured json data.  Start by parsing the message body.  Ensure it's onlly the current message and does not include old email threads or content.  Please ensure you remove all salutations and greetings.  Also remove all odd or unnatural line breaks.  Return the message body in valid markdown format.  

Next, determine if an email signature is present.  If it is, then you are to parse the signature and return it's data in exactly this format.  Please ensure you strip all query params from all URLs.  Please ensure all phone numbers are returned in international format with a valid country code.`

		var prompt strings.Builder

		prompt.WriteString(fmt.Sprintf("Email from: %s\n", request.EmailFrom))
		prompt.WriteString(fmt.Sprintf("From email address: %s\n", request.FromEmailAddress))
		prompt.WriteString(fmt.Sprintf("Email to: %s\n", request.EmailTo))
		prompt.WriteString(fmt.Sprintf("To email address: %s\n", request.FromEmailAddress))
		prompt.WriteString(fmt.Sprintf("Email body text: %s\n", request.EmailBodyText))
		prompt.WriteString(fmt.Sprintf("Email body html: %s\n", request.EmailBodyHTML))
		promptStr := prompt.String()

		temperature := float32(0.2)
		maxOutputTokens := int32(100)

		answer, err := h.services.CommonServices.AIService.AskAIForEmail(ctx, interfaces.AskAIRequest{
			Model:            enum.AIModelGemini,
			SystemPrompt:     &systemPrompt,
			Prompt:           &promptStr,
			ModelTemperature: &temperature,
			MaxOutputTokens:  &maxOutputTokens,
			OutputFormat:     enum.AIOutputJson,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to ask AI"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if answer == nil {
			message := "Empty response from AI"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, AskAIForEmailResponse{
			Answer: *answer,
		})
		return
	}
}

// ValidateAskAIForEmailRequest validates an AskAIForEmailRequest object
func (h *AskAIHandler) validateAskAIForEmailRequest(req *AskAIForEmailRequest) error {
	// Check for required fields
	if req.FromEmailAddress == "" {
		return fmt.Errorf("fromEmailAddress is required")
	}

	if req.ToEmailAddress == "" {
		return fmt.Errorf("toEmailAddress is required")
	}

	// Ensure we have at least one body format (text or HTML)
	if req.EmailBodyText == "" && req.EmailBodyHTML == "" {
		return fmt.Errorf("either emailBodyText or emailBodyHtml is required")
	}

	// Validate email addresses using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(req.FromEmailAddress) {
		return fmt.Errorf("invalid fromEmailAddress format: %s", req.FromEmailAddress)
	}

	if !emailRegex.MatchString(req.ToEmailAddress) {
		return fmt.Errorf("invalid toEmailAddress format: %s", req.ToEmailAddress)
	}

	// Optional: Check maximum lengths for fields
	if len(req.EmailBodyText) > 100000 { // 100KB limit
		return fmt.Errorf("emailBodyText exceeds maximum length")
	}

	if len(req.EmailBodyHTML) > 100000 { // 100KB limit
		return fmt.Errorf("emailBodyHtml exceeds maximum length")
	}

	return nil
}
