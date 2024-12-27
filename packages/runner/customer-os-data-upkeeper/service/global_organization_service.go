package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/biter777/countries"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type GlobalOrganizationService interface {
	SyncDataIntoGlobalOrganizations()
	ScrapinCompanyByWebsite()
}

type globalOrganizationService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.Services
}

func NewGlobalOrganizationService(cfg *config.Config, log logger.Logger, commonServices *commonService.Services) GlobalOrganizationService {
	return &globalOrganizationService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *globalOrganizationService) SyncDataIntoGlobalOrganizations() {
	s.syncScrapinToGlobalOrganization()
	s.syncBrandfetchToGlobalOrganization()
}

func (s *globalOrganizationService) syncScrapinToGlobalOrganization() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.syncScrapinToGlobalOrganization")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 50

	records, err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetToSyncIntoGlobalOrganizations(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to sync"))
		s.log.Errorf("Error getting records to sync: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//process records
	for _, record := range records {

		// mark record as synced initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.MarkSyncedToGlobalOrganizations(ctx, record.ID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
			s.log.Errorf("Error marking record as synced: %s", err.Error())
			continue
		}

		if record.Data == "" {
			continue
		}

		// unmarshal cached data
		data := postgresentity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin data"))
			continue
		}

		if data.Company == nil {
			continue
		}

		// if company website is missing skip processing as we cannot validate primary domain
		if data.Company.WebsiteUrl == "" {
			continue
		}

		// identify primary domain
		_, primaryDomain := domaincheck.PrimaryDomainCheck(data.Company.WebsiteUrl)

		// if primary domain is empty, skip processing
		if primaryDomain == "" {
			continue
		}

		// check if website is accepted
		if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, data.Company.WebsiteUrl) ||
			!s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
			continue
		}

		// if global organization already exists, update otherwise create
		globalOrganization, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting global organization by primary domain"))
			s.log.Errorf("Error getting global organization by primary domain: %s", err.Error())
			continue
		}

		createGlobalOrg := false
		if globalOrganization == nil {
			createGlobalOrg = true
			now := utils.Now()
			globalOrganization = &postgresentity.GlobalOrganization{
				PrimaryDomain: primaryDomain,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
		}

		// populate global organization entity
		if data.Company.Name != "" {
			globalOrganization.Name = data.Company.Name
		}
		if data.Company.Description != "" {
			globalOrganization.Description = data.Company.Description
		}
		if data.Company.Tagline != nil {
			if tagline, ok := data.Company.Tagline.(string); ok {
				globalOrganization.ValueProposition = tagline
			}
		}
		if data.Company.GetEmployeeCount() > 0 {
			globalOrganization.EmployeeCount = data.Company.GetEmployeeCount()
		}
		if data.Company.WebsiteUrl != "" {
			globalOrganization.Website = data.Company.WebsiteUrl
		}
		if data.Company.FoundedOn.Year > 0 {
			globalOrganization.YearFounded = data.Company.FoundedOn.Year
		}
		if data.Company.LinkedInUrl != "" {
			globalOrganization.LinkedInUrl = data.Company.LinkedInUrl
		}
		if data.Company.UniversalName != "" {
			globalOrganization.LinkedInAlias = data.Company.UniversalName
		}
		if data.Company.Logo != "" {
			globalOrganization.LogoUrl = data.Company.Logo
		}
		if data.Company.Headquarter.City != "" {
			globalOrganization.City = data.Company.Headquarter.City
		}
		if data.Company.Headquarter.Country != "" {
			if strings.ToUpper(data.Company.Headquarter.Country) == "OO" {
				globalOrganization.CountryA2 = ""
			} else {
				country := countries.ByName(data.Company.Headquarter.Country)
				if country != countries.Unknown {
					globalOrganization.CountryA2 = country.Alpha2()
				} else {
					globalOrganization.CountryA2 = ""
				}
			}
		}

		if createGlobalOrg {
			// create global organization
			_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Create(ctx, globalOrganization)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error creating global organization"))
				s.log.Errorf("Error creating global organization: %s", err.Error())
				continue
			}
		} else {
			_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Update(ctx, globalOrganization)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error updating global organization"))
				s.log.Errorf("Error updating global organization: %s", err.Error())
				continue
			}
		}
	}
}

