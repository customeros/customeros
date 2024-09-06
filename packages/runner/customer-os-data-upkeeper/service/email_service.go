package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ScrubbyIoResponse struct {
	Email      string `json:"email"`
	Status     string `json:"status"`
	Identifier string `json:"identifier"`
}

type EmailService interface {
	ValidateEmails()
	CheckScrubbyResult()
	CheckEnrowRequestsWithoutResponse()
}

type emailService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.Services
}

func (s *emailService) CheckEnrowRequestsWithoutResponse() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.checkBetterContactRequestsWithoutResponse")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	enrowRequestsWithoutResponse, err := s.commonServices.PostgresRepositories.CacheEmailEnrowRepository.GetWithoutResponses(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, record := range enrowRequestsWithoutResponse {
		// Create HTTP client
		client := &http.Client{}

		// Create POST request
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/email/verify/single?id=%s", s.cfg.EnrowConfig.ApiUrl, url.QueryEscape(record.RequestID)), nil)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			return
		}

		// Set headers
		req.Header.Set("x-api-key", s.cfg.EnrowConfig.ApiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		//Perform the request
		resp, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		if responseBody == nil || string(responseBody) == "" {
			continue
		}

		// Parse the JSON request body
		var enrowResponseBody postgresentity.EnrowResponseBody
		if err = json.Unmarshal(responseBody, &enrowResponseBody); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error unmarshalling request body"))
			continue
		}

		if enrowResponseBody.Qualification != "" {
			err = s.commonServices.PostgresRepositories.CacheEmailEnrowRepository.AddResponse(ctx, enrowResponseBody.Id, enrowResponseBody.Qualification, string(responseBody))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "error saving Enrow response to db"))
				return
			}
		} else {
			err = errors.New("Enrow response qualification is empty")
			tracing.TraceErr(span, err)
		}
	}
}

func NewEmailService(cfg *config.Config, log logger.Logger, commonServices *commonservice.Services) EmailService {
	return &emailService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *emailService) ValidateEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.ValidateEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := s.cfg.Limits.EmailsValidationLimit
	delayFromLastUpdateInMinutes := 2
	delayFromLastValidationAttemptInMinutes := 24 * 60

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.EmailReadRepository.GetEmailsForValidation(ctx, delayFromLastUpdateInMinutes, delayFromLastValidationAttemptInMinutes, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.EmailClient.RequestEmailValidation(ctx, &emailpb.RequestEmailValidationGrpcRequest{
					Tenant:    record.Tenant,
					Id:        record.EmailId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error requesting email validation"))
				s.log.Errorf("Error validating email {%s}: %s", record.EmailId, err.Error())
			}

			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelEmail, record.EmailId, string(neo4jentity.EmailPropertyValidationRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *emailService) CheckScrubbyResult() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.CheckScrubbyResult")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100
	delayFromPreviousCheckInHours := 12

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.PostgresRepositories.CacheEmailScrubbyRepository.GetToCheck(ctx, delayFromPreviousCheckInHours, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		internalCounter := 0
		limitChecksBeforePause := 9
		for _, record := range records {
			if internalCounter >= limitChecksBeforePause {
				internalCounter = 0
				time.Sleep(1 * time.Second)
			}
			internalCounter++
			scrubbyResult, err := s.callScrubbyIo(ctx, record.Email)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error calling scrubby.io"))
				s.log.Errorf("Error calling scrubby.io for email {%s}: %s", record.Email, err.Error())
			}
			if scrubbyResult.Status != string(postgresentity.ScrubbyStatusPending) {
				err = s.commonServices.PostgresRepositories.CacheEmailScrubbyRepository.SetStatus(ctx, record.Email, strings.ToLower(scrubbyResult.Status))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "Error setting scrubby status"))
					s.log.Errorf("Error setting scrubby status for email {%s}: %s", record.Email, err.Error())
				}
			} else {
				_, err = s.commonServices.PostgresRepositories.CacheEmailScrubbyRepository.SetJustChecked(ctx, record.ID)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "Error setting scrubby just checked"))
					s.log.Errorf("Error setting scrubby just checked for email {%s}: %s", record.Email, err.Error())
				}
			}
		}
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *emailService) callScrubbyIo(ctx context.Context, email string) (ScrubbyIoResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.callScrubbyIo")
	defer span.Finish()
	span.LogFields(log.String("email", email))

	encodedEmail := url.QueryEscape(email)
	req, err := http.NewRequest("GET", s.cfg.ScrubbyIoConfig.ApiUrl+"/fetch_email/"+encodedEmail, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return ScrubbyIoResponse{}, err
	}

	// Set the request headers
	req.Header.Set("x-api-key", s.cfg.ScrubbyIoConfig.ApiKey)
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
