package rest

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
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

type EmailVerificationResponse struct {
	Status                string                  `json:"status"`
	Message               string                  `json:"message,omitempty"`
	Email                 string                  `json:"email"`
	Deliverable           string                  `json:"deliverable"`
	Provider              string                  `json:"provider"`
	SecureGatewayProvider string                  `json:"secureGatewayProvider"`
	IsRisky               bool                    `json:"isRisky"`
	IsCatchAll            bool                    `json:"isCatchAll"`
	Risk                  EmailVerificationRisk   `json:"risk"`
	Syntax                EmailVerificationSyntax `json:"syntax"`
	AlternateEmail        string                  `json:"alternateEmail"`
}

type EmailVerificationRisk struct {
	IsFirewalled    bool `json:"isFirewalled"`
	IsRoleMailbox   bool `json:"isRoleMailbox"`
	IsFreeProvider  bool `json:"isFreeProvider"`
	IsMailboxFull   bool `json:"isMailboxFull"`
	IsPrimaryDomain bool `json:"isPrimaryDomain"`
}

type EmailVerificationSyntax struct {
	IsValid bool   `json:"isValid"`
	Domain  string `json:"domain"`
	User    string `json:"user"`
}

func VerifyEmailAddress(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "VerifyEmailAddress", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
		logger := services.Log

		// Check if email address is provided
		emailAddress := c.Query("address")
		if emailAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing address parameter"})
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
		result, err := callApiValidateEmail(ctx, services, emailAddress, verifyCatchAll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal error"})
			return
		}
		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Message})
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
				result.Data.EmailData.IsFreeAccount ||
				result.Data.EmailData.IsMailboxFull ||
				!result.Data.DomainData.IsPrimaryDomain,
			Syntax: EmailVerificationSyntax{
				IsValid: syntaxValidation.IsValid,
				Domain:  syntaxValidation.Domain,
				User:    syntaxValidation.User,
			},
			Risk: EmailVerificationRisk{
				IsFirewalled:    result.Data.DomainData.IsFirewalled,
				IsRoleMailbox:   result.Data.EmailData.IsRoleAccount,
				IsFreeProvider:  result.Data.EmailData.IsFreeAccount,
				IsMailboxFull:   result.Data.EmailData.IsMailboxFull,
				IsPrimaryDomain: result.Data.DomainData.IsPrimaryDomain,
			},
			AlternateEmail: result.Data.EmailData.AlternateEmail,
		}

		c.JSON(http.StatusOK, emailVerificationResponse)
	}
}

func callApiValidateEmail(ctx context.Context, services *service.Services, emailAddress string, verifyCatchAll bool) (*validationmodel.ValidateEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callApiValidateEmail")
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
	req, err := http.NewRequest("POST", services.Cfg.Services.ValidationApi+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, services.Cfg.Services.ValidationApiKey)
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

func BulkUploadEmailsForVerification(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "BulkUploadEmailsForVerification", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, tenant)
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
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to read file"})
			return
		}
		defer file.Close()

		// Validate file type
		if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Only CSV files are accepted"})
			return
		}

		// TODO: Store the file in S3 for logging (placeholder)
		requestID := uuid.New().String()

		// Placeholder: Store the file in S3 (actual implementation to be done later)
		// err = services.S3.UploadFile(requestID, file)
		// if err != nil {
		//     logger.Errorf("Failed to upload file to S3: %v", err)
		//     c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to store file"})
		//     return
		// }

		// Parse the CSV file
		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to parse CSV file"})
			return
		}

		// If emailColumn is provided, ensure it exists in the CSV headers
		var emailIndex int
		if emailColumn != "" {
			emailIndex = -1
			for i, h := range headers {
				if h == emailColumn {
					emailIndex = i
					break
				}
			}
			if emailIndex == -1 {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Column '%s' not found", emailColumn)})
				return
			}
		} else if len(headers) == 1 {
			emailIndex = 0 // Default to first column if only one column exists
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Multiple columns found, please provide 'emailColumn' parameter"})
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
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Error reading CSV file"})
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
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "No records found in the file"})
			return
		}

		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.RegisterRequest(ctx, tenant, requestID, header.Filename, verifyCatchAll, totalEmails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			logger.Errorf("Failed to register request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to register bulk request"})
			return
		}

		// Bulk insert email records into the database
		err = services.Repositories.PostgresRepositories.EmailValidationRecordRepository.BulkInsertRecords(ctx, tenant, requestID, verifyCatchAll, emails)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			logger.Errorf("Failed to insert records: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to store email records"})
			return
		}

		countPendingRequests, err := services.Repositories.PostgresRepositories.EmailValidationRecordRepository.CountPendingRequests(ctx, bulkRequest.Priority, bulkRequest.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to count pending requests"))
			countPendingRequests = 100 // default to 100 records
		}
		countPendingRequests = countPendingRequests + int64(totalEmails)

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{
			"message":               "File uploaded successfully",
			"jobId":                 requestID,
			"resultUrl":             fmt.Sprintf("%s/verify/v1/email/bulk/results/%s", services.Cfg.Services.CustomerOsApiUrl, requestID), // Placeholder for results URL
			"estimatedCompletionTs": calculateEstimatedCompletionTs(countPendingRequests),
		})
	}
}

func GetBulkEmailVerificationResults(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetBulkEmailVerificationResults", c.Request.Header)
		defer span.Finish()

		requestID := c.Param("requestId")
		if requestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing request ID"})
			return
		}
		span.LogKV("requestId", requestID)

		// Fetch the bulk request from the database
		bulkRequest, err := services.Repositories.PostgresRepositories.EmailValidationRequestBulkRepository.GetByRequestID(ctx, requestID)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to insert records"))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve request"})
			return
		}
		if bulkRequest == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Request not found"})
			return
		}

		countPendingRequests, err := services.Repositories.PostgresRepositories.EmailValidationRecordRepository.CountPendingRequests(ctx, bulkRequest.Priority, bulkRequest.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to count pending requests"))
			countPendingRequests = 100 // default to 100 records
		}
		// Check if the processing is completed
		if bulkRequest.Status == entity.EmailValidationRequestBulkStatusProcessing {
			c.JSON(http.StatusOK, gin.H{
				"jobId":                 requestID,
				"status":                "processing",
				"message":               fmt.Sprintf("Completed %s of %s emails", bulkRequest.DeliverableEmails+bulkRequest.UndeliverableEmails, bulkRequest.TotalEmails),
				"estimatedCompletionTs": calculateEstimatedCompletionTs(countPendingRequests),
				"fileName":              bulkRequest.FileName,
				"results":               nil,
			})
			return
		}

		// TODO Placeholder: Generate download URL from S3 (to be implemented later)
		downloadURL := fmt.Sprintf("/path/to/s3/bucket/%s.csv", requestID)

		// Return the results if the processing is completed
		c.JSON(http.StatusOK, gin.H{
			"jobId":    requestID,
			"status":   "completed",
			"fileName": bulkRequest.FileName,
			"results": gin.H{
				"totalEmails":   bulkRequest.TotalEmails,
				"deliverable":   bulkRequest.DeliverableEmails,
				"undeliverable": bulkRequest.UndeliverableEmails,
				"downloadUrl":   downloadURL,
			},
		})
	}
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
