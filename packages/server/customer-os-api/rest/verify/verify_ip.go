// @openapi 3.0.0
package verify

import (
	"net"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// IpIntelligenceResponse represents the IP intelligence response
// @Description Response containing IP intelligence data including threats, geolocation, and network information
type IpIntelligenceResponse struct {
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
// @Failure 400 {object} handlers.ErrorResponse "Invalid IP address format"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /verify/v1/ip [get]
// @Security ApiKeyAuth
func (h *VerifyHandler) IpIntelligence() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "IpIntelligence", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
		logger := h.services.Log

		// Check if address is provided
		ipAddress := c.Query("address")
		if ipAddress == "" {
			message := "Missing parameter: address"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogFields(log.String("address", ipAddress))

		if net.ParseIP(ipAddress) == nil {
			message := "IP address is not valid"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			logger.Warnf("Invalid IP address format: %s", ipAddress)
			return
		}

		var ipIntelligenceResponse IpIntelligenceRecord
		result, err := h.services.CommonServices.VerifyService.LookupIp(ctx, ipAddress)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}
		if result == nil {
			ipIntelligenceResponse = IpIntelligenceRecord{
				IPAddress: ipAddress,
				Threats: IpIntelligenceThreats{
					IsUnallocated: true,
				},
			}
		}

		ipIntelligenceResponse = IpIntelligenceRecord{
			IPAddress: ipAddress,
			Threats: IpIntelligenceThreats{
				IsProxy:       result.Threat.IsProxy,
				IsVpn:         result.Threat.IsVpn,
				IsTor:         result.Threat.IsTor,
				IsUnallocated: result.Threat.IsBogon,
				IsDatacenter:  result.Threat.IsDatacenter,
				IsCloudRelay:  result.Threat.IsIcloudRelay,
				IsMobile:      result.Carrier != nil,
			},
			Geolocation: IpIntelligenceGeolocation{
				City:            result.City,
				Country:         result.CountryName,
				CountryIso:      result.CountryCode,
				IsEuropeanUnion: utils.IsEuropeanUnion(result.CountryCode),
			},
			TimeZone: IpIntelligenceTimeZone{
				Name:        result.TimeZone.Name,
				Abbr:        result.TimeZone.Abbr,
				Offset:      result.TimeZone.Offset,
				IsDst:       result.TimeZone.IsDst,
				CurrentTime: utils.GetCurrentTimeInTimeZone(result.TimeZone.Name),
			},
			Network: IpIntelligenceNetwork{
				ASN:    result.Asn.Asn,
				Name:   result.Asn.Name,
				Domain: result.Asn.Domain,
				Route:  result.Asn.Route,
				Type:   result.Asn.Type,
			},
			Organization: IpIntelligenceOrganization{
				// TBD: Snitcher
				// Name:     TBD,
				// Domain:   TBD,
				// LinkedIn: TBD,
			},
		}

		_, err = h.services.Repositories.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventIpVerificationSuccess,
			postgresrepository.BillableEventDetails{
				ReferenceData: ipAddress,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
		}

		h.responseHandler.HandleSuccess(c, IpIntelligenceResponse{
			IP: ipIntelligenceResponse,
		})
	}
}
