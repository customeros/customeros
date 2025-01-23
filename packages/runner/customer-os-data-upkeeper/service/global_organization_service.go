package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/biter777/countries"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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
}

type globalOrganizationService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
}

func NewGlobalOrganizationService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) GlobalOrganizationService {
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

func (s *globalOrganizationService) EnrichGlobalOrganization() {
	s.enrichName()
	s.enrichIndustry()
	s.enrichDescription()
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

		if s.commonServices.DomainService.IsKnownCompanyHostingUrl(ctx, data.Company.WebsiteUrl) {
			continue
		}

		// identify primary domain
		accessible, _, primaryDomain := s.commonServices.DomainService.CheckDomainWithMailsherpa(ctx, data.Company.WebsiteUrl)
		if !accessible {
			continue
		}

		// if primary domain is empty, skip processing
		if primaryDomain == "" {
			continue
		}

		if !utils.IsValidDomain(primaryDomain) {
			continue
		}

		// check if primary domain is accepted
		if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
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

	// process records
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

		// check if website is accepted
		if s.commonServices.DomainService.IsKnownCompanyHostingUrl(ctx, data.Domain) {
			continue
		}

		// identify primary domain
		accessible, _, primaryDomain := s.commonServices.DomainService.CheckDomainWithMailsherpa(ctx, data.Domain)
		if !accessible {
			continue
		}

		// if primary domain is empty, skip processing
		if primaryDomain == "" {
			continue
		}

		if !utils.IsValidDomain(primaryDomain) {
			continue
		}

		if !s.commonServices.DomainService.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
			continue
		}

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
		for _, link := range data.Links {
			if link.Url != "" && strings.Contains(link.Url, "linkedin.com/company") && globalOrganization.LinkedInUrl == "" {
				globalOrganization.LinkedInUrl = link.Url
				break
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

func (s *globalOrganizationService) enrichIndustry() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichIndustry")
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
		span, ctx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.enrichIndustry.Record")
		defer span.Finish()
		span.LogFields(log.Uint64("record.id", record.ID))

		// mark record as processed initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkIndustryEnrichRequested(ctx, record.ID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking record as processed"))
			s.log.Errorf("Error marking record as processed: %s", err.Error())
			continue
		}

		// prepare Anthropic prompt
		var descLines []string
		descriptions := []string{record.Description, record.SourceDescription1, record.SourceDescription2, record.SourceDescription3, record.SourceDescription4, record.SourceDescription5}
		for i, d := range descriptions {
			if strings.TrimSpace(d) != "" {
				descLines = append(descLines, fmt.Sprintf("Description Line %d: %s", i+1, d))
			}
		}

		// Construct the prompt
		systemPrompt := `
        You are a NAICS industry classification expert. 
        Return only the single most specific and appropriate NAICS code (digits only, e.g. "541511") for the company based on the provided domain, name, and description. 
        No additional text or commentary.

        Important details:
        - Use the latest NAICS codes available.
        - If multiple NAICS codes might apply, choose the best match (the most specific, relevant code).
        - Do not output any text besides the NAICS code itself.
        `

		prompt := map[string]string{
			"domain":      record.PrimaryDomain,
			"name":        record.Name,
			"description": strings.Join(descLines, "\n"),
		}

		// ask AI for NAICS code
		aiOutput, err := s.commonServices.AIService.AskAI(ctx, enum.AIModelAnthropicHaiku, &systemPrompt, &prompt)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error asking AI"))
			continue
		}
		code := utils.IfNotNilString(aiOutput)
		span.LogFields(log.String("result.code", code))

		if code == "" {
			continue
		}

		// get industry for organization
		industryEntity, err := s.commonServices.IndustryService.GetClosestByCode(ctx, code)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error getting industry by code"))
			s.log.Errorf("Error getting industry by code: %s", err.Error())
			continue
		}

		if industryEntity == nil {
			span.LogFields(log.Bool("result.industryFound", false))
			continue
		}

		span.LogFields(log.Bool("result.industryFound", true))

		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetIndustry(ctx, record.ID, industryEntity.Code, industryEntity.Name)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error setting industry"))
			continue
		}
	}
}

func (s *globalOrganizationService) enrichDescription() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichDescription")
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
		recordSpan, recordCtx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.enrichDescription.Record")
		defer recordSpan.Finish()
		recordSpan.LogFields(log.Uint64("record.id", record.ID))

		// mark record as processed initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkDescriptionEnrichRequested(recordCtx, record.ID)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error marking record as processed"))
			s.log.Errorf("Error marking record as processed: %s", err.Error())
			continue
		}

		// prepare Anthropic prompt
		var descLines []string
		descriptions := []string{record.Description, record.SourceDescription1, record.SourceDescription2, record.SourceDescription3, record.SourceDescription4, record.SourceDescription5}
		for i, d := range descriptions {
			if strings.TrimSpace(d) != "" {
				descLines = append(descLines, fmt.Sprintf("Description Line %d: %s", i+1, d))
			}
		}

		// Construct the prompt
		systemPrompt := `
        You are a company analyst who writes clear, direct business descriptions. 
        Return a single paragraph (max 300 characters), in American English, explaining who the company serves and their revenue model. 
        No marketing speak or jargon. 
        Return "N/A" if insufficient information.`

		prompt := map[string]interface{}{
			"name":   record.Name,
			"domain": record.PrimaryDomain,
			"data":   strings.Join(descLines, "\n"),
		}

		// ask AI for concise description
		aiOutput, err := s.commonServices.AIService.AskAI(recordCtx, enum.AIModelAnthropicHaiku, &systemPrompt, &prompt)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error asking AI"))
			continue
		}
		aiDescrition := utils.IfNotNilString(aiOutput)
		recordSpan.LogFields(log.String("result.description", aiDescrition))

		if aiDescrition == "" || aiDescrition == "N/A" {
			aiDescrition = utils.FirstNotEmptyString(record.Description, record.SourceDescription3, record.SourceDescription1, record.SourceDescription4, record.SourceDescription2)
		}

		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetDescription(recordCtx, record.ID, aiDescrition)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error setting description"))
			continue
		}
	}
}

