package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	enrichmentmodel "github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
)

type EnrichmentService interface {
	CallApiFindWorkEmail(ctx context.Context, firstName, lastName, companyName, companyDomain, linkedInUrl string, enrichPhoneNumber bool) (*enrichmentmodel.FindWorkEmailResponse, error)
}

type enrichmentService struct {
	log      logger.Logger
	services *Services
	cfg      *config.Config
}

func NewEnrichmentService(log logger.Logger, services *Services, cfg *config.Config) EnrichmentService {
	return &enrichmentService{
		log:      log,
		services: services,
		cfg:      cfg,
	}
}

func (s *enrichmentService) CallApiFindWorkEmail(ctx context.Context, firstName, lastName, companyName, companyDomain, linkedInUrl string, enrichPhoneNumber bool) (*enrichmentmodel.FindWorkEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.CallApiFindWorkEmail")
	defer span.Finish()

	requestJSON, err := json.Marshal(enrichmentmodel.FindWorkEmailRequest{
		LinkedinUrl:       linkedInUrl,
		FirstName:         firstName,
		LastName:          lastName,
		CompanyName:       companyName,
		CompanyDomain:     companyDomain,
		EnrichPhoneNumber: enrichPhoneNumber,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("GET", s.cfg.InternalServices.EnrichmentApiUrl+"/findWorkEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.InternalServices.EnrichmentApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.status.findWorkEmail", response.StatusCode))

	var findWorkEmailApiResponse enrichmentmodel.FindWorkEmailResponse
	err = json.NewDecoder(response.Body).Decode(&findWorkEmailApiResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode find work email response"))
		return nil, err
	}
	return &findWorkEmailApiResponse, nil
}
