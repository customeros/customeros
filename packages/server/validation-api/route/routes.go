package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net"
	"net/http"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cfg *config.Config, logger logger.Logger) {
	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)
	validateAddress(ctx, r, services)
	validatePhoneNumber(ctx, r, services)
	validateEmailV2(ctx, r, services, cfg, logger)
	validateEmailWithScrubby(ctx, r, services, logger)
	validateEmailWithTrueInbox(ctx, r, services, logger)
	ipLookup(ctx, r, services, logger)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func validateEmailV2(ctx context.Context, r *gin.Engine, services *service.Services, cfg *config.Config, l logger.Logger) {
	r.POST("/validateEmailV2",
		tracing.TracingEnhancer(ctx, "POST /validateEmailV2"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "ValidateEmailV2")
			defer span.Finish()

			var request model.ValidateEmailRequestWithOptions

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.JSON(http.StatusBadRequest, model.ValidateEmailResponse{
					Status:  "error",
					Message: "Invalid request body",
				})
				return
			}
			tracing.LogObjectAsJson(span, "request", request)

			// check email is present
			if request.Email == "" {
				tracing.TraceErr(span, errors.New("Missing email parameter"))
				l.Errorf("Missing email parameter")
				c.JSON(http.StatusBadRequest, model.ValidateEmailResponse{
					Status:  "error",
					Message: "Missing email parameter",
				})
				return
			}
			span.SetTag("email", request.Email)

			emailValidationData, err := services.EmailValidationService.ValidateEmailWithMailSherpa(ctx, request.Email)
			if err != nil {
				tracing.TraceErr(span, err)
				l.Errorf("Error on : %v", err.Error())
				c.JSON(http.StatusInternalServerError, model.ValidateEmailResponse{
					Status:          "error",
					Message:         "Internal server error",
					InternalMessage: err.Error(),
				})
				return
			}

			if emailValidationData != nil && emailValidationData.Syntax.IsValid && !emailValidationData.EmailData.IsRoleAccount && (emailValidationData.EmailData.Deliverable == string(model.EmailDeliverableStatusUnknown)) {
				if request.Options.VerifyCatchAll || emailValidationData.EmailData.RetryValidation == true {
					// Step 1 - try Enrow
					enrowResponseStr := ""
					if cfg.EnrowConfig.Enabled {
						enrowResponseStr, err = services.EmailValidationService.ValidateEmailWithEnrow(ctx, request.Email, request.Options.ExtendedWaitingTime)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "failed to call Enrow"))
							l.Errorf("Error on calling Enrow : %s", err.Error())
						}
						if enrowResponseStr != "" {
							if enrowResponseStr == "valid" {
								emailValidationData.EmailData.Deliverable = string(model.EmailDeliverableStatusDeliverable)
								emailValidationData.EmailData.RetryValidation = false
							} else if enrowResponseStr == "invalid" {
								// accept invalid response only if it's not catch-all
								if !emailValidationData.DomainData.IsCatchAll {
									emailValidationData.EmailData.Deliverable = string(model.EmailDeliverableStatusUndeliverable)
								}
							} else {
								err = fmt.Errorf("Unexpected: Enrow response: {%s}" + enrowResponseStr)
								tracing.TraceErr(span, err)
							}
						}
					}
					// Step 2 - try TrueInbox
					if enrowResponseStr == "" && cfg.TrueinboxConfig.Enabled {
						trueInboxResponse, err := services.EmailValidationService.ValidateEmailWithTrueinbox(ctx, request.Email)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "failed to call TrueInbox"))
							l.Errorf("Error on calling Trueinbox : %s", err.Error())
						} else if trueInboxResponse != nil {
							if trueInboxResponse.Result == "valid" {
								emailValidationData.EmailData.Deliverable = string(model.EmailDeliverableStatusDeliverable)
								if emailValidationData.DomainData.Provider == "" {
									emailValidationData.DomainData.Provider = mapProvider(trueInboxResponse.SmtpProvider)
								}
							} else if trueInboxResponse.Result == "invalid" {
								// accept invalid response only if it's not catch-all
								if !emailValidationData.DomainData.IsCatchAll {
									emailValidationData.EmailData.Deliverable = string(model.EmailDeliverableStatusUndeliverable)
								}
							} else {
								err = fmt.Errorf("Unexpected: Trueinbox response: {%s}" + trueInboxResponse.Result)
								tracing.TraceErr(span, err)
							}
						}
					}
				}
			}

			tracing.LogObjectAsJson(span, "output", emailValidationData)
			c.JSON(http.StatusOK, model.ValidateEmailResponse{
				Status: "success",
				Data:   emailValidationData,
			})
		})
}

