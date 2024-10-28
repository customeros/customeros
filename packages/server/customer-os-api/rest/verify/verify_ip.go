package restverify

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

// IpIntelligenceResponse represents the response for IP intelligence lookup.
// @Description Response structure for IP intelligence lookup.
// @example 200 {object} IpIntelligenceResponse
type IpIntelligenceResponse struct {
	Status       string                     `json:"status" example:"success"`
	Message      string                     `json:"message,omitempty" example:"No threats detected"`
	IP           string                     `json:"ip" example:"192.168.1.1"`
	Threats      IpIntelligenceThreats      `json:"threats"`
	Geolocation  IpIntelligenceGeolocation  `json:"geolocation"`
	TimeZone     IpIntelligenceTimeZone     `json:"time_zone"`
	Network      IpIntelligenceNetwork      `json:"network"`
	Organization IpIntelligenceOrganization `json:"organization"`
}

// IpIntelligenceThreats contains threat intelligence data related to the IP address.
// @Description Threat intelligence data related to the IP address.
type IpIntelligenceThreats struct {
	IsProxy       bool `json:"isProxy" example:"false"`
	IsVpn         bool `json:"isVpn" example:"false"`
	IsTor         bool `json:"isTor" example:"false"`
	IsUnallocated bool `json:"isUnallocated" example:"true"`
	IsDatacenter  bool `json:"isDatacenter" example:"false"`
	IsCloudRelay  bool `json:"isCloudRelay" example:"false"`
	IsMobile      bool `json:"isMobile" example:"true"`
}

// IpIntelligenceGeolocation contains geolocation data related to the IP address.
// @Description Geolocation data related to the IP address.
type IpIntelligenceGeolocation struct {
	City            string `json:"city" example:"Berlin"`
	Country         string `json:"country" example:"Germany"`
	CountryIso      string `json:"countryIso" example:"DE"`
	IsEuropeanUnion bool   `json:"isEuropeanUnion" example:"true"`
}

// IpIntelligenceTimeZone contains timezone data for the IP address.
// @Description Timezone data for the IP address.
type IpIntelligenceTimeZone struct {
	Name        string    `json:"name" example:"Europe/Berlin"`
	Abbr        string    `json:"abbr" example:"CET"`
	Offset      string    `json:"offset" example:"+0100"`
	IsDst       bool      `json:"is_dst" example:"true"`
	CurrentTime time.Time `json:"current_time" example:"2024-09-10T14:00:00+01:00"`
}

// IpIntelligenceNetwork contains network-related data for the IP address.
// @Description Network-related data for the IP address.
type IpIntelligenceNetwork struct {
	ASN    string `json:"asn" example:"AS12345"`
	Name   string `json:"name" example:"ISP Name"`
	Domain string `json:"domain" example:"isp.com"`
	Route  string `json:"route" example:"192.168.0.0/16"`
	Type   string `json:"type" example:"business"`
}

// IpIntelligenceOrganization contains organizational data for the IP address.
// @Description Organizational data for the IP address.
type IpIntelligenceOrganization struct {
	Name     string `json:"name" example:"Company Name"`
	Domain   string `json:"domain" example:"company.com"`
	LinkedIn string `json:"linkedin" example:"https://linkedin.com/company/company"`
}

// @Tags Verify API
// @Summary Get IP Intelligence
// @Description Retrieves threat intelligence and geolocation data for the given IP address.
// @Param address query string true "IP address to check"
// @Success 200 {object} IpIntelligenceResponse
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /verify/v1/ip [get]
func IpIntelligence(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "IpIntelligence", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}
		logger := services.Log

		// Check if address is provided
		ipAddress := c.Query("address")
		if ipAddress == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing address parameter",
				})
			return
		}
		span.LogFields(log.String("address", ipAddress))

		if net.ParseIP(ipAddress) == nil {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Invalid IP address format",
				})
			logger.Warnf("Invalid IP address format: %s", ipAddress)
			return
		}

		requestJSON, err := json.Marshal(validationmodel.IpLookupRequest{
			Ip: ipAddress,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		requestBody := []byte(string(requestJSON))
		req, err := http.NewRequest("POST", services.Cfg.InternalServices.ValidationApi+"/ipLookup", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
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
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
		}
		defer response.Body.Close()

		var result validationmodel.IpLookupResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
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

		_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventIpVerificationSuccess,
			postgresrepository.BillableEventDetails{
				ReferenceData: ipAddress,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
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
