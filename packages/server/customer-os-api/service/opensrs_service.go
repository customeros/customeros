package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

type OpensrsService interface {
	SetupDomainForMailStack(ctx context.Context, tenant, domain string) error
}

type opensrsService struct {
	log      logger.Logger
	services *Services
	cfg      *config.Config
}

// NewOpensrsService initializes the OpensrsService
func NewOpensrsService(log logger.Logger, services *Services, cfg *config.Config) OpensrsService {
	return &opensrsService{
		log:      log,
		services: services,
		cfg:      cfg,
	}
}

func (s *opensrsService) SetupDomainForMailStack(ctx context.Context, tenant, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.SetupDomainForMailStack")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	// step 1: get domain record from the database
	domainRecord, err := s.services.CommonServices.PostgresRepositories.MailStackDomainRepository.GetDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get domain record"))
		s.log.Error("failed to get domain record", err)
		return err
	}
	if domainRecord == nil {
		tracing.TraceErr(span, errors.New("domain record not found"))
		s.log.Errorf("domain record not found for domain")
		return errors.New("domain record not found")
	}
	//dkimPrivateKey := domainRecord.DkimPrivate

	// step 2: Register the domain in Opensrs
	err = s.registerDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to register domain"))
		s.log.Error("failed to register domain")
		return err
	}

	return nil
}

// checkDomainExists checks if the domain is already registered in OpenSRS
func (s *opensrsService) checkDomainExists(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.checkDomainExists")
	defer span.Finish()
	span.LogKV("domain", domain)

	// Construct the OpenSRS API endpoint for checking domain existence
	apiURL := fmt.Sprintf("%s/check_domain/%s", s.cfg.ExternalServices.OpenSRS.Url, url.PathEscape(domain))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		s.log.Error("failed to create request")
		return false, errors.Wrap(err, "failed to create request")
	}

	// Set headers for OpenSRS API authentication
	req.Header.Set("Authorization", "Bearer "+s.cfg.ExternalServices.OpenSRS.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to execute OpenSRS API request"))
		s.log.Error("failed to execute OpenSRS API request")
		return false, errors.Wrap(err, "failed to execute OpenSRS API request")
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body")
		return false, errors.Wrap(err, "failed to read response body")
	}

	span.LogKV("responseBody", string(body))

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		// Domain does not exist
		return false, nil
	} else if resp.StatusCode != http.StatusOK {
		// Unexpected status code, return error
		return false, fmt.Errorf("unexpected status code from OpenSRS API: %d", resp.StatusCode)
	}

	// Assume success means the domain exists
	return true, nil
}

func (s *opensrsService) registerDomain(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.registerDomain")
	defer span.Finish()
	span.LogKV("domain", domain)

	// if domain already exists, return
	exists, err := s.checkDomainExists(ctx, domain)
	if err != nil {
		return errors.Wrap(err, "failed to check domain existence")
	}
	if exists {
		return nil
	}

	apiURL := fmt.Sprintf("%s/add_domain", s.cfg.ExternalServices.OpenSRS.Url)

	// Create the request payload
	payload := map[string]string{
		"domain": domain,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	// Set headers for OpenSRS API authentication
	req.Header.Set("Authorization", "Bearer "+s.cfg.ExternalServices.OpenSRS.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to execute OpenSRS API request"))
		s.log.Error("failed to execute OpenSRS API request")
		return errors.Wrap(err, "failed to execute OpenSRS API request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body")
		return err
	}
	span.LogKV("responseBody", string(body))

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register domain: %s", string(body))
	}

	return nil
}