func validateEmailWithScrubby(ctx context.Context, r *gin.Engine, services *service.Services, l logger.Logger) {
	r.POST("/validateEmailWithScrubby",
		tracing.TracingEnhancer(ctx, "POST /validateEmailWithScrubby"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "ValidateEmailWithScrubby")
			defer span.Finish()

			var request model.ValidateEmailRequest

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.JSON(http.StatusBadRequest, model.ValidateEmailWithScrubbyResponse{
					Status:  "error",
					Message: "Invalid request body",
				})
				return
			}
			span.LogFields(log.String("request.email", request.Email))

			// check email is present
			if request.Email == "" {
				tracing.TraceErr(span, errors.New("Missing email parameter"))
				l.Errorf("Missing email parameter")
				c.JSON(http.StatusBadRequest, model.ValidateEmailWithScrubbyResponse{
					Status:  "error",
					Message: "Missing email parameter",
				})
				return
			}

			validationStatus, err := services.EmailValidationService.ValidateEmailScrubby(ctx, request.Email)
			if err != nil {
				tracing.TraceErr(span, err)
				l.Errorf("Error on : %v", err.Error())
				c.JSON(http.StatusInternalServerError, model.ValidateEmailWithScrubbyResponse{
					Status:          "error",
					Message:         "Internal server error",
					InternalMessage: err.Error(),
				})
				return
			}
			span.LogFields(log.String("result.validationStatus", validationStatus))

			if validationStatus != string(postgresentity.ScrubbyStatusLowercaseValid) &&
				validationStatus != string(postgresentity.ScrubbyStatusLowercaseInvalid) &&
				validationStatus != string(postgresentity.ScrubbyStatusLowercasePending) {
				validationStatus = "unknown"
				l.Errorf("Unknown validation status: %s", validationStatus)
				tracing.TraceErr(span, errors.New("Unknown validation status"))
			}

			c.JSON(http.StatusOK, model.ValidateEmailWithScrubbyResponse{
				Status:         "success",
				EmailIsValid:   validationStatus == string(postgresentity.ScrubbyStatusLowercaseValid),
				EmailIsInvalid: validationStatus == string(postgresentity.ScrubbyStatusLowercaseInvalid),
				EmailIsUnknown: validationStatus == "unknown",
				EmailIsPending: validationStatus == string(postgresentity.ScrubbyStatusLowercasePending),
			})
		})
}

func validateEmailWithTrueInbox(ctx context.Context, r *gin.Engine, services *service.Services, l logger.Logger) {
	r.POST("/validateEmailWithTrueInbox",
		tracing.TracingEnhancer(ctx, "POST /validateEmailWithTrueInbox"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "ValidateEmailWithTrueInbox")
			defer span.Finish()

			var request model.ValidateEmailRequest

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.JSON(http.StatusBadRequest, model.ValidateEmailWithTrueinboxResponse{
					Status:  "error",
					Message: "Invalid request body",
				})
				return
			}
			span.LogFields(log.String("request.email", request.Email))

			// check email is present
			if request.Email == "" {
				tracing.TraceErr(span, errors.New("Missing email parameter"))
				l.Errorf("Missing email parameter")
				c.JSON(http.StatusBadRequest, model.ValidateEmailWithTrueinboxResponse{
					Status:  "error",
					Message: "Missing email parameter",
				})
				return
			}

			validationResult, err := services.EmailValidationService.ValidateEmailWithTrueinbox(ctx, request.Email)
			if err != nil {
				tracing.TraceErr(span, err)
				l.Errorf("Error on : %v", err.Error())
				c.JSON(http.StatusInternalServerError, model.ValidateEmailWithTrueinboxResponse{
					Status:  "error",
					Message: "Internal server error",
				})
				return
			}

			c.JSON(http.StatusOK, model.ValidateEmailWithTrueinboxResponse{
				Status: "success",
				Data:   validationResult,
			})
		})
}

