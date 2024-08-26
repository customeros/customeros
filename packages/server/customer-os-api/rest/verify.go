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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"time"
)

type IpIntelligenceResponse struct {
	Status       string                     `json:"status"`
	Message      string                     `json:"message,omitempty"`
	IP           string                     `json:"ip"`
	Threats      IpIntelligenceThreats      `json:"threats"`
	Geolocation  IpIntelligenceGeolocation  `json:"geolocation"`
	TimeZone     IpIntelligenceTimeZone     `json:"time_zone"`
	Network      IpIntelligenceNetwork      `json:"network"`
	Organization IpIntelligenceOrganization `json:"organization"`
}

type IpIntelligenceThreats struct {
	IsProxy       bool `json:"isProxy"`
	IsVpn         bool `json:"isVpn"`
	IsTor         bool `json:"isTor"`
	IsUnallocated bool `json:"isUnallocated"`
	IsDatacenter  bool `json:"isDatacenter"`
	IsCloudRelay  bool `json:"isCloudRelay"`
	IsMobile      bool `json:"isMobile"`
}

type IpIntelligenceGeolocation struct {
	City            string `json:"city"`
	Country         string `json:"country"`
	CountryIso      string `json:"countryIso"`
	IsEuropeanUnion bool   `json:"isEuropeanUnion"`
}

type IpIntelligenceTimeZone struct {
	Name        string    `json:"name"`
	Abbr        string    `json:"abbr"`
	Offset      string    `json:"offset"`
	IsDst       bool      `json:"is_dst"`
	CurrentTime time.Time `json:"current_time"`
}

type IpIntelligenceNetwork struct {
	ASN    string `json:"asn"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Route  string `json:"route"`
	Type   string `json:"type"`
}

type IpIntelligenceOrganization struct {
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	LinkedIn string `json:"linkedin"`
}

type EmailVerificationResponse struct {
	Status      string                  `json:"status"`
	Message     string                  `json:"message,omitempty"`
	Email       string                  `json:"email"`
	Deliverable string                  `json:"deliverable"`
	Provider    string                  `json:"provider"`
	IsRisky     bool                    `json:"isRisky"`
	IsCatchAll  bool                    `json:"isCatchAll"`
	Risk        EmailVerificationRisk   `json:"risk"`
	Syntax      EmailVerificationSyntax `json:"syntax"`
	Smtp        EmailVerificationSmtp   `json:"smtp"`
	// Deprecated
	IsDeliverable bool `json:"isDeliverable"`
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

type EmailVerificationSmtp struct {
	ResponseCode string `json:"responseCode"`
	ErrorCode    string `json:"errorCode"`
	Description  string `json:"description"`
}

func IpIntelligence(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "IpIntelligence", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
		logger := services.Log

		// Check if address is provided
		ipAddress := c.Query("address")
		if ipAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing address parameter"})
			return
		}
		span.LogFields(log.String("address", ipAddress))

		if net.ParseIP(ipAddress) == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid IP address format"})
			logger.Warnf("Invalid IP address format: %s", ipAddress)
			return
		}

		requestJSON, err := json.Marshal(validationmodel.IpLookupRequest{
			Ip: ipAddress,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		requestBody := []byte(string(requestJSON))
		req, err := http.NewRequest("POST", services.Cfg.Services.ValidationApi+"/ipLookup", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
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
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
		}
		defer response.Body.Close()

		var result validationmodel.IpLookupResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}

		var ipIntelligenceResponse IpIntelligenceResponse
		if result.IpData.StatusCode == 400 {
			ipIntelligenceResponse = IpIntelligenceResponse{
				Status: "success",
				IP:     ipAddress,
				Threats: IpIntelligenceThreats{
					IsUnallocated: true,
				},
			}
		} else {
			ipIntelligenceResponse = IpIntelligenceResponse{
				Status: "success",
				IP:     ipAddress,
				Threats: IpIntelligenceThreats{
					IsProxy:       result.IpData.Threat.IsProxy,
					IsVpn:         result.IpData.Threat.IsVpn,
					IsTor:         result.IpData.Threat.IsTor,
					IsUnallocated: result.IpData.Threat.IsBogon,
					IsDatacenter:  result.IpData.Threat.IsDatacenter,
					IsCloudRelay:  result.IpData.Threat.IsIcloudRelay,
					IsMobile:      result.IpData.Carrier != nil,
				},
				Geolocation: IpIntelligenceGeolocation{
					City:            result.IpData.City,
					Country:         result.IpData.CountryName,
					CountryIso:      result.IpData.CountryCode,
					IsEuropeanUnion: isEuropeanUnion(result.IpData.CountryCode),
				},
				TimeZone: IpIntelligenceTimeZone{
					Name:        result.IpData.TimeZone.Name,
					Abbr:        result.IpData.TimeZone.Abbr,
					Offset:      result.IpData.TimeZone.Offset,
					IsDst:       result.IpData.TimeZone.IsDst,
					CurrentTime: utils.GetCurrentTimeInTimeZone(result.IpData.TimeZone.Name),
				},
				Network: IpIntelligenceNetwork{
					ASN:    result.IpData.Asn.Asn,
					Name:   result.IpData.Asn.Name,
					Domain: result.IpData.Asn.Domain,
					Route:  result.IpData.Asn.Route,
					Type:   result.IpData.Asn.Type,
				},
				Organization: IpIntelligenceOrganization{
					// TBD: Snitcher
					//Name:     TBD,
					//Domain:   TBD,
					//LinkedIn: TBD,
				},
			}
		}

		c.JSON(http.StatusOK, ipIntelligenceResponse)
	}
}

func isEuropeanUnion(countryCodeA2 string) bool {
	switch countryCodeA2 {
	case "AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK":
		return true
	default:
		return false
	}
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
			Status:     "success",
			Email:      emailAddress,
			Provider:   result.Data.DomainData.Provider,
			IsCatchAll: result.Data.DomainData.IsCatchAll,
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
			Smtp: EmailVerificationSmtp{
				Description:  result.Data.EmailData.Description,
				ErrorCode:    result.Data.EmailData.ErrorCode,
				ResponseCode: result.Data.EmailData.ResponseCode,
			},
			Deliverable: result.Data.EmailData.Deliverable,
		}

		c.JSON(http.StatusOK, emailVerificationResponse)
	}
}

func callApiValidateEmail(ctx context.Context, services *service.Services, span opentracing.Span, emailAddress string) (*validationmodel.ValidateEmailResponse, error) {
	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequest{
		Email: emailAddress,
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
		return nil, err
	}
	if validationResponse.Data == nil {
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &validationResponse, nil
}
