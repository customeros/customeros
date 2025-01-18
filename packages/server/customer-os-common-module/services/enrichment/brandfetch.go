package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	AppBrandfetch = "brandfetch"
)

var (
	nonRetryableErrors    = []string{"Invalid Domain Name", "User is not authorized to access this resource with an explicit deny"}
	knownBrandfetchErrors = []string{"Invalid Domain Name", "User is not authorized to access this resource with an explicit deny", "API key quota exceeded", "Endpoint request timed out"}
)

func (s *enrichmentService) GetBrandfetchByDomain(ctx context.Context, domain string) (*entity.BrandfetchResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.GetBrandfetchByDomain")
	defer span.Finish()
	span.LogKV("domain", domain)

	latestEnrichDetailsBrandfetchRecord, err := s.postgres.EnrichDetailsBrandfetchRepository.GetLatestByDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get brandfetch cached data"))
		return nil, err
	}

	var data *entity.BrandfetchResponseBody

	callBrandfetch := false
	success := true

	if latestEnrichDetailsBrandfetchRecord == nil || latestEnrichDetailsBrandfetchRecord.UpdatedAt.AddDate(0, 0, s.config.BrandfetchConfig.TtlDays).Before(utils.Now()) {
		callBrandfetch = true
	} else if latestEnrichDetailsBrandfetchRecord.Success == false {
		// if latest record is not successful
		unmarshalledData := entity.BrandfetchResponseBody{}
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

		_, err = s.postgres.EnrichDetailsBrandfetchRepository.Create(ctx, entity.EnrichDetailsBrandfetch{
			Domain:  domain,
			Data:    string(dataAsString),
			Success: success,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save brandfetch data in db"))
		}
	} else {
		// unmarshal cached data
		unmarshalledData := entity.BrandfetchResponseBody{}
		if err = json.Unmarshal([]byte(latestEnrichDetailsBrandfetchRecord.Data), &unmarshalledData); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch cached data"))
			return nil, err
		}
		data = &unmarshalledData
	}

	// if fresh data not found, check most recent cached data
	if data == nil || !success {
		allSuccessRecords, err := s.postgres.EnrichDetailsBrandfetchRepository.GetAllSuccessByDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get brandfetch data"))
			return nil, err
		}
		if len(allSuccessRecords) > 0 {
			latestEnrichDetailsBrandfetchRecord = &allSuccessRecords[0]
			unmarshalledData := entity.BrandfetchResponseBody{}
			if err = json.Unmarshal([]byte(latestEnrichDetailsBrandfetchRecord.Data), &unmarshalledData); err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal brandfetch cached data"))
				return nil, err
			}
			data = &unmarshalledData
		}
	}
	span.LogFields(log.Bool("result.success", success))

	return data, nil
}

func (s *enrichmentService) callBrandfetch(ctx context.Context, domain string) (*entity.BrandfetchResponseBody, error) {
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

	queryResult := s.postgres.ExternalAppKeysRepository.GetAppKeys(ctx, AppBrandfetch, currentMonth, s.config.BrandfetchConfig.Limit)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		s.log.Errorf("Error getting brandfetch app keys: %s", queryResult.Error)
		return nil, queryResult.Error
	}
	branfetchAppKeys := queryResult.Result.([]entity.ExternalAppKeys)
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
	queryResult = s.postgres.ExternalAppKeysRepository.IncrementUsageCount(ctx, appKey.ID)
	if queryResult.Error != nil {
		tracing.TraceErr(span, queryResult.Error)
		s.log.Errorf("Error incrementing app key usage count: %v", queryResult.Error)
	}

	var brandfetchResponseBody entity.BrandfetchResponseBody
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
