package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	IsProxy      bool `json:"isProxy"`
	IsVpn        bool `json:"isVpn"`
	IsTor        bool `json:"isTor"`
	IsBot        bool `json:"isBot"`
	IsDatacenter bool `json:"isDatacenter"`
	IsCloudRelay bool `json:"isCloudRelay"`
	IsMobile     bool `json:"isMobile"`
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

		// Check if ip_address is provided
		ipAddress := c.Query("ip_address")
		if ipAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing ip_address parameter"})
			return
		}
		span.LogFields(log.String("ip_address", ipAddress))

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

		// This is a placeholder response
		ipIntelligenceResponse := IpIntelligenceResponse{
			Status: "success",
			IP:     ipAddress,
			Threats: IpIntelligenceThreats{
				IsProxy: result.Data.Threat.IsProxy,
				//IsVpn:        TBD,
				IsTor: result.Data.Threat.IsTor,
				//IsBot:        TBD,
				IsDatacenter: result.Data.Threat.IsDatacenter,
				IsCloudRelay: result.Data.Threat.IsIcloudRelay,
				//IsMobile:     TBD,
			},
			Geolocation: IpIntelligenceGeolocation{
				City:            result.Data.City,
				Country:         result.Data.CountryName,
				CountryIso:      result.Data.CountryCode,
				IsEuropeanUnion: isEuropeanUnion(result.Data.CountryCode),
			},
			TimeZone: IpIntelligenceTimeZone{
				Name:        result.Data.TimeZone.Name,
				Abbr:        result.Data.TimeZone.Abbr,
				Offset:      result.Data.TimeZone.Offset,
				IsDst:       result.Data.TimeZone.IsDst,
				CurrentTime: utils.GetCurrentTimeInTimeZone(result.Data.TimeZone.Name),
			},
			Network: IpIntelligenceNetwork{
				ASN:    result.Data.Asn.Asn,
				Name:   result.Data.Asn.Name,
				Domain: result.Data.Asn.Domain,
				Route:  result.Data.Asn.Route,
				Type:   result.Data.Asn.Type,
			},
			Organization: IpIntelligenceOrganization{
				//Name:     TBD,
				//Domain:   TBD,
				//LinkedIn: TBD,
			},
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GenerateEmailTrackingUrls", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		// TODO implement me
		//log := services.Log

		c.JSON(http.StatusOK, "")
	}
}
