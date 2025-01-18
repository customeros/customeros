package enrichment

import (
	"context"
	"errors"
	"strings"

	"github.com/biter777/countries"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type enrichmentService struct {
	log      logger.Logger
	config   *config.ExternalServicesConfig
	postgres *repository.Repositories
}

func NewEnrichmentService(log logger.Logger, config *config.ExternalServicesConfig, postgres *repository.Repositories) interfaces.EnrichmentService {
	return &enrichmentService{
		log:      log,
		config:   config,
		postgres: postgres,
	}
}

func (s *enrichmentService) EnrichPerson(ctx context.Context, person interfaces.PersonSearch) (*uint64, *entity.ScrapInResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.EnrichPerson")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if person.Email == nil && person.Domain == nil && person.LinkedinURL == nil {
		err := errors.New("missing email, domain, or linkedinURL")
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	var errs error

	var personData *entity.ScrapInResponseBody
	// scrapin by linkedinURL
	if person.LinkedinURL != nil {
		_, results, err := s.ScrapInPersonProfile(ctx, *person.LinkedinURL)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		personData = results
	}

	// search by email, domain, and company name if not found in step 1
	if personData != nil || (person.Email == nil && person.Domain == nil && person.CompanyName == nil) {
		return nil, personData, errs
	}

	recordID, response, err := s.ScrapInSearchPerson(
		ctx, *person.Email, *person.FirstName, *person.LastName, *person.Domain, *person.CompanyName,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		errs = multierr.Append(errs, err)
	}
	return &recordID, response, nil
}

func (s *enrichmentService) EnrichOrganization(ctx context.Context, domain, linkedinURL *string) (*interfaces.OrganizationData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.EnrichOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if domain == nil && linkedinURL == nil {
		err := errors.New("missing linkedinURL or domain")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var errs error
	// Scrapin by linkedinURL
	var scrapinResults *entity.ScrapInResponseBody
	if linkedinURL != nil {
		_, results, err := s.ScrapInCompanyProfile(ctx, *linkedinURL)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		scrapinResults = results
	}

	// Scrapin by domain
	if scrapinResults == nil && domain != nil {
		_, results, err := s.ScrapInSearchCompany(ctx, *domain)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		scrapinResults = results
	}

	// Brandfetch
	brandfetchResults, err := s.GetBrandfetchByDomain(ctx, *domain)
	if err != nil {
		tracing.TraceErr(span, err)
		errs = multierr.Append(errs, err)
	}

	if scrapinResults == nil && brandfetchResults == nil {
		return nil, errs
	}

	// Combine data
	return s.combineOrganizationEnrichmentData(ctx, scrapinResults, brandfetchResults, *domain), nil
}

func (s *enrichmentService) combineOrganizationEnrichmentData(ctx context.Context, scrapin *entity.ScrapInResponseBody, brandfetch *entity.BrandfetchResponseBody, domain string) *interfaces.OrganizationData {
	data := interfaces.OrganizationData{}
	s.updateResponseWithScrapinData(&data, scrapin, domain)
	s.updateResponseWithBrandfetchData(&data, brandfetch)
	s.normalizeCountry(&data.Location)
	return &data
}

func (s *enrichmentService) normalizeCountry(data *interfaces.OrganizationLocation) {
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

func (s *enrichmentService) updateResponseWithScrapinData(d *interfaces.OrganizationData, scrapin *entity.ScrapInResponseBody, domain string) {
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
		d.Socials = append(d.Socials, interfaces.OrganizationSocial{
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

func (s *enrichmentService) updateResponseWithBrandfetchData(d *interfaces.OrganizationData, brandfetch *entity.BrandfetchResponseBody) {
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
			d.Socials = append(d.Socials, interfaces.OrganizationSocial{
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