func validatePhoneNumber(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.POST("/validatePhoneNumber",
		tracing.TracingEnhancer(ctx, "POST /validatePhoneNumber"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			var request model.ValidationPhoneNumberRequest

			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			e164, country, err := services.PhoneNumberValidationService.ValidatePhoneNumber(ctx, request.Country, request.PhoneNumber)
			if err != nil {
				errorMessage := err.Error()
				c.JSON(500, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			if e164 == nil {
				errorMessage := "Invalid phone number"
				c.JSON(400, model.MapValidationPhoneNumberResponse(nil, nil, &errorMessage, false))
				return
			}

			c.JSON(200, model.MapValidationPhoneNumberResponse(e164, country, nil, true))
		})
}

func validateAddress(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.POST("/validateAddress",
		tracing.TracingEnhancer(ctx, "POST /validateAddress"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			var request model.ValidationAddressRequest
			if err := c.BindJSON(&request); err != nil {
				errorMessage := "Invalid request body"
				c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
				return
			}
			if request.International {
				internationalAddressLookup, err := services.AddressValidationService.ValidateInternationalAddress(request.Address, request.Country)
				if err != nil {
					errorMessage := err.Error()
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				if internationalAddressLookup == nil {
					errorMessage := "Invalid address"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				addressVerified := false
				for _, v := range internationalAddressLookup.Results {
					if v.Analysis.VerificationStatus == "Verified" {
						addressVerified = true
						break
					}
				}

				if !addressVerified {
					errorMessage := "Address could not be verified"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				c.JSON(200, model.MapValidationInternationalAddressResponse(internationalAddressLookup, nil, true))
			} else {
				validatedAddressResponse, err := services.AddressValidationService.ValidateUsAddress(request.Address)
				if err != nil {
					errorMessage := err.Error()
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				if validatedAddressResponse == nil {
					errorMessage := "Invalid address"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				addressVerified := false
				for _, v := range validatedAddressResponse.Result.Addresses {
					if v.Verified {
						addressVerified = true
						break
					}
				}

				if !addressVerified {
					errorMessage := "Address could not be verified"
					c.JSON(400, model.MapValidationNoAddressResponse(&errorMessage))
					return
				}

				c.JSON(200, model.MapValidationUsAddressResponse(validatedAddressResponse, nil, true))
			}
		})
}

func ipLookup(ctx context.Context, r *gin.Engine, services *service.Services, l logger.Logger) {
	r.POST("/ipLookup",
		tracing.TracingEnhancer(ctx, "POST /ipLookup"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "IpLookup", c.Request.Header)
			defer span.Finish()

			var request model.IpLookupRequest

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.JSON(http.StatusBadRequest, model.IpLookupResponse{
					Status:  "error",
					Message: "Invalid request body",
				})
				return
			}

			// check ip is present
			if request.Ip == "" {
				tracing.TraceErr(span, errors.New("Missing ip parameter"))
				l.Errorf("Missing ip parameter")
				c.JSON(http.StatusBadRequest, model.IpLookupResponse{
					Status:  "error",
					Message: "Missing ip parameter",
				})
				return
			}
			span.LogFields(log.String("request.ip", request.Ip))

			// check ip format
			if net.ParseIP(request.Ip) == nil {
				tracing.TraceErr(span, errors.New("Invalid IP address format"))
				l.Errorf("Invalid IP address format: %s", request.Ip)
				c.JSON(http.StatusBadRequest, model.IpLookupResponse{
					Status:  "error",
					Message: "Invalid IP address format",
				})
				return
			}

			response, err := services.IpIntelligenceService.LookupIp(ctx, request.Ip)
			if err != nil {
				tracing.TraceErr(span, err)
				l.Errorf("Error on : %v", err.Error())
				c.JSON(http.StatusInternalServerError, model.IpLookupResponse{
					Status:  "error",
					Message: "Internal server error",
				})
				return
			}

			c.JSON(http.StatusOK, model.IpLookupResponse{
				Status: "success",
				IpData: response,
			})
		})
}

func mapProvider(input string) string {
	switch input {
	case "Google":
		return "google workspace"
	default:
		return input
	}
}