func (s *globalOrganizationService) syncBrandfetchToGlobalOrganization() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.syncBrandfetchToGlobalOrganization")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 50

	records, err := s.commonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.GetToSyncIntoGlobalOrganizations(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to sync"))
		s.log.Errorf("Error getting records to sync: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//process records
	for _, record := range records {

		// mark record as synced initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.MarkSyncedToGlobalOrganizations(ctx, record.ID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
			s.log.Errorf("Error marking record as synced: %s", err.Error())
			continue
		}

		if record.Data == "" {
			continue
		}

		// unmarshal cached data
		data := postgresentity.BrandfetchResponseBody{}
		if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch data"))
			continue
		}

		// if company domin is missing skip processing as we cannot validate primary domain
		if data.Domain == "" {
			continue
		}

		// identify primary domain
		_, primaryDomain := domaincheck.PrimaryDomainCheck(data.Domain)

		// if primary domain is empty, skip processing
		if primaryDomain == "" {
			continue
		}

		// check if website is accepted
		if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, data.Domain) ||
			!s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
			continue
		}

		globalOrganization, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting global organization by primary domain"))
			s.log.Errorf("Error getting global organization by primary domain: %s", err.Error())
			continue
		}

		// if global organization already exists, skip processing
		if globalOrganization != nil {
			continue
		}

		now := utils.Now()
		globalOrganization = &postgresentity.GlobalOrganization{
			PrimaryDomain: primaryDomain,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// populate global organization entity
		if data.Name != "" {
			globalOrganization.Name = data.Name
		}
		if data.LongDescription != "" {
			globalOrganization.Description = data.LongDescription
		}
		if data.Description != "" {
			globalOrganization.ValueProposition = data.Description
		}
		if data.Company.GetEmployees() > 0 {
			globalOrganization.EmployeeCount = data.Company.GetEmployees()
		}
		if data.Domain != "" {
			globalOrganization.Website = data.Domain
		}
		if data.Company.FoundedYear > 0 {
			globalOrganization.YearFounded = int(data.Company.FoundedYear)
		}
		if len(data.GetLogoUrls()) > 0 {
			globalOrganization.LogoUrl = data.GetLogoUrls()[0]
		}
		if len(data.GetIconUrls()) > 0 {
			globalOrganization.IconUrl = data.GetIconUrls()[0]
		}
		if data.Company.Location.City != "" {
			globalOrganization.City = data.Company.Location.City
		}
		if data.Company.Location.CountryCodeA2 != "" {
			if strings.ToUpper(data.Company.Location.CountryCodeA2) == "OO" {
				globalOrganization.CountryA2 = ""
			} else {
				country := countries.ByName(data.Company.Location.CountryCodeA2)
				if country != countries.Unknown {
					globalOrganization.CountryA2 = country.Alpha2()
				} else {
					globalOrganization.CountryA2 = ""
				}
			}
		}
		for _, link := range data.Links {
			if link.Url != "" && strings.Contains(link.Url, "linkedin.com/company") {
				globalOrganization.LinkedInUrl = link.Url
				break
			}
		}

		// create global organization
		_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Create(ctx, globalOrganization)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error creating global organization"))
			s.log.Errorf("Error creating global organization: %s", err.Error())
			continue
		}
	}
}

func (s *globalOrganizationService) ScrapinCompanyByWebsite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.ScrapinCompanyByWebsite")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 10

	records, err := s.commonServices.PostgresRepositories.GlobalOrganizationWebsiteToProcessRepository.GetWebsitesToProcess(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to process"))
		s.log.Errorf("Error getting records to process: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	//process records
	for _, record := range records {
		// identify primary domain
		_, primaryDomain := domaincheck.PrimaryDomainCheck(record.Website)
		notes := ""
		if primaryDomain == "" {
			notes = "Primary domain not found"
		}

		// mark record as processed initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.GlobalOrganizationWebsiteToProcessRepository.MarkAsProcessed(ctx, record.ID, notes)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
			s.log.Errorf("Error marking record as synced: %s", err.Error())
			continue
		}

		if primaryDomain == "" {
			continue
		}

		// call scrapin for primary domain
		err = s.callApiScrapinOrganization(ctx, primaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error calling enrichment api"))
			return
		}
	}
}

func (s *globalOrganizationService) callApiScrapinOrganization(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.callApiScrapinOrganization")
	defer span.Finish()
	span.LogKV("domain", domain)

	requestJSON, err := json.Marshal(enrichmentmodel.EnrichOrganizationRequest{
		Domain: domain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.EnrichmentApiConfig.Url+"/scrapinOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.EnrichmentApiConfig.ApiKey)

	// Make the HTTP request, retry once if response status is 502
	var response *http.Response
	client := &http.Client{}

	// Make the HTTP request
	response, err = client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return err
	}
	defer response.Body.Close() // Ensures the body is closed only once

	span.LogFields(log.Int("response.statusCode", response.StatusCode))

	if response.StatusCode != http.StatusOK {
		s.log.Errorf("Scrapin organization API response status code is : %d", response.StatusCode)
	}

	return nil
}
