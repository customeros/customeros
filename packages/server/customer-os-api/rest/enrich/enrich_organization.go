package restenrich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// EnrichOrganizationResponse represents the response for the organization enrichment API.
// @Description Response structure for the organization enrichment API.
// @example 200 {object} EnrichOrganizationResponse
type EnrichOrganizationResponse struct {
	// Status of the response.
	// Example: success
	Status string `json:"status" example:"success"`
	// Message for the response.
	// Example: Enrichment completed
	Message string                 `json:"message,omitempty" example:"Enrichment completed"`
	Data    EnrichOrganizationData `json:"data"`
}

// EnrichOrganizationData represents the detailed data about the organization from enrichment.
// @Description Detailed data about an organization from enrichment.
type EnrichOrganizationData struct {
	// Name of the organization.
	// Example: Acme Corporation
	Name string `json:"name" example:"Acme Corporation"`

	// Domain of the organization's website.
	// Example: acme.com
	Domain string `json:"domain" example:"acme.com"`

	// Short description of the organization.
	// Example: A global leader in innovative solutions.
	ShortDescription string `json:"description" example:"A global leader in innovative solutions"`

	// Long description of the organization.
	// Example: Acme Corporation provides cutting-edge technology solutions across the globe.
	LongDescription string `json:"longDescription" example:"Acme Corporation provides cutting-edge technology solutions across the globe."`

	// Website URL of the organization.
	// Example: https://acme.com
	Website string `json:"website" example:"https://acme.com"`

	// Number of employees in the organization.
	// Example: 5000
	Employees int `json:"employees" example:"5000"`

	// Year the organization was founded.
	// Example: 1995
	FoundedYear int `json:"foundedYear" example:"1995"`

	// Indicates whether the organization is publicly traded.
	// Example: true
	Public bool `json:"public,omitempty" example:"true"`

	// List of logo URLs for the organization.
	// Example: ["https://acme.com/logo.png"]
	Logos []string `json:"logos" example:"https://acme.com/logo.png"`

	// List of icon URLs for the organization.
	// Example: ["https://acme.com/icon.png"]
	Icons []string `json:"icons" example:"https://acme.com/icon.png"`

	// Industry in which the organization operates.
	Industry EnrichOrganizationIndustry `json:"industry"`

	// List of social media URLs for the organization.
	// Example: ["https://linkedin.com/company/acme"]
	Socials []string `json:"socials" example:"https://linkedin.com/company/acme"`

	// Location information about the organization.
	Location EnrichOrganizationLocation `json:"location"`
}

type EnrichOrganizationIndustry struct {
	// Industry in which the organization operates.
	// Example: Technology
	Industry string `json:"industry" example:"Technology"`
}

// EnrichOrganizationLocation represents the location details of an organization.
// @Description Location details of an organization.
type EnrichOrganizationLocation struct {
	// Indicates if the location is the headquarters.
	// Example: true
	IsHeadquarter bool `json:"isHeadquarter" example:"true"`

	// Country of the organization.
	// Example: United States
	Country string `json:"country" example:"United States"`

	// ISO Alpha-2 code of the country.
	// Example: US
	CountryCodeA2 string `json:"countryCodeA2" example:"US"`

	// City or locality of the organization.
	// Example: San Francisco
	City string `json:"city" example:"San Francisco"`

	// Region or state of the organization.
	// Example: California
	Region string `json:"region" example:"California"`

	// Postal code of the organization's location.
	// Example: 94105
	PostalCode string `json:"postalCode" example:"94105"`

	// Address line 1 of the organization.
	// Example: 123 Main St
	AddressLine1 string `json:"addressLine1" example:"123 Main St"`

	// Address line 2 of the organization (optional).
	// Example: Suite 100
	AddressLine2 string `json:"addressLine2" example:"Suite 100"`
}

// @Summary Enrich Organization Information
// @Description Enriches an organization's information using the domain or other details.
// @Tags Enrichment API
// @Param linkedinUrl query string false "LinkedIn profile URL of the organization"
// @Param domain query string false "Domain of the organization"
// @Success 200 {object} EnrichOrganizationResponse "Enrichment results including organizational data"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /enrich/v1/organization [get]
func EnrichOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichOrganization", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)
		commontracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}

		linkedinUrl := c.Query("linkedinUrl")
		domain := c.Query("domain")

		// check linked in or email params are present
		if strings.TrimSpace(linkedinUrl) == "" && strings.TrimSpace(domain) == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required parameters linkedinUrl or domain",
				})
			return
		}
		span.LogFields(
			log.String("request.domain", domain),
			log.String("request.linkedinUrl", linkedinUrl))

		// Call enrichPerson API
		enrichOrganizationApiResponse, err := callApiEnrichOrganization(ctx, services, span, linkedinUrl, domain)
		if err != nil || enrichOrganizationApiResponse.Status == "error" {
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Internal error",
				})
			return
		}
		if enrichOrganizationApiResponse.Success == false {
			c.JSON(http.StatusOK,
				rest.ErrorResponse{
					Status:  "warning",
					Message: "Organization not found",
				})
			return
		}
		// Compose the response
		response := EnrichOrganizationResponse{
			Status:  "success",
			Message: "Enrichment completed",
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
				Socials: enrichOrganizationApiResponse.Data.Socials,
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
			_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, postgresentity.BillableEventEnrichOrganizationSuccess, "", "",
				fmt.Sprintf("LinkedIn URL: %s, Domain: %s", linkedinUrl, domain))
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
	req, err := http.NewRequest("GET", services.Cfg.InternalServices.EnrichmentApiUrl+"/enrichOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.InternalServices.EnrichmentApiKey)
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
