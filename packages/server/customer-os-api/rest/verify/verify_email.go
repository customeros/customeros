package restverify

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
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
	"io"
	"net/http"
	"strings"
	"time"
)

const singleEmailVerificationAproxDurationInSeconds = float64(3.5)
const threadsToVerifyBulkEmails = 6

// EmailVerificationResponse represents the response returned after verifying an email address
// @Description The response structure for email verification, providing detailed validation results.
// @example 200 {object} EmailVerificationResponse
type EmailVerificationResponse struct {
	// Status indicates the status of the verification (e.g., "success" or "failure")
	Status string `json:"status" example:"success"`

	// Message contains any additional information or errors related to the verification
	Message string `json:"message,omitempty" example:"Email verified successfully"`

	// Email is the email address that was verified
	Email string `json:"email" example:"example@example.com"`

	// Deliverable indicates whether the email is deliverable (e.g., "true", "false", "unknown")
	Deliverable string `json:"deliverable" example:"true"`

	// Provider is the email service provider (e.g., Gmail, Outlook)
	Provider string `json:"provider" example:"gmail"`

	// SecureGatewayProvider is the secure gateway provider (e.g., Proofpoint, Mimecast)
	SecureGatewayProvider string `json:"secureGatewayProvider" example:"Proofpoint"`

	// IsRisky indicates whether the email address is risky (e.g., used in spam or phishing)
	IsRisky bool `json:"isRisky" example:"false"`

	// IsCatchAll indicates if the email address is a catch-all address
	IsCatchAll bool `json:"isCatchAll" example:"false"`

	// Risk provides detailed risk factors associated with the email address
	Risk EmailVerificationRisk `json:"risk"`

	// Syntax provides details on the syntax validation of the email
	Syntax EmailVerificationSyntax `json:"syntax"`

	// AlternateEmail provides an alternate email if available
	AlternateEmail string `json:"alternateEmail" example:"alternate@example.com"`
}

// EmailVerificationRisk provides details on potential risks associated with the email address
type EmailVerificationRisk struct {
	// IsFirewalled indicates whether the email is protected by a firewall
	IsFirewalled bool `json:"isFirewalled" example:"false"`

	// IsRoleMailbox indicates if the email belongs to a role (e.g., info@, support@)
	IsRoleMailbox bool `json:"isRoleMailbox" example:"false"`

	// IsSystemGenerated indicates if the email is system-generated
	IsSystemGenerated bool `json:"isSystemGenerated" example:"false"`

	// IsFreeProvider indicates if the email uses a free provider like Gmail or Yahoo
	IsFreeProvider bool `json:"isFreeProvider" example:"true"`

	// IsMailboxFull indicates if the mailbox is full
	IsMailboxFull bool `json:"isMailboxFull" example:"false"`

	// IsPrimaryDomain indicates if the email belongs to a primary domain (not an alias)
	IsPrimaryDomain bool `json:"isPrimaryDomain" example:"true"`
}

// EmailVerificationSyntax provides details on the syntax validation of the email address
type EmailVerificationSyntax struct {
	// IsValid indicates if the syntax of the email is valid
	IsValid bool `json:"isValid" example:"true"`

	// Domain represents the domain part of the email address
	Domain string `json:"domain" example:"example.com"`

	// User represents the local part (before the @) of the email address
	User string `json:"user" example:"example"`
}

// BulkUploadResponse represents the response for the bulk email upload API.
// @Description Response structure for bulk email upload, containing job ID, result URL, and estimated completion time.
// @example 200 {object} BulkUploadResponse
type BulkUploadResponse struct {
	Message               string  `json:"message" example:"File uploaded successfully"`
	JobID                 string  `json:"jobId" example:"550e8400-e29b-41d4-a716-446655440000"`
	ResultURL             string  `json:"resultUrl" example:"https://api.customeros.ai/verify/v1/email/bulk/results/550e8400-e29b-41d4-a716-446655440000"`
	EstimatedCompletionTs float64 `json:"estimatedCompletionTs" example:"1694030400"` // Epoch timestamp
}

