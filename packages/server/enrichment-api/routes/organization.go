package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/biter777/countries"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
)

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

		var scrapinResponseBody *entity.ScrapInResponseBody
		var brandfetchResponseBody *entity.BrandfetchResponseBody
		domain := request.Domain

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			_, response, err := services.ScrapeInService.ScrapInCompanyProfile(ctx, request.LinkedinUrl)
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

		// Step 2 - Scrapin by domain
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
		if domain == "" && scrapinResponseBody != nil && scrapinResponseBody.Company != nil {
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
		} else if (scrapinResponseBody != nil && scrapinResponseBody.Success == false) && (brandfetchResponseBody != nil && brandfetchResponseBody.IsEmpty()) {
			span.LogKV("result", "No data found")
			c.JSON(http.StatusNotFound, model.EnrichOrganizationResponse{
				Status:  "Not found",
				Message: "No data found",
				Success: false,
			})
			return
		}
		primaryEnrichSource := SCRAPIN
		if scrapinResponseBody == nil || scrapinResponseBody.Success == false {
			primaryEnrichSource = BRANDFETCH
		}

		// combine data
		combinedData := combineData(scrapinResponseBody, brandfetchResponseBody, domain)

		output := model.EnrichOrganizationResponse{
			Status:              "success",
			Data:                &combinedData,
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

func scrapinOrganization(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "scrapinOrganization")
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

		var scrapinResponseBody *entity.ScrapInResponseBody
		domain := request.Domain

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			_, response, err := services.ScrapeInService.ScrapInCompanyProfile(ctx, request.LinkedinUrl)
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

		// Step 2 - Scrapin by domain
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

		if scrapinResponseBody == nil || scrapinResponseBody.Success == false {
			output := model.EnrichOrganizationResponse{
				Status:              "error",
				Success:             false,
				PrimaryEnrichSource: SCRAPIN,
			}
			c.JSON(http.StatusNotFound, output)
			return
		}

		responseData := model.EnrichOrganizationResponseData{}
		updateResponseWithScrapinData(&responseData, scrapinResponseBody, domain)
		normalizeCountry(&responseData.Location)
		output := model.EnrichOrganizationResponse{
			Status:              "success",
			Data:                &responseData,
			Success:             true,
			PrimaryEnrichSource: SCRAPIN,
		}
		c.JSON(http.StatusOK, output)
	}
}

func combineData(scrapin *entity.ScrapInResponseBody, brandfetch *entity.BrandfetchResponseBody, domain string) model.EnrichOrganizationResponseData {
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

func updateResponseWithScrapinData(d *model.EnrichOrganizationResponseData, scrapin *entity.ScrapInResponseBody, domain string) {
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

func updateResponseWithBrandfetchData(d *model.EnrichOrganizationResponseData, brandfetch *entity.BrandfetchResponseBody) {
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
