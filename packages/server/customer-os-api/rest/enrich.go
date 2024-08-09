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
)

func EnrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
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