// BulkResultsResponse represents the response for the bulk email results API.
// @Description Response structure for returning bulk email verification results after processing.
// @example 200 {object} BulkResultsResponse
type BulkResultsResponse struct {
	JobID                 string              `json:"jobId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status                string              `json:"status" example:"completed"`
	FileName              string              `json:"fileName" example:"emails.csv"`
	Message               string              `json:"message" example:"Completed 1000 of 1000 emails"`
	Results               *BulkResultsDetails `json:"results"`
	EstimatedCompletionTs int64               `json:"estimatedCompletionTs" example:"1694030400"` // Epoch timestamp
}

// BulkResultsDetails contains the details of the results of the bulk verification.
// @Description Detailed results of the bulk email verification.
type BulkResultsDetails struct {
	TotalEmails   int    `json:"totalEmails" example:"1000"`
	Deliverable   int    `json:"deliverable" example:"950"`
	Undeliverable int    `json:"undeliverable" example:"45"`
	DownloadURL   string `json:"downloadUrl" example:"https://api.customeros.ai/verify/v1/email/bulk/results/550e8400-e29b-41d4-a716-446655440000/download"`
}

// @Summary Verify Single Email Address
// @Description Checks the validity and various characteristics of the given email address
// @Tags Verify API
// @Param address query string true "Email address to verify"
// @Param verifyCatchAll query string false "Verify catch-all domain" default(true)"
// @Success 200 {object} EmailVerificationResponse "Successful response"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Accept json
// @Router /verify/v1/email [get]
func VerifyEmailAddress(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "VerifyEmailAddress", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, rest.ErrorResponse{Status: "error", Message: "Missing tenant context"})
			return
		}
		logger := services.Log

		// Check if email address is provided
		emailAddress := c.Query("address")
		if emailAddress == "" {
			c.JSON(http.StatusBadRequest, rest.ErrorResponse{Status: "error", Message: "Missing address parameter"})
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
				Status: "success",
				Email:  emailAddress,
				Syntax: EmailVerificationSyntax{
					IsValid: false,
				},
			})
			logger.Warnf("Invalid email address format: %s", emailAddress)
			return
		}

		// call validation api
		result, err := CallApiValidateEmail(ctx, services, emailAddress, verifyCatchAll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, rest.ErrorResponse{Status: "error", Message: "Internal error"})
			return
		}
		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, rest.ErrorResponse{Status: "error", Message: result.Message})
			return
		}

		emailVerificationResponse := EmailVerificationResponse{
			Status:                "success",
			Email:                 emailAddress,
			Deliverable:           result.Data.EmailData.Deliverable,
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

		if emailVerificationResponse.Deliverable != "unknown" && emailVerificationResponse.Deliverable != "" {
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
		c.JSON(http.StatusOK, emailVerificationResponse)
	}
}

// @Summary Bulk Upload Emails for Verification
// @Description Uploads a CSV file with email addresses for bulk verification.
// @Tags Verify API
// @Param file formData file true "CSV file containing email addresses to verify"
// @Param emailColumn formData string false "The column name in the CSV that contains the email addresses (optional if only one column exists)"
// @Param verifyCatchAll formData string false "Verify catch-all domain" default(true)"
// @Success 200 {object} BulkUploadResponse "File uploaded successfully, with job ID and result URL"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Accept multipart/form-data
// @Router /verify/v1/email/bulk [post]
func BulkUploadEmailsForVerification(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BulkUploadEmailsForVerification", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing tenant context",
				})
			return
		}
		logger := services.Log

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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Failed to read file",
				})
			return
		}
		defer file.Close()

		// Validate file type
		if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Only CSV files are accepted",
				})
			return
		}

		requestID := uuid.New().String()

		// Parse the CSV file
		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Failed to parse CSV file",
				})
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
				c.JSON(http.StatusBadRequest,
					rest.ErrorResponse{
						Status:  "error",
						Message: fmt.Sprintf("Column '%s' not found", emailColumn),
					})
				return
			}
		} else if len(headers) == 1 {
			emailIndex = 0 // Default to first column if only one column exists
		} else {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Multiple columns found, please provide 'emailColumn' parameter",
				})
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
				c.JSON(http.StatusBadRequest,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Error reading CSV file",
					})
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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "No records found in the file",
				})
			return
		}

		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.RegisterRequest(ctx, tenant, requestID, header.Filename, verifyCatchAll, totalEmails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			logger.Errorf("Failed to register request: %v", err)
			c.JSON(http.StatusInternalServerError, rest.ErrorResponse{
				Status:  "error",
				Message: "Failed to register bulk request",
			})
		}

		// Bulk insert email records into the database
		err = services.Repositories.PostgresRepositories.EmailValidationRecordRepository.BulkInsertRecords(ctx, tenant, requestID, verifyCatchAll, emails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			logger.Errorf("Failed to insert records: %v", err)
			c.JSON(http.StatusInternalServerError, rest.ErrorResponse{
				Status:  "error",
				Message: "Failed to store email records",
			})
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

// @Summary Get Bulk Email Verification Results
// @Description Retrieves the results of bulk email verification if the processing is completed.
// @Tags Verify API
// @Param requestId path string true "Job ID of the bulk email verification"
// @Success 200 {object} BulkResultsResponse "Bulk email verification results if processing is completed"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce json
// @Router /verify/v1/email/bulk/results/{requestId} [get]
func GetBulkEmailVerificationResults(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetBulkEmailVerificationResults", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		requestID := c.Param("requestId")
		if requestID == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing request ID",
				})
			return
		}
		span.LogKV("requestId", requestID)

		// Fetch the bulk request from the database
		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Failed to retrieve request",
				})
			return
		}
		if bulkRequest == nil {
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Request not found",
				})
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

// @Summary Download Bulk Email Verification Results
// @Description Downloads the CSV file containing the results of bulk email verification if the processing is completed.
// @Tags Verify API
// @Param requestId path string true "Job ID of the bulk email verification"
// @Success 200 "CSV file containing the results of the bulk email verification"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Produce text/csv
// @Router /verify/v1/email/bulk/results/{requestId}/download [get]
func DownloadBulkEmailVerificationResults(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetBulkEmailVerificationResults", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// Extract requestID from the path parameter
		requestID := c.Param("requestId")
		if requestID == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing request ID",
				})
			return
		}

		// Fetch the bulk request to ensure it exists and is completed
		// Fetch the bulk request from the database
		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Failed to retrieve request",
				})
			return
		}
		if bulkRequest == nil {
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Request not found",
				})
			return
		}

		// Check if the bulk request is completed before proceeding
		if bulkRequest.Status != postgresentity.EmailValidationRequestBulkStatusCompleted {
			c.JSON(http.StatusAccepted, gin.H{
				"status":  "processing",
				"message": "The bulk request is still being processed. Please try again later.",
			})
			return
		}

		if bulkRequest.FileStoreId == "" {
			c.JSON(http.StatusNotFound,
				rest.ErrorResponse{
					Status:  "error",
					Message: "File not found"})
			return
		}

		fileDTO, fileContent, err := services.FileStoreApiService.GetFile(bulkRequest.Tenant, bulkRequest.FileStoreId, span)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get file using file store api"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Failed to fetch the file"})
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
