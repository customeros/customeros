package route

import (
	"context"
	"encoding/json"
	"github.com/biter777/countries"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.GET("/enrichPerson",
		tracing.TracingEnhancer(ctx, "GET /enrichPerson"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API, security.WithCache(services.CommonServices.Cache)),
		enrichPerson(services))
	r.GET("/findWorkEmail",
		tracing.TracingEnhancer(ctx, "GET /findWorkEmail"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API, security.WithCache(services.CommonServices.Cache)),
		findWorkEmail(services))
	r.GET("/enrichOrganizationWithScrapin",
		tracing.TracingEnhancer(ctx, "GET /enrichOrganizationWithScrapin"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API, security.WithCache(services.CommonServices.Cache)),
		enrichOrganizationWithScrapin(services))
	r.GET("/enrichOrganization",
		tracing.TracingEnhancer(ctx, "GET /enrichOrganization"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API, security.WithCache(services.CommonServices.Cache)),
		enrichOrganization(services))
}

func enrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "enrichPerson")
		defer span.Finish()

		var request model.EnrichPersonRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichPersonScrapinResponse{
				Status:      "error",
				Message:     "Invalid request body",
				PersonFound: false,
			})
			return
		}
		request.Normalize()

		tracing.LogObjectAsJson(span, "request", request)

		// validate mandatory parameters
		if request.LinkedinUrl == "" && request.Email == "" && request.Domain == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin, email and domain parameters"))
			services.Logger.Errorf("Missing linkedin, email and domain parameters")
			c.JSON(http.StatusBadRequest, model.EnrichPersonScrapinResponse{
				Status:      "error",
				Message:     "Missing linkedin, email and domain parameters",
				PersonFound: false,
			})
			return
		}

		var scrapinRecordId uint64
		var enrichPersonData *model.EnrichedPersonData

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			recordId, response, err := services.ScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInPersonProfile"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonScrapinResponse{
					Status:      "error",
					Message:     "Internal server error",
					PersonFound: false,
				})
				return
			}
			enrichPersonData = &model.EnrichedPersonData{
				PersonProfile: response,
			}
			scrapinRecordId = recordId
		}

		foundByLinkedInUrl := enrichPersonData != nil && enrichPersonData.PersonProfile != nil && enrichPersonData.PersonProfile.Person != nil

		// Step 2 - Search by email and domain
		if !foundByLinkedInUrl && (request.Email != "" || request.Domain != "") {
			recordId, response, err := services.ScrapeInService.ScrapInSearchPerson(ctx, request.Email, request.FirstName, request.LastName, request.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInSearchPerson"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonScrapinResponse{
					Status:      "error",
					Message:     "Internal server error",
					PersonFound: false,
				})
				return
			}
			enrichPersonData = &model.EnrichedPersonData{
				PersonProfile: response,
			}
			scrapinRecordId = recordId
		}

		c.JSON(http.StatusOK, model.EnrichPersonScrapinResponse{
			Status:      "success",
			RecordId:    scrapinRecordId,
			PersonFound: enrichPersonData.PersonProfile != nil && enrichPersonData.PersonProfile.Person != nil,
			Data:        enrichPersonData,
		})
	}
}

func findWorkEmail(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "findWorkEmail")
		defer span.Finish()

		var request model.FindWorkEmailRequest

		if err := c.BindJSON(&request); err != nil {
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.FindWorkEmailResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
			return
		}
		tracing.LogObjectAsJson(span, "request", request)

		recordId, requestId, response, err := services.BettercontactService.FindWorkEmail(ctx, request.LinkedinUrl, request.FirstName, request.LastName, request.CompanyName, request.CompanyDomain, request.EnrichPhoneNumber)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, model.FindWorkEmailResponse{
				Status:  "error",
				Message: "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, model.FindWorkEmailResponse{
			Status:                 "success",
			RecordId:               recordId,
			BetterContactRequestId: requestId,
			Data:                   response,
		})
	}
}

