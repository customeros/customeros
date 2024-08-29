package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

var knowIpDataBadResponseMessages = []string{"is a reserved IP address"}

type IpIntelligenceService interface {
	LookupIp(ctx context.Context, ip string) (*postgresentity.IPDataResponseBody, error)
}

type ipIntelligenceService struct {
	config   *config.Config
	Services *Services
	log      logger.Logger
}

func NewIpIntelligenceService(config *config.Config, services *Services, log logger.Logger) IpIntelligenceService {
	return &ipIntelligenceService{
		config:   config,
		Services: services,
		log:      log,
	}
}

func (s *ipIntelligenceService) LookupIp(ctx context.Context, ip string) (*postgresentity.IPDataResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IpIntelligenceService.LookupIp")
	defer span.Finish()
	span.LogFields(log.String("ip", ip))

	cacheIpData, err := s.Services.CommonServices.PostgresRepositories.CacheIpDataRepository.Get(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return nil, err
	}
	var data *postgresentity.IPDataResponseBody
	// if cached data is missing or last time fetched > 90 days ago
	if cacheIpData == nil || cacheIpData.UpdatedAt.AddDate(0, 0, s.config.IpDataConfig.IpDataCacheTtlDays).Before(utils.Now()) {
		// get data from IPData
		if data, err = s.askIpData(ctx, ip); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get IPData"))
			return nil, err
		}
		// save to db
		dataAsString, err := json.Marshal(data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal data"))
			return nil, err
		}
		s.Services.CommonServices.PostgresRepositories.CacheIpDataRepository.Save(ctx, postgresentity.CacheIpData{
			Ip:   ip,
			Data: string(dataAsString),
		})
	} else {
		// unmarshal cached data
		data = &postgresentity.IPDataResponseBody{}
		if err = json.Unmarshal([]byte(cacheIpData.Data), data); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal cache data"))
			return nil, err
		}
	}

	data.Ip = ip

	return data, nil
}

func (s *ipIntelligenceService) askIpData(ctx context.Context, ip string) (*postgresentity.IPDataResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IpIntelligenceService.askIpData")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create IPData request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.config.IpDataConfig.ApiUrl, ip, s.config.IpDataConfig.ApiKey), nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create GET request for IPData"))
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to perform GET request for IPData")
		tracing.TraceErr(span, wrappedErr)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	knownBadResponse := false
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			for _, msg := range knowIpDataBadResponseMessages {
				if strings.Contains(string(responseBody), msg) {
					knownBadResponse = true
					break
				}
			}
		}
		if !knownBadResponse {
			span.LogFields(log.String("response.body", string(responseBody)))
			tracing.TraceErr(span, errors.Errorf("IPData returned status code %d", resp.StatusCode))
		}
	}

	// Parse the JSON request body
	var ipDataResponseBody postgresentity.IPDataResponseBody
	if err = json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		return nil, err
	}
	ipDataResponseBody.StatusCode = resp.StatusCode

	return &ipDataResponseBody, nil
}
