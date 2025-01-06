package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"

	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type VerifyService interface {
	Threats(ctx context.Context, ipAddress string) (*IpThreats, error)
	IdentifyCompanyDomain(ctx context.Context, ipAddress string) (*string, error)
}

type verifyService struct {
	cfg      *config.GlobalConfig
	services *Services
}

func NewVerifyService(services *Services, cfg *config.GlobalConfig) VerifyService {
	return &verifyService{
		services: services,
		cfg:      cfg,
	}
}

type IpThreats struct {
	IsThreat      bool
	IsAnonymous   bool
	IsBogon       bool
	IsDatacenter  bool
	IsICloudRelay bool
	IsKnownAbuser bool
	IsProxy       bool
	IsTor         bool
	IsVpn         bool
}

func (s *verifyService) IdentifyCompanyDomain(ctx context.Context, ipAddress string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.IsBot")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	// lookup company in enrich details tracking table
	data, err := s.services.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ipAddress)
	if err != nil {
		return nil, err
	}

	if data.CompanyDomain != nil {
		return data.CompanyDomain, nil
	}

	// call snitcher if domain not known
	snitcherData, err := s.services.EnrichmentService.Snitcher(ctx, ipAddress)
	if err != nil {
		return nil, err
	}

	if !snitcherData.CompanyFound {
		return nil, nil
	}

	return &snitcherData.Data.Domain, nil
}

func (s *verifyService) Threats(ctx context.Context, ipAddress string) (*IpThreats, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.IsBot")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

	ipData, err := s.callVerifyAPIForIpData(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	results := IpThreats{
		IsAnonymous:   ipData.IpData.Threat.IsAnonymous,
		IsBogon:       ipData.IpData.Threat.IsBogon,
		IsDatacenter:  ipData.IpData.Threat.IsDatacenter,
		IsICloudRelay: ipData.IpData.Threat.IsIcloudRelay,
		IsKnownAbuser: ipData.IpData.Threat.IsKnownAbuser,
		IsProxy:       ipData.IpData.Threat.IsProxy,
		IsTor:         ipData.IpData.Threat.IsTor,
		IsVpn:         ipData.IpData.Threat.IsVpn,
	}

	if results.IsAnonymous ||
		results.IsBogon ||
		results.IsDatacenter ||
		results.IsICloudRelay ||
		results.IsKnownAbuser ||
		results.IsProxy ||
		results.IsTor ||
		results.IsVpn {
		results.IsThreat = true
	}

	return &results, nil
}

func (s *verifyService) callVerifyAPIForIpData(ctx context.Context, ipAddress string) (*validationmodel.IpLookupResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.callVerifyAPIForIpData")
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

	if net.ParseIP(ipAddress) == nil {
		err := errors.New("invalid IP address")
		tracing.TraceErr(span, err)
		return nil, err
	}

	requestJSON, err := json.Marshal(validationmodel.IpLookupRequest{
		Ip: ipAddress,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", s.cfg.InternalServices.ValidationApiConfig.Url+"/ipLookup", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.InternalServices.ValidationApiConfig.ApiKey)
	req.Header.Set(security.TenantHeader, "")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()

	var result validationmodel.IpLookupResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return nil, err
	}

	return &result, nil
}
