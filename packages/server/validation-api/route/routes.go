package route

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/service"
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
	validateEmailWithReacher(ctx, r, services, logger)
	validateEmailV2(ctx, r, services, logger)
	ipLookup(ctx, r, services, logger)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func validateEmailWithReacher(ctx context.Context, r *gin.Engine, services *service.Services, l logger.Logger) {
	r.POST("/validateEmail",
		tracing.TracingEnhancer(ctx, "POST /validateEmail"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		func(c *gin.Context) {
			var request model.ValidateEmailRequest

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			response, err := services.EmailValidationService.ValidateEmailWithReacher(ctx, request.Email)
			if err != nil {
				l.Errorf("Error validating email: %v", err.Error())
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, model.MapValidationEmailResponse(response, nil))
		})
}

func validateEmailV2(ctx context.Context, r *gin.Engine, services *service.Services, l logger.Logger) {
	r.POST("/validateEmailV2",
		tracing.TracingEnhancer(ctx, "POST /validateEmailV2"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		func(c *gin.Context) {
			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "ValidateEmailV2", c.Request.Header)
			defer span.Finish()

			var request model.ValidateEmailRequest

			if err := c.BindJSON(&request); err != nil {
				l.Errorf("Fail reading request: %v", err.Error())
				c.JSON(http.StatusBadRequest, model.ValidateEmailResponse{
					Status:  "error",
					Message: "Invalid request body",
				})
				return
			}

			// check ip is present
			if request.Email == "" {
				tracing.TraceErr(span, errors.New("Missing email parameter"))
				l.Errorf("Missing email parameter")
				c.JSON(http.StatusBadRequest, model.ValidateEmailResponse{
					Status:  "error",
					Message: "Missing email parameter",
				})
				return
			}
			span.LogFields(log.String("request.email", request.Email))

			response, err := services.EmailValidationService.ValidateEmailWithMailsherpa(ctx, request.Email)
			if err != nil {
				tracing.TraceErr(span, err)
				l.Errorf("Error on : %v", err.Error())
				c.JSON(http.StatusInternalServerError, model.ValidateEmailResponse{
					Status:  "error",
					Message: "Internal server error",
				})
				return
			}

			c.JSON(http.StatusOK, model.ValidateEmailResponse{
				Status: "success",
				Data:   response,
			})
		})
}

func validatePhoneNumber(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.POST("/validatePhoneNumber",
		tracing.TracingEnhancer(ctx, "POST /validatePhoneNumber"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
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
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
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
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
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
				Data:   response,
			})
		})
}
