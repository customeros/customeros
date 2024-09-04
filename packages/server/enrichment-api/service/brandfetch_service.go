package service

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/constants"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"math/rand"
	"net/http"
)

var nonRetryableErrors = []string{"Invalid Domain Name", "User is not authorized to access this resource with an explicit deny"}
var knownBrandfetchErrors = []string{"Invalid Domain Name", "User is not authorized to access this resource with an explicit deny", "API key quota exceeded", "Endpoint request timed out"}

type BrandfetchService interface {
	GetByDomain(ctx context.Context, domain string) (*postgresentity.BrandfetchResponseBody, error)
}

type brandfetchService struct {
	config   *config.Config
	services *Services
	log      logger.Logger
}

func NewBrandfetchService(config *config.Config, services *Services, log logger.Logger) BrandfetchService {
	return &brandfetchService{
		config:   config,
		services: services,
		log:      log,
	}
}

func (s *brandfetchService) GetByDomain(ctx context.Context, domain string) (*postgresentity.BrandfetchResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandfetchService.GetByDomain")
	defer span.Finish()
	span.LogKV("domain", domain)

	latestEnrichDetailsBrandfetchRecord, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.GetLatestByDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get brandfetch cached data"))
		return nil, err
	}

	var data *postgresentity.BrandfetchResponseBody

	callBrandfetch := false
	success := true

	if latestEnrichDetailsBrandfetchRecord == nil || latestEnrichDetailsBrandfetchRecord.UpdatedAt.AddDate(0, 0, s.config.BrandfetchConfig.TtlDays).Before(utils.Now()) {
		callBrandfetch = true
	} else if latestEnrichDetailsBrandfetchRecord.Success == false {
		// if latest record is not successful
		unmarshalledData := postgresentity.BrandfetchResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsBrandfetchRecord.Data), &unmarshalledData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch cached data"))
		}
		data = &unmarshalledData
		if utils.Contains(nonRetryableErrors, data.Message) {
			callBrandfetch = false
		}
	}

	// if cached data is missing or last time fetched > ttl refresh
	if callBrandfetch {
		// get data from brandfetch
		if data, err = s.callBrandfetch(ctx, domain); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to call brandfetch"))
			return nil, err
		}

		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal data"))
			return nil, err
		}

		if utils.Contains(knownBrandfetchErrors, data.Message) {
			success = false
		}

		_, err = s.services.CommonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.Create(ctx, postgresentity.EnrichDetailsBrandfetch{
			Domain:  domain,
			Data:    string(dataAsString),
			Success: success,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save brandfetch data in db"))
		}
	} else {
		// unmarshal cached data
		unmarshalledData := postgresentity.BrandfetchResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsBrandfetchRecord.Data), &unmarshalledData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch cached data"))
			return nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data
	if data == nil || !success {
		allSuccessRecords, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsBrandfetchRepository.GetAllSuccessByDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get brandfetch data"))
			return nil, err
		}
		if len(allSuccessRecords) > 0 {
			latestEnrichDetailsBrandfetchRecord = &allSuccessRecords[0]
			unmarshalledData := postgresentity.BrandfetchResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsBrandfetchRecord.Data), &unmarshalledData); err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch cached data"))
				return nil, err
			}
			data = &unmarshalledData
		}
	}

	return data, nil
}

func (s *brandfetchService) callBrandfetch(ctx context.Context, domain string) (*postgresentity.BrandfetchResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandfetchService.callBrandfetch")
	defer span.Finish()
	span.LogKV("domain", domain)

	brandfetchUrl := s.config.BrandfetchConfig.Url

	if brandfetchUrl == "" {
		err := errors.New("Brandfetch URL not set")
		tracing.TraceErr(span, err)
		s.log.Errorf("Brandfetch URL not set")
		return nil, err
	}

	// get current month in format yyyy-mm
	currentMonth := utils.Now().Format("2006-01")

	queryResult := s.services.CommonServices.PostgresRepositories.ExternalAppKeysRepository.GetAppKeys(ctx, constants.AppBrandfetch, currentMonth, h.cfg.Services.BrandfetchLimit)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		s.log.Errorf("Error getting brandfetch app keys: %s", queryResult.Error)
		return nil, queryResult.Error
	}
	branfetchAppKeys := queryResult.Result.([]postgresentity.ExternalAppKeys)
	if len(branfetchAppKeys) == 0 {
		err := errors.New(fmt.Sprintf("no brandfetch app keys available for %s", currentMonth))
		tracing.TraceErr(span, err)
		s.log.Errorf("No brandfetch app keys available for %s", currentMonth)
		return nil, err
	}
	// pick random app key from list
	appKey := branfetchAppKeys[rand.Intn(len(branfetchAppKeys))]

	body, err := makeBrandfetchHTTPRequest(brandfetchUrl, appKey.AppKey, domain)

	// Increment usage count of the app key
	queryResult = s.services.CommonServices.PostgresRepositories.ExternalAppKeysRepository.IncrementUsageCount(ctx, appKey.ID)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		s.log.Errorf("Error incrementing app key usage count: %v", queryResult.Error)
	}

	var brandfetchResponseBody postgresentity.BrandfetchResponseBody
	err = json.Unmarshal(body, &brandfetchResponseBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch response"))
		s.log.Errorf("Error unmarshalling brandfetch response: %s", err.Error())
		return nil, err
	}

	return &brandfetchResponseBody, nil
}

func makeBrandfetchHTTPRequest(baseUrl, apiKey, domain string) ([]byte, error) {
	url := baseUrl + "/" + domain

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}
