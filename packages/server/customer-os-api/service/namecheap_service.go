package service

import (
	"encoding/xml"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
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

type NamecheapDomainInfo struct {
	DomainName  string   `json:"domainName"`
	CreatedDate string   `json:"createdDate"`
	ExpiredDate string   `json:"expiredDate"`
	Nameservers []string `json:"nameservers"`
	WhoisGuard  bool     `json:"whoisGuard"`
}

type NamecheapService interface {
	CheckDomainAvailability(ctx context.Context, domain string) (bool, bool, error)
	PurchaseDomain(ctx context.Context, tenant, domain string) error
	GetDomainPrice(ctx context.Context, domain string) (float64, error)
	GetDomainInfo(ctx context.Context, tenant, domain string) (NamecheapDomainInfo, error)
	UpdateNameservers(ctx context.Context, tenant, domain string, nameservers []string) error
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

	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.check")
	params.Add("DomainList", domain)

	resp, err := http.PostForm(s.cfg.ExternalServices.Namecheap.Url, params)
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
		if string(responseBody) == "error code: 522" {
			return false, false, coserrors.ErrConnectionTimeout
		}
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

	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.create")
	params.Add("DomainName", domain)
	params.Add("Years", strconv.Itoa(s.cfg.ExternalServices.Namecheap.Years))
	params.Add("AddFreeWhoisguard", "yes")

	params.Add("RegistrantFirstName", s.cfg.ExternalServices.Namecheap.RegistrantFirstName)
	params.Add("RegistrantLastName", s.cfg.ExternalServices.Namecheap.RegistrantLastName)
	params.Add("RegistrantJobTitle", s.cfg.ExternalServices.Namecheap.RegistrantJobTitle)
	params.Add("RegistrantAddress1", s.cfg.ExternalServices.Namecheap.RegistrantAddress1)
	params.Add("RegistrantOrganizationName", s.cfg.ExternalServices.Namecheap.RegistrantCompanyName)
	params.Add("RegistrantCity", s.cfg.ExternalServices.Namecheap.RegistrantCity)
	params.Add("RegistrantStateProvince", s.cfg.ExternalServices.Namecheap.RegistrantState)
	params.Add("RegistrantPostalCode", s.cfg.ExternalServices.Namecheap.RegistrantZIP)
	params.Add("RegistrantCountry", s.cfg.ExternalServices.Namecheap.RegistrantCountry)
	params.Add("RegistrantPhone", s.cfg.ExternalServices.Namecheap.RegistrantPhoneNumber)
	params.Add("RegistrantEmailAddress", s.cfg.ExternalServices.Namecheap.RegistrantEmail)

	params.Add("TechFirstName", s.cfg.ExternalServices.Namecheap.RegistrantFirstName)
	params.Add("TechLastName", s.cfg.ExternalServices.Namecheap.RegistrantLastName)
	params.Add("TechJobTitle", s.cfg.ExternalServices.Namecheap.RegistrantJobTitle)
	params.Add("TechAddress1", s.cfg.ExternalServices.Namecheap.RegistrantAddress1)
	params.Add("TechOrganizationName", s.cfg.ExternalServices.Namecheap.RegistrantCompanyName)
	params.Add("TechCity", s.cfg.ExternalServices.Namecheap.RegistrantCity)
	params.Add("TechStateProvince", s.cfg.ExternalServices.Namecheap.RegistrantState)
	params.Add("TechPostalCode", s.cfg.ExternalServices.Namecheap.RegistrantZIP)
	params.Add("TechCountry", s.cfg.ExternalServices.Namecheap.RegistrantCountry)
	params.Add("TechPhone", s.cfg.ExternalServices.Namecheap.RegistrantPhoneNumber)
	params.Add("TechEmailAddress", s.cfg.ExternalServices.Namecheap.RegistrantEmail)

	params.Add("AdminFirstName", s.cfg.ExternalServices.Namecheap.RegistrantFirstName)
	params.Add("AdminLastName", s.cfg.ExternalServices.Namecheap.RegistrantLastName)
	params.Add("AdminJobTitle", s.cfg.ExternalServices.Namecheap.RegistrantJobTitle)
	params.Add("AdminAddress1", s.cfg.ExternalServices.Namecheap.RegistrantAddress1)
	params.Add("AdminOrganizationName", s.cfg.ExternalServices.Namecheap.RegistrantCompanyName)
	params.Add("AdminCity", s.cfg.ExternalServices.Namecheap.RegistrantCity)
	params.Add("AdminStateProvince", s.cfg.ExternalServices.Namecheap.RegistrantState)
	params.Add("AdminPostalCode", s.cfg.ExternalServices.Namecheap.RegistrantZIP)
	params.Add("AdminCountry", s.cfg.ExternalServices.Namecheap.RegistrantCountry)
	params.Add("AdminPhone", s.cfg.ExternalServices.Namecheap.RegistrantPhoneNumber)
	params.Add("AdminEmailAddress", s.cfg.ExternalServices.Namecheap.RegistrantEmail)

	params.Add("AuxBillingFirstName", s.cfg.ExternalServices.Namecheap.RegistrantFirstName)
	params.Add("AuxBillingLastName", s.cfg.ExternalServices.Namecheap.RegistrantLastName)
	params.Add("AuxBillingJobTitle", s.cfg.ExternalServices.Namecheap.RegistrantJobTitle)
	params.Add("AuxBillingAddress1", s.cfg.ExternalServices.Namecheap.RegistrantAddress1)
	params.Add("AuxBillingOrganizationName", s.cfg.ExternalServices.Namecheap.RegistrantCompanyName)
	params.Add("AuxBillingCity", s.cfg.ExternalServices.Namecheap.RegistrantCity)
	params.Add("AuxBillingStateProvince", s.cfg.ExternalServices.Namecheap.RegistrantState)
	params.Add("AuxBillingPostalCode", s.cfg.ExternalServices.Namecheap.RegistrantZIP)
	params.Add("AuxBillingCountry", s.cfg.ExternalServices.Namecheap.RegistrantCountry)
	params.Add("AuxBillingPhone", s.cfg.ExternalServices.Namecheap.RegistrantPhoneNumber)
	params.Add("AuxBillingEmailAddress", s.cfg.ExternalServices.Namecheap.RegistrantEmail)

	// Execute the request
	resp, err := http.PostForm(s.cfg.ExternalServices.Namecheap.Url, params)
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
				Registered    bool   `xml:"Registered,attr"`
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
	if !result.CommandResponse.DomainCreateResult.Registered {
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

	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.users.getPricing")
	params.Add("ProductType", "DOMAIN")
	params.Add("ProductCategory", "REGISTER")
	params.Add("ProductName", tld)

	resp, err := http.PostForm(s.cfg.ExternalServices.Namecheap.Url, params)
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

func (s *namecheapService) GetDomainInfo(ctx context.Context, tenant, domain string) (NamecheapDomainInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NamecheapService.GetDomainInfo")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	// Check if domain belongs to the tenant in PostgreSQL and is active
	exists, err := s.repositories.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check domain ownership in postgres"))
		s.log.Error("failed to check domain ownership in postgres", err)
		return NamecheapDomainInfo{}, err
	}
	if !exists {
		err := fmt.Errorf("domain %s does not belong to tenant %s or is not active", domain, tenant)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return NamecheapDomainInfo{}, err
	}

	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.getInfo")
	params.Add("DomainName", domain)

	// Execute the request
	resp, err := http.PostForm(s.cfg.ExternalServices.Namecheap.Url, params)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call Namecheap API for domain info"))
		s.log.Error("failed to call Namecheap API for domain info", err)
		return NamecheapDomainInfo{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	span.LogFields(tracingLog.String("responseBody", string(responseBody)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read Namecheap response"))
		s.log.Error("failed to read Namecheap response", err)
		return NamecheapDomainInfo{}, err
	}

	// Define XML response structure for Namecheap domain info
	type NamecheapDomainInfoResult struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error []struct {
				Number  string `xml:"Number,attr"`
				Message string `xml:",chardata"`
			} `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			DomainGetInfoResult struct {
				Status        string `xml:"Status,attr"`
				ID            string `xml:"ID,attr"`
				DomainName    string `xml:"DomainName,attr"`
				OwnerName     string `xml:"OwnerName,attr"`
				IsOwner       bool   `xml:"IsOwner,attr"`
				IsPremium     bool   `xml:"IsPremium,attr"`
				DomainDetails struct {
					CreatedDate string `xml:"CreatedDate"`
					ExpiredDate string `xml:"ExpiredDate"`
					NumYears    int    `xml:"NumYears"`
				} `xml:"DomainDetails"`
				WhoisGuard struct {
					Enabled     bool   `xml:"Enabled,attr"`
					ID          string `xml:"ID"`
					ExpiredDate string `xml:"ExpiredDate"`
				} `xml:"Whoisguard"`
				PremiumDnsSubscription struct {
					UseAutoRenew   bool   `xml:"UseAutoRenew"`
					SubscriptionId string `xml:"SubscriptionId"`
					CreatedDate    string `xml:"CreatedDate"`
					ExpirationDate string `xml:"ExpirationDate"`
					IsActive       bool   `xml:"IsActive"`
				} `xml:"PremiumDnsSubscription"`
				DnsDetails struct {
					ProviderType     string   `xml:"ProviderType,attr"`
					IsUsingOurDNS    bool     `xml:"IsUsingOurDNS,attr"`
					HostCount        int      `xml:"HostCount,attr"`
					EmailType        string   `xml:"EmailType,attr"`
					DynamicDNSStatus bool     `xml:"DynamicDNSStatus,attr"`
					IsFailover       bool     `xml:"IsFailover,attr"`
					Nameservers      []string `xml:"Nameserver"`
				} `xml:"DnsDetails"`
			} `xml:"DomainGetInfoResult"`
		} `xml:"CommandResponse"`
	}

	var result NamecheapDomainInfoResult
	if err = xml.Unmarshal(responseBody, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse Namecheap XML response"))
		s.log.Error("failed to parse Namecheap XML response", err)
		return NamecheapDomainInfo{}, err
	}

	// Check if any errors exist
	if len(result.Errors.Error) > 0 {
		for _, e := range result.Errors.Error {
			errMsg := fmt.Sprintf("Error %s: %s", e.Number, e.Message)
			tracing.TraceErr(span, fmt.Errorf(errMsg))
			s.log.Errorf("Namecheap API returned an error: %s", errMsg)
		}
		return NamecheapDomainInfo{}, fmt.Errorf("Namecheap API returned errors")
	}

	// Populate NamecheapDomainInfo
	domainInfo := NamecheapDomainInfo{
		DomainName:  result.CommandResponse.DomainGetInfoResult.DomainName,
		CreatedDate: result.CommandResponse.DomainGetInfoResult.DomainDetails.CreatedDate,
		ExpiredDate: result.CommandResponse.DomainGetInfoResult.DomainDetails.ExpiredDate,
		Nameservers: result.CommandResponse.DomainGetInfoResult.DnsDetails.Nameservers,
		WhoisGuard:  result.CommandResponse.DomainGetInfoResult.WhoisGuard.Enabled,
	}

	// Log retrieved domain info
	span.LogKV("domainInfo", domainInfo)

	return domainInfo, nil
}

func (s *namecheapService) UpdateNameservers(ctx context.Context, tenant, domain string, nameservers []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NamecheapService.UpdateNameservers")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain, "nameservers", nameservers)

	// Check if domain belongs to the tenant in PostgreSQL and is active
	exists, err := s.repositories.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check domain ownership in postgres"))
		s.log.Error("failed to check domain ownership in postgres", err)
		return err
	}
	if !exists {
		err := fmt.Errorf("domain %s does not belong to tenant %s or is not active", domain, tenant)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return err
	}

	sld := strings.Split(domain, ".")[0]
	tld := strings.Split(domain, ".")[1]

	// Prepare the parameters for the Namecheap API call
	params := url.Values{}
	params.Add("ApiKey", s.cfg.ExternalServices.Namecheap.ApiKey)
	params.Add("ApiUser", s.cfg.ExternalServices.Namecheap.ApiUser)
	params.Add("UserName", s.cfg.ExternalServices.Namecheap.ApiUsername)
	params.Add("ClientIp", s.cfg.ExternalServices.Namecheap.ApiClientIp)
	params.Add("Command", "namecheap.domains.dns.setCustom")
	params.Add("SLD", sld)
	params.Add("TLD", tld)
	params.Add("Nameservers", strings.Join(nameservers, ","))

	// Execute the request
	resp, err := http.PostForm(s.cfg.ExternalServices.Namecheap.Url, params)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call Namecheap API for setting custom nameservers"))
		s.log.Error("failed to call Namecheap API for setting custom nameservers", err)
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

	// Define XML response structure for Namecheap setCustom DNS response
	type NamecheapSetCustomDNSResult struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error []struct {
				Number  string `xml:"Number,attr"`
				Message string `xml:",chardata"`
			} `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			DomainDNSSetCustomResult struct {
				Domain  string `xml:"Domain,attr"`
				Updated bool   `xml:"Updated,attr"`
			} `xml:"DomainDNSSetCustomResult"`
		} `xml:"CommandResponse"`
	}

	var result NamecheapSetCustomDNSResult
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

	// Check if the operation was successful
	if !result.CommandResponse.DomainDNSSetCustomResult.Updated {
		err := fmt.Errorf("failed to set custom nameservers for domain %s", domain)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return err
	}

	// Log success
	span.LogKV("result", "success")
	s.log.Infof("Successfully set custom nameservers for domain %s", domain)

	return nil
}
