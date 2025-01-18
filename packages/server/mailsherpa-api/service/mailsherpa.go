package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/verify"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/mailsherpa-api/config"
)

type MailSherpaService struct {
	log      logger.Logger
	cfg      *config.Config
	postgres *postgres_repository.Repositories
}

func NewMailSherpaService(log logger.Logger, config *config.Config, postgres *postgres_repository.Repositories) *MailSherpaService {
	return &MailSherpaService{
		log:      log,
		cfg:      config,
		postgres: postgres,
	}
}

const MaxDurationCallMailSherpa = 10 * time.Second

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
		result.EmailData.Deliverable = string(verify.EmailDeliverableStatusUnknown)
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
		result.EmailData.Deliverable = string(verify.EmailDeliverableStatusUnknown)
		return result, nil
	}

	var providersToSkip []string
	if s.cfg.EmailConfig.EmailValidationSkipProvidersCommaSeparated != "" {
		providersToSkip = strings.Split(s.cfg.EmailConfig.EmailValidationSkipProvidersCommaSeparated, ",")
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
			result.EmailData.Deliverable = string(verify.EmailDeliverableStatusUnknown)
			return result, nil
		}
		alternateEmailValidation := postgres_entity.CacheEmailValidation{}
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
		if emailValidation.AlternateEmail != "" && alternateEmailValidation.Deliverable == string(verify.EmailDeliverableStatusDeliverable) {
			result.EmailData.AlternateEmail = emailValidation.AlternateEmail
		}
	} else {
		result.EmailData.SkippedValidation = true
		result.EmailData.Deliverable = string(verify.EmailDeliverableStatusUnknown)
	}

	return result, nil
}

func (s *MailSherpaService) getDomainValidationWithTimeout(ctx context.Context, domain, email string) (postgres_entity.CacheEmailValidationDomain, bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, MaxDurationCallMailSherpa)
	defer cancel()

	// Create a new independent background context for validation
	validationCtx := context.Background()

	domainValidation := postgres_entity.CacheEmailValidationDomain{}
	errChan := make(chan error, 1)
	timeout := false

	// Run getEmailValidation in a goroutine to handle timeout
	go func() {
		var err error
		domainValidation, err = s.getDomainValidation(validationCtx, domain, email)
		errChan <- err
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return domainValidation, false, err
		}
	case <-ctxWithTimeout.Done():
		// Timeout occurred, set default values
		timeout = true
		domainValidation = postgres_entity.CacheEmailValidationDomain{}
	}
	return domainValidation, timeout, nil
}

func (s *MailSherpaService) getDomainValidation(ctx context.Context, domain, email string) (postgres_entity.CacheEmailValidationDomain, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getDomainValidation")
	defer span.Finish()
	span.LogKV("domain", domain, "email", email)

	cacheDomain, err := s.postgres.CacheEmailValidationDomainRepository.Get(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	if cacheDomain == nil || cacheDomain.IsPrimaryDomain == nil || cacheDomain.UpdatedAt.AddDate(0, 0, s.cfg.EmailConfig.EmailDomainValidationCacheTtlDays).Before(utils.Now()) {
		// get domain data with mailsherpa
		domainValidation := mailsherpa.ValidateDomain(mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.cfg.EmailConfig.EmailValidationFromDomain,
		})
		jsonData, err := json.Marshal(domainValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal domain validation data"))
		}
		cacheDomain, err = s.postgres.CacheEmailValidationDomainRepository.Save(ctx, postgres_entity.CacheEmailValidationDomain{
			Domain:              domain,
			Provider:            domainValidation.Provider,
			Firewall:            domainValidation.SecureGatewayProvider,
			IsCatchAll:          domainValidation.IsCatchAll,
			IsFirewalled:        domainValidation.IsFirewalled,
			HasMXRecord:         domainValidation.HasMXRecord,
			HasSPFRecord:        domainValidation.HasSPFRecord,
			Error:               domainValidation.Error,
			CanConnectSMTP:      domainValidation.SmtpResponse.CanConnectSMTP,
			TLSRequired:         domainValidation.SmtpResponse.TLSRequired,
			ResponseCode:        domainValidation.SmtpResponse.ResponseCode,
			ErrorCode:           domainValidation.SmtpResponse.ErrorCode,
			Description:         utils.SanitizeUTF8(domainValidation.SmtpResponse.Description),
			HealthFromEmail:     domainValidation.MailServerHealth.FromEmail,
			HealthServerIP:      domainValidation.MailServerHealth.ServerIP,
			HealthIsGreylisted:  domainValidation.MailServerHealth.IsGreylisted,
			HealthIsBlacklisted: domainValidation.MailServerHealth.IsBlacklisted,
			HealthRetryAfter:    domainValidation.MailServerHealth.RetryAfter,
			IsPrimaryDomain:     &domainValidation.IsPrimaryDomain,
			PrimaryDomain:       domainValidation.PrimaryDomain,
			Data:                utils.SanitizeUTF8(string(jsonData)),
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save domain data"))
			return postgres_entity.CacheEmailValidationDomain{}, err
		}
	}

	return *cacheDomain, nil
}

func (s *MailSherpaService) getEmailValidationWithTimeout(ctx context.Context, email string, syntaxValidation mailsherpa.SyntaxValidation, isPrimaryDomain bool, primaryDomain string) (postgres_entity.CacheEmailValidation, bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, MaxDurationCallMailSherpa)
	defer cancel()

	// Create a new independent background context for validation
	validationCtx := context.Background()

	emailValidation := postgres_entity.CacheEmailValidation{}
	errChan := make(chan error, 1)
	timeout := false

	// Run getEmailValidation in a goroutine to handle timeout
	go func() {
		var err error
		emailValidation, err = s.getEmailValidation(validationCtx, email, syntaxValidation, isPrimaryDomain, primaryDomain)
		errChan <- err
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return emailValidation, false, err
		}
	case <-ctxWithTimeout.Done():
		// Timeout occurred, set default values
		timeout = true
		emailValidation = postgres_entity.CacheEmailValidation{
			Deliverable: string(verify.EmailDeliverableStatusUndeliverable),
		}
	}
	return emailValidation, timeout, nil
}

