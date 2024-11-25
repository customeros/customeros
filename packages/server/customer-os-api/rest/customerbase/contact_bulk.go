package customerbase

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type BulkResponse struct {
	rest.BaseResponse
	Summary BulkSummary      `json:"summary,omitempty"`
	Details BulkErrorDetails `json:"details,omitempty"`
}

type BulkResponseMultipleErrors struct {
	rest.BaseResponse
	Summary BulkSummary        `json:"summary,omitempty"`
	Details []BulkErrorDetails `json:"details,omitempty"`
}

type BulkErrorDetails struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type BulkSummary struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
}

func CreateBulkContacts(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		handleBulkJSONRequest(httpContext)
	}
}

func ImportContacts(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrUnsupportedContentType)
		}

		handleCSVUpload(httpContext)
	}
}

func handleBulkJSONRequest(ctx rest.HTTPContext) {
	var multipleContacts []ContactRecord
	if err := ctx.GinContext.ShouldBindJSON(&multipleContacts); err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest)
		return
	}

	if len(multipleContacts) == 0 {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No contacts provided"))
		return
	}

	var fail int
	var total int
	var validationErrors []BulkErrorDetails

	for _, contact := range multipleContacts {
		err, errValue := validateContact(&contact)
		total++
		if err != nil {
			fail++
			errDetails := BulkErrorDetails{
				Value:       errValue,
				Description: fmt.Sprintf("%v", err),
			}

			validationErrors = append(validationErrors, errDetails)
		}
		contact.ContactId = processContact(ctx, contact)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	case fail == 1:
		resp := BulkResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: BulkErrorDetails{
				Value:       validationErrors[0].Value,
				Description: validationErrors[0].Description,
			},
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	case fail == total:
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No valid contacts found in request"))
	default:
		resp := BulkResponseMultipleErrors{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: validationErrors,
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	}
}

func handleCSVUpload(ctx rest.HTTPContext) {
	file, err := rest.ValidateAndOpenCsvFile(ctx.GinContext)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if err := validateFileHeaders(ctx.GinContext, ctx.Span, reader); err != nil {
		return
	}

	processCSVRecords(ctx, reader)
}

func validateFileHeaders(c *gin.Context, span opentracing.Span, reader *csv.Reader) error {
	headers, err := reader.Read()
	if err != nil {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to read file"))
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
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing required headers: email, linkedin_url"))
		return errors.New("invalid headers")
	}
	return nil
}

func processCSVRecords(ctx rest.HTTPContext, reader *csv.Reader) {
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
				Value:       fmt.Sprintf("%s", record),
				Description: "Unable to read record",
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

		err, errVal := validateContact(&contactRecord)
		if err != nil {
			fail++
			csvErrors = append(csvErrors, BulkErrorDetails{
				Value:       errVal,
				Description: fmt.Sprintf("%s", err),
			})
		}
		contactRecord.ContactId = processContact(ctx, contactRecord)
	}

	switch {
	case fail == 0:
		resp := BulkResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total,
				Failed:  fail,
			},
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	case fail == 1:
		resp := BulkResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: BulkErrorDetails{
				Value:       csvErrors[0].Value,
				Description: csvErrors[0].Description,
			},
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	case fail == total:
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("No valid contacts found in request"))
	default:
		resp := BulkResponseMultipleErrors{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Summary: BulkSummary{
				Total:   total,
				Success: total - fail,
				Failed:  fail,
			},
			Details: csvErrors,
		}
		ctx.GinContext.JSON(http.StatusCreated, resp)
	}
}