func enrichOrganizationWithScrapin(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "enrichOrganizationWithScrapin")
		defer span.Finish()

		var request model.EnrichOrganizationRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichOrganizationScrapinResponse{
				Status:            "error",
				Message:           "Invalid request body",
				OrganizationFound: false,
			})
			return
		}
		request.Normalize()

		tracing.LogObjectAsJson(span, "request", request)

		// validate mandatory parameters
		if request.LinkedinUrl == "" && request.Domain == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin and domain parameters"))
			services.Logger.Errorf("Missing linkedin and domain parameters")
			c.JSON(http.StatusBadRequest, model.EnrichOrganizationScrapinResponse{
				Status:            "error",
				Message:           "Missing linkedin and domain parameters",
				OrganizationFound: false,
			})
			return
		}

		var scrapinRecordId uint64
		var scrapinResponseBody *postgresEntity.ScrapInResponseBody

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			recordId, response, err := services.ScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInCompanyProfile"))
				c.JSON(http.StatusInternalServerError, model.EnrichOrganizationScrapinResponse{
					Status:            "error",
					Message:           "Internal server error",
					OrganizationFound: false,
				})
				return
			}
			scrapinResponseBody = response
			scrapinRecordId = recordId
		}

		foundByLinkedInUrl := scrapinResponseBody != nil && scrapinResponseBody.Company != nil

		// Step 2 - Scrapin by email
		if !foundByLinkedInUrl && request.Domain != "" {
			recordId, response, err := services.ScrapeInService.ScrapInSearchCompany(ctx, request.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInSearchCompany"))
				c.JSON(http.StatusInternalServerError, model.EnrichOrganizationScrapinResponse{
					Status:            "error",
					Message:           "Internal server error",
					OrganizationFound: false,
				})
				return
			}
			scrapinResponseBody = response
			scrapinRecordId = recordId
		}

		c.JSON(http.StatusOK, model.EnrichOrganizationScrapinResponse{
			Status:            "success",
			RecordId:          scrapinRecordId,
			OrganizationFound: scrapinResponseBody != nil && scrapinResponseBody.Company != nil,
			Data:              scrapinResponseBody,
		})
	}
}

func enrichOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "enrichOrganization")
		defer span.Finish()

		var request model.EnrichOrganizationRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichOrganizationResponse{
				Status:  "error",
				Message: "Invalid request body",
				Success: false,
			})
			return
		}
		request.Normalize()

		tracing.LogObjectAsJson(span, "request", request)

		// validate mandatory parameters
		if request.LinkedinUrl == "" && request.Domain == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin and domain parameters"))
			services.Logger.Errorf("Missing linkedin and domain parameters")
			c.JSON(http.StatusBadRequest, model.EnrichOrganizationResponse{
				Status:  "error",
				Message: "Missing linkedin and domain parameters",
				Success: false,
			})
			return
		}

		var scrapinResponseBody *postgresEntity.ScrapInResponseBody
		var brandfetchResponseBody *postgresEntity.BrandfetchResponseBody
		domain := request.Domain

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			_, response, err := services.ScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInCompanyProfile"))
				c.JSON(http.StatusInternalServerError, model.EnrichOrganizationResponse{
					Status:  "error",
					Message: "Internal server error",
					Success: false,
				})
				return
			}
			scrapinResponseBody = response
		}

		foundByLinkedInUrl := scrapinResponseBody != nil && scrapinResponseBody.Company != nil

		// Step 2 - Scrapin by email
		if !foundByLinkedInUrl && domain != "" {
			_, response, err := services.ScrapeInService.ScrapInSearchCompany(ctx, domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInSearchCompany"))
				c.JSON(http.StatusInternalServerError, model.EnrichOrganizationResponse{
					Status:  "error",
					Message: "Internal server error",
					Success: false,
				})
				return
			}
			scrapinResponseBody = response
		}

		// Step3 - Brandfetch
		if domain == "" {
			domain = utils.ExtractDomain(scrapinResponseBody.Company.WebsiteUrl)
		}
		if domain != "" {
			response, err := services.BrandfetchService.GetByDomain(ctx, domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Brandfetch by domain"))
				c.JSON(http.StatusInternalServerError, model.EnrichOrganizationResponse{
					Status:  "error",
					Message: "Internal server error",
					Success: false,
				})
				return
			}
			brandfetchResponseBody = response
		}

		if scrapinResponseBody == nil && brandfetchResponseBody == nil {
			span.LogKV("result", "No data found, both scrapin and brandfetch are nil")
			c.JSON(http.StatusNotFound, model.EnrichOrganizationResponse{
				Status:  "Not found",
				Message: "No data found",
				Success: false,
			})
			return
		} else if scrapinResponseBody.Success == false && brandfetchResponseBody.IsEmpty() {
			span.LogKV("result", "No data found")
			c.JSON(http.StatusNotFound, model.EnrichOrganizationResponse{
				Status:  "Not found",
				Message: "No data found",
				Success: false,
			})
			return
		}
		primaryEnrichSource := "SCRAPIN"
		if scrapinResponseBody == nil || scrapinResponseBody.Success == false {
			primaryEnrichSource = "BRANDFETCH"
		}

		// combine data
		combinedData := combineData(scrapinResponseBody, brandfetchResponseBody, domain)

		output := model.EnrichOrganizationResponse{
			Status:              "success",
			Data:                combinedData,
			Success:             true,
			PrimaryEnrichSource: primaryEnrichSource,
		}
		outputStr, err := json.Marshal(output)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal response"))
			c.JSON(http.StatusInternalServerError, model.EnrichOrganizationResponse{
				Status:  "error",
				Message: "Internal server error",
				Success: false,
			})
			return
		}
		tracing.LogObjectAsJson(span, "response.body", outputStr)

		c.JSON(http.StatusOK, output)
	}
}

