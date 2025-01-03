// @openapi 3.0.0
package customerbase

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// BulkResponse represents the response for bulk operations with single error
// @Description Response structure for bulk operations with single error detail
type BulkResponse struct {
	// Inherits standard response fields
	enum.BaseResponse
	// Summary of the bulk operation
	Summary BulkSummary `json:"summary,omitempty"`
	// Error details if any
	Details BulkErrorDetails `json:"details,omitempty"`
}

// BulkResponseMultipleErrors represents the response for bulk operations with multiple errors
// @Description Response structure for bulk operations with multiple error details
type BulkResponseMultipleErrors struct {
	// Inherits standard response fields
	enum.BaseResponse
	// Summary of the bulk operation
	Summary BulkSummary `json:"summary,omitempty"`
	// List of error details
	Details []BulkErrorDetails `json:"details,omitempty"`
}

// BulkErrorDetails represents error details for bulk operations
// @Description Error details for failed operations in bulk processing
type BulkErrorDetails struct {
	// The value that caused the error
	// example: invalid@email..com
	Value string `json:"value"`
	// The location of the record that caused the error
	// example: 17
	RecordNumber int `json:"recordNumber"`
	// Description of the error
	// example: invalid email format
	Description string `json:"description"`
}

// BulkSummary represents the summary of a bulk operation
// @Description Summary statistics for bulk operations
type BulkSummary struct {
	// Total number of records processed
	// example: 100
	Total int `json:"total"`
	// Number of successfully processed records
	// example: 95
	Success int `json:"success"`
	// Number of failed records
	// example: 5
	Failed int `json:"failed"`
}

// @Summary Create multiple contacts
// @Description Creates multiple contacts from JSON input
// @Tags CustomerBASE API
// @Accept json
// @Produce json
// @Param contacts body []ContactRecord true "Array of contacts to create"
// @Success 201 {object} BulkResponse "All contacts created successfully"
// @Success 207 {object} BulkResponseMultipleErrors "Contacts created with some failures"
// @Failure 400 {object} rest.BaseResponse "Invalid request data"
// @Failure 401 {object} rest.BaseResponse "Unauthorized"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/contacts/bulk [post]
// @Security ApiKeyAuth
func CreateBulkContacts(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.CreateBulkContacts", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		handleBulkJSONRequest(c, s)
	}
}

// @Summary Import contacts from CSV
// @Description Creates multiple contacts from CSV file upload
// @Tags CustomerBASE API
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file with contact data (required headers: email, linkedin_url)"
// @Success 201 {object} BulkResponse "All contacts imported successfully"
// @Success 207 {object} BulkResponseMultipleErrors "Contacts imported with some failures"
// @Failure 400 {object} rest.BaseResponse "Invalid file format or data"
// @Failure 401 {object} rest.BaseResponse "Unauthorized"
// @Failure 415 {object} rest.BaseResponse "Unsupported content type"
// @Failure 500 {object} rest.BaseResponse "Internal server error"
// @Router /customerbase/v1/contacts/import [post]
// @Security ApiKeyAutl
func ImportContacts(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.ImportContacts", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrUnsupportedContentType)
			return
		}

		handleCSVUpload(c, s)
	}
}

func handleBulkJSONRequest(c *gin.Context, s *service.Services) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.handleBulkJSONRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var multipleContacts []ContactRecord
	if err := c.ShouldBindJSON(&multipleContacts); err != nil {
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
		tracing.TraceErr(span, err)
		return
	}

	if len(multipleContacts) == 0 {
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("No contacts provided"))
		return
	}

	var fail int
	var total int
	var validationErrors []BulkErrorDetails

	for i, contact := range multipleContacts {
		err, errValue := validateContactRecord(&contact)
		total++
		if err != nil {
			fail++
			errDetails := BulkErrorDetails{
				Value:        errValue,
				RecordNumber: i,
				Description:  fmt.Sprintf("%v", err),
			}

			validationErrors = append(validationErrors, errDetails)
		}
		contact.ContactId = processContact(c, s, contact)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		c.JSON(http.StatusCreated, resp)
	case fail == 1:
		resp := BulkResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusPartialSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: BulkErrorDetails{
				Value:        validationErrors[0].Value,
				RecordNumber: validationErrors[0].RecordNumber,
				Description:  validationErrors[0].Description,
			},
		}
		c.JSON(http.StatusPartialContent, resp)
	case fail == total:
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("No valid contacts found in request"))
	default:
		resp := BulkResponseMultipleErrors{
			BaseResponse: enum.BuildBaseResponse(enum.StatusPartialSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: validationErrors,
		}
		c.JSON(http.StatusCreated, resp)
	}
}

func handleCSVUpload(c *gin.Context, s *service.Services) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.handleCSVUpload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	file, err := rest.ValidateAndOpenCsvFile(c, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if err := validateFileHeaders(c, reader); err != nil {
		return
	}

	processCSVRecords(c, s, reader)
}

func validateFileHeaders(c *gin.Context, reader *csv.Reader) error {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.validateFileHeaders")
	defer span.Finish()
	tracing.TagComponentRest(span)

	headers, err := reader.Read()
	if err != nil {
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to read file"))
		return err
	}

	hasEmail := false
	hasLinkedIn := false
	for _, header := range headers {
		if header == "email" {
			hasEmail = true
		}
		if header == "linkedin_url" {
			hasLinkedIn = true
		}
	}

	if !hasEmail || !hasLinkedIn {
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing required headers: email, linkedin_url"))
		return errors.New("invalid headers")
	}
	return nil
}

func processCSVRecords(c *gin.Context, s *service.Services, reader *csv.Reader) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Customerbase.processCSVRecords")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var csvErrors []BulkErrorDetails
	var total int
	var fail int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			total++
			fail++
			csvErrors = append(csvErrors, BulkErrorDetails{
				Value:        fmt.Sprintf("%s", record),
				RecordNumber: total,
				Description:  "Unable to read record",
			})
			continue
		}

		contactRecord := ContactRecord{
			Email:       record[0],
			LinkedInURL: record[1],
		}

		if contactRecord.Email == "" && contactRecord.LinkedInURL == "" {
			continue
		}

		total++

		err, errVal := validateContactRecord(&contactRecord)
		if err != nil {
			fail++
			csvErrors = append(csvErrors, BulkErrorDetails{
				Value:        errVal,
				RecordNumber: total,
				Description:  fmt.Sprintf("%s", err),
			})
		}
		contactRecord.ContactId = processContact(c, s, contactRecord)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		c.JSON(http.StatusCreated, resp)
	case fail == 1:
		resp := BulkResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: BulkErrorDetails{
				Value:        csvErrors[0].Value,
				RecordNumber: csvErrors[0].RecordNumber,
				Description:  csvErrors[0].Description,
			},
		}
		c.JSON(http.StatusCreated, resp)
	case fail == total:
		rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("No valid contacts found in request"))
	default:
		resp := BulkResponseMultipleErrors{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: csvErrors,
		}
		c.JSON(http.StatusCreated, resp)
	}
}
