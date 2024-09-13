package service

import (
	"encoding/xml"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	coserrors "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Namecheap supported commands: https://www.namecheap.com/support/api/methods/

type NamecheapService interface {
	CheckDomainAvailability(ctx context.Context, domain string) (bool, bool, error)
	PurchaseDomain(ctx context.Context, tenant, domain string) error
	GetDomainPrice(ctx context.Context, domain string) (float64, error)
}

type namecheapService struct {
	log          logger.Logger
	cfg          *config.Config
	repositories *repository.Repositories
}

func NewNamecheapService(log logger.Logger, cfg *config.Config, repositories *repository.Repositories) NamecheapService {
	return &namecheapService{
		log:          log,
		cfg:          cfg,
		repositories: repositories,
	}
}

// CheckDomainAvailability checks if the domain is available using Namecheap API
func (s *namecheapService) CheckDomainAvailability(ctx context.Context, domain string) (bool, bool, error) {
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
		return false, false, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	span.LogFields(tracingLog.String("responseBody", string(responseBody)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read Namecheap response"))
		s.log.Error("failed to read Namecheap response", err)
		return false, false, err
	}

	// Define namecheap XML struct for domain check
	type NamecheapCheckResult struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error []struct {
				Number  string `xml:"Number,attr"`
				Message string `xml:",chardata"`
			} `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			DomainCheckResult struct {
				Domain                   string `xml:"Domain,attr"`
				Available                bool   `xml:"Available,attr"`
				IsPremiumName            bool   `xml:"IsPremiumName,attr"`
				PremiumRegistrationPrice string `xml:"PremiumRegistrationPrice,attr"`
			} `xml:"DomainCheckResult"`
		} `xml:"CommandResponse"`
	}
	var result NamecheapCheckResult

	if err = xml.Unmarshal(responseBody, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse Namecheap XML response"))
		s.log.Error("failed to parse Namecheap XML response", err)
		return false, false, err
	}
	// Check if any errors exist
	if len(result.Errors.Error) > 0 {
		for _, e := range result.Errors.Error {
			errMsg := fmt.Sprintf("Error %s: %s", e.Number, e.Message)
			tracing.TraceErr(span, fmt.Errorf(errMsg))
			s.log.Errorf("Namecheap API returned an error: %s", errMsg)
		}
		return false, false, fmt.Errorf("Namecheap API returned errors")
	}

	// Check availability
	span.LogFields(tracingLog.Bool("result.available", result.CommandResponse.DomainCheckResult.Available))
	span.LogFields(tracingLog.Bool("result.premium", result.CommandResponse.DomainCheckResult.IsPremiumName))

	return result.CommandResponse.DomainCheckResult.Available, result.CommandResponse.DomainCheckResult.IsPremiumName, nil
}

// PurchaseDomain purchases the domain using Namecheap API
func (s *namecheapService) PurchaseDomain(ctx context.Context, tenant, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NamecheapService.PurchaseDomain")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	apiURL := "https://api.namecheap.com/xml.response"
	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.create")
	params.Add("DomainName", domain)
	params.Add("Years", strconv.Itoa(s.cfg.ExternalServices.Namecheap.Years))
	params.Add("AddFreeWhoisGuard", "yes")
	params.Add("AutoRenew", utils.BoolToString(s.cfg.ExternalServices.Namecheap.AutoRenew))

	// Execute the request
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call Namecheap API for domain purchase"))
		s.log.Error("failed to call Namecheap API for domain purchase", err)
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	span.LogFields(tracingLog.String("responseBody", string(responseBody)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read Namecheap response"))
		s.log.Error("failed to read Namecheap response", err)
		return err
	}

	// Define namecheap XML struct for domain registration result
	type NamecheapPurchaseResult struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error []struct {
				Number  string `xml:"Number,attr"`
				Message string `xml:",chardata"`
			} `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			DomainCreateResult struct {
				Domain        string `xml:"Domain,attr"`
				Success       bool   `xml:"Success,attr"`
				OrderID       string `xml:"OrderID,attr"`
				TransactionID string `xml:"TransactionID,attr"`
				ChargedAmount string `xml:"ChargedAmount,attr"`
			} `xml:"DomainCreateResult"`
		} `xml:"CommandResponse"`
	}
	var result NamecheapPurchaseResult

	if err = xml.Unmarshal(responseBody, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse Namecheap XML response"))
		s.log.Error("failed to parse Namecheap XML response", err)
		return err
	}
	// Check if any errors exist
	if len(result.Errors.Error) > 0 {
		for _, e := range result.Errors.Error {
			errMsg := fmt.Sprintf("Error %s: %s", e.Number, e.Message)
			tracing.TraceErr(span, fmt.Errorf(errMsg))
			s.log.Errorf("Namecheap API returned an error: %s", errMsg)
		}
		return fmt.Errorf("Namecheap API returned errors")
	}

	// Check if the purchase was successful
	if !result.CommandResponse.DomainCreateResult.Success {
		err = fmt.Errorf("failed to register domain %s: Namecheap API returned unsuccessful status", domain)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return err
	}

	// Log and store the purchase details
	span.LogFields(
		tracingLog.String("result.domain", result.CommandResponse.DomainCreateResult.Domain),
		tracingLog.String("result.orderID", result.CommandResponse.DomainCreateResult.OrderID),
		tracingLog.String("result.transactionID", result.CommandResponse.DomainCreateResult.TransactionID),
		tracingLog.String("result.chargedAmount", result.CommandResponse.DomainCreateResult.ChargedAmount),
	)
	s.log.Infof("Domain purchased successfully: %s, Order ID: %s, Transaction ID: %s",
		result.CommandResponse.DomainCreateResult.Domain,
		result.CommandResponse.DomainCreateResult.OrderID,
		result.CommandResponse.DomainCreateResult.TransactionID,
	)

	// Store domain
	_, err = s.repositories.PostgresRepositories.MailStackDomainRepository.RegisterDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store mailstack domain in postgres"))
		s.log.Error("failed to store domain in postgres", err)
		return nil
	}

	// Return purchase result details
	return nil
}

