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
	ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (*postgresentity.ScrapInPersonResponse, error)
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

func (s scrapinPersonService) ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (*postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapinPersonService.ScrapInPersonProfile")
	defer span.Finish()
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	scrapinPersonProfileData, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow(ctx, linkedInUrl, postgresentity.ScrapInFlowPersonProfile)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get scrapin data"))
		return nil, err
	}

	var data *postgresentity.ScrapInPersonResponse

	// if cached data is missing or last time fetched > ttl refresh
	if scrapinPersonProfileData == nil || scrapinPersonProfileData.UpdatedAt.AddDate(0, 0, s.config.ScrapinConfig.ScrapInTtlDays).Before(utils.Now()) {

		// get data from scrapin
		if data, err = s.callScrapinPersonProfile(ctx, linkedInUrl); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to call scrapin"))
			return nil, err
		}

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal data"))
			return nil, err
		}
		paramsAsString, err := json.Marshal(ScrapInPersonSearchRequestParams{
			LinkedInUrl: linkedInUrl,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal params"))
		}
		_, err = s.services.CommonServices.PostgresRepositories.EnrichDetailsScrapInRepository.Create(ctx, postgresentity.EnrichDetailsScrapIn{
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
	} else {
		// unmarshal cached data
		if err = json.Unmarshal([]byte(scrapinPersonProfileData.Data), &data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal scrapin cached data"))
			return nil, err
		}
	}

	return data, nil
}

func (s *scrapinPersonService) callScrapinPersonProfile(ctx context.Context, linkedInUrl string) (*postgresentity.ScrapInPersonResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapInService.callScrapinPersonProfile")
	defer span.Finish()
	span.LogFields(log.String("linkedInUrl", linkedInUrl))

	baseUrl := s.config.ScrapinConfig.ScrapInApiUrl
	if baseUrl == "" {
		err := errors.New("ScrapIn URL not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("ScrapIn URL not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}
	scrapInApiKey := s.config.ScrapinConfig.ScrapInApiKey
	if scrapInApiKey == "" {
		err := errors.New("Scrapin Api key not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("Scrapin Api key not set")
		return &postgresentity.ScrapInPersonResponse{}, err
	}

	url := baseUrl + "/enrichment/profile" + "?apikey=" + scrapInApiKey + "&linkedInUrl=" + linkedInUrl
	span.LogFields(log.String("url", url))

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
