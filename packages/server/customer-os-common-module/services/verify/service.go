package verify

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/nyaruka/phonenumbers"
	"github.com/opentracing/opentracing-go"
)

type verifyService struct {
	log        logger.Logger
	postgres   *repository.Repositories
	cfg        *config.CommonConfig
	enrichment interfaces.EnrichmentService
}

func NewVerifyService(log logger.Logger, postgres *repository.Repositories, config *config.CommonConfig, enrichment interfaces.EnrichmentService) interfaces.VerifyService {
	return &verifyService{
		log:        log,
		postgres:   postgres,
		cfg:        config,
		enrichment: enrichment,
	}
}

func (s *verifyService) ValidateEmail(ctx context.Context, email string) (*interfaces.ValidateEmailMailSherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.ValidateEmail")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	if email == "" {
		err := errors.New("email is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.SetTag("email", email)

	// call mailsherpa
	emailValidationData, err := s.ValidateEmailWithMailSherpa(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if emailValidationData == nil {
		return nil, nil
	}

	// check if mailsherpa data is complete
	if s.isMailsherpaDataComplete(ctx, emailValidationData) {
		return emailValidationData, nil
	}

	// try Enrow
	enrowData, err := s.ValidateEmailWithEnrow(ctx, email, false)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if enrowData != "" {
		span.LogKV("enrowResponse", enrowData)
		switch enrowData {
		case "valid":
			emailValidationData.EmailData.Deliverable = string(EmailDeliverableStatusDeliverable)
			emailValidationData.EmailData.RetryValidation = false
			return emailValidationData, nil
		case "invalid":
			// do nothing
		default:
			err := errors.New("unexpected Enrow response: " + enrowData)
			tracing.TraceErr(span, err)
		}
	}

	// try TrueInbox
	trueInbox, err := s.ValidateEmailWithTrueinbox(ctx, email)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	if trueInbox == nil {
		return emailValidationData, nil
	}
	span.LogKV("trueInboxResponse", trueInbox.Result)
	switch trueInbox.Result {
	case "valid":
		emailValidationData.EmailData.Deliverable = string(EmailDeliverableStatusDeliverable)
		emailValidationData.EmailData.RetryValidation = false
		return emailValidationData, nil
	case "invalid":
		// do nothing
	default:
		err := errors.New("unexpected TrueInbox response: " + enrowData)
		tracing.TraceErr(span, err)
	}

	return emailValidationData, nil
}

func (s *verifyService) isMailsherpaDataComplete(ctx context.Context, data *interfaces.ValidateEmailMailSherpaData) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "VerifyService.isMailsherpaDataComplete")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	if data == nil {
		return true
	}

	if data.EmailData.Deliverable != string(EmailDeliverableStatusUnknown) {
		return true
	}

	if data.EmailData.IsRoleAccount ||
		data.EmailData.IsSystemGenerated ||
		data.EmailData.IsMailboxFull {
		return true
	}
	return false
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
	snitcherData, err := s.enrichment.IPIdentity(ctx, ipAddress)
	if err != nil {
		return nil, err
	}

	if !snitcherData.CompanyFound() {
		return nil, nil
	}

	domain := snitcherData.CompanyDomain()
	return &domain, nil
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
