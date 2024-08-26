package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type ScrubbyIoRequest struct {
	Email       string `json:"email"`
	CallbackUrl string `json:"callback_url"`
	Identifier  string `json:"identifier"`
}

type ScrubbyIoResponse struct {
	Email      string `json:"email"`
	Status     string `json:"status"`
	Identifier string `json:"identifier"`
}

type EmailValidationService interface {
	ValidateEmailWithMailsherpa(ctx context.Context, email string) (*model.ValidateEmailMailSherpaData, error)
	ValidateEmailScrubby(ctx context.Context, email string) (string, error)
}

type emailValidationService struct {
	config   *config.Config
	Services *Services
	log      logger.Logger
}

func NewEmailValidationService(config *config.Config, services *Services, log logger.Logger) EmailValidationService {
	return &emailValidationService{
		config:   config,
		Services: services,
		log:      log,
	}
}

func (s *emailValidationService) ValidateEmailWithMailsherpa(ctx context.Context, email string) (*model.ValidateEmailMailSherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithMailsherpa")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	result := &model.ValidateEmailMailSherpaData{
		Email: email,
	}

	syntaxValidation := mailsherpa.ValidateEmailSyntax(email)
	result.Syntax.IsValid = syntaxValidation.IsValid
	result.Syntax.Domain = syntaxValidation.Domain
	result.Syntax.User = syntaxValidation.User
	result.Syntax.CleanEmail = syntaxValidation.CleanEmail

	// if syntax is not valid, return
	if !syntaxValidation.IsValid {
		return result, nil
	}

	domainValidation, err := s.getDomainValidation(ctx, syntaxValidation.Domain, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to validate email domain"))
		return nil, err
	}

	result.DomainData.IsFirewalled = domainValidation.IsFirewalled
	result.DomainData.IsCatchAll = domainValidation.IsCatchAll
	result.DomainData.CanConnectSMTP = domainValidation.CanConnectSMTP
	result.DomainData.Provider = domainValidation.Provider
	result.DomainData.Firewall = domainValidation.Firewall
	result.DomainData.HasMXRecord = domainValidation.HasMXRecord
	result.DomainData.HasSPFRecord = domainValidation.HasSPFRecord
	result.DomainData.TLSRequired = domainValidation.TLSRequired
	result.DomainData.ResponseCode = domainValidation.ResponseCode
	result.DomainData.ErrorCode = domainValidation.ErrorCode
	result.DomainData.Description = domainValidation.Description

	if domainValidation.HealthIsGreylisted || domainValidation.HealthIsBlacklisted {
		result.EmailData.Deliverable = string(model.EmailDeliverableStatusUnknown)
		return result, nil
	}

	var providersToSkip []string
	if s.config.EmailConfig.EmailValidationSkipProvidersCommaSeparated != "" {
		providersToSkip = strings.Split(s.config.EmailConfig.EmailValidationSkipProvidersCommaSeparated, ",")
		// remove spaces
		for i, provider := range providersToSkip {
			providersToSkip[i] = strings.TrimSpace(provider)
		}
	}

	// Check for providers that are marked for skip
	if len(providersToSkip) == 0 || !utils.Contains(providersToSkip, domainValidation.Provider) {
		emailValidation, err := s.getEmailValidation(ctx, email, syntaxValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
			return nil, err
		}
		result.EmailData.Deliverable = emailValidation.Deliverable
		result.EmailData.IsMailboxFull = emailValidation.IsMailboxFull
		result.EmailData.IsRoleAccount = emailValidation.IsRoleAccount
		result.EmailData.IsFreeAccount = emailValidation.IsFreeAccount
		result.EmailData.SmtpSuccess = emailValidation.SmtpSuccess
		result.EmailData.ResponseCode = emailValidation.ResponseCode
		result.EmailData.ErrorCode = emailValidation.ErrorCode
		result.EmailData.Description = emailValidation.Description
		result.EmailData.RetryValidation = emailValidation.RetryValidation
		result.EmailData.TLSRequired = emailValidation.TLSRequired
	} else {
		result.EmailData.SkippedValidation = true
		result.EmailData.Deliverable = string(model.EmailDeliverableStatusUnknown)
	}

	return result, nil
}

func (s *emailValidationService) getDomainValidation(ctx context.Context, domain, email string) (*postgresentity.CacheEmailValidationDomain, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getDomainValidation")
	defer span.Finish()

	cacheDomain, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationDomainRepository.Get(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	if cacheDomain == nil || cacheDomain.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailDomainValidationCacheTtlDays).Before(utils.Now()) {
		// get domain data with mailsherpa
		domainValidation := mailsherpa.ValidateDomain(mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.config.EmailConfig.EmailValidationFromDomain,
		})
		jsonData, err := json.Marshal(domainValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal domain validation data"))
		}
		cacheDomain, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationDomainRepository.Save(ctx, postgresentity.CacheEmailValidationDomain{
			Domain:              domain,
			Provider:            domainValidation.Provider,
			Firewall:            domainValidation.Firewall,
			IsCatchAll:          domainValidation.IsCatchAll,
			IsFirewalled:        domainValidation.IsFirewalled,
			HasMXRecord:         domainValidation.HasMXRecord,
			HasSPFRecord:        domainValidation.HasSPFRecord,
			Error:               domainValidation.Error,
			CanConnectSMTP:      domainValidation.SmtpResponse.CanConnectSMTP,
			TLSRequired:         domainValidation.SmtpResponse.TLSRequired,
			ResponseCode:        domainValidation.SmtpResponse.ResponseCode,
			ErrorCode:           domainValidation.SmtpResponse.ErrorCode,
			Description:         domainValidation.SmtpResponse.Description,
			HealthFromEmail:     domainValidation.MailServerHealth.FromEmail,
			HealthServerIP:      domainValidation.MailServerHealth.ServerIP,
			HealthIsGreylisted:  domainValidation.MailServerHealth.IsGreylisted,
			HealthIsBlacklisted: domainValidation.MailServerHealth.IsBlacklisted,
			HealthRetryAfter:    domainValidation.MailServerHealth.RetryAfter,
			Data:                string(jsonData),
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save domain data"))
			return nil, err
		}
	}

	return cacheDomain, nil
}

