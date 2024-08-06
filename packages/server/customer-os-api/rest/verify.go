package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/opentracing/opentracing-go/log"
	"net"
	"net/http"
	"time"
)

type IpIntelligenceResponse struct {
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

		// TODO: Implement the actual logic to gather IP intelligence

		// This is a placeholder response
		response := IpIntelligenceResponse{
			IP: ipAddress,
			Threats: IpIntelligenceThreats{
				IsProxy:      true,
				IsVpn:        false,
				IsTor:        false,
				IsBot:        true,
				IsDatacenter: false,
				IsCloudRelay: false,
				IsMobile:     true,
			},
			Geolocation: IpIntelligenceGeolocation{
				City:            "Berlin",
				Country:         "Germany",
				CountryIso:      "DEU",
				IsEuropeanUnion: true,
			},
			TimeZone: IpIntelligenceTimeZone{
				Name:        "Europe/London",
				Abbr:        "BST",
				Offset:      "+0100",
				IsDst:       true,
				CurrentTime: time.Now(),
			},
			Network: IpIntelligenceNetwork{
				ASN:    "AS5089",
				Name:   "Virgin Media Limited",
				Domain: "virginmedia.com",
				Route:  "92.238.0.0/15",
				Type:   "business",
			},
			Organization: IpIntelligenceOrganization{
				Name:     "Knowsley",
				Domain:   "knowsley.gov.uk",
				LinkedIn: "https://linkedin.com/in/knowsley",
			},
		}

		c.JSON(http.StatusOK, response)
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

//func tp_be_deleted() {
//	findEmailRequest := FindEmailRequest{
//		FirstName: firstName,
//		LastName:  lastName,
//		Domain:    (*domainEntitiesPtr)[0].Domain,
//	}
//
//	efJSON, err := json.Marshal(findEmailRequest)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		graphql.AddErrorf(ctx, "Internal error")
//		return "", nil
//	}
//	requestBody := []byte(string(efJSON))
//	req, err := http.NewRequest("POST", r.cfg.Services.ValidationApi+"/findEmail", bytes.NewBuffer(requestBody))
//	if err != nil {
//		tracing.TraceErr(span, err)
//		graphql.AddErrorf(ctx, "Internal error")
//		return "", nil
//	}
//	// Inject span context into the HTTP request
//	req = commonTracing.InjectSpanContextIntoHTTPRequest(req, span)
//
//	// Set the request headers
//	req.Header.Set(security.ApiKeyHeader, r.cfg.Services.ValidationApiKey)
//	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))
//
//	// Make the HTTP request
//	client := &http.Client{}
//	response, err := client.Do(req)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		graphql.AddErrorf(ctx, "Internal error")
//		return "", nil
//	}
//	defer response.Body.Close()
//	var result FindEmailResponse
//	err = json.NewDecoder(response.Body).Decode(&result)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		graphql.AddErrorf(ctx, "Internal error")
//		return "", nil
//	}
//}
