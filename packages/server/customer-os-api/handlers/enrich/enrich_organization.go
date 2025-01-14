// @openapi 3.0.0
package restenrich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
)

// EnrichOrganizationResponse represents the response for organization enrichment
// @Description Response structure for organization enrichment operations
type EnrichOrganizationResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// Enriched organization data
	// required: true
	Data EnrichOrganizationData `json:"data"`
}

// EnrichOrganizationData represents enriched organization information
// @Description Detailed enriched information about an organization
type EnrichOrganizationData struct {
	// Organization name
	// required: true
	// example: Acme Corporation
	Name string `json:"name"`

	// Organization's primary domain
	// required: true
	// example: acme.com
	Domain string `json:"domain"`

	// Brief description of the organization
	// required: false
	// example: A global leader in innovative solutions
	ShortDescription string `json:"description,omitempty"`

	// Detailed description of the organization
	// required: false
	// example: Acme Corporation provides cutting-edge technology solutions across the globe
	LongDescription string `json:"longDescription,omitempty"`

	// Organization's website URL
	// required: true
	// example: https://acme.com
	// format: uri
	Website string `json:"website"`

	// Number of employees
	// required: false
	// minimum: 0
	// example: 5000
	Employees int `json:"employees,omitempty"`

	// Year the organization was founded
	// required: false
	// minimum: 1800
	// maximum: 2100
	// example: 1995
	FoundedYear int `json:"foundedYear,omitempty"`

	// Indicates if the organization is publicly traded
	// required: false
	// example: true
	Public bool `json:"public,omitempty"`

	// URLs to organization logos
	// required: false
	// example: ["https://acme.com/logo.png"]
	Logos []string `json:"logos,omitempty"`

	// URLs to organization icons
	// required: false
	// example: ["https://acme.com/icon.png"]
	Icons []string `json:"icons,omitempty"`

	// Industry classification
	// required: false
	Industry EnrichOrganizationIndustry `json:"industry,omitempty"`

	// Social media presence
	// required: false
	// example: ["https://linkedin.com/company/acme"]
	Socials []string `json:"socials,omitempty"`

	// Organization location information
	// required: false
	Location EnrichOrganizationLocation `json:"location,omitempty"`
}

// EnrichOrganizationIndustry represents industry classification
// @Description Industry classification information
type EnrichOrganizationIndustry struct {
	// Primary industry category
	// required: true
	// example: Technology
	Industry string `json:"industry"`
}

// EnrichOrganizationLocation represents location information
// @Description Detailed location information for an organization
type EnrichOrganizationLocation struct {
	// Indicates if this is the headquarters location
	// required: true
	// example: true
	IsHeadquarter bool `json:"isHeadquarter"`

	// Country name
	// required: true
	// example: United States
	Country string `json:"country"`

	// ISO 3166-1 alpha-2 country code
	// required: true
	// example: US
	// pattern: ^[A-Z]{2}$
	CountryCodeA2 string `json:"countryCodeA2"`

	// City name
	// required: false
	// example: San Francisco
	City string `json:"city,omitempty"`

	// State or region
	// required: false
	// example: California
	Region string `json:"region,omitempty"`

	// Postal code
	// required: false
	// example: 94105
	PostalCode string `json:"postalCode,omitempty"`

	// Primary address line
	// required: false
	// example: 123 Main St
	AddressLine1 string `json:"addressLine1,omitempty"`

	// Secondary address line
	// required: false
	// example: Suite 100
	AddressLine2 string `json:"addressLine2,omitempty"`
}

