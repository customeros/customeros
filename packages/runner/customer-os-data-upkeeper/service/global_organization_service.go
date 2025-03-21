package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/biter777/countries"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonconstants "github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/webscraper"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type EnrichOrganizationRequest struct {
	Domain      string `json:"domain"`
	LinkedinUrl string `json:"linkedinUrl"`
}

func (e *EnrichOrganizationRequest) Normalize() {
	e.LinkedinUrl = strings.TrimSpace(e.LinkedinUrl)
	e.Domain = strings.TrimSpace(e.Domain)
}

type GlobalOrganizationService interface {
	SyncDataIntoGlobalOrganizations()
	ScrapinCompanyByWebsite()
	EnrichGlobalOrganization()
	SyncGlobalOrgsToTenantOrganizations()
	ScrapeGlobalOrgs()
	ExtractWebpageLinks()
}

type globalOrganizationService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.CommonServices
}

func NewGlobalOrganizationService(cfg *config.Config, log logger.Logger, commonServices *commonservice.CommonServices) GlobalOrganizationService {
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

	// process records
	for _, record := range records {
		s.syncScrapinToGlobalOrganizationRecord(ctx, record)
	}
}

func (s *globalOrganizationService) syncScrapinToGlobalOrganizationRecord(ctx context.Context, record *postgres_entity.EnrichDetailsScrapIn) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.syncScrapinToGlobalOrganizationRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// mark record as synced initially to not process same record again, even if error occurs
	err := s.commonServices.PostgresRepositories.EnrichDetailsScrapInRepository.MarkSyncedToGlobalOrganizations(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
		s.log.Errorf("Error marking record as synced: %s", err.Error())
		return
	}

	if record.Data == "" {
		return
	}

	// unmarshal cached data
	data := postgresentity.ScrapInResponseBody{}
	if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin data"))
		return
	}

	if data.Company == nil {
		return
	}

	// if company website is missing skip processing as we cannot validate primary domain
	if data.Company.WebsiteUrl == "" {
		return
	}

	if s.commonServices.DomainService.IsKnownCompanyHostingUrl(ctx, data.Company.WebsiteUrl) {
		return
	}

	// identify primary domain
	accessible, _, primaryDomain := s.commonServices.DomainService.CheckDomainWithMailsherpa(ctx, data.Company.WebsiteUrl)
	if !accessible {
		return
	}

	// if primary domain is empty, skip processing
	if primaryDomain == "" {
		return
	}

	if !utils.IsValidDomain(primaryDomain) {
		return
	}

	// check if primary domain is accepted
	if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
		return
	}

	// if global organization already exists, update otherwise create
	globalOrganization, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting global organization by primary domain"))
		s.log.Errorf("Error getting global organization by primary domain: %s", err.Error())
		return
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
	if data.Company.Name != "" && globalOrganization.Name == "" {
		name := strings.TrimSpace(data.Company.Name)
		name = utils.SanitizeUTF8(name)
		globalOrganization.Name = name
	}
	if data.Company.Description != "" {
		globalOrganization.SourceDescription1 = data.Company.Description
	}
	if data.Company.Tagline != nil {
		if tagline, ok := data.Company.Tagline.(string); ok {
			globalOrganization.SourceDescription2 = tagline
		}
	}
	if data.Company.GetEmployeeCount() > 0 && globalOrganization.EmployeeCount == 0 {
		globalOrganization.EmployeeCount = data.Company.GetEmployeeCount()
	}
	if data.Company.WebsiteUrl != "" && globalOrganization.Website == "" {
		globalOrganization.Website = data.Company.WebsiteUrl
	}
	if data.Company.FoundedOn.Year > 0 && globalOrganization.YearFounded == 0 {
		globalOrganization.YearFounded = data.Company.FoundedOn.Year
	}
	if data.Company.LinkedInUrl != "" && globalOrganization.LinkedInUrl == "" {
		globalOrganization.LinkedInUrl = data.Company.LinkedInUrl
	}
	if data.Company.UniversalName != "" && globalOrganization.LinkedInAlias == "" {
		globalOrganization.LinkedInAlias = data.Company.UniversalName
	}
	if data.Company.Logo != "" && globalOrganization.LogoUrl == "" {
		globalOrganization.LogoUrl = data.Company.Logo
	}
	if data.Company.Headquarter.City != "" && globalOrganization.City == "" {
		globalOrganization.City = data.Company.Headquarter.City
	}
	if data.Company.Headquarter.GeographicArea != "" && globalOrganization.Region == "" {
		globalOrganization.Region = data.Company.Headquarter.GeographicArea
	}
	if data.Company.Headquarter.Country != "" && globalOrganization.CountryA2 == "" {
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
			return
		}
		_, err := s.commonServices.WebscraperService.Scrape(ctx, "https://"+globalOrganization.PrimaryDomain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	} else {
		_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Update(ctx, globalOrganization)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error updating global organization"))
			s.log.Errorf("Error updating global organization: %s", err.Error())
			return
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

	// process records
	for _, record := range records {
		s.syncBrandfetchToGlobalOrganizationRecord(ctx, record)
	}
}

