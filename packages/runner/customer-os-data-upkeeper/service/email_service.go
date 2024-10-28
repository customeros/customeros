package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	validationcsv "github.com/openline-ai/openline-customer-os/packages/server/validation-api/csv"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScrubbyIoResponse struct {
	Email      string `json:"email"`
	Status     string `json:"status"`
	Identifier string `json:"identifier"`
}

type EmailService interface {
	ValidateEmails()
	ValidateEmailsFromBulkRequests()
	CheckScrubbyResult()
	CheckEnrowRequestsWithoutResponse()
	CleanEmails()
	SendEmails()
	ProcessSentEmails()
}

type emailService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.Services
}

func (s *emailService) CheckEnrowRequestsWithoutResponse() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.CheckEnrowRequestsWithoutResponse")
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
	delayFromLastUpdateInSeconds := 10
	delayFromLastValidationAttemptInMinutes := 30 // 30 minutes

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.EmailReadRepository.GetEmailsForValidation(ctx, delayFromLastUpdateInSeconds, delayFromLastValidationAttemptInMinutes, limit)
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

func (s *emailService) ValidateEmailsFromBulkRequests() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.ValidateEmailsFromBulkRequests")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 200
	workers := s.cfg.Limits.BulkEmailsValidationThreads

	records, err := s.commonServices.PostgresRepositories.EmailValidationRecordRepository.GetUnprocessedEmailRecords(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	// if records found, process them
	if len(records) > 0 {
		// Create a worker pool
		recordChan := make(chan postgresentity.EmailValidationRecord, workers)
		wg := sync.WaitGroup{}

		// Start workers
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for record := range recordChan {
					// Call the email validation method (Placeholder)
					validationResult, err := s.callEmailValidation(ctx, record.Tenant, record.Email, record.VerifyCatchAll)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "Error calling email validation"))
						s.log.Errorf("Error validating email: %v", err)
						continue
					}

					dataObj := validationResult.Data
					if dataObj == nil {
						dataObj = &validationmodel.ValidateEmailMailSherpaData{}
					}
					data, err := json.Marshal(dataObj)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "Error marshalling data"))
						s.log.Errorf("Error marshalling data: %v", err)
						continue
					}

					// Update the record with the validation result
					err = s.commonServices.PostgresRepositories.EmailValidationRecordRepository.UpdateEmailRecord(ctx, record.ID, string(data))
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "Error updating email record"))
						s.log.Errorf("Failed to update email record: %s", err.Error())
						continue
					}

					// create billable event
					if dataObj.EmailData.Deliverable != "unknown" && dataObj.EmailData.Deliverable != "" {
						billableEvent := postgresentity.BillableEventEmailVerifiedNotCatchAll
						if dataObj.DomainData.IsCatchAll {
							billableEvent = postgresentity.BillableEventEmailVerifiedCatchAll
						}
						_, err = s.commonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, record.Tenant, billableEvent,
							postgresrepository.BillableEventDetails{
								ExternalID:    strconv.FormatUint(record.ID, 10),
								ReferenceData: record.Email,
							})
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
						}
					}

					// Update bulk request based on deliverable or undeliverable result
					if strings.ToLower(dataObj.EmailData.Deliverable) == "true" {
						err = s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.IncrementDeliverableEmails(ctx, record.RequestID)
					} else {
						err = s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.IncrementUndeliverableEmails(ctx, record.RequestID)
					}
					if err != nil {
						s.log.Errorf("Failed to increment email count for bulk request: %v", err)
					}
				}
			}()
		}

		// Feed the records to the workers
		for _, record := range records {
			recordChan <- record
		}

		close(recordChan)
		wg.Wait()
	}

	// After processing, check if all records for each request are processed
	requestsToCheck := make(map[string]struct{}) // Track unique Request IDs
	// Gather all unique Request IDs from the processed records
	for _, record := range records {
		requestsToCheck[record.RequestID] = struct{}{}
	}

	// additionally include oldest 5 uncompleted requests for check
	oldestUncompletedRequests, err := s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.GetOldestUncompletedRequests(ctx, 5)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting oldest uncompleted requests"))
	}
	for _, request := range oldestUncompletedRequests {
		requestsToCheck[request.RequestID] = struct{}{}
	}

	if len(requestsToCheck) > 0 {
		s.checkAndUpdateBulkRequests(ctx, requestsToCheck)
	}
}

