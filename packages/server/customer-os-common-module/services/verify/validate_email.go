package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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

type MailsherpaRequest struct {
	Email string `json:"email"`
}
type MailsherpaResponse struct {
	Status          string                                  `json:"status"`
	Message         string                                  `json:"message"`
	InternalMessage string                                  `json:"internalMessage"`
	Data            *interfaces.ValidateEmailMailSherpaData `json:"data"`
}

func (s *verifyService) ValidateEmailWithMailSherpa(ctx context.Context, email string) (*interfaces.ValidateEmailMailSherpaData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithMailSherpa")
	defer span.Finish()
	span.LogKV("email", email)

	// check if mailsherpa is configured
	if s.cfg.Internal.MailSherpaApiConfig.MailsherpaApiUrl == "" || s.cfg.Internal.MailSherpaApiConfig.MailsherpaApiKey == "" {
		err := errors.New("MailSherpa is not configured")
		tracing.TraceErr(span, err)
		s.log.Errorf("MailSherpa is not configured")
		return nil, err
	}

	// Construct the URL with the email as a query parameter
	requestUrl := fmt.Sprintf("%s/validateEmail", s.cfg.Internal.MailSherpaApiConfig.MailsherpaApiUrl)

	request := MailsherpaRequest{
		Email: email,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewBuffer(payload))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}

	// Set the request headers
	req.Header.Set("x-Openline-API-KEY", s.cfg.Internal.MailSherpaApiConfig.MailsherpaApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	span.LogFields(log.Int("response.mailsherpa.status", response.StatusCode))
	body, err := io.ReadAll(response.Body)
	span.LogFields(log.String("response.mailsherpa.body", string(body)))

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("Mailsherpa returned %d status code", response.StatusCode)
		tracing.TraceErr(span, errors.Wrap(err, "failed to get response from Enrow"))
		return nil, err
	}

	var mailsherpaResponse MailsherpaResponse
	err = json.Unmarshal(body, &mailsherpaResponse)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode Mailsherpa response"))
		s.log.Errorf("failed to decode Mailsherpa response: %s", err.Error())
		return nil, err
	}

	if mailsherpaResponse.Data == nil {
		span.LogKV("MailsherpaResponse", "NoData")
	}

	return mailsherpaResponse.Data, nil
}