func (s *globalOrganizationService) syncBrandfetchToGlobalOrganizationRecord(ctx context.Context, record *postgres_entity.EnrichDetailsBrandfetch) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.syncBrandfetchToGlobalOrganizationRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// mark record as synced initially to not process same record again, even if error occurs
	err := s.commonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.MarkSyncedToGlobalOrganizations(ctx, record.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking record as synced"))
		s.log.Errorf("Error marking record as synced: %s", err.Error())
		return
	}

	if record.Data == "" {
		return
	}

	// unmarshal cached data
	data := postgresentity.BrandfetchResponseBody{}
	if err = json.Unmarshal([]byte(record.Data), &data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch data"))
		return
	}

	// if company domin is missing skip processing as we cannot validate primary domain
	if data.Domain == "" {
		return
	}

	// check if website is accepted
	if s.commonServices.DomainService.IsKnownCompanyHostingUrl(ctx, data.Domain) {
		return
	}

	// identify primary domain
	accessible, _, primaryDomain := s.commonServices.DomainService.CheckDomainWithMailsherpa(ctx, data.Domain)
	if !accessible {
		return
	}

	// if primary domain is empty, skip processing
	if primaryDomain == "" {
		return
	}

	if !utils.IsValidDomain(primaryDomain) {
		return
	}

	if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
		return
	}

	globalOrganization, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting global organization by primary domain"))
		s.log.Errorf("Error getting global organization by primary domain: %s", err.Error())
		return
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
	if data.Name != "" && globalOrganization.Name == "" {
		name := strings.TrimSpace(data.Name)
		name = utils.SanitizeUTF8(name)
		globalOrganization.Name = name
	}
	if data.LongDescription != "" {
		globalOrganization.SourceDescription3 = data.LongDescription
	}
	if data.Description != "" {
		globalOrganization.SourceDescription4 = data.Description
	}
	if data.Company.GetEmployees() > 0 && globalOrganization.EmployeeCount == 0 {
		globalOrganization.EmployeeCount = data.Company.GetEmployees()
	}
	if data.Domain != "" && globalOrganization.Website == "" {
		globalOrganization.Website = data.Domain
	}
	if data.Company.FoundedYear > 0 && globalOrganization.YearFounded == 0 {
		globalOrganization.YearFounded = int(data.Company.FoundedYear)
	}
	if len(data.GetLogoUrls()) > 0 && globalOrganization.LogoUrl == "" {
		globalOrganization.LogoUrl = data.GetLogoUrls()[0]
	}
	if len(data.GetIconUrls()) > 0 && globalOrganization.IconUrl == "" {
		globalOrganization.IconUrl = data.GetIconUrls()[0]
	}
	if data.Company.Location.City != "" && globalOrganization.City == "" {
		globalOrganization.City = data.Company.Location.City
	}
	if data.Company.Location.Region != "" && globalOrganization.Region == "" {
		globalOrganization.Region = data.Company.Location.Region
	}
	if data.Company.Location.CountryCodeA2 != "" && globalOrganization.CountryA2 == "" {
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
	otherSocials := []string{}
	for _, link := range data.Links {
		if link.Url != "" && strings.Contains(link.Url, "linkedin.com/company") && globalOrganization.LinkedInUrl == "" {
			_, scrapinResponse, err := s.commonServices.EnrichmentService.ScrapInCompanyProfile(ctx, link.Url)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error calling scrapin for company profile"))
				s.log.Errorf("Error calling scrapin for company profile: %s", err.Error())
				continue
			}

			if scrapinResponse != nil && scrapinResponse.Company != nil {
				globalOrganization.LinkedInUrl = scrapinResponse.Company.LinkedInUrl
				if scrapinResponse.Company.UniversalName != "" && globalOrganization.LinkedInAlias == "" {
					globalOrganization.LinkedInAlias = scrapinResponse.Company.UniversalName
				}
				if scrapinResponse.Company.Logo != "" && globalOrganization.LogoUrl == "" {
					globalOrganization.LogoUrl = scrapinResponse.Company.Logo
				}
				if scrapinResponse.Company.Headquarter.City != "" && globalOrganization.City == "" {
					globalOrganization.City = scrapinResponse.Company.Headquarter.City
				}
				if scrapinResponse.Company.Headquarter.GeographicArea != "" && globalOrganization.Region == "" {
					globalOrganization.Region = scrapinResponse.Company.Headquarter.GeographicArea
				}
				if scrapinResponse.Company.Headquarter.Country != "" && globalOrganization.CountryA2 == "" {
					if strings.ToUpper(scrapinResponse.Company.Headquarter.Country) == "OO" {
						globalOrganization.CountryA2 = ""
					} else {
						country := countries.ByName(scrapinResponse.Company.Headquarter.Country)
						if country != countries.Unknown {
							globalOrganization.CountryA2 = country.Alpha2()
						} else {
							globalOrganization.CountryA2 = ""
						}
					}
				}
			}
		} else if link.Url != "" && !strings.Contains(link.Url, "linkedin.com") {
			otherSocials = append(otherSocials, link.Url)
		}
	}

	if len(data.Company.Industries) > 0 {
		industryDesc := "Business area: "
		for _, industry := range data.Company.Industries {
			industryDesc += industry.Name + "; "
		}
		globalOrganization.SourceDescription5 = industryDesc
	}

	if createGlobalOrg {
		// create global organization
		_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Create(ctx, globalOrganization)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error creating global organization"))
			s.log.Errorf("Error creating global organization: %s", err.Error())
			return
		}
		_, err := s.commonServices.WebscraperService.Scrape(ctx, "https://"+globalOrganization.PrimaryDomain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	} else {
		_, err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Update(ctx, globalOrganization)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error updating global organization"))
			s.log.Errorf("Error updating global organization: %s", err.Error())
			return
		}
	}
	if len(otherSocials) > 0 {
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.AddOtherSocials(ctx, globalOrganization.PrimaryDomain, otherSocials)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error adding other socials"))
			s.log.Errorf("Error adding other socials: %s", err.Error())
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

	// process records
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

	requestJSON, err := json.Marshal(EnrichOrganizationRequest{
		Domain: domain,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequestWithContext(ctx, "GET", s.cfg.Common.Internal.CustomerOsApi.ApiUrl+"/scrapinOrganization", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.Common.Internal.CustomerOsApi.ApiKey)

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

func (s *globalOrganizationService) ExtractWebpageLinks() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.ExtractWebpageLinks")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100
	webpages, err := s.commonServices.PostgresRepositories.ScrapedWebpageRepository.GetWebpagesWithoutLinks(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to scrape"))
		return
	}

	if len(webpages) == 0 {
		span.LogFields(log.Int("result.count", len(webpages)))
		return
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, 10) // Adjust the number based on your needs

	for _, page := range webpages {
		wg.Add(1)
		// Acquire semaphore
		semaphore <- struct{}{}

		go func(scrapedWebpage *postgres_entity.ScrapedWebpage) {
			// Create timeout context for this goroutine
			childCtx, childCancel := context.WithTimeout(ctx, 90*time.Second)
			defer childCancel()

			childSpan, childCtx := tracing.StartTracerSpan(childCtx, "GlobalOrganizationService.ExtractWebpageLinks")
			defer childSpan.Finish()
			childSpan.LogFields(log.String("url", scrapedWebpage.Url))

			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore when done

			content, links := s.commonServices.WebscraperService.ProcessWebContent(childCtx, scrapedWebpage.Content)
			err := s.commonServices.PostgresRepositories.ScrapedWebpageRepository.SetLinks(childCtx, scrapedWebpage.Url, links)
			if err != nil {
				tracing.TraceErr(childSpan, err)
			}
			err = s.commonServices.PostgresRepositories.ScrapedWebpageRepository.SetContent(childCtx, scrapedWebpage.Url, content)
			if err != nil {
				tracing.TraceErr(childSpan, err)
			}
		}(page)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (s *globalOrganizationService) ScrapeGlobalOrgs() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.ScrapeGlobalOrgs")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100
	orgs, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToScrape(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to scrape"))
		return
	}

	if len(orgs) == 0 {
		return
	}

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, 10) // Adjust the number based on your needs

	for _, org := range orgs {
		wg.Add(1)
		// Acquire semaphore
		semaphore <- struct{}{}

		go func(org *postgres_entity.GlobalOrganization) {
			// Create timeout context for this goroutine
			childCtx, childCancel := context.WithTimeout(ctx, 90*time.Second)
			defer childCancel()

			childSpan, childCtx := opentracing.StartSpanFromContext(childCtx, "GlobalOrganizationService.ScrapeGlobalOrg")
			defer childSpan.Finish()
			tracing.TagEntity(childSpan, org.PrimaryDomain)

			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore when done

			page := "https://" + org.PrimaryDomain
			contents, err := s.commonServices.WebscraperService.Scrape(childCtx, page)
			if err != nil {
				switch {
				case errors.Is(err, webscraper.ErrUnprocessable):
					childSpan.LogKV("error", "Unprocessable content")
					childSpan.LogKV("url", page)
				default:
					tracing.TraceErr(childSpan, errors.Wrap(err, "error scraping global org primary domain"))
				}

				err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetScrapeStatus(childCtx, org.ID, enum.ScrapeError)
				if err != nil {
					tracing.TraceErr(childSpan, errors.Wrap(err, "error updating global org scraped status"))
				}
			} else if contents == "" {
				err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetScrapeStatus(childCtx, org.ID, enum.ScrapeError)
				if err != nil {
					tracing.TraceErr(childSpan, errors.Wrap(err, "error updating global org scraped status"))
				}
			}

			if contents == "" {
				// delete global org after 5 attempts
				if org.ScrapeAttempt > 5 {
					err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.Delete(childCtx, org.ID)
					if err != nil {
						tracing.TraceErr(childSpan, errors.Wrap(err, "error deleting global org"))
					}
				}
			}

			if contents != "" {
				err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetScrapeStatus(childCtx, org.ID, enum.ScrapeCompleted)
				if err != nil {
					tracing.TraceErr(childSpan, errors.Wrap(err, "error updating global org scraped status"))
					return
				}
			}
		}(org)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func (s *globalOrganizationService) EnrichGlobalOrganization() {
	s.enrichNames()
	s.enrichIndustries()
	s.enrichDescriptions()
}

func (s *globalOrganizationService) enrichIndustries() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichIndustries")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 30
	hoursFromPreviousAttempt := 24
	maxAttempts := 3

	records, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToEnrichIndustry(ctx, hoursFromPreviousAttempt, maxAttempts, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to process"))
		s.log.Errorf("Error getting records to process: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	// process records
	for _, record := range records {
		s.enrichIndustry(ctx, record)
	}
}

func (s *globalOrganizationService) enrichIndustry(ctx context.Context, globalOrganization *postgresentity.GlobalOrganization) {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichIndustry")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagEntity(span, strconv.FormatUint(globalOrganization.ID, 10))

	// mark globalOrganization as processed initially to not process same globalOrganization again, even if error occurs
	err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkIndustryEnrichRequested(ctx, globalOrganization.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking globalOrganization as processed"))
		s.log.Errorf("Error marking globalOrganization as processed: %s", err.Error())
		return
	}

	url := "https://" + globalOrganization.PrimaryDomain
	scrapedPage, err := s.commonServices.WebscraperService.Scrape(ctx, url)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// Construct the prompt
	systemPrompt := `I'm going to provide you metadata about a company, including the company name, website url, description, and content scraped from their homepage (if available).
        Your job is to classify the NAICS industry code the company belongs to based on the data provided.
        Return only the single most specific and appropriate NAICS code (digits only, e.g. "541511") for the company. 

        Important details:
        - Use the latest NAICS codes available.
        - If multiple NAICS codes might apply, choose the best match (the most specific, relevant code).`

	var p strings.Builder
	p.WriteString(fmt.Sprintf("Company name: %s\n", globalOrganization.Name))
	p.WriteString(fmt.Sprintf("Website: %s\n", url))
	p.WriteString(fmt.Sprintf("Company description: %s\n", globalOrganization.Description))
	if scrapedPage != "" {
		p.WriteString("---Homepage content---\n")
		p.WriteString(scrapedPage)
	}
	prompt := p.String()

	temperature := float32(0.1)
	aiOutput, err := s.commonServices.AIService.AskAIForIndustryCode(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelGemini,
		SystemPrompt:     &systemPrompt,
		Prompt:           &prompt,
		ModelTemperature: &temperature,
		OutputFormat:     enum.AIOutputJson,
		MaxOutputTokens:  utils.Int32Ptr(250),
		Retries:          utils.IntPtr(2),
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error asking AI"))
		return
	}

	span.LogFields(log.Float64("result.confidence", aiOutput.Confidence))

	if aiOutput.Confidence < 0.5 {
		return
	}

	span.LogFields(log.String("result.code", aiOutput.Code))

	// get industry for organization
	industryEntity, err := s.commonServices.IndustryService.GetClosestByCode(ctx, aiOutput.Code)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting industry by code"))
		s.log.Errorf("Error getting industry by code: %s", err.Error())
		return
	}

	if industryEntity == nil {
		span.LogFields(log.Bool("result.industryFound", false))
		return
	}

	span.LogFields(log.Bool("result.industryFound", true))

	err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetIndustry(ctx, globalOrganization.ID, industryEntity.Code, industryEntity.Name)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error setting industry"))
		return
	}
}

func (s *globalOrganizationService) enrichDescriptions() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichDescriptions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 30
	hoursFromPreviousAttempt := 24
	maxAttempts := 3

	records, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToEnrichDescription(ctx, hoursFromPreviousAttempt, maxAttempts, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to process"))
		s.log.Errorf("Error getting records to process: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	// process records
	for _, record := range records {
		s.enrichDescription(ctx, record)
	}
}

func (s *globalOrganizationService) enrichDescription(ctx context.Context, globalOrganization *postgresentity.GlobalOrganization) {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichDescription")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagEntity(span, strconv.FormatUint(globalOrganization.ID, 10))

	// mark record as processed initially to not process same record again, even if error occurs
	err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkDescriptionEnrichRequested(ctx, globalOrganization.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking record as processed"))
		s.log.Errorf("Error marking record as processed: %s", err.Error())
		return
	}

	url := "https://" + globalOrganization.PrimaryDomain
	pageContent, err := s.commonServices.WebscraperService.Scrape(ctx, url)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// prepare Anthropic prompt
	var descLines []string
	descriptions := []string{globalOrganization.Description, globalOrganization.SourceDescription1, globalOrganization.SourceDescription2, globalOrganization.SourceDescription3, globalOrganization.SourceDescription4, globalOrganization.SourceDescription5}
	for i, d := range descriptions {
		if strings.TrimSpace(d) != "" {
			descLines = append(descLines, fmt.Sprintf("Description Line %d: %s", i+1, d))
		}
	}

	// Construct the prompt
	systemPrompt := `I am going to provide you metadata about a company, including the company name, website url, various descriptions from social media, and scraped content from their homepage (if available).  Your job is to write a clear, direct description of the company that explains who they serve and their revenue model.  

        Please return a single paragraph (max 300 characters), in American English.
        No marketing speak or jargon.`

	var p strings.Builder
	p.WriteString(fmt.Sprintf("Company name: %s\n", globalOrganization.Name))
	p.WriteString(fmt.Sprintf("Company website: %s\n", url))
	for _, desc := range descLines {
		p.WriteString(desc)
	}
	if pageContent != "" {
		p.WriteString("---Scraped homepage content---\n")
		p.WriteString(pageContent)
	}
	prompt := p.String()

	// ask AI for concise description
	temperature := float32(1.0)
	aiOutput, err := s.commonServices.AIService.AskAIForCompanyDescription(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelAnthropicHaiku,
		SystemPrompt:     &systemPrompt,
		Prompt:           &prompt,
		ModelTemperature: &temperature,
		OutputFormat:     enum.AIOutputJson,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error asking AI"))
		return
	}

	if aiOutput.Description == "" {
		aiOutput.Description = utils.FirstNotEmptyString(globalOrganization.Description, globalOrganization.SourceDescription3, globalOrganization.SourceDescription1, globalOrganization.SourceDescription4, globalOrganization.SourceDescription2)
	}

	span.LogFields(log.String("result.description", aiOutput.Description))

	err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDescription(ctx, globalOrganization.ID, aiOutput.Description)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error setting description"))
		return
	}
}

func (s *globalOrganizationService) enrichNames() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichNames")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 30
	hoursFromPreviousAttempt := 24
	maxAttempts := 3

	records, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetOrganizationsToEnrichName(ctx, hoursFromPreviousAttempt, maxAttempts, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to process"))
		s.log.Errorf("Error getting records to process: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	// process records
	for _, record := range records {
		s.enrichName(ctx, record)
	}
}

func (s *globalOrganizationService) enrichName(ctx context.Context, globalOrganization *postgres_entity.GlobalOrganization) {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichName")
	defer span.Finish()
	tracing.TagComponentCronJob(span)
	tracing.TagEntity(span, strconv.FormatUint(globalOrganization.ID, 10))

	// mark record as processed initially to not process same record again, even if error occurs
	err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkNameEnrichRequested(ctx, globalOrganization.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking record as processed"))
		s.log.Errorf("Error marking record as processed: %s", err.Error())
		return
	}

	url := "https://" + globalOrganization.PrimaryDomain
	pageContent, err := s.commonServices.WebscraperService.Scrape(ctx, url)

	var sp strings.Builder
	sp.WriteString("I am going to provide you with metadata about a company, which will include:\n")
	sp.WriteString("•	A current/partial company name\n")
	sp.WriteString("•	The company's website\n")
	sp.WriteString("•	Scraped page content from the company's website (if available)\n")
	sp.WriteString("Your task is to identify and return the commonly recognized company name in English, along with a confidence score between 0 and 1.\n")
	sp.WriteString("The confidence score should reflect how certain you are that the name your return accurately reflects the commonly recognized company name in English, with a score of 1 indicating absolute confidence.\n")
	sp.WriteString("SPECIAL INSTRUCTIONS:\n")
	sp.WriteString("If the business is more commonly recognized by a brand name e.g., Apple instead of Apple Inc., return that simpler, branded name.\n")
	sp.WriteString("Remove all suffixes like Inc, Ltd, LLC, Corp, etc\n")
	sp.WriteString("Return the name in title case EXCEPT when the recognized brand name is spelled in uppercase e.g. UPS or HSBC, or the commonly recognized brand starts with a lower case e.g. ebay\n")

	systemPrompt := sp.String()

	var p strings.Builder
	p.WriteString(fmt.Sprintf("Company name as we currently have it: %s\n", globalOrganization.Name))
	p.WriteString(fmt.Sprintf("Webpage url: %s\n", url))
	if pageContent != "" {
		p.WriteString("--- Webpage content --- \n")
		p.WriteString(pageContent)
	}
	prompt := p.String()

	// ask AI for concise description
	temperature := float32(0.1)
	aiOutput, err := s.commonServices.AIService.AskAIForCompanyName(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelLlama8B,
		SystemPrompt:     &systemPrompt,
		Prompt:           &prompt,
		ModelTemperature: &temperature,
		OutputFormat:     enum.AIOutputJson,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error asking AI"))
		return
	}

	span.LogFields(log.String("result.confidence", fmt.Sprintf("%f", aiOutput.Confidence)))
	if aiOutput.Confidence < 0.5 {
		return
	}

	span.LogFields(log.String("result.name", aiOutput.Name))
	err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetName(ctx, globalOrganization.ID, aiOutput.Name)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error setting name"))
		return
	}
}

func (s *globalOrganizationService) SyncGlobalOrgsToTenantOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.SyncGlobalOrgsToTenantOrganizations")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	daysFromPreviousSync := 1

	records, err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.GetGlobalOrganizationsToSyncIntoTenantOrganizations(ctx, daysFromPreviousSync, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting records to process"))
		s.log.Errorf("Error getting records to process: %s", err.Error())
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	// process records
	for _, record := range records {
		s.syncGlobalOrganizationToTenantOrganizationRecord(ctx, record)
	}
}

func (s *globalOrganizationService) syncGlobalOrganizationToTenantOrganizationRecord(ctx context.Context, globalOrganization *postgresentity.GlobalOrganization) {
	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.syncGlobalOrganizationToTenantOrganizationRecord")
	defer span.Finish()
	span.LogFields(log.Uint64("record.id", globalOrganization.ID), log.String("record.primaryDomain", globalOrganization.PrimaryDomain))
	tracing.TagEntity(span, globalOrganization.PrimaryDomain)

	// mark record as processed initially to not process same record again, even if error occurs
	err := s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkGlobalOrganizationSyncedToNeo(ctx, globalOrganization.ID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error marking record as processed"))
		s.log.Errorf("Error marking record as processed: %s", err.Error())
		return
	}

	// Find organizations by domain across all tenants
	tenantWithOrgId, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsByDomainAcrossAllTenants(ctx, globalOrganization.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting organizations by domain"))
		s.log.Errorf("Error getting organizations by domain: %s", err.Error())
		return
	}

	for _, tenantOrg := range tenantWithOrgId {
		innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenantOrg.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})
		s.syncOrganizationToTenant(innerCtx, globalOrganization, tenantOrg)
	}
}

func (s *globalOrganizationService) syncOrganizationToTenant(ctx context.Context, globalOrganization *postgresentity.GlobalOrganization, tenantOrg neo4jrepository.TenantAndOrganizationId) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.syncOrganizationToTenant")
	defer span.Finish()
	tenant := tenantOrg.Tenant
	organizationId := tenantOrg.OrganizationId
	tracing.TagTenant(span, tenant)
	span.LogKV("industryNaicsCode", globalOrganization.IndustryNaicsCode, "description", globalOrganization.Description, "name", globalOrganization.Name)

	// sync organization
	dataFields := data_fields.OrganizationFields{}
	if globalOrganization.IndustryNaicsCode != "" {
		dataFields.IndustryCode = utils.StringPtr(globalOrganization.IndustryNaicsCode)
	}
	if globalOrganization.Description != "" {
		dataFields.Description = utils.StringPtr(globalOrganization.Description)
	}
	if globalOrganization.Name != "" {
		dataFields.Name = utils.StringPtr(globalOrganization.Name)
	}
	if globalOrganization.LogoPath != "" {
		dataFields.LogoUrl = utils.StringPtr(commonconstants.S3ImagesCDN + globalOrganization.LogoPath)
	} else if globalOrganization.LogoUrl != "" {
		dataFields.LogoUrl = utils.StringPtr(globalOrganization.LogoUrl)
	}
	if globalOrganization.IconPath != "" {
		dataFields.IconUrl = utils.StringPtr(commonconstants.S3ImagesCDN + globalOrganization.IconPath)
	} else if globalOrganization.IconUrl != "" {
		dataFields.IconUrl = utils.StringPtr(globalOrganization.IconUrl)
	}
	_, err := s.commonServices.OrganizationService.Save(ctx, nil, utils.StringPtr(organizationId), dataFields)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error syncing organization"))
		s.log.Errorf("Error syncing organization: %s", err.Error())
	}

	// sync Linked In
	socialEntities, err := s.commonServices.SocialService.GetAllForEntities(ctx, tenant, model.ORGANIZATION, []string{organizationId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting social entities"))
		s.log.Errorf("Error getting social entities: %s", err.Error())
		return
	}

	if globalOrganization.LinkedInUrl != "" {
		// check social entites contains no linked in before adding
		skipAddingLinkedIn := false
		for _, socialEntity := range *socialEntities {
			if socialEntity.IsLinkedin() {
				skipAddingLinkedIn = true
			}
		}
		if !skipAddingLinkedIn {
			// Check no other org contains current linked in
			linkedInIdentifier := neo4jentity.SocialEntity{Url: globalOrganization.LinkedInUrl}.ExtractLinkedinCompanyIdentifierFromUrl()
			orgsWithLinkedIn, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, linkedInIdentifier, globalOrganization.LinkedInAlias, linkedInIdentifier)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error getting organizations by linked in"))
				s.log.Errorf("Error getting organizations by linked in: %s", err.Error())
				return
			}
			if len(orgsWithLinkedIn) > 0 {
				skipAddingLinkedIn = true
			}
			if !skipAddingLinkedIn {
				s.commonServices.SocialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
					Type: model.ORGANIZATION,
					Id:   tenantOrg.OrganizationId,
				}, neo4jentity.SocialEntity{
					Url:   globalOrganization.LinkedInUrl,
					Alias: globalOrganization.LinkedInAlias,
				})
			}
		}
	}

	// sync other socials
	for _, otherSocialUrl := range globalOrganization.OtherSocials {
		s.commonServices.SocialService.AddSocialToEntity(ctx, nil, common_srv.LinkWith{
			Type: model.ORGANIZATION,
			Id:   organizationId,
		}, neo4jentity.SocialEntity{
			Url: otherSocialUrl,
		})
	}
}