func combineData(scrapin *postgresEntity.ScrapInResponseBody, brandfetch *postgresEntity.BrandfetchResponseBody, domain string) model.EnrichOrganizationResponseData {
	data := model.EnrichOrganizationResponseData{}

	updateResponseWithScrapinData(&data, scrapin, domain)
	updateResponseWithBrandfetchData(&data, brandfetch)
	normalizeCountry(&data.Location)

	return data
}

func normalizeCountry(data *model.EnrichOrganizationResponseDataLocation) {
	// handle special cases
	if strings.ToUpper(data.Country) == "OO" {
		data.Country = ""
	}
	countryFields := []string{data.CountryCodeA3, data.CountryCodeA2, data.Country}
	for _, code := range countryFields {
		if code != "" {
			country := countries.ByName(code)
			if country != countries.Unknown {
				data.Country = country.String()
				data.CountryCodeA2 = country.Alpha2()
				data.CountryCodeA3 = country.Alpha3()
				return
			}
		}
	}
}

func updateResponseWithScrapinData(d *model.EnrichOrganizationResponseData, scrapin *postgresEntity.ScrapInResponseBody, domain string) {
	if scrapin == nil {
		return
	}
	if scrapin.Company == nil {
		return
	}

	if d.Employees == 0 {
		d.Employees = scrapin.Company.GetEmployeeCount()
	}
	if d.FoundedYear == 0 {
		d.FoundedYear = int64(scrapin.Company.FoundedOn.Year)
	}
	if d.Name == "" {
		d.Name = scrapin.Company.Name
	}
	if d.ShortDescription == "" {
		if scrapin.Company.Tagline != nil {
			if tagline, ok := scrapin.Company.Tagline.(string); ok {
				d.ShortDescription = tagline
			}
		}
	}
	if d.LongDescription == "" {
		d.LongDescription = scrapin.Company.Description
	}
	if d.Domain == "" {
		if domain != "" {
			d.Domain = domain
		} else {
			d.Domain = utils.ExtractDomain(scrapin.Company.WebsiteUrl)
		}
	}
	if d.Website == "" {
		d.Website = scrapin.Company.WebsiteUrl
	}
	if scrapin.Company.Logo != "" {
		d.Logos = append(d.Logos, scrapin.Company.Logo)
	}
	if d.Industry == "" {
		d.Industry = scrapin.Company.Industry
	}
	if scrapin.Company.LinkedInUrl != "" {
		d.Socials = append(d.Socials, model.EnrichOrganizationResponseSocial{
			Url:   scrapin.Company.LinkedInUrl,
			Alias: scrapin.Company.UniversalName,
			Id:    scrapin.Company.LinkedInId,
		})
	}

	if !scrapin.Company.HeadquarterIsEmpty() {
		d.Location.IsHeadquarter = utils.BoolPtr(true)
		if d.Location.Country == "" {
			d.Location.Country = scrapin.Company.Headquarter.Country
		}
		if d.Location.Locality == "" {
			d.Location.Locality = scrapin.Company.Headquarter.City
		}
		if d.Location.Region == "" {
			d.Location.Region = scrapin.Company.Headquarter.GeographicArea
		}
		if d.Location.PostalCode == "" {
			d.Location.PostalCode = scrapin.Company.Headquarter.PostalCode
		}
		if d.Location.AddressLine1 == "" {
			d.Location.AddressLine1 = scrapin.Company.Headquarter.Street1
		}
		if d.Location.AddressLine2 == "" {
			if scrapin.Company.Headquarter.Street2 != nil {
				if street2, ok := scrapin.Company.Headquarter.Street2.(string); ok {
					d.Location.AddressLine2 = street2
				}
			}
		}
	}
}

