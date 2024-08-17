package service

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type ScrapinPersonService interface {
	ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (uint64, *postgresentity.ScrapInPersonResponse, error)
	ScrapInSearchPerson(ctx context.Context, email, fistName, lastName, domain string) (uint64, *postgresentity.ScrapInPersonResponse, error)
}

type ScrapInPersonSearchRequestParams struct {
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	CompanyDomain string `json:"companyDomain,omitempty"`
	Email         string `json:"email,omitempty"`
	LinkedInUrl   string `json:"linkedInUrl,omitempty"`
}

type scrapinPersonService struct {
	config   *config.Config
	services *Services
	log      logger.Logger
}

func NewPersonScrapeInService(config *config.Config, services *Services, log logger.Logger) ScrapinPersonService {
	return &scrapinPersonService{
		config:   config,
		services: services,
		log:      log,
	}
}

func (s scrapinPersonService) ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (uint64, *postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapinPersonService.ScrapInPersonProfile")
	defer span.Finish()
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	latestEnrichDetailsScrapInRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, linkedInUrl, postgresentity.ScrapInFlowPersonProfile)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgresentity.ScrapInPersonResponse
	var recordId uint64

	callScrapin := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		callScrapin = true
	} else if latestEnrichDetailsScrapInRecord.PersonFound == false {
		// if latest record has not found person response, check root cause
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
		}
		// if credits left is 0 or last attempt was > 1 day ago, call scrapin
		if data.CreditsLeft == 0 {
			callScrapin = true
		} else if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapin = true
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callScrapin {
		// get data from scrapin
		if data, err = s.callScrapinPersonProfile(ctx, linkedInUrl); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInPersonSearchRequestParams{
			LinkedInUrl: linkedInUrl,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.Create(ctx, postgresentity.EnrichDetailsScrapIn{
			Param1:        linkedInUrl,
			Flow:          postgresentity.ScrapInFlowPersonProfile,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   data.Person != nil,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
	}

	// if fresh data not found, check most recent cached data with person found
	if data == nil || (data != nil && data.Person == nil) {
		latestEnrichDetailsScrapInRecordWithPersonFound, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithPersonFound(ctx, linkedInUrl, postgresentity.ScrapInFlowPersonProfile)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithPersonFound != nil {
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithPersonFound.Data), &data); err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			return latestEnrichDetailsScrapInRecordWithPersonFound.ID, data, nil
		}
	}

	return recordId, data, nil
}

func (s *scrapinPersonService) callScrapinPersonProfile(ctx context.Context, linkedInUrl string) (*postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.callScrapinPersonProfile")
	defer span.Finish()
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	url := baseUrl + "/enrichment/profile" + "?apikey=" + scrapInApiKey + "&linkedInUrl=" + linkedInUrl

	body, err := makeScrapInHTTPRequest(url)

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	var scrapinResponse postgresentity.ScrapInPersonResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
		span.LogFields(log.String("response.body", string(body)))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	return &scrapinResponse, nil
}

func (s scrapinPersonService) ScrapInSearchPerson(ctx context.Context, email, firstName, lastName, domain string) (uint64, *postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapinPersonService.ScrapInSearchPerson")
	defer span.Finish()
	span.LogFields(log.String("email", email), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	latestEnrichDetailsScrapInRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlow(ctx, email, firstName, lastName, domain, postgresentity.ScrapInFlowPersonSearch)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin data"))
		return 0, nil, err
	}

	var data *postgresentity.ScrapInPersonResponse
	var recordId uint64

	callScrapin := false

	if latestEnrichDetailsScrapInRecord == nil || latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.TtlDays).Before(utils.Now()) {
		callScrapin = true
	} else if latestEnrichDetailsScrapInRecord.PersonFound == false {
		// if latest record has not found person response, check root cause
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
		}
		// if credits left is 0 or last attempt was > 1 day ago, call scrapin
		if data.CreditsLeft == 0 {
			callScrapin = true
		} else if latestEnrichDetailsScrapInRecord.UpdatedAt.AddDate(0, 0, 1).Before(utils.Now()) {
			callScrapin = true
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callScrapin {
		// get data from scrapin
		if data, err = s.callScrapinPersonSearch(ctx, email, firstName, lastName, domain); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to call scrapin"))
			return 0, nil, err
		}

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal data"))
			return 0, nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInPersonSearchRequestParams{
			Email:         email,
			FirstName:     firstName,
			LastName:      lastName,
			CompanyDomain: domain,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal params"))
		}
		dbRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.Create(ctx, postgresentity.EnrichDetailsScrapIn{
			Param1:        email,
			Param2:        firstName,
			Param3:        lastName,
			Param4:        domain,
			Flow:          postgresentity.ScrapInFlowPersonSearch,
			AllParamsJson: string(paramsAsString),
			Data:          string(dataAsString),
			PersonFound:   data.Person != nil,
			CompanyFound:  data.Company != nil,
			Success:       data.Success,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save scrapin data in db"))
		}
		if dbRecord != nil {
			recordId = dbRecord.ID
		}
	} else {
		recordId = latestEnrichDetailsScrapInRecord.ID
		// unmarshal cached data
		if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecord.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return 0, nil, err
		}
	}

	// if fresh data not found, check most recent cached data with person found
	if data == nil || (data != nil && data.Person == nil) {
		latestEnrichDetailsScrapInRecordWithPersonFound, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlowWithPersonFound(ctx, email, firstName, lastName, domain, postgresentity.ScrapInFlowPersonSearch)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin data"))
			return 0, nil, err
		}
		if latestEnrichDetailsScrapInRecordWithPersonFound != nil {
			if err = json.Unmarshal([]byte(latestEnrichDetailsScrapInRecordWithPersonFound.Data), &data); err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
				return 0, nil, err
			}
			return latestEnrichDetailsScrapInRecordWithPersonFound.ID, data, nil
		}
	}

	return recordId, data, nil
}

func (s *scrapinPersonService) callScrapinPersonSearch(ctx context.Context, email, firstName, lastName, domain string) (*postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.callScrapinPersonSearch")
	defer span.Finish()
	span.LogFields(log.String("email", email), log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	baseUrl := s.config.ScrapinConfig.Url
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	url := baseUrl + "/enrichment" + "?apikey=" + scrapInApiKey + "&email=" + email
	if firstName != "" {
		url += "&firstName=" + firstName
	}
	if lastName != "" {
		url += "&lastName=" + lastName
	}
	if domain != "" {
		url += "&companyDomain=" + domain
	}

	body, err := makeScrapInHTTPRequest(url)

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "makeScrapInHTTPRequest"))
		s.log.Errorf("Error making scrapin HTTP request: %s", err.Error())
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	var scrapinResponse postgresentity.ScrapInPersonResponse
	err = json.Unmarshal(body, &scrapinResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
		span.LogFields(log.String("response.body", string(body)))
		s.log.Errorf("Error unmarshalling scrapin response: %s", err.Error())
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	return &scrapinResponse, nil
}

func makeScrapInHTTPRequest(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}
