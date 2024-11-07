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
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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

type EnrowRequest struct {
	Email    string `json:"email"`
	Settings struct {
		Webhook string `json:"webhook"`
	} `json:"settings"`
}

type EnrowResponse struct {
	Id          string  `json:"id"`
	CreditsUsed float64 `json:"credits_used"`
	Message     string  `json:"message"`
}

const MaxDurationCallMailSherpa = 10 * time.Second

type EmailValidationService interface {
	ValidateEmailWithMailSherpa(ctx context.Context, email string) (*model.ValidateEmailMailSherpaData, error)
	ValidateEmailScrubby(ctx context.Context, email string) (string, error)
	ValidateEmailWithTrueinbox(ctx context.Context, email string) (*postgresentity.TrueInboxResponseBody, error)
	ValidateEmailWithEnrow(ctx context.Context, email string, extendedWaitingTimeForResponse bool) (string, error)
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

func (s *emailValidationService) ValidateEmailWithMailSherpa(ctx context.Context, email string) (*model.ValidateEmailMailSherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithMailSherpa")
	defer span.Finish()
	span.LogKV("email", email)

	result := &model.ValidateEmailMailSherpaData{
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
		result.EmailData.Deliverable = string(model.EmailDeliverableStatusUnknown)
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
		emailValidation, emailCheckTimeoutOccurred, err := s.getEmailValidationWithTimeout(ctx, email, syntaxValidation, utils.BoolDefaultIfNil(domainValidation.IsPrimaryDomain, true), domainValidation.PrimaryDomain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
			return nil, err
		}
		if emailCheckTimeoutOccurred {
			s.log.Warnf("Timeout occurred while validating email %s", email)
			result.EmailData.Deliverable = string(model.EmailDeliverableStatusUnknown)
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
		if emailValidation.AlternateEmail != "" && alternateEmailValidation.Deliverable == string(model.EmailDeliverableStatusDeliverable) {
			result.EmailData.AlternateEmail = emailValidation.AlternateEmail
		}
	} else {
		result.EmailData.SkippedValidation = true
		result.EmailData.Deliverable = string(model.EmailDeliverableStatusUnknown)
	}

	return result, nil
}

func (s *emailValidationService) getDomainValidationWithTimeout(ctx context.Context, domain, email string) (postgresentity.CacheEmailValidationDomain, bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, MaxDurationCallMailSherpa)
	defer cancel()

	// Create a new independent background context for validation
	validationCtx := context.Background()

	domainValidation := postgresentity.CacheEmailValidationDomain{}
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
		domainValidation = postgresentity.CacheEmailValidationDomain{}
	}
	return domainValidation, timeout, nil
}

func (s *emailValidationService) getDomainValidation(ctx context.Context, domain, email string) (postgresentity.CacheEmailValidationDomain, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getDomainValidation")
	defer span.Finish()
	span.LogKV("domain", domain, "email", email)

	cacheDomain, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationDomainRepository.Get(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	if cacheDomain == nil || cacheDomain.IsPrimaryDomain == nil || cacheDomain.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailDomainValidationCacheTtlDays).Before(utils.Now()) {
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
			return postgresentity.CacheEmailValidationDomain{}, err
		}
	}

	return *cacheDomain, nil
}

func (s *emailValidationService) getEmailValidationWithTimeout(ctx context.Context, email string, syntaxValidation mailsherpa.SyntaxValidation, isPrimaryDomain bool, primaryDomain string) (postgresentity.CacheEmailValidation, bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, MaxDurationCallMailSherpa)
	defer cancel()

	// Create a new independent background context for validation
	validationCtx := context.Background()

	emailValidation := postgresentity.CacheEmailValidation{}
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
		emailValidation = postgresentity.CacheEmailValidation{
			Deliverable: string(model.EmailDeliverableStatusUndeliverable),
		}
	}
	return emailValidation, timeout, nil
}

