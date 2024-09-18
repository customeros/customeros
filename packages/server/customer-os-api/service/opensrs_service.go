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
	"time"
)

type OpenSRSResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	ErrorNumber int    `json:"error_number,omitempty"`
}

type MailboxDetails struct {
	Email             string   `json:"email"`
	ForwardingEnabled bool     `json:"forwardingEnabled"`
	ForwardingTo      []string `json:"forwardingTo"`
	WebmailEnabled    bool     `json:"webmailEnabled"`
}

type OpensrsService interface {
	SetupDomainForMailStack(ctx context.Context, tenant, domain string) error
	SetMailbox(ctx context.Context, tenant, domain, username, password string, forwardingEnabled bool, forwardingTo []string, webmailEnabled bool) error
	GetMailboxDetails(ctx context.Context, email string) (MailboxDetails, error)
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

	// step 2: Configure the domain in OpenSRS
	err = s.setEmailDomainInOpenSRS(ctx, domain, domainRecord.DkimPrivate)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to configure email domain in open SRS"))
		s.log.Error("failed to configure email domain in open SRS", err)
		return err
	}

	return nil
}

func (s *opensrsService) setEmailDomainInOpenSRS(ctx context.Context, domain, dkimPrivateKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.setEmailDomainInOpenSRS")
	defer span.Finish()
	span.LogKV("domain", domain)

	// Define the API endpoint (replace with your environment's URL)
	apiURL := s.cfg.ExternalServices.OpenSRS.Url + "/api/change_domain"

	// Prepare the request body
	requestBody := map[string]interface{}{
		"credentials": map[string]string{
			"user":     s.cfg.ExternalServices.OpenSRS.Username,
			"password": s.cfg.ExternalServices.OpenSRS.ApiKey,
		},
		"domain": domain,
		"attributes": map[string]interface{}{
			"dkim_selector": "dkim",
			"dkim_key":      dkimPrivateKey,
		},
	}

	// Convert the request body to JSON
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		s.log.Error("failed to marshal request body", err)
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to make API request"))
		s.log.Error("failed to make API request", err)
		return fmt.Errorf("failed to make API request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return errors.Wrap(err, "failed to read response body")
	}
	span.LogKV("responseBody", string(body))

	// Check for a successful response
	var response OpenSRSResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response"))
		s.log.Error("failed to unmarshal response", err)
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Check if the response indicates success
	if !response.Success {
		tracing.TraceErr(span, errors.New(response.Error))
		s.log.Error("API request failed", response.Error)
		return fmt.Errorf("API request failed: %s", response.Error)
	}

	return nil
}

func (s *opensrsService) SetMailbox(ctx context.Context, tenant, domain, username, password string, forwardingEnabled bool, forwardingTo []string, webmailEnabled bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.SetMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain, "username", username)

	// Define the API endpoint for adding a mailbox (replace with your environment's URL)
	apiURL := s.cfg.ExternalServices.OpenSRS.Url + "/api/change_user"

	attributes := map[string]interface{}{
		"type":           "mailbox",
		"password":       password,
		"delivery_local": true, // Store mail locally
	}

	if webmailEnabled {
		attributes["service_webmail"] = "enabled"
	} else {
		attributes["service_webmail"] = "disabled"
	}
	// Add forwarding options if enabled
	if forwardingEnabled && len(forwardingTo) > 0 {
		attributes["delivery_forward"] = true
		attributes["forward_recipients"] = forwardingTo
	}

	// Create the requestBody with the extracted attributes
	requestBody := map[string]interface{}{
		"credentials": map[string]string{
			"user":     s.cfg.ExternalServices.OpenSRS.Username,
			"password": s.cfg.ExternalServices.OpenSRS.ApiKey,
		},
		"user":       username + "@" + domain,
		"attributes": attributes,
	}

	// Convert the request body to JSON
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		s.log.Error("failed to marshal request body", err)
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return fmt.Errorf("failed to create HTTP request: %s", err.Error())
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to make API request"))
		s.log.Error("failed to make API request", err)
		return fmt.Errorf("failed to make API request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tracing.TraceErr(span, errors.New("API request failed"))
		s.log.Error("API request failed", err)
		return fmt.Errorf("API request failed")
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return err
	}
	span.LogKV("responseBody", string(body))

	// Check for a successful response
	var response OpenSRSResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response"))
		s.log.Error("failed to unmarshal response", err)
		return err
	}

	// Check if the response indicates success
	if !response.Success {
		tracing.TraceErr(span, errors.New(response.Error))
		s.log.Error("API request failed", response.Error)
		return err
	}

	return nil
}

func (s *opensrsService) GetMailboxDetails(ctx context.Context, email string) (MailboxDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensrsService.GetMailboxDetails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("email", email)

	// Define the API endpoint for getting mailbox information
	apiURL := s.cfg.ExternalServices.OpenSRS.Url + "/api/get_user"

	// Create the request body
	requestBody := map[string]interface{}{
		"credentials": map[string]string{
			"user":     s.cfg.ExternalServices.OpenSRS.Username,
			"password": s.cfg.ExternalServices.OpenSRS.ApiKey,
		},
		"user": email,
	}

	// Convert the request body to JSON
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request body"))
		s.log.Error("failed to marshal request body", err)
		return MailboxDetails{}, fmt.Errorf("failed to marshal request body: %s", err.Error())
	}

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return MailboxDetails{}, fmt.Errorf("failed to create HTTP request: %s", err.Error())
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to make API request"))
		s.log.Error("failed to make API request", err)
		return MailboxDetails{}, fmt.Errorf("failed to make API request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tracing.TraceErr(span, errors.New("API request failed"))
		s.log.Error("API request failed")
		return MailboxDetails{}, fmt.Errorf("API request failed")
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return MailboxDetails{}, err
	}
	span.LogKV("responseBody", string(body))

	// Define a map to parse the response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response"))
		s.log.Error("failed to unmarshal response", err)
		return MailboxDetails{}, err
	}

	// Check if the response indicates success
	if success, ok := response["success"].(bool); !ok || !success {
		errMessage := response["error"].(string)
		tracing.TraceErr(span, errors.New(errMessage))
		s.log.Error("API request failed", errMessage)
		return MailboxDetails{}, fmt.Errorf("API request failed: %s", errMessage)
	}

	// Extract the mailbox details: creation date and attributes
	attributes := response["attributes"].(map[string]interface{})
	mailboxDetails := MailboxDetails{
		Email:             email,
		ForwardingEnabled: attributes["delivery_forward"].(bool),
	}
	recipients := make([]string, 0)
	for _, recipient := range attributes["forward_recipients"].([]interface{}) {
		if str, ok := recipient.(string); ok {
			recipients = append(recipients, str)
		}
	}
	mailboxDetails.ForwardingTo = recipients
	if (attributes["service_webmail"].(string)) == "enabled" {
		mailboxDetails.WebmailEnabled = true
	}

	return mailboxDetails, nil
}