func (s *globalOrganizationService) enrichName() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.enrichName")
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
		recordSpan, recordCtx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.enrichName.Record")
		defer recordSpan.Finish()
		recordSpan.LogFields(log.Uint64("record.id", record.ID))

		// mark record as processed initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkNameEnrichRequested(recordCtx, record.ID)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error marking record as processed"))
			s.log.Errorf("Error marking record as processed: %s", err.Error())
			continue
		}

		// prepare Anthropic prompt
		systemPrompt := `
You are a world-class naming assistant. The user will provide:
	•	A current/partial company name
	•	The company’s primary domain
	•	The company’s website

Your task:
	1.	Identify the standard recognized company name.
	2.	If the business is more commonly recognized by a brand name (e.g., “Apple” instead of “Apple Inc.”), return that simpler, branded name.
	3.	If you find a longer official name that’s different from the brand, remove suffixes like “Inc,” “Ltd,” “LLC,” “Corp,” etc., but keep the rest of the formal name.
    4. 	If you cannot identify a single valid name, return “N/A”
	5.	Only output the name itself or “N/A,” with no explanations, disclaimers, or additional text.
	6.  Output in english.

Important Samples:
	•	If input suggests “Verizon Communications Inc.,” return “Verizon”
	•	If input suggests “Apple Inc,” return “Apple”
	•	If brand differs from legal name (e.g., “Nestlé S.A.” vs. “Nescafé”), return the official company name “Nestle”
	•	If unsure, return “N/A”

Provide the final name or N/A as your entire response.
`

		prompt := map[string]interface{}{
			"name":    record.Name,
			"domain":  record.PrimaryDomain,
			"website": record.Website,
		}

		// ask AI for concise description
		aiOutput, err := s.commonServices.AIService.AskAI(recordCtx, enum.AIModelAnthropicHaiku, &systemPrompt, &prompt)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error asking AI"))
			continue
		}
		aiName := utils.IfNotNilString(aiOutput)
		recordSpan.LogFields(log.String("result.name", aiName))

		if aiName == "" || aiName == "N/A" {
			continue
		}

		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.SetName(recordCtx, record.ID, aiName)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error setting name"))
			continue
		}
	}
}

func (s *globalOrganizationService) SyncGlobalOrgsToTenantOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "GlobalOrganizationService.SyncGlobalOrgsToTenantOrganizations")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	daysFromPreviousSync := 30

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
		recordSpan, recordCtx := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationService.SyncGlobalOrgsToTenantOrganizations.Record")
		defer recordSpan.Finish()
		recordSpan.LogFields(log.Uint64("record.id", record.ID), log.String("record.primaryDomain", record.PrimaryDomain))
		tracing.TagEntity(recordSpan, record.PrimaryDomain)

		// mark record as processed initially to not process same record again, even if error occurs
		err = s.commonServices.PostgresRepositories.GlobalOrganizationRepository.MarkGlobalOrganizationSyncedToNeo(recordCtx, record.ID)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error marking record as processed"))
			s.log.Errorf("Error marking record as processed: %s", err.Error())
			continue
		}

		// Find organizations by domain across all tenants
		tenantWithOrgId, err := s.commonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsByDomainAcrossAllTenants(recordCtx, record.PrimaryDomain)
		if err != nil {
			tracing.TraceErr(recordSpan, errors.Wrap(err, "error getting organizations by domain"))
			s.log.Errorf("Error getting organizations by domain: %s", err.Error())
			continue
		}

		for _, tenantOrgs := range tenantWithOrgId {
			innerCtx := common.WithCustomContext(recordCtx, &common.CustomContext{
				Tenant:    tenantOrgs.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			if record.IndustryNaicsCode == "" && record.Description == "" {
				continue
			}

			// sync organization
			dataFields := data_fields.OrganizationFields{}
			if record.IndustryNaicsCode != "" {
				dataFields.IndustryCode = utils.StringPtr(record.IndustryNaicsCode)
			}
			if record.Description != "" {
				dataFields.Description = utils.StringPtr(record.Description)
			}
			if record.Name != "" {
				dataFields.Name = utils.StringPtr(record.Name)
			}
			_, err = s.commonServices.OrganizationService.Save(innerCtx, nil, utils.StringPtr(tenantOrgs.OrganizationId), dataFields)
			if err != nil {
				tracing.TraceErr(recordSpan, errors.Wrap(err, "error syncing organization"))
				s.log.Errorf("Error syncing organization: %s", err.Error())
				continue
			}
		}
	}
}