func (s *emailValidationService) getEmailValidation(ctx context.Context, email string, syntaxValidation mailsherpa.SyntaxValidation, isPrimaryDomain bool, primaryDomain string) (postgresentity.CacheEmailValidation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getEmailValidation")
	defer span.Finish()
	span.LogFields(
		log.String("email", email),
		log.Bool("isPrimaryDomain", isPrimaryDomain),
		log.String("primaryDomain", primaryDomain))

	cachedEmail, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Get(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
	}

	// if no cached data found, or last time fetched > 90 days ago, or is retry validation
	if cachedEmail == nil ||
		cachedEmail.RetryValidation ||
		cachedEmail.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailValidationCacheTtlDays).Before(utils.Now()) {
		// get email data with mailsherpa
		emailValidationRequest := mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.config.EmailConfig.EmailValidationFromDomain,
			DomainValidationParams: &mailsherpa.DomainValidationParams{
				IsPrimaryDomain: isPrimaryDomain,
				PrimaryDomain:   primaryDomain,
			},
		}

		emailValidation := mailsherpa.ValidateEmail(emailValidationRequest)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email data with mailsherpa"))
			return postgresentity.CacheEmailValidation{}, err
		}
		jsonData, err := json.Marshal(emailValidation)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal email validation data"))
		}

		cacheEmailValidationEntity := postgresentity.CacheEmailValidation{
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
		cachedEmail, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Save(ctx, cacheEmailValidationEntity)
		if err != nil {
			// retry saving once
			cachedEmail, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Save(ctx, cacheEmailValidationEntity)
		}
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save email data"))
			return postgresentity.CacheEmailValidation{}, err
		}
	}

	return *cachedEmail, nil
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

func (s *emailValidationService) ValidateEmailWithTrueinbox(ctx context.Context, email string) (*postgresentity.TrueInboxResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithTrueinbox")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedTrueInboxRecord, err := s.Services.CommonServices.PostgresRepositories.CacheEmailTrueinboxRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return nil, err
	}

	var data *postgresentity.TrueInboxResponseBody
	if cachedTrueInboxRecord == nil || cachedTrueInboxRecord.CreatedAt.AddDate(0, 0, s.config.TrueinboxConfig.CacheTtlDays).Before(utils.Now()) {
		trueInboxResponse, err := s.callTrueinboxToValidateEmail(ctx, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email with trueinbox"))
			s.log.Errorf("failed to validate email with trueinbox: %s", err.Error())
			return nil, err
		}
		responseJson, err := json.Marshal(trueInboxResponse)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to marshal trueinbox response"))
			s.log.Errorf("failed to marshal trueinbox response: %s", err.Error())
			return nil, err
		}
		_, err = s.Services.CommonServices.PostgresRepositories.CacheEmailTrueinboxRepository.Create(ctx, postgresentity.CacheEmailTrueinbox{
			Email:  email,
			Data:   string(responseJson),
			Result: trueInboxResponse.Result,
		})
		data = &trueInboxResponse
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save trueinbox data"))
			s.log.Errorf("failed to save trueinbox data: %s", err.Error())
			return nil, err
		}
	} else {
		data = &postgresentity.TrueInboxResponseBody{}
		err = json.Unmarshal([]byte(cachedTrueInboxRecord.Data), data)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal trueinbox data"))
			s.log.Errorf("failed to unmarshal trueinbox data: %s", err.Error())
			return nil, err
		}
	}
	return data, nil
}

