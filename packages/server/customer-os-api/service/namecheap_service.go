package service

import (
	"encoding/xml"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/url"
)

// Define namecheap XML struct
type NamecheapResult struct {
	XMLName         xml.Name `xml:"ApiResponse"`
	Status          string   `xml:"Status,attr"`
	CommandResponse struct {
		DomainCheckResult struct {
			Domain                   string `xml:"Domain,attr"`
			Available                bool   `xml:"Available,attr"`
			IsPremiumName            bool   `xml:"IsPremiumName,attr"`
			PremiumRegistrationPrice string `xml:"PremiumRegistrationPrice,attr"`
		} `xml:"DomainCheckResult"`
	} `xml:"CommandResponse"`
}

type NamecheapService interface {
	CheckDomainAvailability(ctx context.Context, domain string) (bool, error)
	//PurchaseDomain(ctx context.Context, tenant, domain string) error
}

type namecheapService struct {
	log logger.Logger
	cfg *config.Config
}

func NewNamecheapService(log logger.Logger, cfg *config.Config) NamecheapService {
	return &namecheapService{
		log: log,
		cfg: cfg,
	}
}

// CheckDomainAvailability checks if the domain is available using Namecheap API
func (s *namecheapService) CheckDomainAvailability(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NamecheapService.CheckDomainAvailability")
	defer span.Finish()
	span.LogKV("domain", domain)

	apiURL := "https://api.namecheap.com/xml.response"
	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.check")
	params.Add("DomainList", domain)

	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call Namecheap API"))
		s.log.Error("failed to call Namecheap API", err)
		return false, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.LogObjectAsJson(span, "responseBody", responseBody)
		tracing.TraceErr(span, errors.Wrap(err, "failed to read Namecheap response"))
		s.log.Error("failed to read Namecheap response", err)
		return false, err
	}

	var result NamecheapResult

	if err = xml.Unmarshal(responseBody, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse Namecheap XML response"))
		s.log.Error("failed to parse Namecheap XML response", err)
		return false, err
	}

	// Check availability
	if result.CommandResponse.DomainCheckResult.Available {
		span.LogFields(tracingLog.Bool("result.available", true))
		return true, nil
	} else {
		span.LogFields(tracingLog.Bool("result.available", false))
		return false, nil
	}
}

//// PurchaseDomain purchases a domain using the Namecheap API and stores it in the DB
//func (s *namecheapService) PurchaseDomain(ctx context.Context, tenant, domain string) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "PurchaseDomain")
//	defer span.Finish()
//	tracing.TagTenant(span, tenant)
//	span.LogKV("domain", domain)
//
//	// Step 1: Check domain availability
//	isAvailable, err := s.CheckDomainAvailability(ctx, domain)
//	if err != nil {
//		return fmt.Errorf("failed to check domain availability: %w", err)
//	}
//
//	if !isAvailable {
//		return fmt.Errorf("domain %s is unavailable", domain)
//	}
//
//	// Step 2: Register domain
//	apiURL := "https://api.namecheap.com/xml.response"
//	params := url.Values{}
//	params.Add("ApiKey", s.apiKey)
//	params.Add("ApiUser", s.apiUser)
//	params.Add("UserName", s.apiUsername)
//	params.Add("ClientIp", s.clientIp)
//	params.Add("Command", "namecheap.domains.create")
//	params.Add("DomainName", domain)
//
//	resp, err := http.PostForm(apiURL, params)
//	if err != nil {
//		return fmt.Errorf("failed to call Namecheap API for domain registration: %w", err)
//	}
//	defer resp.Body.Close()
//
//	var result struct {
//		ApiResponse struct {
//			CommandResponse struct {
//				DomainCreateResult struct {
//					Domain  string `json:"Domain"`
//					Success bool   `json:"Success"`
//					Errors  string `json:"Errors"`
//				} `json:"DomainCreateResult"`
//			} `json:"CommandResponse"`
//		} `json:"ApiResponse"`
//	}
//
//	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
//		return fmt.Errorf("failed to parse Namecheap response: %w", err)
//	}
//
//	if !result.ApiResponse.CommandResponse.DomainCreateResult.Success {
//		return fmt.Errorf("failed to register domain: %s", result.ApiResponse.CommandResponse.DomainCreateResult.Errors)
//	}
//
//	// Step 3: Store domain in PostgreSQL after successful purchase
//	domainRecord := entity.MailStackDomain{
//		Tenant:        "some-tenant-id", // Assuming tenant is available in context or input
//		DomainName:    domain,
//		Configuration: "", // Placeholder for now
//		CreatedAt:     time.Now(),
//		UpdatedAt:     time.Now(),
//	}
//	if err := s.repo.CreateDomain(ctx, &domainRecord); err != nil {
//		return fmt.Errorf("failed to store domain record: %w", err)
//	}
//
//	span.LogKV("result", "Domain successfully purchased and stored")
//	return nil
//}