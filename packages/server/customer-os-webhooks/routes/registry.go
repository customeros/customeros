// routes/registry.go
package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonlogger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	webhook "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/webhooks"
)

type Registry struct {
	engine   *gin.Engine
	services *service.Services
	config   *config.Config
	logger   commonlogger.Logger
	cache    *commoncaches.Cache
}

func NewRegistry(
	engine *gin.Engine,
	services *service.Services,
	config *config.Config,
	logger commonlogger.Logger,
	cache *commoncaches.Cache,
) *Registry {
	return &Registry{
		engine:   engine,
		services: services,
		config:   config,
		logger:   logger,
		cache:    cache,
	}
}

func (r *Registry) RegisterRoutes(ctx context.Context) {
	r.registerWebhooks(ctx)
	r.registerBasicRoutes()
}

func (r *Registry) registerWebhooks(ctx context.Context) {
	// add all webhook routes here
	r.registerExternalSystemRoutes(ctx)
	r.registerInvoiceRoutes(ctx)
	r.registerOrganizationRoutes(ctx)
	r.registerPostmarkRoutes(ctx)
	r.registerUserRoutes(ctx)
}

func (r *Registry) registerInvoiceRoutes(ctx context.Context) {
	r.engine.POST("/sync/invoice",
		tracing.TracingEnhancer(ctx, "/sync/invoice"),
		security.ApiKeyCheckerHTTP(
			r.services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			r.services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.CUSTOMER_OS_WEBHOOKS,
			security.WithCache(r.cache),
		),
		webhook.NewInvoiceHandler(r.services, r.logger).Handle)
}

func (r *Registry) registerExternalSystemRoutes(ctx context.Context) {
	r.engine.POST("/sync/external-system",
		tracing.TracingEnhancer(ctx, "/sync/external-system"),
		security.ApiKeyCheckerHTTP(
			r.services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			r.services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.CUSTOMER_OS_WEBHOOKS,
			security.WithCache(r.cache),
		),
		webhook.NewExternalSystemHandler(r.services, r.logger).Handle)
}

func (r *Registry) registerOrganizationRoutes(ctx context.Context) {
	// Organization Routes
	r.engine.POST("/sync/organization",
		tracing.TracingEnhancer(ctx, "/sync/organization"),
		security.ApiKeyCheckerHTTP(
			r.services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			r.services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.CUSTOMER_OS_WEBHOOKS,
			security.WithCache(r.cache),
		),
		webhook.NewOrganizationHandler(r.services, r.logger).Handle)
}

func (r *Registry) registerPostmarkRoutes(ctx context.Context) {
	r.engine.POST("/sync/postmark-interaction-event",
		tracing.TracingEnhancer(ctx, "/sync/postmark-interaction-event"),
		webhook.NewInteractionEventHandler(r.services, r.config, r.logger).Handle)
}

func (r *Registry) registerUserRoutes(ctx context.Context) {
	r.engine.POST("/sync/user",
		tracing.TracingEnhancer(ctx, "/sync/user"),
		security.ApiKeyCheckerHTTP(
			r.services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			r.services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.CUSTOMER_OS_WEBHOOKS,
			security.WithCache(r.cache),
		),
		webhook.NewUserHandler(r.services, r.logger).Handle)
}
