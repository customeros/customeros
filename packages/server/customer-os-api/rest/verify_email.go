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
	"io"
	"net/http"
	"strings"
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
	AlternateEmail        string                  `json:"alternateEmail"`
}

type EmailVerificationRisk struct {
	IsFirewalled    bool `json:"isFirewalled"`
	IsRoleMailbox   bool `json:"isRoleMailbox"`
	IsFreeProvider  bool `json:"isFreeProvider"`
	IsMailboxFull   bool `json:"isMailboxFull"`
	IsPrimaryDomain bool `json:"isPrimaryDomain"`
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
		span.LogKV("request.address", emailAddress)
		span.LogKV("request.verifyCatchAll", c.Query("verifyCatchAll"))
		// check if verifyCatchAll param exists, defaulted to true
		verifyCatchAll := true
		if strings.ToLower(c.Query("verifyCatchAll")) == "false" {
			verifyCatchAll = false
		}

		syntaxValidation := mailsherpa.ValidateEmailSyntax(emailAddress)
		if !syntaxValidation.IsValid {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid email address syntax"})
			logger.Warnf("Invalid email address format: %s", emailAddress)
			return
		}

		// call validation api
		result, err := callApiValidateEmail(ctx, services, span, emailAddress, verifyCatchAll)
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
				result.Data.EmailData.IsMailboxFull ||
				!result.Data.DomainData.IsPrimaryDomain,
			Syntax: EmailVerificationSyntax{
				IsValid: syntaxValidation.IsValid,
				Domain:  syntaxValidation.Domain,
				User:    syntaxValidation.User,
			},
			Risk: EmailVerificationRisk{
				IsFirewalled:    result.Data.DomainData.IsFirewalled,
				IsRoleMailbox:   result.Data.EmailData.IsRoleAccount,
				IsFreeProvider:  result.Data.EmailData.IsFreeAccount,
				IsMailboxFull:   result.Data.EmailData.IsMailboxFull,
				IsPrimaryDomain: result.Data.DomainData.IsPrimaryDomain,
			},
			AlternateEmail: result.Data.EmailData.AlternateEmail,
		}

		c.JSON(http.StatusOK, emailVerificationResponse)
	}
}

func callApiValidateEmail(ctx context.Context, services *service.Services, span opentracing.Span, emailAddress string, verifyCatchAll bool) (*validationmodel.ValidateEmailResponse, error) {
	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequestWithOptions{
		Email: emailAddress,
		Options: validationmodel.ValidateEmailRequestOptions{
			VerifyCatchAll: verifyCatchAll,
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
	span.LogFields(log.Int("response.status", response.StatusCode))
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	var validationResponse validationmodel.ValidateEmailResponse
	err = json.Unmarshal(responseBody, &validationResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(responseBody)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return nil, err
	}
	if validationResponse.Data == nil {
		tracing.LogObjectAsJson(span, "response", validationResponse)
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &validationResponse, nil
}