// @Summary Enrich organization information
// @Description Enriches organization information using either domain or LinkedIn URL
// @Tags Enrichment API
// @Accept json
// @Produce json
// @Param linkedinUrl query string false "Organization's LinkedIn URL" example(https://linkedin.com/company/acme)
// @Param domain query string false "Organization's domain" example(acme.com)
// @Success 200 {object} EnrichOrganizationResponse "Successfully retrieved enriched data"
// @Success 200 {object} rest.ErrorResponse "Organization not found (status: warning)"
// @Failure 400 {object} rest.BaseResponse "Missing or invalid parameters"
// @Failure 401 {object} rest.BaseResponse "Missing or invalid API key"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /enrich/v1/organization [get]
// @Security ApiKeyAuth
func EnrichOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichOrganization", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		linkedinUrl := c.Query("linkedinUrl")
		domain := c.Query("domain")

		// check linked in or email params are present
		if strings.TrimSpace(linkedinUrl) == "" && strings.TrimSpace(domain) == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing linkedinUrl or domain"))
			return
		}

		span.LogFields(
			log.String("request.domain", domain),
			log.String("request.linkedinUrl", linkedinUrl))

		// Call enrichPerson API
		enrichOrganizationApiResponse, err := callApiEnrichOrganization(ctx, services, span, linkedinUrl, domain)
		if err != nil || enrichOrganizationApiResponse.Status == "error" {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if enrichOrganizationApiResponse.Success == false || enrichOrganizationApiResponse.Data == nil {
			handlers.SendError(c, span, http.StatusNotFound, &enum.ErrorResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusWarning),
				Message:      "Organization not found",
			})
			return
		}
		// Compose the response
		socialUrls := make([]string, 0)
		for _, social := range enrichOrganizationApiResponse.Data.Socials {
			socialUrls = append(socialUrls, social.Url)
		}
		response := EnrichOrganizationResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Data: EnrichOrganizationData{
				Name:             enrichOrganizationApiResponse.Data.Name,
				Domain:           enrichOrganizationApiResponse.Data.Domain,
				ShortDescription: enrichOrganizationApiResponse.Data.ShortDescription,
				LongDescription:  enrichOrganizationApiResponse.Data.LongDescription,
				Website:          enrichOrganizationApiResponse.Data.Website,
				Employees:        int(enrichOrganizationApiResponse.Data.Employees),
				FoundedYear:      int(enrichOrganizationApiResponse.Data.FoundedYear),
				Public:           utils.BoolDefaultIfNil(enrichOrganizationApiResponse.Data.Public, false),
				Logos:            enrichOrganizationApiResponse.Data.Logos,
				Icons:            enrichOrganizationApiResponse.Data.Icons,
				Industry: EnrichOrganizationIndustry{
					Industry: enrichOrganizationApiResponse.Data.Industry,
				},
				Socials: socialUrls,
				Location: EnrichOrganizationLocation{
					IsHeadquarter: utils.BoolDefaultIfNil(enrichOrganizationApiResponse.Data.Location.IsHeadquarter, false),
					Country:       enrichOrganizationApiResponse.Data.Location.Country,
					CountryCodeA2: enrichOrganizationApiResponse.Data.Location.CountryCodeA2,
					City:          enrichOrganizationApiResponse.Data.Location.Locality,
					Region:        enrichOrganizationApiResponse.Data.Location.Region,
					PostalCode:    enrichOrganizationApiResponse.Data.Location.PostalCode,
					AddressLine1:  enrichOrganizationApiResponse.Data.Location.AddressLine1,
					AddressLine2:  enrichOrganizationApiResponse.Data.Location.AddressLine2,
				},
			},
		}

		if enrichOrganizationApiResponse.Success == true {
			_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichOrganizationSuccess,
				postgresrepository.BillableEventDetails{
					ReferenceData: fmt.Sprintf("LinkedIn URL: %s, Domain: %s", linkedinUrl, domain),
				})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func callApiEnrichOrganization(ctx context.Context, services *service.Services, span opentracing.Span, linkedinUrl, domain string) (*enrichmentmodel.EnrichOrganizationResponse, error) {
	requestJSON, err := json.Marshal(enrichmentmodel.EnrichOrganizationRequest{
		Domain:      domain,
		LinkedinUrl: linkedinUrl,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", services.Cfg.InternalServices.EnrichmentApiConfig.Url+"/enrichOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.InternalServices.EnrichmentApiConfig.ApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.enrichOrganization", response.StatusCode))

	var enrichOrganizationApiResponse enrichmentmodel.EnrichOrganizationResponse
	err = json.NewDecoder(response.Body).Decode(&enrichOrganizationApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode enrich organization response"))
		return nil, err
	}
	return &enrichOrganizationApiResponse, nil
}
