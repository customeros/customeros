package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
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
	"time"
)

type EmailValidationService interface {
	ValidateEmailWithReacher(ctx context.Context, email string) (*model.RancherEmailResponseDTO, error)
	ValidateEmailWithMailsherpa(ctx context.Context, email string) (*model.ValidateEmailMailsherpaData, error)
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

func (s *emailValidationService) ValidateEmailWithReacher(ctx context.Context, email string) (*model.RancherEmailResponseDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.OnEmailCreate")
	defer span.Finish()

	message := map[string]string{"to_email": email}
	bytesRepresentation, _ := json.Marshal(message)

	client := http.Client{}
	// Create the request
	req, err := http.NewRequest("POST", s.config.ReacherConfig.ApiPath, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on creating request: %v", err.Error())
		return nil, err
	}
	req.Header.Set("x-reacher-secret", s.config.ReacherConfig.Secret)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on sending request: %v", err.Error())
		return nil, err
	}
	// Process the response
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on reading body: %v", err.Error())
		return nil, err
	}
	if resp.StatusCode == 200 {
		d := new(model.RancherEmailResponseDTO)

		err = json.Unmarshal(body, &d)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on unmarshalling body: %v", err.Error())
			return nil, err
		}
		return d, nil
	} else {
		return nil, errors.New(fmt.Sprintf("validation error: %s", body))
	}
}

func (s *emailValidationService) ValidateEmailWithMailsherpa(ctx context.Context, email string) (*model.ValidateEmailMailsherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithMailsherpa")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	result := &model.ValidateEmailMailsherpaData{
		Email: email,
	}

	syntaxValidation := mailsherpa.ValidateEmailSyntax(email)
	result.Syntax.IsValid = syntaxValidation.IsValid
	result.Syntax.Domain = syntaxValidation.Domain
	result.Syntax.User = syntaxValidation.User

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

	emailValidation, err := s.getEmailValidation(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to validate email"))
		return nil, err
	}
	result.EmailData.IsDeliverable = emailValidation.IsDeliverable
	result.EmailData.IsMailboxFull = emailValidation.IsMailboxFull
	result.EmailData.IsRoleAccount = emailValidation.IsRoleAccount
	result.EmailData.IsFreeAccount = emailValidation.IsFreeAccount
	result.EmailData.SmtpSuccess = emailValidation.SmtpSuccess
	result.EmailData.ResponseCode = emailValidation.ResponseCode
	result.EmailData.ErrorCode = emailValidation.ErrorCode
	result.EmailData.Description = emailValidation.Description
	result.EmailData.RetryValidation = emailValidation.RetryValidation
	result.EmailData.SmtpResponse = emailValidation.SmtpResponse

	return result, nil
}

func (s *emailValidationService) getDomainValidation(ctx context.Context, domain, email string) (*postgresentity.CacheEmailValidationDomain, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getDomainValidation")
	defer span.Finish()

	cacheDomain, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationDomainRepository.Get(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return nil, err
	}

	if cacheDomain == nil || cacheDomain.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailDomainValidationCacheTtlDays).Before(utils.Now()) {
		// get domain data with mailsherpa
		domainValidation, err := mailsherpa.ValidateDomain(mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.config.EmailConfig.EmailValidationFromDomain,
		}, true)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get domain data"))
			return nil, err
		}
		cacheDomain, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationDomainRepository.Save(ctx, postgresentity.CacheEmailValidationDomain{
			Domain:         domain,
			IsCatchAll:     domainValidation.IsCatchAll,
			IsFirewalled:   domainValidation.IsFirewalled,
			CanConnectSMTP: domainValidation.CanConnectSMTP,
			Provider:       domainValidation.Provider,
			Firewall:       domainValidation.Firewall,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save domain data"))
			return nil, err
		}
	}

	return cacheDomain, nil
}

func (s *emailValidationService) getEmailValidation(ctx context.Context, email string) (*postgresentity.CacheEmailValidation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.getEmailValidation")
	defer span.Finish()

	cacheEmail, err := s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Get(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return nil, err
	}

	// if no cached data found, or last time fetched > 90 days ago, or is retry validation and last time fetched > 1 hour ago
	if cacheEmail == nil ||
		cacheEmail.UpdatedAt.AddDate(0, 0, s.config.EmailConfig.EmailValidationCacheTtlDays).Before(utils.Now()) ||
		(cacheEmail.RetryValidation && cacheEmail.UpdatedAt.Add(time.Hour).Before(utils.Now())) {
		// get email data with mailsherpa
		emailValidation, err := mailsherpa.ValidateEmail(mailsherpa.EmailValidationRequest{
			Email:      email,
			FromDomain: s.config.EmailConfig.EmailValidationFromDomain,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get email data"))
			return nil, err
		}
		cacheEmail, err = s.Services.CommonServices.PostgresRepositories.CacheEmailValidationRepository.Save(ctx, postgresentity.CacheEmailValidation{
			Email:           email,
			IsDeliverable:   emailValidation.IsDeliverable,
			IsMailboxFull:   emailValidation.IsMailboxFull,
			IsRoleAccount:   emailValidation.IsRoleAccount,
			IsFreeAccount:   emailValidation.IsFreeAccount,
			SmtpSuccess:     emailValidation.SmtpSuccess,
			ResponseCode:    emailValidation.ResponseCode,
			ErrorCode:       emailValidation.ErrorCode,
			Description:     emailValidation.Description,
			RetryValidation: emailValidation.RetryValidation,
			SmtpResponse:    emailValidation.SmtpResponse,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save email data"))
			return nil, err
		}
	}

	return cacheEmail, nil
}
