package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type DNSConfig struct {
	RecordType string
	Name       string
	Content    string
	Proxied    bool
	TTL        int
	Priority   *int
}

type CloudflareService interface {
	SetupDomainForMailStack(ctx context.Context, tenant, domain string) ([]string, error)
}

type cloudflareService struct {
	log      logger.Logger
	services *Services
	cfg      *config.Config
}

// NewCloudflareService initializes the CloudflareService
func NewCloudflareService(log logger.Logger, services *Services, cfg *config.Config) CloudflareService {
	return &cloudflareService{
		log:      log,
		services: services,
		cfg:      cfg,
	}
}

func (s *cloudflareService) SetupDomainForMailStack(ctx context.Context, tenant, domain string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.SetupDomainForMailStack")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	// step 1: Check if the domain exists in Cloudflare
	domainExists, zoneID, err := s.checkDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check domain existence"))
		s.log.Error("failed to check domain existence")
		return nil, err
	}

	// step 2: Delete all DNS records for the domain
	if domainExists {
		err = s.deleteAllDNSRecords(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to delete DNS records"))
			s.log.Error("failed to delete DNS records")
			return nil, err
		}
	}

	// step 3: Add the domain to Cloudflare
	if !domainExists {
		zoneID, err = s.addDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to add domain"))
			s.log.Errorf("failed to add domain %s", domain)
			return nil, err
		}
		domainExists = true
	}

	// step 4: Configure DNS records
	dnsConfigs := dnsConfigsForMailStack(domain)
	for _, dnsConfig := range dnsConfigs {
		err = s.addDNSRecord(ctx, zoneID, dnsConfig.RecordType, dnsConfig.Name, dnsConfig.Content, dnsConfig.TTL, dnsConfig.Proxied, dnsConfig.Priority)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to add DNS record"))
			s.log.Errorf("failed to add DNS record %s %s -> %s", dnsConfig.RecordType, dnsConfig.Name, dnsConfig.Content)
			return nil, err
		}
	}

	nameservers, err := s.getNameservers(ctx, zoneID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get nameservers"))
		s.log.Error("failed to get nameservers")
		return nil, err
	}

	return nameservers, nil
}

func (s *cloudflareService) deleteAllDNSRecords(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.deleteAllDNSRecords")
	defer span.Finish()
	span.LogKV("domain", domain)

	domainExists, zoneID, err := s.checkDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check domain existence"))
		s.log.Error("failed to check domain existence")
		return err
	}

	if !domainExists {
		span.LogFields(tracingLog.String("result", "Domain does not exist"))
		return nil
	}

	cloudflareUrl := fmt.Sprintf("%s/zones/%s/dns_records", s.cfg.ExternalServices.Cloudflare.Url, zoneID)
	req, err := http.NewRequestWithContext(ctx, "GET", cloudflareUrl, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		s.log.Error("failed to create request")
		return err
	}

	req.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get DNS records"))
		s.log.Error("failed to get DNS records")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body")
		return err
	}

	var recordsResponse struct {
		Success bool `json:"success"`
		Result  []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
	}

	if err = json.Unmarshal(body, &recordsResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		s.log.Error("failed to unmarshal response body")
		return err
	}

	if !recordsResponse.Success {
		err := fmt.Errorf("failed to fetch DNS records from Cloudflare")
		tracing.TraceErr(span, err)
		return err
	}

	// Step 3: Delete DNS records that match the given domain or its subdomains
	for _, record := range recordsResponse.Result {
		if record.Name == domain || strings.HasSuffix(record.Name, "."+domain) {
			delURL := fmt.Sprintf("%s/zones/%s/dns_records/%s", s.cfg.ExternalServices.Cloudflare.Url, zoneID, record.ID)
			deleteReq, err := http.NewRequestWithContext(ctx, "DELETE", delURL, nil)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to create delete request"))
				return err
			}

			deleteReq.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
			deleteReq.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
			deleteReq.Header.Set("Content-Type", "application/json")

			delResp, err := http.DefaultClient.Do(deleteReq)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to delete DNS record"))
				return err
			}
			defer delResp.Body.Close()

			// Check if the deletion was successful
			if delResp.StatusCode != http.StatusOK {
				delBody, _ := io.ReadAll(delResp.Body)
				err := fmt.Errorf("failed to delete DNS record: %s", string(delBody))
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}