// Check if all records for each request are processed and update bulk request status
func (s *emailService) checkAndUpdateBulkRequests(ctx context.Context, requestsToCheck map[string]struct{}) {
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.checkAndUpdateBulkRequests")
	defer span.Finish()

	// For each unique Request ID, check if all records are processed
	for requestID := range requestsToCheck {
		unprocessedCount, err := s.commonServices.PostgresRepositories.EmailValidationRecordRepository.CountPendingRequestsByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error counting pending records"))
			s.log.Errorf("Failed to count pending records for request %s: %v", requestID, err)
			continue
		}

		// get request by requestID
		request, err := s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error getting request by requestID"))
			s.log.Errorf("Failed to get request %s: %v", requestID, err)
			continue
		}

		// If there are no unprocessed records, mark the request as completed
		if unprocessedCount == 0 && request.Status != postgresentity.EmailValidationRequestBulkStatusCompleted {
			// generate csv result file
			csvContent, err := s.generateBulkEmailValidationResponseCSVFileContent(ctx, requestID)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error generating CSV content"))
				s.log.Errorf("Failed to generate CSV content for request %s: %v", requestID, err.Error())
				continue
			}

			// Upload result file to S3
			basePath := fmt.Sprintf("/EMAIL_VALIDATION/BULK/%d", utils.Now().Year())
			filesStoreService := fsc.NewFileStoreApiService(&s.cfg.FileStoreApiConfig)

			fileDTO, err := filesStoreService.UploadSingleFileBytes(request.Tenant, basePath, requestID, requestID+".csv", csvContent, span)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "UploadSingleFileBytes"))
				continue
			}

			if fileDTO.Id == "" {
				tracing.TraceErr(span, errors.New("fileDTO.Id is empty"))
				continue
			}

			err = s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.MarkRequestAsCompleted(ctx, requestID, fileDTO.Id)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error marking request as completed"))
				s.log.Errorf("Failed to mark request %s as completed: %v", requestID, err)
			}
		}
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

func (s *emailService) callEmailValidation(ctx context.Context, tenant, email string, verifyCatchAll bool) (validationmodel.ValidateEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.callEmailValidation")
	defer span.Finish()
	span.LogFields(log.String("email", email), log.Bool("verifyCatchAll", verifyCatchAll))
	tracing.TagTenant(span, tenant)

	emptyResponse := validationmodel.ValidateEmailResponse{}

	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequestWithOptions{
		Email: email,
		Options: validationmodel.ValidateEmailRequestOptions{
			VerifyCatchAll:      verifyCatchAll,
			ExtendedWaitingTime: true,
		},
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return emptyResponse, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", s.cfg.ValidationApi.Url+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return emptyResponse, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, s.cfg.ValidationApi.ApiKey)
	req.Header.Set(security.TenantHeader, tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return emptyResponse, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return emptyResponse, err
	}

	// if response status is 504 retry once
	if response.StatusCode == http.StatusGatewayTimeout {
		span.LogFields(log.Int("response.status.firstAttempt", response.StatusCode))
		response, err = client.Do(req)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
			return emptyResponse, err
		}
		defer response.Body.Close()
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
			return emptyResponse, err
		}
	}

	span.LogFields(log.Int("response.statusCode", response.StatusCode))

	if response.StatusCode == http.StatusGatewayTimeout {
		err = errors.New("validation api returned 504 status code")
		tracing.TraceErr(span, err)
		return emptyResponse, err
	}

	var validationResponse validationmodel.ValidateEmailResponse
	err = json.Unmarshal(responseBody, &validationResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(responseBody)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return emptyResponse, err
	}
	if validationResponse.Data == nil {
		tracing.LogObjectAsJson(span, "response", validationResponse)
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return emptyResponse, err
	}
	return validationResponse, nil
}

