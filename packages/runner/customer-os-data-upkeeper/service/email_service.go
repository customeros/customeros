package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
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
	commonServices *commonservice.CommonServices
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
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/email/verify/single?id=%s", s.cfg.Common.External.EnrowConfig.ApiUrl, url.QueryEscape(record.RequestID)), nil)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
			return
		}

		// Set headers
		req.Header.Set("x-api-key", s.cfg.Common.External.EnrowConfig.ApiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Perform the request
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

func NewEmailService(cfg *config.Config, log logger.Logger, commonServices *commonservice.CommonServices) EmailService {
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

	limit := s.cfg.App.Limits.EmailsValidationLimit
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

		// no record found, exit
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			err = s.commonServices.EmailService.RequestEmailValidation(innerCtx, record.EmailId)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error requesting email validation"))
				s.log.Errorf("Error publishing email validation request: %s", err.Error())
			}

			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelEmail, record.EmailId, string(neo4jentity.EmailPropertyValidationRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
			}

			// pause 1 second before next request
			time.Sleep(1 * time.Second)
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
	workers := s.cfg.App.Limits.BulkEmailsValidationThreads

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
					validationResult, err := s.commonServices.VerifyService.ValidateEmail(ctx, record.Email)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "Error validating email"))
						s.log.Errorf("Error validating email: %v", err)
						continue
					}

					data, err := json.Marshal(validationResult)
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
					if validationResult.EmailData.Deliverable != "unknown" && validationResult.EmailData.Deliverable != "" {
						billableEvent := postgresentity.BillableEventEmailVerifiedNotCatchAll
						if validationResult.DomainData.IsCatchAll {
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
					if strings.ToLower(validationResult.EmailData.Deliverable) == "true" {
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

			// log csv content size
			span.LogFields(log.Int("csvContent.size.request."+requestID, len(csvContent)))

			// Upload result file to S3
			basePath := fmt.Sprintf("/EMAIL_VALIDATION/BULK/%d", utils.Now().Year())

			fileDTO, err := s.commonServices.FileService.UploadSingleFileBytes(ctx, basePath, requestID, requestID+".csv", &csvContent, false)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "UploadSingleFileBytes"))
				continue
			}

			if fileDTO.ID == "" {
				tracing.TraceErr(span, errors.New("fileDTO.Id is empty"))
				continue
			}

			err = s.commonServices.PostgresRepositories.EmailValidationRequestBulkRepository.MarkRequestAsCompleted(ctx, requestID, fileDTO.ID)
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
	req, err := http.NewRequest("GET", s.cfg.Common.External.ScrubbyIoConfig.ApiUrl+"/fetch_email/"+encodedEmail, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return ScrubbyIoResponse{}, err
	}

	// Set the request headers
	req.Header.Set("x-api-key", s.cfg.Common.External.ScrubbyIoConfig.ApiKey)
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

func (s *emailService) generateBulkEmailValidationResponseCSVFileContent(ctx context.Context, requestId string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.generateBulkEmailValidationResponseCSVFileContent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("requestId", requestId)

	// Create an in-memory buffer to write the CSV content
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Write the CSV header
	header, _ := GenerateCSVRow(interfaces.ValidateEmailMailSherpaData{})
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
			var validationData interfaces.ValidateEmailMailSherpaData
			if record.Data == "" {
				tracing.TraceErr(span, fmt.Errorf("validation data is empty for email %s and requestId %s", record.Email, record.RequestID))
				continue
			}
			if err := json.Unmarshal([]byte(record.Data), &validationData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal validation data for email %s: %v", record.Email, err)
			}

			// Generate CSV row
			_, row := GenerateCSVRow(validationData)
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

func GenerateCSVRow(data interfaces.ValidateEmailMailSherpaData) (header []string, row []string) {
	// Define the header fields (without the excluded ones)
	header = []string{
		"Email", "SyntaxIsValid", "User", "Domain", "CleanEmail",
		"IsFirewalled", "Provider", "SecureGatewayProvider", "IsCatchAll", "CanConnectSMTP",
		"HasMXRecord", "HasSPFRecord", "TLSRequired", "IsPrimaryDomain", "PrimaryDomain",
		"SkippedValidation", "Deliverable", "IsMailboxFull", "IsRoleAccount", "IsSystemGenerated",
		"IsFreeAccount", "SmtpSuccess", "RetryValidation", "TLSRequired", "AlternateEmail",
	}

	// Create the corresponding row for the given data
	row = []string{
		data.Email,
		fmt.Sprintf("%v", data.Syntax.IsValid),
		data.Syntax.User,
		data.Syntax.Domain,
		data.Syntax.CleanEmail,
		fmt.Sprintf("%v", data.DomainData.IsFirewalled),
		data.DomainData.Provider,
		data.DomainData.SecureGatewayProvider,
		fmt.Sprintf("%v", data.DomainData.IsCatchAll),
		fmt.Sprintf("%v", data.DomainData.CanConnectSMTP),
		fmt.Sprintf("%v", data.DomainData.HasMXRecord),
		fmt.Sprintf("%v", data.DomainData.HasSPFRecord),
		fmt.Sprintf("%v", data.DomainData.TLSRequired),
		fmt.Sprintf("%v", data.DomainData.IsPrimaryDomain),
		data.DomainData.PrimaryDomain,
		fmt.Sprintf("%v", data.EmailData.SkippedValidation),
		data.EmailData.Deliverable,
		fmt.Sprintf("%v", data.EmailData.IsMailboxFull),
		fmt.Sprintf("%v", data.EmailData.IsRoleAccount),
		fmt.Sprintf("%v", data.EmailData.IsSystemGenerated),
		fmt.Sprintf("%v", data.EmailData.IsFreeAccount),
		fmt.Sprintf("%v", data.EmailData.SmtpSuccess),
		fmt.Sprintf("%v", data.EmailData.RetryValidation),
		fmt.Sprintf("%v", data.EmailData.TLSRequired),
		data.EmailData.AlternateEmail,
	}

	return header, row
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
	delayFromLastUpdateInHours := 24 // 24 hours

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
			innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    record.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})

			err = s.commonServices.EmailService.DeleteOrphanEmail(innerCtx, record.EmailId)
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
