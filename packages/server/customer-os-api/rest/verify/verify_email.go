// @openapi 3.0.0
package restverify

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

const (
	singleEmailVerificationAproxDurationInSeconds = float64(3.5)
	threadsToVerifyBulkEmails                     = 6
)

// Email verification responses and types

// EmailVerificationResponse represents the email verification response
// @Description Response for single email verification including detailed validation results
type EmailVerificationResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// Email verification details
	// required: true
	Email EmailVerificationRecord `json:"email,omitempty"`
}

// EmailVerificationRecord represents detailed email verification results
// @Description Detailed validation results for an email address
type EmailVerificationRecord struct {
	// Email address that was verified
	// required: true
	// format: email
	EmailAddress string `json:"emailAddress" example:"example@example.com"`

	// Deliverability status
	// required: true
	// enum: true,false,unknown
	Deliverable string `json:"deliverable" example:"true"`

	// Email service provider
	// required: false
	Provider string `json:"provider" example:"gmail"`

	// Security gateway provider
	// required: false
	SecureGatewayProvider string `json:"secureGatewayProvider" example:"Proofpoint"`

	// Indicates if email is considered risky
	// required: true
	IsRisky bool `json:"isRisky" example:"false"`

	// Indicates if domain is catch-all
	// required: true
	IsCatchAll bool `json:"isCatchAll" example:"false"`

	// Risk assessment details
	// required: true
	Risk EmailVerificationRisk `json:"risk"`

	// Syntax validation details
	// required: true
	Syntax EmailVerificationSyntax `json:"syntax"`

	// Alternative email address if available
	// required: false
	// format: email
	AlternateEmail string `json:"alternateEmail,omitempty" example:"alternate@example.com"`
}

// EmailVerificationRisk represents risk assessment details
// @Description Risk factors associated with the email address
type EmailVerificationRisk struct {
	// Indicates if email is behind a firewall
	// required: true
	IsFirewalled bool `json:"isFirewalled" example:"false"`

	// Indicates if email is a role account
	// required: true
	IsRoleMailbox bool `json:"isRoleMailbox" example:"false"`

	// Indicates if email is system-generated
	// required: true
	IsSystemGenerated bool `json:"isSystemGenerated" example:"false"`

	// Indicates if email uses a free provider
	// required: true
	IsFreeProvider bool `json:"isFreeProvider" example:"true"`

	// Indicates if mailbox is full
	// required: true
	IsMailboxFull bool `json:"isMailboxFull" example:"false"`

	// Indicates if domain is primary
	// required: true
	IsPrimaryDomain bool `json:"isPrimaryDomain" example:"true"`
}

// EmailVerificationSyntax represents syntax validation results
// @Description Email syntax validation details
type EmailVerificationSyntax struct {
	// Indicates if email syntax is valid
	// required: true
	IsValid bool `json:"isValid" example:"true"`

	// Domain part of email
	// required: true
	Domain string `json:"domain" example:"example.com"`

	// Local part of email
	// required: true
	User string `json:"user" example:"example"`
}

// Bulk verification responses

// BulkUploadResponse represents bulk verification job initiation response
// @Description Response after initiating bulk email verification
type BulkUploadResponse struct {
	// Status message
	// required: true
	Message string `json:"message" example:"File uploaded successfully"`

	// Unique job identifier
	// required: true
	// format: uuid
	JobID string `json:"jobId" example:"550e8400-e29b-41d4-a716-446655440000"`

	// URL to check verification results
	// required: true
	// format: uri
	ResultURL string `json:"resultUrl" example:"https://api.customeros.ai/verify/v1/email/bulk/results/550e8400-e29b-41d4-a716-446655440000"`

	// Estimated completion timestamp
	// required: true
	EstimatedCompletionTs float64 `json:"estimatedCompletionTs" example:"1694030400"`
}

