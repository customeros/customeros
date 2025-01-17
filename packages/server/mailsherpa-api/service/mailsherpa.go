package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/mailsherpa-api/config"
	"github.com/opentracing/opentracing-go"
)

type MailSherpaService struct {
	log    logger.Logger
	config *config.Config
}

func NewMailSherpaService(log logger.Logger, config *config.Config) *MailSherpaService {
	return &MailSherpaService{
		log:    log,
		config: config,
	}
}

func (s *MailSherpaService) ValidateEmailWithMailSherpa(ctx context.Context, email string) (*interfaces.ValidateEmailMailSherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailsherpaAPI.ValidateEmailWithMailSherpa")
	defer span.Finish()
	span.LogKV("email", email)

	result := &interfaces.ValidateEmailMailSherpaData{
		Email: email,
	}

	syntaxValidation := mailsherpa.ValidateEmailSyntax(email)
	result.Syntax.IsValid = syntaxValidation.IsValid
	result.Syntax.Domain = syntaxValidation.Domain
	result.Syntax.User = syntaxValidation.User
	result.Syntax.CleanEmail = syntaxValidation.CleanEmail
	result.EmailData.IsFreeAccount = syntaxValidation.IsFreeAccount
	result.EmailData.IsRoleAccount = syntaxValidation.IsRoleAccount
	result.EmailData.IsSystemGenerated = syntaxValidation.IsSystemGenerated

	// if syntax is not valid, return
	if !syntaxValidation.IsValid {
		return result, nil
	}

	domainValidation, domainCheckTimeoutOccurred, err := s.getDomainValidationWithTimeout(ctx, syntaxValidation.Domain, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to validate email domain"))
		return nil, err
	}
	if domainCheckTimeoutOccurred {
		s.log.Warnf("Timeout occurred while validating domain %s", syntaxValidation.Domain)
		result.EmailData.Deliverable = string(EmailDeliverableStatusUnknown)
		return result, nil
	}

	result.DomainData.IsFirewalled = domainValidation.IsFirewalled
	result.DomainData.IsCatchAll = domainValidation.IsCatchAll
	result.DomainData.CanConnectSMTP = domainValidation.CanConnectSMTP
	result.DomainData.Provider = domainValidation.Provider
	result.DomainData.SecureGatewayProvider = domainValidation.Firewall
	result.DomainData.HasMXRecord = domainValidation.HasMXRecord
	result.DomainData.HasSPFRecord = domainValidation.HasSPFRecord
	result.DomainData.TLSRequired = domainValidation.TLSRequired
	result.DomainData.ResponseCode = domainValidation.ResponseCode
	result.DomainData.ErrorCode = domainValidation.ErrorCode
	result.DomainData.Description = domainValidation.Description
	result.DomainData.IsPrimaryDomain = *domainValidation.IsPrimaryDomain
	result.DomainData.PrimaryDomain = domainValidation.PrimaryDomain

	if domainValidation.HealthIsGreylisted || domainValidation.HealthIsBlacklisted || domainValidation.IsCatchAll {
		result.EmailData.Deliverable = string(EmailDeliverableStatusUnknown)
		return result, nil
	}

	var providersToSkip []string
	if s.cfg.InternalServices.EmailConfig.EmailValidationSkipProvidersCommaSeparated != "" {
		providersToSkip = strings.Split(s.cfg.InternalServices.EmailConfig.EmailValidationSkipProvidersCommaSeparated, ",")
		// remove spaces
		for i, provider := range providersToSkip {
			providersToSkip[i] = strings.TrimSpace(provider)
		}
	}

	// Check for providers that are marked for skip
	if len(providersToSkip) == 0 || !utils.Contains(providersToSkip, domainValidation.Provider) {
		emailValidation, emailCheckTimeoutOccurred, err := s.getEmailValidationWithTimeout(ctx, email, syntaxValidation, utils.BoolDefaultIfNil(domainValidation.IsPrimaryDomain, true), domainValidation.PrimaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
			return nil, err
		}
		if emailCheckTimeoutOccurred {
			s.log.Warnf("Timeout occurred while validating email %s", email)
			result.EmailData.Deliverable = string(EmailDeliverableStatusUnknown)
			return result, nil
		}
		alternateEmailValidation := postgresentity.CacheEmailValidation{}
		if emailValidation.AlternateEmail != "" {
			alternateEmailValidation, _, err = s.getEmailValidationWithTimeout(ctx, emailValidation.AlternateEmail, syntaxValidation, true, "")
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to validate alternate email"))
			}
		}
		result.EmailData.Deliverable = emailValidation.Deliverable
		result.EmailData.IsMailboxFull = emailValidation.IsMailboxFull
		result.EmailData.SmtpSuccess = emailValidation.SmtpSuccess
		result.EmailData.ResponseCode = emailValidation.ResponseCode
		result.EmailData.ErrorCode = emailValidation.ErrorCode
		result.EmailData.Description = emailValidation.Description
		result.EmailData.RetryValidation = emailValidation.RetryValidation
		result.EmailData.TLSRequired = emailValidation.TLSRequired
		if emailValidation.AlternateEmail != "" && alternateEmailValidation.Deliverable == string(EmailDeliverableStatusDeliverable) {
			result.EmailData.AlternateEmail = emailValidation.AlternateEmail
		}
	} else {
		result.EmailData.SkippedValidation = true
		result.EmailData.Deliverable = string(EmailDeliverableStatusUnknown)
	}

	return result, nil
}
