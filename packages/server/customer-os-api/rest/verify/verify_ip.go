// @openapi 3.0.0
package restverify

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// IpIntelligenceResponse represents the IP intelligence response
// @Description Response containing IP intelligence data including threats, geolocation, and network information
type IpIntelligenceResponse struct {
	// Inherits standard response fields
	rest.BaseResponse
	// IP intelligence details
	// required: true
	IP IpIntelligenceRecord `json:"ip,omitempty"`
}

// IpIntelligenceRecord represents detailed IP information
// @Description Comprehensive information about an IP address
type IpIntelligenceRecord struct {
	// IP address being analyzed
	// required: true
	// pattern: ^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$
	IPAddress string `json:"ipAddress" example:"192.168.1.1"`

	// Threat intelligence information
	// required: true
	Threats IpIntelligenceThreats `json:"threats"`

	// Geolocation information
	// required: true
	Geolocation IpIntelligenceGeolocation `json:"geolocation"`

	// Timezone information
	// required: true
	TimeZone IpIntelligenceTimeZone `json:"time_zone"`

	// Network information
	// required: true
	Network IpIntelligenceNetwork `json:"network"`

	// Organization information
	// required: false
	Organization IpIntelligenceOrganization `json:"organization"`
}

// IpIntelligenceThreats represents threat intelligence data
// @Description Threat intelligence indicators for the IP address
type IpIntelligenceThreats struct {
	// Indicates if IP is a proxy
	// required: true
	IsProxy bool `json:"isProxy" example:"false"`

	// Indicates if IP is a VPN
	// required: true
	IsVpn bool `json:"isVpn" example:"false"`

	// Indicates if IP is a TOR exit node
	// required: true
	IsTor bool `json:"isTor" example:"false"`

	// Indicates if IP is unallocated
	// required: true
	IsUnallocated bool `json:"isUnallocated" example:"false"`

	// Indicates if IP belongs to a datacenter
	// required: true
	IsDatacenter bool `json:"isDatacenter" example:"false"`

	// Indicates if IP is a cloud relay
	// required: true
	IsCloudRelay bool `json:"isCloudRelay" example:"false"`

	// Indicates if IP belongs to a mobile network
	// required: true
	IsMobile bool `json:"isMobile" example:"false"`
}

// IpIntelligenceGeolocation represents geolocation data
// @Description Geographic location information for the IP address
type IpIntelligenceGeolocation struct {
	// City name
	// required: false
	City string `json:"city" example:"Berlin"`

	// Country name
	// required: true
	Country string `json:"country" example:"Germany"`

	// ISO 3166-1 alpha-2 country code
	// required: true
	// pattern: ^[A-Z]{2}$
	CountryIso string `json:"countryIso" example:"DE"`

	// Indicates if country is in the European Union
	// required: true
	IsEuropeanUnion bool `json:"isEuropeanUnion" example:"true"`
}

// IpIntelligenceTimeZone represents timezone data
// @Description Timezone information for the IP address location
type IpIntelligenceTimeZone struct {
	// IANA timezone name
	// required: true
	// example: Europe/Berlin
	Name string `json:"name"`

	// Timezone abbreviation
	// required: true
	// example: CET
	Abbr string `json:"abbr"`

	// UTC offset
	// required: true
	// pattern: ^[+-]\d{4}$
	Offset string `json:"offset" example:"+0100"`

	// Indicates if daylight saving time is active
	// required: true
	IsDst bool `json:"is_dst" example:"true"`

	// Current time in the timezone
	// required: true
	// format: date-time
	CurrentTime time.Time `json:"current_time" example:"2024-09-10T14:00:00+01:00"`
}

// IpIntelligenceNetwork represents network data
// @Description Network information for the IP address
type IpIntelligenceNetwork struct {
	// Autonomous System Number
	// required: true
	// pattern: ^AS\d+$
	ASN string `json:"asn" example:"AS12345"`

	// Network name
	// required: true
	Name string `json:"name" example:"ISP Name"`

	// Network domain
	// required: false
	Domain string `json:"domain" example:"isp.com"`

	// Network route (CIDR notation)
	// required: true
	// pattern: ^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$
	Route string `json:"route" example:"192.168.0.0/16"`

	// Network type
	// required: true
	// enum: business,hosting,isp,education,government
	Type string `json:"type" example:"business"`
}

// IpIntelligenceOrganization represents organization data
// @Description Organization information associated with the IP address
type IpIntelligenceOrganization struct {
	// Organization name
	// required: false
	Name string `json:"name" example:"Company Name"`

	// Organization domain
	// required: false
	Domain string `json:"domain" example:"company.com"`

	// LinkedIn profile URL
	// required: false
	// format: uri
	LinkedIn string `json:"linkedin" example:"https://linkedin.com/company/company"`
}

// @Summary Get IP intelligence data
// @Description Retrieves comprehensive information about an IP address including threats, geolocation, and network details
// @Tags IP Intelligence
// @Accept json
// @Produce json
// @Param address query string true "IP address to analyze" pattern(^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$)
// @Success 200 {object} IpIntelligenceResponse "IP intelligence data retrieved successfully"
// @Failure 400 {object} rest.ErrorResponse "Invalid IP address format"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /verify/v1/ip [get]
// @Security ApiKeyAuth
func IpIntelligence(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "IpIntelligence", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}
		logger := services.Log

		// Check if address is provided
		ipAddress := c.Query("address")
		if ipAddress == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing parameter: address"))
			return
		}
		span.LogFields(log.String("address", ipAddress))

		if net.ParseIP(ipAddress) == nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("IP address is not valid"))
			logger.Warnf("Invalid IP address format: %s", ipAddress)
			return
		}

		requestJSON, err := json.Marshal(validationmodel.IpLookupRequest{
			Ip: ipAddress,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}
		requestBody := []byte(string(requestJSON))
		req, err := http.NewRequest("POST", services.Cfg.InternalServices.ValidationApi+"/ipLookup", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}
		// Inject span context into the HTTP request
		req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

		// Set the request headers
		req.Header.Set(security.ApiKeyHeader, services.Cfg.InternalServices.ValidationApiKey)
		req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

		// Make the HTTP request
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
		}
		defer response.Body.Close()

		var result validationmodel.IpLookupResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		var ipIntelligenceResponse IpIntelligenceRecord
		if result.IpData.StatusCode == 400 {
			ipIntelligenceResponse = IpIntelligenceRecord{
				IPAddress: ipAddress,
				Threats: IpIntelligenceThreats{
					IsUnallocated: true,
				},
			}
		} else {
			ipIntelligenceResponse = IpIntelligenceRecord{
				IPAddress: ipAddress,
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
					// Name:     TBD,
					// Domain:   TBD,
					// LinkedIn: TBD,
				},
			}
		}

		_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventIpVerificationSuccess,
			postgresrepository.BillableEventDetails{
				ReferenceData: ipAddress,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
		}

		c.JSON(http.StatusOK, IpIntelligenceResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			IP:           ipIntelligenceResponse,
		})
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