func (s *verifyService) ValidateEmailScrubby(ctx context.Context, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailScrubby")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedScrubbyRecord, err := s.postgres.CacheEmailScrubbyRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return "", err
	}

	validationStatus := ""

	if cachedScrubbyRecord == nil ||
		cachedScrubbyRecord.Status == "" ||
		cachedScrubbyRecord.CheckedAt.AddDate(0, 0, s.cfg.External.ScrubbyIoConfig.CacheTtlDays).Before(utils.Now()) {
		identifier := uuid.New().String()
		scrubbyResponse, err := s.callScrubbyIo(ctx, identifier, email)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to validate email with scrubby"))
		} else {
			savedRecord, err := s.postgres.CacheEmailScrubbyRepository.Save(ctx, postgresentity.CacheEmailScrubby{
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
		allCachedRecords, err := s.postgres.CacheEmailScrubbyRepository.GetAllByEmail(ctx, email)
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

func (s *verifyService) callScrubbyIo(ctx context.Context, identifier, email string) (ScrubbyIoResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callScrubbyIo")
	defer span.Finish()
	span.LogFields(log.String("email", email), log.String("identifier", identifier))

	// validate if scrubby is configured
	if s.cfg.External.ScrubbyIoConfig.ApiKey == "" || s.cfg.External.ScrubbyIoConfig.ApiUrl == "" {
		err := errors.New("scrubby.io is not configured")
		tracing.TraceErr(span, err)
		s.log.Errorf("scrubby.io is not configured")
		return ScrubbyIoResponse{}, err
	}

	requestJSON, err := json.Marshal(ScrubbyIoRequest{
		Email:       email,
		Identifier:  identifier,
		CallbackUrl: s.cfg.External.ScrubbyIoConfig.CallbackUrl,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return ScrubbyIoResponse{}, err
	}

	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", s.cfg.External.ScrubbyIoConfig.ApiUrl+"/add_email", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return ScrubbyIoResponse{}, err
	}

	// Set the request headers
	req.Header.Set("x-api-key", s.cfg.External.ScrubbyIoConfig.ApiKey)
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

func (s *verifyService) ValidateEmailWithTrueinbox(ctx context.Context, email string) (*postgresentity.TrueInboxResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailWithTrueinbox")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedTrueInboxRecord, err := s.postgres.CacheEmailTrueinboxRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return nil, err
	}

	var data *postgresentity.TrueInboxResponseBody
	if cachedTrueInboxRecord == nil || cachedTrueInboxRecord.CreatedAt.AddDate(0, 0, s.cfg.External.TrueInboxConfig.CacheTtlDays).Before(utils.Now()) {
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
		_, err = s.postgres.CacheEmailTrueinboxRepository.Create(ctx, postgresentity.CacheEmailTrueinbox{
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

func (s *verifyService) callTrueinboxToValidateEmail(ctx context.Context, email string) (postgresentity.TrueInboxResponseBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callTrueinboxToValidateEmail")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	// validate trueinbox is configured
	if s.cfg.External.TrueInboxConfig.ApiKey == "" || s.cfg.External.TrueInboxConfig.ApiUrl == "" {
		err := errors.New("TrueInbox is not configured")
		tracing.TraceErr(span, err)
		s.log.Errorf("TrueInbox is not configured")
		return postgresentity.TrueInboxResponseBody{}, err
	}

	// Construct the URL with the email as a query parameter
	requestUrl := fmt.Sprintf("%s/v1/api/verify-single-email?email=%s", s.cfg.External.TrueInboxConfig.ApiUrl, url.QueryEscape(email))

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return postgresentity.TrueInboxResponseBody{}, err
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+s.cfg.External.TrueInboxConfig.ApiKey)
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

func (s *verifyService) ValidateEmailWithEnrow(ctx context.Context, email string, extendedWaitingTimeForResponse bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.ValidateEmailEnrow")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	cachedEnrowRecord, err := s.postgres.CacheEmailEnrowRepository.GetLatestByEmail(ctx, email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get cache data"))
		return "", err
	}

	if cachedEnrowRecord == nil || cachedEnrowRecord.CreatedAt.AddDate(0, 0, s.cfg.External.EnrowConfig.CacheTtlDays).Before(utils.Now()) {
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
		_, err = s.postgres.CacheEmailEnrowRepository.RegisterRequest(ctx, postgresentity.CacheEmailEnrow{
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
	waitingTimeSec := s.cfg.External.EnrowConfig.MaxWaitResultsSeconds
	if extendedWaitingTimeForResponse {
		waitingTimeSec = waitingTimeSec * 3
	}

	for i := 0; i < waitingTimeSec; i++ {
		cachedEnrowRecord, err = s.postgres.CacheEmailEnrowRepository.GetLatestByEmail(ctx, email)
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

func (s *verifyService) callEnrowToValidateEmail(ctx context.Context, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.callEnrowToValidateEmail")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	// validate enrow is configured
	if s.cfg.External.EnrowConfig.ApiKey == "" || s.cfg.External.EnrowConfig.ApiUrl == "" {
		err := errors.New("Enrow is not configured")
		tracing.TraceErr(span, err)
		s.log.Errorf("Enrow is not configured")
		return "", err
	}

	// Construct the URL with the email as a query parameter
	requestUrl := fmt.Sprintf("%s/email/verify/single", s.cfg.External.EnrowConfig.ApiUrl)

	request := EnrowRequest{
		Email: email,
		Settings: struct {
			Webhook string `json:"webhook"`
		}{
			Webhook: s.cfg.External.EnrowConfig.CallbackUrl,
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
	req.Header.Set("x-api-key", s.cfg.External.EnrowConfig.ApiKey)
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