func updateResponseWithBrandfetchData(d *model.EnrichOrganizationResponseData, brandfetch *postgresEntity.BrandfetchResponseBody) {
	if brandfetch == nil {
		return
	}

	if d.Employees == 0 {
		d.Employees = brandfetch.Company.GetEmployees()
	}
	if d.FoundedYear == 0 {
		d.FoundedYear = brandfetch.Company.FoundedYear
	}
	if d.Name == "" {
		d.Name = brandfetch.Name
	}
	if d.ShortDescription == "" {
		d.ShortDescription = brandfetch.Description
	}
	if d.LongDescription == "" {
		d.LongDescription = brandfetch.LongDescription
	}
	if d.Domain == "" {
		d.Domain = brandfetch.Domain
	}
	if d.Website == "" {
		d.Website = brandfetch.Domain
	}
	if brandfetch.Company.Kind != "" {
		if brandfetch.Company.Kind == "PUBLIC_COMPANY" {
			d.Public = utils.BoolPtr(true)
		} else {
			d.Public = utils.BoolPtr(false)
		}
	}
	if len(brandfetch.Logos) > 0 {
		for _, logo := range brandfetch.Logos {
			if logo.Type == "icon" {
				d.Icons = append(d.Icons, logo.Formats[0].Src)
			} else if logo.Type == "symbol" {
				d.Icons = append(d.Icons, logo.Formats[0].Src)
			} else if logo.Type == "logo" {
				d.Logos = append(d.Logos, logo.Formats[0].Src)
			} else if logo.Type == "other" {
				d.Logos = append(d.Logos, logo.Formats[0].Src)
			}
		}
	}
	if d.Industry == "" {
		industryMaxScore := float64(0)
		if len(brandfetch.Company.Industries) > 0 {
			for _, industry := range brandfetch.Company.Industries {
				if industry.Name != "" && industry.Score > industryMaxScore {
					d.Industry = industry.Name
					industryMaxScore = industry.Score
				}
			}
		}
	}

	for _, link := range brandfetch.Links {
		if link.Url != "" {
			d.Socials = append(d.Socials, model.EnrichOrganizationResponseSocial{
				Url: link.Url,
			})
		}
	}

	if !brandfetch.Company.LocationIsEmpty() {
		if d.Location.CountryCodeA2 == "" && d.Location.Country == "" {
			d.Location.CountryCodeA2 = brandfetch.Company.Location.CountryCodeA2
			d.Location.Country = brandfetch.Company.Location.Country
		} else if d.Location.Country == brandfetch.Company.Location.Country && d.Location.CountryCodeA2 == "" {
			d.Location.CountryCodeA2 = brandfetch.Company.Location.CountryCodeA2
		}
		if d.Location.Locality == "" {
			d.Location.Locality = brandfetch.Company.Location.City
		}
		if d.Location.Region == "" {
			d.Location.Region = brandfetch.Company.Location.Region
		}
	}
}