// getZoneID fetches the zone ID for the given domain using the Cloudflare API
func (s *cloudflareService) checkDomain(ctx context.Context, domain string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.checkDomain")
	defer span.Finish()

	cloudflareUrl := fmt.Sprintf("%s/zones?name=%s", s.cfg.ExternalServices.Cloudflare.Url, domain)
	req, err := http.NewRequestWithContext(ctx, "GET", cloudflareUrl, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		s.log.Error("failed to create request")
		return false, "", err
	}

	req.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check domain existence in Cloudflare"))
		s.log.Error("failed to check domain existence in Cloudflare")
		return false, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return false, "", err
	}

	var zonesResponse struct {
		Success bool `json:"success"`
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
	}

	if err = json.Unmarshal(body, &zonesResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		s.log.Error("failed to unmarshal response body")
		return false, "", err
	}

	// Check if the domain exists
	if zonesResponse.Success && len(zonesResponse.Result) > 0 {
		span.LogFields(tracingLog.Bool("result.exists", true))
		return true, zonesResponse.Result[0].ID, nil
	}

	span.LogFields(tracingLog.Bool("result.exists", false))

	return false, "", nil
}

func (s *cloudflareService) addDomain(ctx context.Context, domain string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.addDomain")
	defer span.Finish()
	span.LogKV("domain", domain)

	cloudflareUrl := fmt.Sprintf("%s/zones", s.cfg.ExternalServices.Cloudflare.Url)

	// Create the request payload
	payload := map[string]interface{}{
		"name":       domain,
		"jump_start": true,   // Automatically scan for common DNS records
		"type":       "full", // "full" for full domain management
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal payload"))
		s.log.Error("failed to marshal payload", err)
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", cloudflareUrl, bytes.NewBuffer(payloadData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return "", err
	}

	req.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to add domain to Cloudflare"))
		s.log.Error("failed to add domain to Cloudflare", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return "", err
	}
	span.LogKV("responseBody", string(body))

	// Define the response structure
	var addDomainResponse struct {
		Success bool `json:"success"`
		Result  struct {
			ID string `json:"id"`
		} `json:"result"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err = json.Unmarshal(body, &addDomainResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		s.log.Error("failed to unmarshal response body", err)
		return "", err
	}

	if !addDomainResponse.Success {
		errMsg := "failed to add domain to Cloudflare"
		if len(addDomainResponse.Errors) > 0 {
			errMsg = addDomainResponse.Errors[0].Message
		}
		err := fmt.Errorf(errMsg)
		tracing.TraceErr(span, err)
		s.log.Error("Cloudflare API error: ", errMsg)
		return "", err
	}

	// Log success and return the Zone ID
	zoneID := addDomainResponse.Result.ID
	span.LogKV("zoneID", zoneID)
	s.log.Infof("Successfully added domain to Cloudflare. Zone ID: %s", zoneID)

	return zoneID, nil
}

func (s *cloudflareService) addDNSRecord(ctx context.Context, zoneID, recordType, name, content string, ttl int, proxied bool, priority *int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.addDNSRecord")
	defer span.Finish()
	span.LogKV("zoneID", zoneID, "recordType", recordType, "name", name)
	span.LogFields(tracingLog.String("content", content), tracingLog.Int("ttl", ttl), tracingLog.Bool("proxied", proxied))

	cloudflareUrl := fmt.Sprintf("%s/zones/%s/dns_records", s.cfg.ExternalServices.Cloudflare.Url, zoneID)

	// Create the request payload
	payload := map[string]interface{}{
		"type":    recordType,
		"name":    name,
		"content": content,
		"ttl":     ttl,
		"proxied": proxied,
	}
	// Include priority if the record type is MX
	if recordType == "MX" && priority != nil {
		payload["priority"] = *priority
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal payload"))
		s.log.Error("failed to marshal payload", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", cloudflareUrl, bytes.NewBuffer(payloadData))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return err
	}

	req.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to add DNS record to Cloudflare"))
		s.log.Error("failed to add DNS record to Cloudflare", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return err
	}

	// Define the response structure
	var addDNSResponse struct {
		Success bool `json:"success"`
		Errors  []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err = json.Unmarshal(body, &addDNSResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		s.log.Error("failed to unmarshal response body", err)
		return err
	}

	if !addDNSResponse.Success {
		errMsg := "failed to add DNS record to Cloudflare"
		if len(addDNSResponse.Errors) > 0 {
			errMsg = addDNSResponse.Errors[0].Message
		}
		err := fmt.Errorf(errMsg)
		tracing.TraceErr(span, err)
		s.log.Error("Cloudflare API error: ", errMsg)
		return err
	}

	// Log success
	s.log.Infof("Successfully added DNS record to Cloudflare: %s %s -> %s", recordType, name, content)

	return nil
}

func (s *cloudflareService) getNameservers(ctx context.Context, zoneID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.getNameservers")
	defer span.Finish()
	span.LogKV("zoneID", zoneID)

	cloudflareUrl := fmt.Sprintf("%s/zones/%s", s.cfg.ExternalServices.Cloudflare.Url, zoneID)

	req, err := http.NewRequestWithContext(ctx, "GET", cloudflareUrl, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create HTTP request"))
		s.log.Error("failed to create HTTP request", err)
		return nil, err
	}

	req.Header.Set("X-Auth-Email", s.cfg.ExternalServices.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", s.cfg.ExternalServices.Cloudflare.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get domain info from Cloudflare"))
		s.log.Error("failed to get domain info from Cloudflare", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		s.log.Error("failed to read response body", err)
		return nil, err
	}

	// Define the response structure
	var domainInfoResponse struct {
		Success bool `json:"success"`
		Result  struct {
			NameServers []string `json:"name_servers"`
		} `json:"result"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err = json.Unmarshal(body, &domainInfoResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal response body"))
		s.log.Error("failed to unmarshal response body", err)
		return nil, err
	}

	if !domainInfoResponse.Success {
		errMsg := "failed to get nameservers from Cloudflare"
		if len(domainInfoResponse.Errors) > 0 {
			errMsg = domainInfoResponse.Errors[0].Message
		}
		err = fmt.Errorf(errMsg)
		tracing.TraceErr(span, err)
		s.log.Error("Cloudflare API error: ", errMsg)
		return nil, err
	}

	// Log and return the nameservers
	nameservers := domainInfoResponse.Result.NameServers
	span.LogKV("result.nameservers", nameservers)
	s.log.Infof("Retrieved nameservers: %v", nameservers)

	return nameservers, nil
}

func dnsConfigsForMailStack(domain string) []DNSConfig {
	return []DNSConfig{
		{RecordType: "A", Name: "@", Content: "192.0.2.1", Proxied: true, TTL: 1},
		{RecordType: "CNAME", Name: "www", Content: domain, Proxied: true, TTL: 1},
		{RecordType: "CNAME", Name: "mail", Content: "mail.customerosmail.com", Proxied: false, TTL: 1},
		{RecordType: "MX", Name: "@", Content: fmt.Sprintf("mx.%s.cust.a.hostedemail.com", domain), Proxied: false, TTL: 1, Priority: utils.IntPtr(10)},
		{RecordType: "TXT", Name: "@", Content: "v=spf1 include:_spf.hostedemail.com -all", Proxied: false, TTL: 1},
		{RecordType: "TXT", Name: "_dmarc", Content: "v=DMARC1; p=reject; aspf=s; adkim=s; sp=reject; pct=100; ruf=mailto:dmarc@customerosmail.com; rua=mailto:monitor@customerosmail.com; fo=1; ri=86400", Proxied: false, TTL: 1},
		{RecordType: "TXT", Name: "dkim._domainkey", Content: "Copy the output value from the script dkim.py", Proxied: false, TTL: 1},
	}
}
