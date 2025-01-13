package verify

import (
	"context"

	"github.com/nyaruka/phonenumbers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type verifyService struct {
	log        logger.Logger
	cfg        *config.VerifyServiceConfig
	postgres   *repository.Repositories
	enrichment interfaces.EnrichmentService
	USClient   *extract.Client
	IntlClient *international_street.Client
}

func NewVerifyService(log logger.Logger, cfg *config.VerifyServiceConfig, postgres *repository.Repositories, enrichment interfaces.EnrichmentService) interfaces.VerifyService {
	return &verifyService{
		log:        log,
		cfg:        cfg,
		postgres:   postgres,
		enrichment: enrichment,
		USClient:   wireup.BuildUSExtractAPIClient(wireup.SecretKeyCredential(cfg.SmartyConfig.AuthId, cfg.SmartyConfig.AuthToken)),
		IntlClient: wireup.BuildInternationalStreetAPIClient(wireup.SecretKeyCredential(cfg.SmartyConfig.AuthId, cfg.SmartyConfig.AuthToken)),
	}
}

func (s *verifyService) IdentifyCompanyDomain(ctx context.Context, ipAddress string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.IsBot")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	// lookup company in enrich details tracking table
	data, err := s.postgres.EnrichDetailsTrackingRepository.GetByIP(ctx, ipAddress)
	if err != nil {
		return nil, err
	}

	if data.CompanyDomain != nil {
		return data.CompanyDomain, nil
	}

	// call snitcher if domain not known
	snitcherData, err := s.enrichment.Snitcher(ctx, ipAddress)
	if err != nil {
		return nil, err
	}

	if !snitcherData.CompanyFound {
		return nil, nil
	}

	return &snitcherData.Data.Domain, nil
}

func (s *verifyService) Threats(ctx context.Context, ipAddress string) (*interfaces.IpThreats, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.IsBot")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

	ipData, err := s.LookupIp(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	results := interfaces.IpThreats{
		IsAnonymous:   ipData.Threat.IsAnonymous,
		IsBogon:       ipData.Threat.IsBogon,
		IsDatacenter:  ipData.Threat.IsDatacenter,
		IsICloudRelay: ipData.Threat.IsIcloudRelay,
		IsKnownAbuser: ipData.Threat.IsKnownAbuser,
		IsProxy:       ipData.Threat.IsProxy,
		IsTor:         ipData.Threat.IsTor,
		IsVpn:         ipData.Threat.IsVpn,
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

func (s *verifyService) ValidatePhoneNumber(ctx context.Context, countryCodeA2 string, phoneNumber string) (*string, *string, error) {
	num, err := phonenumbers.Parse(phoneNumber, countryCodeA2)
	if err != nil {
		return nil, nil, err
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, nil, nil
	} else {
		e164 := phonenumbers.Format(num, phonenumbers.E164)
		extractedCountryCodeA2 := phonenumbers.GetRegionCodeForNumber(num)
		return &e164, &extractedCountryCodeA2, nil
	}
}
