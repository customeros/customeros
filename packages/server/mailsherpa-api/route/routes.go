package route

import (
	"github.com/gin-gonic/gin"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/config"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/logger"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/model"
	"github.com/customeros/customeros/packages/server/mailsherpa-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cfg *config.Config, logger logger.Logger) {
	r.GET("/health", healthCheckHandler)
	r.GET("/readiness", healthCheckHandler)
	validateEmail(ctx, r, services, cfg, logger)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func validateEmail(ctx context.Context, r *gin.Engine, services *service.Services, cfg *config.Config, l logger.Logger) {
	r.POST("/validateEmail",
		tracing.TracingEnhancer(ctx, "POST /validateEmail"),
		security.ApiKeyCheckerHTTP(services.PostgresRepositories.TenantWebhookApiKeyRepository, services.PostgresRepositories.AppKeyRepository, security.MAILSHEPRA_API, security.WithCache(services.Cache)),
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

			emailValidationData, err := services.MailsherpaService.ValidateEmailWithMailSherpa(ctx, request.Email)
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

			tracing.LogObjectAsJson(span, "output", emailValidationData)
			c.JSON(http.StatusOK, model.ValidateEmailResponse{
				Status: "success",
				Data:   emailValidationData,
			})
		})
}