func (s *emailService) generateBulkEmailValidationResponseCSVFileContent(ctx context.Context, requestId string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.generateBulkEmailValidationResponseCSVFileContent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("requestId", requestId)

	// Create an in-memory buffer to write the CSV content
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Write the CSV header
	header, _ := validationcsv.GenerateCSVRow(validationmodel.ValidateEmailMailSherpaData{})
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write header: %v", err)
	}

	chunkSize := 1000
	offset := 0

	for {
		// Fetch a chunk of records
		records, err := s.commonServices.PostgresRepositories.EmailValidationRecordRepository.GetEmailRecordsInChunks(ctx, requestId, chunkSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch email records: %v", err)
		}

		// Break the loop if no more records are returned
		if len(records) == 0 {
			break
		}

		// Process each record and generate the CSV row
		for _, record := range records {
			// Parse the record data (assuming it's in JSON format) into ValidateEmailMailSherpaData
			var validationData validationmodel.ValidateEmailMailSherpaData
			if record.Data == "" {
				tracing.TraceErr(span, fmt.Errorf("validation data is empty for email %s and requestId %s", record.Email, record.RequestID))
				continue
			}
			if err := json.Unmarshal([]byte(record.Data), &validationData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal validation data for email %s: %v", record.Email, err)
			}

			// Generate CSV row
			_, row := validationcsv.GenerateCSVRow(validationData)
			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write row for email %s: %v", record.Email, err)
			}
		}

		// Flush the data to the buffer after processing each chunk
		writer.Flush()

		// Move to the next chunk
		offset += chunkSize
	}

	// Return the CSV content as a byte slice
	return buffer.Bytes(), nil
}

func (s *emailService) CleanEmails() {
	s.deleteOrphanEmails()
}

func (s *emailService) deleteOrphanEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.deleteOrphanEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	delayFromLastUpdateInHours := 7 * 24

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.EmailReadRepository.GetOrphanEmailNodes(ctx, limit, delayFromLastUpdateInHours)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error getting orphan emails"))
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			err = s.commonServices.EmailService.DeleteOrphanEmail(ctx, record.Tenant, record.EmailId, constants.AppSourceDataUpkeeper)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error deleting orphan email"))
				s.log.Errorf("Error deleting orphan email {%s}: %s", record.EmailId, err.Error())
			}
		}
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *emailService) SendEmails() {
	s.sendEmails()
}

func (s *emailService) sendEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.sendEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	emailMessages, err := s.commonServices.PostgresRepositories.EmailMessageRepository.GetForSending(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return // return if error
	}

	if len(emailMessages) == 0 {
		return
	}

	for _, emailMessage := range emailMessages {
		localCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    emailMessage.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		err := s.commonServices.MailService.SendMail(localCtx, emailMessage)
		if err != nil {
			tracing.TraceErr(span, err)

			s2 := err.Error()

			emailMessage.Status = postgresentity.EmailMessageStatusError
			emailMessage.Error = &s2

			err := s.commonServices.PostgresRepositories.EmailMessageRepository.Store(ctx, emailMessage.Tenant, emailMessage)
			if err != nil {
				tracing.TraceErr(span, err)
				break
			}

			continue
		}
	}
}

func (s *emailService) ProcessSentEmails() {
	s.processSentEmails()
}

func (s *emailService) processSentEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.processSentEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	emailMessages, err := s.commonServices.PostgresRepositories.EmailMessageRepository.GetForProcessing(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return // return if error
	}

	if len(emailMessages) == 0 {
		return
	}

	for _, emailMessage := range emailMessages {
		localCtx := common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    emailMessage.Tenant,
			AppSource: constants.AppSourceDataUpkeeper,
		})

		_, err := s.commonServices.MailService.ProcessSentEmail(localCtx, nil, emailMessage)
		if err != nil {
			tracing.TraceErr(span, err)

			s2 := err.Error()

			emailMessage.Status = postgresentity.EmailMessageStatusError
			emailMessage.Error = &s2

			err := s.commonServices.PostgresRepositories.EmailMessageRepository.Store(ctx, emailMessage.Tenant, emailMessage)
			if err != nil {
				tracing.TraceErr(span, err)
				break
			}

			continue
		}
	}
}