// BulkResultsResponse represents bulk verification results
// @Description Response containing bulk verification results or status
type BulkResultsResponse struct {
	// Unique job identifier
	// required: true
	// format: uuid
	JobID string `json:"jobId" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Processing status
	// required: true
	// enum: processing,completed
	Status string `json:"status" example:"completed"`

	// Original filename
	// required: true
	FileName string `json:"fileName" example:"emails.csv"`

	// Progress message
	// required: true
	Message string `json:"message" example:"Completed 1000 of 1000 emails"`

	// Verification results if completed
	// required: false
	Results *BulkResultsDetails `json:"results,omitempty"`

	// Estimated completion timestamp
	// required: true
	EstimatedCompletionTs int64 `json:"estimatedCompletionTs" example:"1694030400"`
}

// BulkResultsDetails represents detailed bulk verification results
// @Description Detailed statistics for bulk verification results
type BulkResultsDetails struct {
	// Total number of emails processed
	// required: true
	// minimum: 0
	TotalEmails int `json:"totalEmails" example:"1000"`

	// Number of deliverable emails
	// required: true
	// minimum: 0
	Deliverable int `json:"deliverable" example:"950"`

	// Number of undeliverable emails
	// required: true
	// minimum: 0
	Undeliverable int `json:"undeliverable" example:"45"`

	// URL to download detailed results
	// required: true
	// format: uri
	DownloadURL string `json:"downloadUrl" example:"https://api.customeros.ai/verify/v1/email/bulk/results/550e8400-e29b-41d4-a716-446655440000/download"`
}

// @Summary Verify single email address
// @Description Performs comprehensive validation of a single email address
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param address query string true "Email address to verify" format(email)
// @Param verifyCatchAll query bool false "Verify catch-all domain" default(true)
// @Success 200 {object} EmailVerificationResponse "Email verification results"
// @Failure 400 {object} rest.ErrorResponse "Invalid email format or missing parameters"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /verify/v1/email [get]
// @Security ApiKeyAuth
func VerifyEmailAddress(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "VerifyEmailAddress", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}
		logger := services.Log

		// Check if email address is provided
		emailAddress := c.Query("address")
		if emailAddress == "" {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing parameter: address"))
			return
		}
		span.LogKV("request.address", emailAddress)
		span.LogKV("request.verifyCatchAll", c.Query("verifyCatchAll"))
		// check if verifyCatchAll param exists, defaulted to true
		verifyCatchAll := true
		if strings.ToLower(c.Query("verifyCatchAll")) == "false" {
			verifyCatchAll = false
		}

		syntaxValidation := mailsherpa.ValidateEmailSyntax(emailAddress)
		if !syntaxValidation.IsValid {
			c.JSON(http.StatusOK, EmailVerificationResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Email: EmailVerificationRecord{
					EmailAddress: emailAddress,
					Syntax: EmailVerificationSyntax{
						IsValid: false,
					},
				},
			})
			logger.Warnf("Invalid email address format: %s", emailAddress)
			return
		}

		// call validation api
		result, err := CallApiValidateEmail(ctx, services, emailAddress, verifyCatchAll)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if result.Status != "success" {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage(result.Message))
			return
		}

		deliverable := result.Data.EmailData.Deliverable
		if deliverable == "" {
			deliverable = "unknown"
		}

		emailVerificationResponse := EmailVerificationRecord{
			EmailAddress:          emailAddress,
			Deliverable:           deliverable,
			Provider:              result.Data.DomainData.Provider,
			SecureGatewayProvider: result.Data.DomainData.SecureGatewayProvider,
			IsCatchAll:            result.Data.DomainData.IsCatchAll,
			IsRisky: result.Data.DomainData.IsFirewalled ||
				result.Data.EmailData.IsRoleAccount ||
				result.Data.EmailData.IsSystemGenerated ||
				result.Data.EmailData.IsFreeAccount ||
				result.Data.EmailData.IsMailboxFull ||
				!result.Data.DomainData.IsPrimaryDomain,
			Syntax: EmailVerificationSyntax{
				IsValid: syntaxValidation.IsValid,
				Domain:  syntaxValidation.Domain,
				User:    syntaxValidation.User,
			},
			Risk: EmailVerificationRisk{
				IsFirewalled:      result.Data.DomainData.IsFirewalled,
				IsRoleMailbox:     result.Data.EmailData.IsRoleAccount,
				IsSystemGenerated: result.Data.EmailData.IsSystemGenerated,
				IsFreeProvider:    result.Data.EmailData.IsFreeAccount,
				IsMailboxFull:     result.Data.EmailData.IsMailboxFull,
				IsPrimaryDomain:   result.Data.DomainData.IsPrimaryDomain,
			},
			AlternateEmail: result.Data.EmailData.AlternateEmail,
		}

		if emailVerificationResponse.Deliverable != "unknown" {
			billableEvent := postgresentity.BillableEventEmailVerifiedNotCatchAll
			if emailVerificationResponse.IsCatchAll {
				billableEvent = postgresentity.BillableEventEmailVerifiedCatchAll
			}
			_, err = services.CommonServices.PostgresRepositories.ApiBillableEventRepository.RegisterEvent(ctx, tenant, billableEvent,
				postgresrepository.BillableEventDetails{
					ReferenceData: emailAddress,
				})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to register billable event"))
			}
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, EmailVerificationResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Email:        emailVerificationResponse,
		})
	}
}

// @Summary Upload emails for bulk verification
// @Description Initiates bulk verification process for emails from CSV file
// @Tags Email Verification
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing email addresses"
// @Param emailColumn formData string false "CSV column containing emails" default(first column)
// @Param verifyCatchAll formData bool false "Verify catch-all domains" default(true)
// @Success 200 {object} BulkUploadResponse "Bulk verification initiated"
// @Failure 400 {object} rest.ErrorResponse "Invalid file format or missing parameters"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /verify/v1/email/bulk [post]
// @Security ApiKeyAuth
func BulkUploadEmailsForVerification(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BulkUploadEmailsForVerification", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		// Get email column param (optional)
		emailColumn := c.DefaultPostForm("emailColumn", "")
		verifyCatchAllParam := c.DefaultPostForm("verifyCatchAll", "true")
		// check if verifyCatchAll param exists, defaulted to true
		verifyCatchAll := true
		if strings.ToLower(verifyCatchAllParam) == "false" {
			verifyCatchAll = false
		}
		span.LogKV("emailColumn", emailColumn)
		span.LogKV("verifyCatchAll", verifyCatchAllParam)

		// Parse the uploaded CSV file
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to read csv file"))
			return
		}
		defer file.Close()

		// Validate file type
		if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("File type not permitted.  Please send .csv file"))
			return
		}

		requestID := uuid.New().String()

		// Parse the CSV file
		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to parse csv file"))
			return
		}

		// If emailColumn is provided, ensure it exists in the CSV headers
		var emailIndex int
		if emailColumn != "" {
			emailIndex = -1
			for i, h := range headers {
				if strings.ToLower(h) == strings.ToLower(emailColumn) {
					emailIndex = i
					break
				}
			}
			if emailIndex == -1 {
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage(fmt.Sprintf("Column '%s' not found", emailColumn)))
				return
			}
		} else if len(headers) == 1 {
			emailIndex = 0 // Default to first column if only one column exists
		} else {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Please provide emailColumn parameter"))
			return
		}

		// Initialize slices for storing emails and validation data
		var emails []string

		// Read and validate each email
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to read csv file"))
				return
			}

			email := record[emailIndex]
			// Clean and skip empty or duplicate emails
			if email == "" || utils.Contains(emails, email) {
				continue
			}
			emails = append(emails, email)
		}

		// Register the bulk request in the database
		totalEmails := len(emails)
		if totalEmails == 0 {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("No records found in the csv file"))
			return
		}

		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.RegisterRequest(ctx, tenant, requestID, header.Filename, verifyCatchAll, totalEmails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to process bulk request"))
		}

		// Bulk insert email records into the database
		err = services.Repositories.PostgresRepositories.EmailValidationRecordRepository.BulkInsertRecords(ctx, tenant, requestID, verifyCatchAll, emails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert email records"))
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to process bulk request"))
			return
		}

		countPendingRequests, err := services.Repositories.PostgresRepositories.EmailValidationRecordRepository.CountPendingRequests(ctx, bulkRequest.Priority, bulkRequest.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to count pending requests"))
			countPendingRequests = 100 // default to 100 records
		}
		countPendingRequests = countPendingRequests + int64(totalEmails)

		// Respond with success message
		c.JSON(http.StatusOK,
			BulkUploadResponse{
				Message:               "File uploaded successfully",
				JobID:                 requestID,
				ResultURL:             fmt.Sprintf("%s/verify/v1/email/bulk/results/%s", services.Cfg.InternalServices.CustomerOsApiUrl, requestID), // Placeholder for results URL
				EstimatedCompletionTs: float64(calculateEstimatedCompletionTs(countPendingRequests)),
			})
	}
}

// @Summary Get bulk verification results
// @Description Retrieves results or status of bulk verification job
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param requestId path string true "Bulk verification job ID" format(uuid)
// @Success 200 {object} BulkResultsResponse "Verification results or status"
// @Failure 400 {object} rest.ErrorResponse "Invalid job ID"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} rest.ErrorResponse "Job not found"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /verify/v1/email/bulk/results/{requestId} [get]
// @Security ApiKeyAuth
func GetBulkEmailVerificationResults(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetBulkEmailVerificationResults", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		requestID := c.Param("requestId")
		if requestID == "" {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing parameter: ID"))
			return
		}
		span.LogKV("requestId", requestID)

		// Fetch the bulk request from the database
		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("invalid requestId"))
			return
		}
		if bulkRequest == nil {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}

		countPendingRequests, err := services.Repositories.PostgresRepositories.EmailValidationRecordRepository.CountPendingRequests(ctx, bulkRequest.Priority, bulkRequest.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to count pending requests"))
			countPendingRequests = 100 // default to 100 records
		}
		// Check if the processing is completed
		if bulkRequest.Status == postgresentity.EmailValidationRequestBulkStatusProcessing {
			c.JSON(http.StatusOK,
				BulkResultsResponse{
					Status:                "processing",
					JobID:                 requestID,
					FileName:              bulkRequest.FileName,
					Message:               fmt.Sprintf("Completed %d of %d emails", bulkRequest.DeliverableEmails+bulkRequest.UndeliverableEmails, bulkRequest.TotalEmails),
					Results:               nil,
					EstimatedCompletionTs: calculateEstimatedCompletionTs(countPendingRequests),
				})
			return
		}

		// Return the results if the processing is completed
		c.JSON(http.StatusOK, BulkResultsResponse{
			JobID:    requestID,
			Status:   "completed",
			FileName: bulkRequest.FileName,
			Results: &BulkResultsDetails{
				TotalEmails:   bulkRequest.TotalEmails,
				Deliverable:   bulkRequest.DeliverableEmails,
				Undeliverable: bulkRequest.UndeliverableEmails,
				DownloadURL:   fmt.Sprintf("%s/verify/v1/email/bulk/results/%s/download", services.Cfg.InternalServices.CustomerOsApiUrl, requestID),
			},
		})
	}
}

// @Summary Download bulk verification results
// @Description Downloads CSV file containing detailed verification results
// @Tags Email Verification
// @Accept json
// @Produce text/csv
// @Param requestId path string true "Bulk verification job ID" format(uuid)
// @Success 200 {file} csv "CSV file containing verification results"
// @Failure 400 {object} rest.ErrorResponse "Invalid job ID"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} rest.ErrorResponse "Results not found"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /verify/v1/email/bulk/results/{requestId}/download [get]
// @Security ApiKeyAuth
func DownloadBulkEmailVerificationResults(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetBulkEmailVerificationResults", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// Extract requestID from the path parameter
		requestID := c.Param("requestId")
		if requestID == "" {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing requestId"))
			return
		}

		// Fetch the bulk request to ensure it exists and is completed
		// Fetch the bulk request from the database
		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to retrieve request"))
			return
		}
		if bulkRequest == nil {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("Unable to find request"))
			return
		}

		// Check if the bulk request is completed before proceeding
		if bulkRequest.Status != postgresentity.EmailValidationRequestBulkStatusCompleted {
			c.JSON(http.StatusAccepted, gin.H{
				"status":    string(enum.StatusProcessing),
				"requestId": requestID,
				"message":   "The bulk request is still being processed. Please try again later.",
			})
			return
		}

		if bulkRequest.FileStoreId == "" {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("csv file not found"))
			return
		}

		fileDTO, fileContent, err := services.FileStoreApiService.GetFile(bulkRequest.Tenant, bulkRequest.FileStoreId, span)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get file using file store api"))
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to fetch the csv file"))
			return
		}

		// Set the response headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileDTO.FileName))
		c.Header("Content-Type", "text/csv")

		// Return the CSV content as a response
		c.Writer.Write(*fileContent)
	}
}

func CallApiValidateEmail(ctx context.Context, services *service.Services, emailAddress string, verifyCatchAll bool) (*validationmodel.ValidateEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CallApiValidateEmail")
	defer span.Finish()

	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequestWithOptions{
		Email: emailAddress,
		Options: validationmodel.ValidateEmailRequestOptions{
			VerifyCatchAll: verifyCatchAll,
		},
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", services.Cfg.InternalServices.ValidationApi+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.InternalServices.ValidationApiKey)
	req.Header.Set(security.TenantHeader, common.GetTenantFromContext(ctx))

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	// if response status is 504 retry once
	if response.StatusCode == http.StatusGatewayTimeout {
		span.LogFields(log.Int("response.status.firstAttempt", response.StatusCode))
		response, err = client.Do(req)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
			return nil, err
		}
		defer response.Body.Close()
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
			return nil, err
		}
	}

	span.LogFields(log.Int("response.statusCode", response.StatusCode))

	if response.StatusCode == http.StatusGatewayTimeout {
		err = errors.New("validation api returned 504 status code")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var validationResponse validationmodel.ValidateEmailResponse
	err = json.Unmarshal(responseBody, &validationResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(responseBody)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return nil, err
	}
	if validationResponse.Data == nil {
		tracing.LogObjectAsJson(span, "response", validationResponse)
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &validationResponse, nil
}

func calculateEstimatedCompletionTs(pendingRequests int64) int64 {
	// Total estimated time in seconds for pending requests
	totalTime := float64(pendingRequests) * singleEmailVerificationAproxDurationInSeconds / threadsToVerifyBulkEmails

	// Add 10 seconds to account for delay between cron runs
	totalTime += 10

	// Get the current time and add the estimated time
	estimatedCompletionTime := utils.Now().Add(time.Duration(totalTime) * time.Second)

	// Return the epoch timestamp (in seconds)
	return estimatedCompletionTime.Unix()
}