func (s *namecheapService) GetDomainPrice(ctx context.Context, domain string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NamecheapService.GetDomainPrice")
	defer span.Finish()
	span.LogKV("domain", domain)

	// Extract the TLD from the domain (e.g., "com" from "example.com")
	tld := strings.Split(domain, ".")[1]

	apiURL := "https://api.namecheap.com/xml.response"
	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.users.getPricing")
	params.Add("ProductType", "DOMAIN")
	params.Add("ProductCategory", "REGISTER")
	params.Add("ProductName", tld)

	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call Namecheap API for domain pricing"))
		s.log.Error("failed to call Namecheap API for domain pricing", err)
		return 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	span.LogFields(tracingLog.String("responseBody", string(responseBody)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read Namecheap response"))
		s.log.Error("failed to read Namecheap response", err)
		return 0, err
	}

	// Define the XML struct for domain pricing response
	type NamecheapPricingResult struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error []struct {
				Number  string `xml:"Number,attr"`
				Message string `xml:",chardata"`
			} `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			UserGetPricingResult struct {
				ProductType struct {
					Name            string `xml:"Name,attr"`
					ProductCategory []struct {
						Name    string `xml:"Name,attr"`
						Product []struct {
							Name  string `xml:"Name,attr"`
							Price []struct {
								Duration       string `xml:"Duration,attr"`
								DurationType   string `xml:"DurationType,attr"`
								Price          string `xml:"Price,attr"`
								PricingType    string `xml:"PricingType,attr"`
								YourPrice      string `xml:"YourPrice,attr"`
								AdditionalCost string `xml:"AdditionalCost,attr"`
								Currency       string `xml:"Currency,attr"`
							} `xml:"Price"`
						} `xml:"Product"`
					} `xml:"ProductCategory"`
				} `xml:"ProductType"`
			} `xml:"UserGetPricingResult"`
		} `xml:"CommandResponse"`
	}
	var result NamecheapPricingResult

	if err = xml.Unmarshal(responseBody, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse Namecheap XML response"))
		s.log.Error("failed to parse Namecheap XML response", err)
		return 0, err
	}
	// Check if any errors exist
	if len(result.Errors.Error) > 0 {
		for _, e := range result.Errors.Error {
			errMsg := fmt.Sprintf("Error %s: %s", e.Number, e.Message)
			tracing.TraceErr(span, fmt.Errorf(errMsg))
			s.log.Errorf("Namecheap API returned an error: %s", errMsg)
		}
		return 0, fmt.Errorf("Namecheap API returned errors")
	}

	// Search for the TLD pricing information
	for _, category := range result.CommandResponse.UserGetPricingResult.ProductType.ProductCategory {
		if category.Name == "register" {
			for _, product := range category.Product {
				if product.Name == tld {
					for _, price := range product.Price {
						if price.Duration == "1" && price.DurationType == "YEAR" {
							// Parse the price and return it
							parsedPrice, err := strconv.ParseFloat(price.YourPrice, 64)
							if err != nil {
								tracing.TraceErr(span, errors.Wrap(err, "failed to parse registration price"))
								s.log.Error("failed to parse registration price", err)
								return 0, err
							}
							span.LogKV("result.price", parsedPrice)
							return parsedPrice, nil
						}
					}
				}
			}
		}
	}

	return 0, coserrors.ErrDomainPriceNotFound
}
