package rest

import (
	"bytes"
	"encoding/json"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

type EmailVerificationResponse struct {
	Status                string                  `json:"status"`
	Message               string                  `json:"message,omitempty"`
	Email                 string                  `json:"email"`
	Deliverable           string                  `json:"deliverable"`
	Provider              string                  `json:"provider"`
	SecureGatewayProvider string                  `json:"secureGatewayProvider"`
	IsRisky               bool                    `json:"isRisky"`
	IsCatchAll            bool                    `json:"isCatchAll"`
	Risk                  EmailVerificationRisk   `json:"risk"`
	Syntax                EmailVerificationSyntax `json:"syntax"`
}

type EmailVerificationRisk struct {
	IsFirewalled   bool `json:"isFirewalled"`
	IsRoleMailbox  bool `json:"isRoleMailbox"`
	IsFreeProvider bool `json:"isFreeProvider"`
	IsMailboxFull  bool `json:"isMailboxFull"`
}

type EmailVerificationSyntax struct {
	IsValid bool   `json:"isValid"`
	Domain  string `json:"domain"`
	User    string `json:"user"`
}

func VerifyEmailAddress(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "VerifyEmailAddress", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
		logger := services.Log

		// Check if email address is provided
		emailAddress := c.Query("address")
		if emailAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing address parameter"})
			return
		}
		span.LogFields(log.String("address", emailAddress))

		syntaxValidation := mailsherpa.ValidateEmailSyntax(emailAddress)
		if !syntaxValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid email address syntax"})
			logger.Warnf("Invalid email address format: %s", emailAddress)
			return
		}

		// call validation api
		result, err := callApiValidateEmail(ctx, services, span, emailAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Message})
			return
		}

		emailVerificationResponse := EmailVerificationResponse{
			Status:                "success",
			Email:                 emailAddress,
			Deliverable:           result.Data.EmailData.Deliverable,
			Provider:              result.Data.DomainData.Provider,
			SecureGatewayProvider: result.Data.DomainData.SecureGatewayProvider,
			IsCatchAll:            result.Data.DomainData.IsCatchAll,
			IsRisky: result.Data.DomainData.IsFirewalled ||
				result.Data.EmailData.IsRoleAccount ||
				result.Data.EmailData.IsFreeAccount ||
				result.Data.EmailData.IsMailboxFull,
			Syntax: EmailVerificationSyntax{
				IsValid: syntaxValidation.IsValid,
				Domain:  syntaxValidation.Domain,
				User:    syntaxValidation.User,
			},
			Risk: EmailVerificationRisk{
				IsFirewalled:   result.Data.DomainData.IsFirewalled,
				IsRoleMailbox:  result.Data.EmailData.IsRoleAccount,
				IsFreeProvider: result.Data.EmailData.IsFreeAccount,
				IsMailboxFull:  result.Data.EmailData.IsMailboxFull,
			},
		}

		c.JSON(http.StatusOK, emailVerificationResponse)
	}
}

func callApiValidateEmail(ctx context.Context, services *service.Services, span opentracing.Span, emailAddress string) (*validationmodel.ValidateEmailResponse, error) {
	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequestWithOptions{
		Email: emailAddress,
		Options: validationmodel.ValidateEmailRequestOptions{
			CallTrueInbox: true,
		},
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", services.Cfg.Services.ValidationApi+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.ValidationApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()

	var validationResponse validationmodel.ValidateEmailResponse
	err = json.NewDecoder(response.Body).Decode(&validationResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		responseBody := new(bytes.Buffer)
		_, readErr := responseBody.ReadFrom(response.Body)
		if readErr != nil {
			tracing.TraceErr(span, errors.Wrap(readErr, "failed to read response body"))
		}
		span.LogFields(log.String("response.body", responseBody.String()))
		return nil, err
	}
	if validationResponse.Data == nil {
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &validationResponse, nil
}