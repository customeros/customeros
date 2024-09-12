package service

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type CloudflareService interface {
	SetupDomain(ctx context.Context, domain string) error
}

type cloudflareService struct {
	log      logger.Logger
	services *Services
	Api      *cloudflare.API
	ZoneID   string
}

func NewCloudflareService(log logger.Logger, services *Services) CloudflareService {
	return &cloudflareService{
		log:      log,
		services: services,
	}
}

//func (s *cloudflareService) getCloudflareAPI() (*cloudflare.API, error) {
//	api, err := cloudflare.New(apiKey, email)
//	if err != nil {
//
//		return nil, errors.Wrap(err, "failed to create Cloudflare API client")
//	}
//
//	return api, nil
//}

func (s *cloudflareService) SetupDomain(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CloudflareService.SetupDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Fetch zone ID for the domain
	zoneID, err := s.Api.ZoneIDByName(domain) // Assumes domain has no trailing dot
	if err != nil {
		return errors.Wrap(err, "failed to get Zone ID")
	}
	s.ZoneID = zoneID

	//// Add A Record
	//err = s.addARecord(domain)
	//if err != nil {
	//	return err
	//}
	//
	//// Add CNAME Records
	//err = s.addCNAMERecord("www", domain)
	//if err != nil {
	//	return err
	//}
	//
	//err = s.addCNAMERecord("mail", "mail.customerosmail.com")
	//if err != nil {
	//	return err
	//}
	//
	//err = s.addCNAMERecord("stats", "custosmetrics.com")
	//if err != nil {
	//	return err
	//}
	//
	//// Add MX Record
	//err = s.addMXRecord(domain, fmt.Sprintf("mx.%s.cust.a.hostedemail.com", domain))
	//if err != nil {
	//	return err
	//}
	//
	//// Add TXT Records for SPF, DMARC, DKIM
	//err = s.addTXTRecord(domain, "v=spf1 include:_spf.hostedemail.com -all")
	//if err != nil {
	//	return err
	//}
	//
	//err = s.addTXTRecord("_dmarc", "v=DMARC1; p=reject; aspf=s; adkim=s; sp=reject; pct=100; ruf=mailto:dmarc@customerosmail.com; rua=mailto:monitor@customerosmail.com; fo=1; ri=86400")
	//if err != nil {
	//	return err
	//}
	//
	//// Add DKIM Record (Example with a dummy key)
	//err = s.addTXTRecord("dkim._domainkey", "Your DKIM Value from dkim.py")
	//if err != nil {
	//	return err
	//}

	return nil
}

//// addARecord adds an A Record to Cloudflare
//func (s *cloudflareService) addARecord(domain string) error {
//	record := cloudflare.DNSRecord{
//		Type:    "A",
//		Name:    domain,
//		Content: "192.0.2.1",
//		TTL:     1, // TTL "automatic" in Cloudflare
//		Proxied: true,
//	}
//	_, err := s.Api.CreateDNSRecord(s.ZoneID, record)
//	if err != nil {
//		return errors.Wrap(err, "failed to create A record")
//	}
//	return nil
//}

//// addCNAMERecord adds a CNAME Record to Cloudflare
//func (s *cloudflareService) addCNAMERecord(name, target string) error {
//	record := cloudflare.DNSRecord{
//		Type:    "CNAME",
//		Name:    name,
//		Content: target,
//		TTL:     1,
//		Proxied: name == "www", // Proxy only "www"
//	}
//	_, err := s.Api.CreateDNSRecord(s.ZoneID, record)
//	if err != nil {
//		return errors.Wrap(err, fmt.Sprintf("failed to create CNAME record for %s", name))
//	}
//	return nil
//}
//
//// addMXRecord adds an MX Record to Cloudflare
//func (s *cloudflareService) addMXRecord(domain, mailServer string) error {
//	record := cloudflare.DNSRecord{
//		Type:     "MX",
//		Name:     domain,
//		Content:  mailServer,
//		TTL:      1,
//		Priority: 1,
//	}
//	_, err := s.Api.CreateDNSRecord(s.ZoneID, record)
//	if err != nil {
//		return errors.Wrap(err, "failed to create MX record")
//	}
//	return nil
//}
//
//// addTXTRecord adds a TXT Record to Cloudflare
//func (s *cloudflareService) addTXTRecord(name, content string) error {
//	record := cloudflare.DNSRecord{
//		Type:    "TXT",
//		Name:    name,
//		Content: content,
//		TTL:     1,
//	}
//	_, err := s.Api.CreateDNSRecord(s.ZoneID, record)
//	if err != nil {
//		return errors.Wrap(err, fmt.Sprintf("failed to create TXT record for %s", name))
//	}
//	return nil
//}
