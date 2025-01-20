// @openapi 3.0.0
package customerbase

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// BulkResponse represents the response for bulk operations with single error
// @Description Response structure for bulk operations with single error detail
type BulkResponse struct {
	// Summary of the bulk operation
	Summary BulkSummary `json:"summary,omitempty"`
	// Error details if any
	Details BulkErrorDetails `json:"details,omitempty"`
}

// BulkResponseMultipleErrors represents the response for bulk operations with multiple errors
// @Description Response structure for bulk operations with multiple error details
type BulkResponseMultipleErrors struct {
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
func (h *ContactHandler) CreateBulkContacts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.CreateBulkContacts", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		h.handleBulkJSONRequest(c)
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
func (h *ContactHandler) ImportContacts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Customerbase.ImportContacts", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			message := "Unsupported Content-Type"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		h.handleCSVUpload(c)
	}
}

func (h *ContactHandler) handleBulkJSONRequest(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Customerbase.handleBulkJSONRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var multipleContacts []ContactRecord
	if err := c.ShouldBindJSON(&multipleContacts); err != nil {
		message := "Unable to parse request"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		tracing.TraceErr(span, err)
		return
	}

	if len(multipleContacts) == 0 {
		message := "No contacts provided"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	var fail int
	var total int
	var validationErrors []BulkErrorDetails

	for i, contact := range multipleContacts {
		err, errValue := h.validateContactRecord(&contact)
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
		contact.ContactId = h.processContact(c.Request.Context(), contact)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		h.responseHandler.HandleSuccess(c, resp)
	case fail == 1:
		resp := BulkResponse{
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
		h.responseHandler.HandleSuccess(c, resp)
	case fail == total:
		message := "No valid contacts found in request"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
	default:
		resp := BulkResponseMultipleErrors{
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: validationErrors,
		}
		h.responseHandler.HandleCreated(c, resp)
		c.JSON(http.StatusCreated, resp)
	}
}

func (h *ContactHandler) handleCSVUpload(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Customerbase.handleCSVUpload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	file, err := h.validateAndOpenCsvFile(c)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if err := h.validateFileHeaders(c, reader); err != nil {
		return
	}

	h.processCSVRecords(c, reader)
}

func (h *ContactHandler) validateFileHeaders(c *gin.Context, reader *csv.Reader) error {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Customerbase.validateFileHeaders")
	defer span.Finish()
	tracing.TagComponentRest(span)

	headers, err := reader.Read()
	if err != nil {
		message := "Unable to read file"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
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
		message := "Missing required headers: email, linkedin_url"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return errors.New("invalid headers")
	}
	return nil
}

func (h *ContactHandler) processCSVRecords(c *gin.Context, reader *csv.Reader) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Customerbase.processCSVRecords")
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

		err, errVal := h.validateContactRecord(&contactRecord)
		if err != nil {
			fail++
			csvErrors = append(csvErrors, BulkErrorDetails{
				Value:        errVal,
				RecordNumber: total,
				Description:  fmt.Sprintf("%s", err),
			})
		}
		contactRecord.ContactId = h.processContact(c.Request.Context(), contactRecord)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		h.responseHandler.HandleCreated(c, resp)
	case fail == 1:
		resp := BulkResponse{
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
		h.responseHandler.HandleCreated(c, resp)
	case fail == total:
		message := "No valid contacts found in request"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
	default:
		resp := BulkResponseMultipleErrors{
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: csvErrors,
		}
		h.responseHandler.HandleCreated(c, resp)
	}
}

func (h *ContactHandler) validateAndOpenCsvFile(c *gin.Context) (multipart.File, error) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Customerbase.validateAndOpenCsvFile")
	defer span.Finish()
	tracing.TagComponentRest(span)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		message := "Unable to parse file"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		message := "Invalid file type"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return nil, errors.New("invalid file type")
	}

	return file, nil
}
