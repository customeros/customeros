package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"io"
	"net/http"
	"net/url"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/customeros/mailsherpa/mailvalidate"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type ScrapInSearchRequestParams struct {
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	CompanyDomain string `json:"companyDomain,omitempty"`
	Email         string `json:"email,omitempty"`
	LinkedInUrl   string `json:"linkedInUrl,omitempty"`
}

func (s *enrichmentService) ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (uint64, *postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.ScrapInPersonProfile")
	defer spans.Finish()

	spans.LogKV("linkedInUrl", linkedInUrl)

	latestEnrichDetailsScrapInRecord, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, linkedInUrl, postgres_entity.ScrapInFlowPersonProfile)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgres_entity.ScrapInResponseBody
	var recordId uint64

	callScrapInNow := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		// if no cached data found, or cache data is older than ttl (90 days), call scrapin
		callScrapInNow = true
	} else if latestEnrichDetailsScrapInRecord.PersonFound == false {
		// if last attempt was > 1 day ago, call scrapin
		if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapInNow = true
		}
	}

	if callScrapInNow {
		// get data from scrapin
		if data, err = s.callScrapinPersonProfile(ctx, linkedInUrl); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		traceIfCreditsDepleting(data, *spans)

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInSearchRequestParams{
			LinkedInUrl: linkedInUrl,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.postgres.EnrichDetailsScrapInRepository.Create(ctx, postgres_entity.EnrichDetailsScrapIn{
			Param1:        linkedInUrl,
			Flow:          postgres_entity.ScrapInFlowPersonProfile,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   data.Person != nil,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		unmarshalledData := postgres_entity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &unmarshalledData); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data with person found
	if data == nil || data.Person == nil {
		latestEnrichDetailsScrapInRecordWithPersonFound, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithPersonFound(ctx, linkedInUrl, postgres_entity.ScrapInFlowPersonProfile)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to extract cached scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithPersonFound != nil {
			unmarshalledData := postgres_entity.ScrapInResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithPersonFound.Data), &unmarshalledData); err != nil {
				spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			data = &unmarshalledData
			return latestEnrichDetailsScrapInRecordWithPersonFound.ID, data, nil
		}
	}

	return recordId, data, nil
}

func (s *enrichmentService) callScrapinPersonProfile(ctx context.Context, linkedInUrl string) (*postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.callScrapinPersonProfile")
	defer spans.Finish()

	spans.LogKV("linkedInUrl", linkedInUrl)

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		spans.TraceError(err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		spans.TraceError(err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	params := url.Values{}
	params.Add("apikey", scrapInApiKey)
	params.Add("linkedInUrl", linkedInUrl)

	scrapinStatusCode, body, err := makeScrapInHTTPRequest(baseUrl + "/enrichment/profile" + "?" + params.Encode())
	if err != nil {
		spans.TraceError(errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	spans.LogKV("scrapin.statusCode", scrapinStatusCode)

	var scrapinResponse postgres_entity.ScrapInResponseBody
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "json.Unmarshal"))
		spans.LogKV("response.body", string(body))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	return &scrapinResponse, nil
}

func (s *enrichmentService) ScrapInSearchPerson(ctx context.Context, email, firstName, lastName, domain, companyName string) (uint64, *postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.ScrapInSearchPerson")
	defer spans.Finish()

	spans.LogKV("email", email)
	spans.LogKV("firstName", firstName)
	spans.LogKV("lastName", lastName)
	spans.LogKV("domain", domain)
	spans.LogKV("companyName", companyName)

	latestEnrichDetailsScrapInRecord, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlow(ctx, email, firstName, lastName, domain, companyName, postgres_entity.ScrapInFlowPersonSearch)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgres_entity.ScrapInResponseBody
	var recordId uint64

	callScrapInNow := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		callScrapInNow = true
	} else if latestEnrichDetailsScrapInRecord.PersonFound == false {
		// if last attempt was > 1 day ago, call scrapin
		if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapInNow = true
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callScrapInNow {
		// get data from scrapin
		if data, err = s.callScrapinPersonSearch(ctx, email, firstName, lastName, domain, companyName); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		traceIfCreditsDepleting(data, *spans)

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInSearchRequestParams{
			Email:         email,
			FirstName:     firstName,
			LastName:      lastName,
			CompanyDomain: domain,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.postgres.EnrichDetailsScrapInRepository.Create(ctx, postgres_entity.EnrichDetailsScrapIn{
			Param1:        email,
			Param2:        firstName,
			Param3:        lastName,
			Param4:        domain,
			Flow:          postgres_entity.ScrapInFlowPersonSearch,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   data.Person != nil,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		unmarshalledData := postgres_entity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &unmarshalledData); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data with person found
	if data == nil || data.Person == nil {
		latestEnrichDetailsScrapInRecordWithPersonFound, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlowWithPersonFound(ctx, email, firstName, lastName, domain, companyName, postgres_entity.ScrapInFlowPersonSearch)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithPersonFound != nil {
			unmarshalledData := postgres_entity.ScrapInResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithPersonFound.Data), &unmarshalledData); err != nil {
				spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			data = &unmarshalledData
			recordId = latestEnrichDetailsScrapInRecordWithPersonFound.ID
		}
	}

	// validate correct person identified by checking domains of email and company
	if data != nil && data.Person != nil {
		inputPrimaryDomains := []string{}
		outputPrimaryDomains := []string{}
		// collect primary domains from input params
		if email != "" {
			syntaxValidation := mailvalidate.ValidateEmailSyntax(email)
			if syntaxValidation.IsValid {
				_, primaryDomain := domaincheck.PrimaryDomainCheck(syntaxValidation.Domain)
				inputPrimaryDomains = append(inputPrimaryDomains, primaryDomain)
			}
		}
		if domain != "" {
			_, primaryDomain := domaincheck.PrimaryDomainCheck(domain)
			inputPrimaryDomains = append(inputPrimaryDomains, primaryDomain)
		}
		// collect primary domains from scrapin results
		for data.Company != nil {
			_, primaryDomain := domaincheck.PrimaryDomainCheck(data.Company.WebsiteUrl)
			outputPrimaryDomains = append(outputPrimaryDomains, primaryDomain)
			break
		}
		// check if any of the input primary domains match with output primary domains
		matchFound := false
		if len(inputPrimaryDomains) > 0 && len(outputPrimaryDomains) > 0 {
			for _, inputPrimaryDomain := range inputPrimaryDomains {
				for _, outputPrimaryDomain := range outputPrimaryDomains {
					if inputPrimaryDomain == outputPrimaryDomain {
						matchFound = true
						break
					}
				}
				if matchFound {
					break
				}
			}
		}
		if !matchFound {
			spans.LogKV("inputPrimaryDomains", fmt.Sprintf("%v", inputPrimaryDomains), "outputPrimaryDomains", fmt.Sprintf("%v", outputPrimaryDomains))
			spans.LogKV("result.error", "Person identified does not match with input params")
			return 0, nil, nil
		}
	}

	return recordId, data, nil
}

func (s *enrichmentService) callScrapinPersonSearch(ctx context.Context, email, firstName, lastName, domain, companyName string) (*postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.callScrapinPersonSearch")
	defer spans.Finish()

	spans.LogKV("email", email, "firstName", firstName, "lastName", lastName, "domain", domain, "companyName", companyName)

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		spans.TraceError(err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		spans.TraceError(err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	params := url.Values{}
	params.Add("apikey", scrapInApiKey)
	params.Add("email", email)
	if firstName != "" {
		params.Add("firstName", firstName)
	}
	if lastName != "" {
		params.Add("lastName", lastName)
	}
	if domain != "" {
		params.Add("companyDomain", domain)
	}
	if companyName != "" {
		params.Add("companyName", companyName)
	}

	scrapinStatusCode, body, err := makeScrapInHTTPRequest(baseUrl + "/enrichment" + "?" + params.Encode())
	if err != nil {
		spans.TraceError(errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	spans.LogKV("scrapin.statusCode", scrapinStatusCode)

	var scrapinResponse postgres_entity.ScrapInResponseBody
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "json.Unmarshal"))
		spans.LogKV("response.body", string(body))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	return &scrapinResponse, nil
}

func (s *enrichmentService) ScrapInCompanyProfile(ctx context.Context, linkedInUrl string) (uint64, *postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.ScrapInCompanyProfile")
	defer spans.Finish()

	spans.LogKV("linkedInUrl", linkedInUrl)

	latestEnrichDetailsScrapInRecord, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, linkedInUrl, postgres_entity.ScrapInFlowCompanyProfile)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgres_entity.ScrapInResponseBody
	var recordId uint64

	callScrapInNow := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		callScrapInNow = true
	} else if latestEnrichDetailsScrapInRecord.CompanyFound == false {
		// if last attempt was > 1 day ago, call scrapin
		if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapInNow = true
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callScrapInNow {
		// get data from scrapin
		if data, err = s.callScrapinCompanyProfile(ctx, linkedInUrl); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		traceIfCreditsDepleting(data, *spans)

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInSearchRequestParams{
			LinkedInUrl: linkedInUrl,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.postgres.EnrichDetailsScrapInRepository.Create(ctx, postgres_entity.EnrichDetailsScrapIn{
			Param1:        linkedInUrl,
			Flow:          postgres_entity.ScrapInFlowCompanyProfile,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   false,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		unmarshalledData := postgres_entity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &unmarshalledData); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data with company found
	if data == nil || data.Company == nil {
		latestEnrichDetailsScrapInRecordWithCompanyFound, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithCompanyFound(ctx, linkedInUrl, postgres_entity.ScrapInFlowCompanyProfile)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithCompanyFound != nil {
			unmarshalledData := postgres_entity.ScrapInResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithCompanyFound.Data), &unmarshalledData); err != nil {
				spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			data = &unmarshalledData
			return latestEnrichDetailsScrapInRecordWithCompanyFound.ID, data, nil
		}
	}

	return recordId, data, nil
}

func (s *enrichmentService) callScrapinCompanyProfile(ctx context.Context, linkedInUrl string) (*postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ScrapInService.callScrapinCompanyProfile")
	defer spans.Finish()
	spans.LogKV("linkedInUrl", linkedInUrl)

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		spans.TraceError(err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		spans.TraceError(err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	params := url.Values{}
	params.Add("apikey", scrapInApiKey)
	params.Add("linkedInUrl", linkedInUrl)

	scrapinStatusCode, body, err := makeScrapInHTTPRequest(baseUrl + "/enrichment/company" + "?" + params.Encode())
	if err != nil {
		spans.TraceError(errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	spans.LogKV("scrapin.statusCode", scrapinStatusCode)

	var scrapinResponse postgres_entity.ScrapInResponseBody
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "json.Unmarshal"))
		spans.LogKV("response.body", string(body))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	return &scrapinResponse, nil
}

func (s *enrichmentService) ScrapInSearchCompany(ctx context.Context, domain string) (uint64, *postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ScrapinService.ScrapInSearchCompany")
	defer spans.Finish()
	spans.LogKV("domain", domain)

	latestEnrichDetailsScrapInRecord, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, domain, postgres_entity.ScrapInFlowCompanySearch)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgres_entity.ScrapInResponseBody
	var recordId uint64

	callScrapInNow := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		callScrapInNow = true
	} else if latestEnrichDetailsScrapInRecord.CompanyFound == false {
		// if last attempt was > 1 day ago, call scrapin
		if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapInNow = true
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callScrapInNow {
		// get data from scrapin
		if data, err = s.callScrapinCompanySearch(ctx, domain); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		traceIfCreditsDepleting(data, *spans)

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInSearchRequestParams{
			CompanyDomain: domain,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.postgres.EnrichDetailsScrapInRepository.Create(ctx, postgres_entity.EnrichDetailsScrapIn{
			Param1:        domain,
			Flow:          postgres_entity.ScrapInFlowCompanySearch,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   false,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		unmarshalledData := postgres_entity.ScrapInResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &unmarshalledData); err != nil {
			spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data with company found
	if data == nil || data.Company == nil {
		latestEnrichDetailsScrapInRecordWithCompanyFound, err := s.postgres.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithCompanyFound(ctx, domain, postgres_entity.ScrapInFlowCompanySearch)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to get scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithCompanyFound != nil {
			unmarshalledData := postgres_entity.ScrapInResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithCompanyFound.Data), &unmarshalledData); err != nil {
				spans.TraceError(errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			data = &unmarshalledData

			if data.Company != nil && data.Company.WebsiteUrl != "" {
				// check primary domain matches
				inputIsPrimary, inputAltPrimaryDomain := domaincheck.PrimaryDomainCheck(domain)
				outputIsPrimary, outputAltPrimaryDomain := domaincheck.PrimaryDomainCheck(data.Company.WebsiteUrl)
				inputPrimaryDomain := utils.ExtractDomain(inputAltPrimaryDomain)
				if inputIsPrimary {
					inputPrimaryDomain = utils.ExtractDomain(domain)
				}
				outputPrimaryDomain := utils.ExtractDomain(outputAltPrimaryDomain)
				if outputIsPrimary {
					outputPrimaryDomain = utils.ExtractDomain(data.Company.WebsiteUrl)
				}
				spans.LogKV("result.inputPrimaryDomain", inputPrimaryDomain, "result.outputPrimaryDomain", outputPrimaryDomain)
				if inputPrimaryDomain != "" && inputPrimaryDomain == outputPrimaryDomain {
					return latestEnrichDetailsScrapInRecordWithCompanyFound.ID, data, nil
				}
			}
			return 0, nil, nil
		}
	}

	if data != nil && data.Company != nil && data.Company.WebsiteUrl != "" {
		// check primary domain matches
		inputIsPrimary, inputAltPrimaryDomain := domaincheck.PrimaryDomainCheck(domain)
		outputIsPrimary, outputAltPrimaryDomain := domaincheck.PrimaryDomainCheck(data.Company.WebsiteUrl)
		inputPrimaryDomain := utils.ExtractDomain(inputAltPrimaryDomain)
		if inputIsPrimary {
			inputPrimaryDomain = utils.ExtractDomain(domain)
		}
		outputPrimaryDomain := utils.ExtractDomain(outputAltPrimaryDomain)
		if outputIsPrimary {
			outputPrimaryDomain = utils.ExtractDomain(data.Company.WebsiteUrl)
		}
		spans.LogKV("result.inputPrimaryDomain", inputPrimaryDomain, "outputPrimaryDomain", outputPrimaryDomain)
		if inputPrimaryDomain != "" && inputPrimaryDomain == outputPrimaryDomain {
			return recordId, data, nil
		} else {
			spans.LogKV("result.info", fmt.Sprintf("Scrapin retuned different company domain, expected: %s, got: %s", inputPrimaryDomain, outputPrimaryDomain))
		}
	}

	return 0, nil, nil
}

func (s *enrichmentService) callScrapinCompanySearch(ctx context.Context, domain string) (*postgres_entity.ScrapInResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ScrapInService.callScrapinCompanySearch")
	defer spans.Finish()
	spans.LogKV("domain", domain)

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		spans.TraceError(err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		spans.TraceError(err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	params := url.Values{}
	params.Add("apikey", scrapInApiKey)
	params.Add("domain", domain)
	scrapinUrl := baseUrl + "/enrichment/company/domain" + "?" + params.Encode()

	scrapinStatusCode, body, err := makeScrapInHTTPRequest(scrapinUrl)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}
	spans.LogKV("scrapin.statusCode", scrapinStatusCode)

	var scrapinResponse postgres_entity.ScrapInResponseBody
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "json.Unmarshal"))
		spans.LogKV("response.body", string(body))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgres_entity.ScrapInResponseBody{}, err
	}

	return &scrapinResponse, nil
}

func makeScrapInHTTPRequest(url string) (int, []byte, error) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	var body []byte
	var err error
	var res *http.Response
	var statusCode int

	// try few times to get data from scrapin
	for i := 0; i < 2; i++ {
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return -1, nil, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		statusCode = res.StatusCode

		if statusCode == http.StatusTooManyRequests {
			// sleep 1 second and try again
			time.Sleep(1 * time.Second)
			continue
		} else {
			return statusCode, body, err
		}
	}
	return statusCode, body, err
}

func traceIfCreditsDepleting(data *postgres_entity.ScrapInResponseBody, spans telemetry.Spans) {
	if data != nil {
		if data.CreditsLeft > 0 && data.CreditsLeft < 50 {
			spans.TraceError(errors.New(fmt.Sprintf("ScrapIn credits are depleting, only %d credits left", data.CreditsLeft)))
		}
		if data.RateLimitLeft > 0 && data.RateLimitLeft < 50 {
			spans.TraceError(errors.New(fmt.Sprintf("ScrapIn rate limit is depleting, only %d requests left", data.RateLimitLeft)))
		}
	}
}