func (s *emailValidationService) getEmailValidation(ctx context.Context, email string, syntaxValidation mailsherpa.SyntaxValidation) (*postgresentity.CacheEmailValidation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getEmailValidation")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedEmail, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Get(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	// if no cached data found, or last time fetched > 90 days ago, or is retry validation
	if cachedEmail == nil ||
		cachedEmail.RetryValidation ||
		cachedEmail.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailValidationCacheTtlDays).Before(utils.Now()) {
		// get email data with mailsherpa
		emailValidation := mailsherpa.ValidateEmail(mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.config.EmailConfig.EmailValidationFromDomain,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email data with mailsherpa"))
			return nil, err
		}
		jsonData, err := json.Marshal(emailValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal email validation data"))
		}
		cachedEmail, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Save(ctx, postgresentity.CacheEmailValidation{
			Email:               email,
			Deliverable:         emailValidation.IsDeliverable,
			IsMailboxFull:       emailValidation.IsMailboxFull,
			IsRoleAccount:       emailValidation.IsRoleAccount,
			IsFreeAccount:       emailValidation.IsFreeAccount,
			RetryValidation:     emailValidation.RetryValidation,
			Error:               emailValidation.Error,
			Data:                string(jsonData),
			HealthIsGreylisted:  emailValidation.MailServerHealth.IsGreylisted,
			HealthIsBlacklisted: emailValidation.MailServerHealth.IsBlacklisted,
			HealthServerIP:      emailValidation.MailServerHealth.ServerIP,
			HealthFromEmail:     emailValidation.MailServerHealth.FromEmail,
			HealthRetryAfter:    emailValidation.MailServerHealth.RetryAfter,
			SmtpSuccess:         emailValidation.SmtpResponse.CanConnectSMTP,
			ResponseCode:        emailValidation.SmtpResponse.ResponseCode,
			ErrorCode:           emailValidation.SmtpResponse.ErrorCode,
			Description:         emailValidation.SmtpResponse.Description,
			TLSRequired:         emailValidation.SmtpResponse.TLSRequired,
			Username:            syntaxValidation.User,
			NormalizedEmail:     syntaxValidation.CleanEmail,
			Domain:              syntaxValidation.Domain,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save email data"))
			return nil, err
		}
	}

	return cachedEmail, nil
}

func (s *emailValidationService) ValidateEmailScrubby(ctx context.Context, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailScrubby")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedScrubbyRecord, err := s.Services.CommonServices.PostgresRepositories.CacheEmailScrubbyRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return "", err
	}

	validationStatus := ""

	if cachedScrubbyRecord == nil ||
		cachedScrubbyRecord.Status == "" ||
		cachedScrubbyRecord.CheckedAt.AddDate(0, 0, s.config.ScrubbyIoConfig.CacheTtlDays).Before(utils.Now()) {
		identifier := uuid.New().String()
		scrubbyResponse, err := s.callScrubbyIo(ctx, identifier, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email with scrubby"))
		} else {
			savedRecord, err := s.Services.CommonServices.PostgresRepositories.CacheEmailScrubbyRepository.Save(ctx, postgresentity.CacheEmailScrubby{
				ID:        identifier,
				Email:     email,
				Status:    strings.ToLower(scrubbyResponse.Status),
				CheckedAt: utils.Now(),
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to save scrubby data"))
				return "", err
			}
			validationStatus = savedRecord.Status
		}
	}

	if validationStatus == "" || validationStatus == "pending" {
		allCachedRecords, err := s.Services.CommonServices.PostgresRepositories.CacheEmailScrubbyRepository.GetAllByEmail(ctx, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get all scrubby records"))
			return validationStatus, err
		}
		for _, record := range allCachedRecords {
			if record.Status == "valid" {
				validationStatus = "valid"
				break
			} else if record.Status == "invalid" {
				validationStatus = "invalid"
				break
			} else if record.Status != "" {
				validationStatus = record.Status
			}
		}
	}
	return validationStatus, nil
}

func (s *emailValidationService) callScrubbyIo(ctx context.Context, identifier, email string) (ScrubbyIoResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callScrubbyIo")
	defer span.Finish()
	span.LogFields(log.String("email", email), log.String("identifier", identifier))

	requestJSON, err := json.Marshal(ScrubbyIoRequest{
		Email:       email,
		Identifier:  identifier,
		CallbackUrl: s.config.ScrubbyIoConfig.CallbackUrl,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return ScrubbyIoResponse{}, err
	}

	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", s.config.ScrubbyIoConfig.ApiUrl+"/add_email", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return ScrubbyIoResponse{}, err
	}

	// Set the request headers
	req.Header.Set("x-api-key", s.config.ScrubbyIoConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return ScrubbyIoResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("scrubby.io returned %d status code", response.StatusCode))
		tracing.TraceErr(span, err)
		return ScrubbyIoResponse{}, err
	}

	var scrubbyIoResponse ScrubbyIoResponse
	err = json.NewDecoder(response.Body).Decode(&scrubbyIoResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode scrubby.io response"))
		return ScrubbyIoResponse{}, err
	}
	tracing.LogObjectAsJson(span, "response.scrubby", scrubbyIoResponse)

	return scrubbyIoResponse, nil
}