func (s *emailValidationService) callTrueinboxToValidateEmail(ctx context.Context, email string) (postgresentity.TrueInboxResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callTrueinboxToValidateEmail")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	// Construct the URL with the email as a query parameter
	requestUrl := fmt.Sprintf("%s/v1/api/verify-single-email?email=%s", s.config.TrueinboxConfig.ApiUrl, url.QueryEscape(email))

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return postgresentity.TrueInboxResponseBody{}, err
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+s.config.TrueinboxConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return postgresentity.TrueInboxResponseBody{}, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.statusCode", response.StatusCode))
	body, err := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		span.LogFields(log.String("response.body", string(body)))
		err = fmt.Errorf("TrueInbox returned %d status code", response.StatusCode)
		tracing.TraceErr(span, errors.Wrap(err, "failed to get response from TrueInbox"))
		return postgresentity.TrueInboxResponseBody{}, err
	}

	var trueInboxResponse postgresentity.TrueInboxResponseBody
	err = json.Unmarshal(body, &trueInboxResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode TrueInbox response"))
		s.log.Errorf("failed to decode TrueInbox response: %s", err.Error())
		return trueInboxResponse, err
	}
	tracing.LogObjectAsJson(span, "response.trueinbox", trueInboxResponse)

	return trueInboxResponse, nil
}

func (s *emailValidationService) ValidateEmailWithEnrow(ctx context.Context, email string, extendedWaitingTimeForResponse bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailEnrow")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedEnrowRecord, err := s.Services.CommonServices.PostgresRepositories.CacheEmailEnrowRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return "", err
	}

	if cachedEnrowRecord == nil || cachedEnrowRecord.CreatedAt.AddDate(0, 0, s.config.EnrowConfig.CacheTtlDays).Before(utils.Now()) {
		enrowRequestId, err := s.callEnrowToValidateEmail(ctx, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to call enrow"))
			s.log.Errorf("failed to call enrow: %s", err.Error())
			return "", err
		}
		if enrowRequestId == "" {
			err = errors.New("enrow request id is empty")
			tracing.TraceErr(span, err)
			s.log.Errorf("enrow request id is empty")
			return "", err
		}
		_, err = s.Services.CommonServices.PostgresRepositories.CacheEmailEnrowRepository.RegisterRequest(ctx, postgresentity.CacheEmailEnrow{
			Email:     email,
			RequestID: enrowRequestId,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to register enrow request"))
			s.log.Errorf("failed to register enrow request: %s", err.Error())
			return "", err
		}
	}

	result := ""
	waitingTimeSec := s.config.EnrowConfig.MaxWaitResultsSeconds
	if extendedWaitingTimeForResponse {
		waitingTimeSec = waitingTimeSec * 3
	}

	for i := 0; i < waitingTimeSec; i++ {
		cachedEnrowRecord, err = s.Services.CommonServices.PostgresRepositories.CacheEmailEnrowRepository.GetLatestByEmail(ctx, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
			return "", err
		}
		if cachedEnrowRecord.Qualification != "" {
			result = cachedEnrowRecord.Qualification
			break
		}
		time.Sleep(time.Second)
	}

	return result, nil
}

func (s *emailValidationService) callEnrowToValidateEmail(ctx context.Context, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callEnrowToValidateEmail")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	// Construct the URL with the email as a query parameter
	requestUrl := fmt.Sprintf("%s/email/verify/single", s.config.EnrowConfig.ApiUrl)

	request := EnrowRequest{
		Email: email,
		Settings: struct {
			Webhook string `json:"webhook"`
		}{
			Webhook: s.config.EnrowConfig.CallbackUrl,
		},
	}
	payload, err := json.Marshal(request)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return "", err
	}

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewBuffer(payload))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return "", err
	}

	// Set the request headers
	req.Header.Set("x-api-key", s.config.EnrowConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return "", err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.enrow.status", response.StatusCode))
	body, err := io.ReadAll(response.Body)
	span.LogFields(log.String("response.enrow.body", string(body)))

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("Enrow returned %d status code", response.StatusCode)
		tracing.TraceErr(span, errors.Wrap(err, "failed to get response from Enrow"))
		return "", err
	}

	var enrowResponse EnrowResponse
	err = json.Unmarshal(body, &enrowResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode Enrow response"))
		s.log.Errorf("failed to decode Enrow response: %s", err.Error())
		return "", err
	}

	return enrowResponse.Id, nil
}