func (s *MailSherpaService) getEmailValidation(ctx context.Context, email string, syntaxValidation mailsherpa.SyntaxValidation, isPrimaryDomain bool, primaryDomain string) (postgres_entity.CacheEmailValidation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getEmailValidation")
	defer span.Finish()
	span.LogFields(
		log.String("email", email),
		log.Bool("isPrimaryDomain", isPrimaryDomain),
		log.String("primaryDomain", primaryDomain))

	cachedEmail, err := s.postgres.CacheEmailValidationRepository.Get(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	// if no cached data found, or last time fetched > 90 days ago, or is retry validation
	if cachedEmail == nil ||
		cachedEmail.RetryValidation ||
		cachedEmail.UpdatedAt.AddDate(0, 0, s.cfg.EmailConfig.EmailValidationCacheTtlDays).Before(utils.Now()) {
		// get email data with mailsherpa
		emailValidationRequest := mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.cfg.EmailConfig.EmailValidationFromDomain,
			DomainValidationParams: &mailsherpa.DomainValidationParams{
				IsPrimaryDomain: isPrimaryDomain,
				PrimaryDomain:   primaryDomain,
			},
		}

		emailValidation := mailsherpa.ValidateEmail(emailValidationRequest)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email data with mailsherpa"))
			return postgres_entity.CacheEmailValidation{}, err
		}
		jsonData, err := json.Marshal(emailValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal email validation data"))
		}

		cacheEmailValidationEntity := postgres_entity.CacheEmailValidation{
			Email:               email,
			Deliverable:         emailValidation.IsDeliverable,
			IsMailboxFull:       emailValidation.IsMailboxFull,
			IsRoleAccount:       emailValidation.IsRoleAccount,
			IsFreeAccount:       emailValidation.IsFreeAccount,
			RetryValidation:     emailValidation.RetryValidation,
			Error:               emailValidation.Error,
			Data:                utils.SanitizeUTF8(string(jsonData)),
			HealthIsGreylisted:  emailValidation.MailServerHealth.IsGreylisted,
			HealthIsBlacklisted: emailValidation.MailServerHealth.IsBlacklisted,
			HealthServerIP:      emailValidation.MailServerHealth.ServerIP,
			HealthFromEmail:     emailValidation.MailServerHealth.FromEmail,
			HealthRetryAfter:    emailValidation.MailServerHealth.RetryAfter,
			SmtpSuccess:         emailValidation.SmtpResponse.CanConnectSMTP,
			ResponseCode:        emailValidation.SmtpResponse.ResponseCode,
			ErrorCode:           emailValidation.SmtpResponse.ErrorCode,
			Description:         utils.SanitizeUTF8(emailValidation.SmtpResponse.Description),
			TLSRequired:         emailValidation.SmtpResponse.TLSRequired,
			Username:            syntaxValidation.User,
			NormalizedEmail:     syntaxValidation.CleanEmail,
			Domain:              syntaxValidation.Domain,
			IsSystemGenerated:   syntaxValidation.IsSystemGenerated,
			AlternateEmail:      emailValidation.AlternateEmail.Email,
		}
		cachedEmail, err = s.postgres.CacheEmailValidationRepository.Save(ctx, cacheEmailValidationEntity)
		if err != nil {
			// retry saving once
			cachedEmail, err = s.postgres.CacheEmailValidationRepository.Save(ctx, cacheEmailValidationEntity)
		}
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save email data"))
			return postgres_entity.CacheEmailValidation{}, err
		}
	}

	return *cachedEmail, nil
}
