package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
)

type TenantService interface {
	ResolveTenant(ctx context.Context, data *model.PostmarkEmailWebhookData) (string, error)
	ValidateTenant(ctx context.Context, name string) (string, error)
	GetTenantFromEmail(email string) (string, error)
}

type tenantService struct {
	services *Services
	logger   logger.Logger
}

func NewTenantService(services *Services) TenantService {
	return &tenantService{
		services: services,
		logger:   logger.Logger,
	}
}

// ResolveTenant extracts and validates the tenant from the webhook data
func (s *tenantService) ResolveTenant(ctx context.Context, data *model.PostmarkEmailWebhookData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.ResolveTenant")
	defer span.Finish()

	// First try to get tenant from BCC addresses
	for _, email := range data.BccFull {
		if tenant, err := s.GetTenantFromEmail(email.Email); err == nil && tenant != "" {
			validTenant, err := s.ValidateTenant(ctx, tenant)
			if err != nil {
				span.LogFields(tracingLog.Bool("tenant.found", false))
				return "", err
			}

			span.LogFields(
				tracingLog.Bool("tenant.found", true),
				tracingLog.String("tenant.name", validTenant),
			)
			return validTenant, nil
		}
	}

	span.LogFields(tracingLog.Bool("tenant.found", false))
	return "", fmt.Errorf("no valid tenant found in webhook data")
}

// ValidateTenant verifies the tenant exists in the system and returns the normalized name
func (s *tenantService) ValidateTenant(ctx context.Context, name string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.ValidateTenant")
	defer span.Finish()

	// Get tenant from repository
	tenant, err := s.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, name)
	if err != nil {
		span.LogFields(tracingLog.Bool("tenant.found", false))
		return "", fmt.Errorf("error looking up tenant: %w", err)
	}

	if tenant == nil {
		span.LogFields(tracingLog.Bool("tenant.found", false))
		return "", fmt.Errorf("tenant not found: %s", name)
	}

	// Map to entity and return normalized name
	tenantEntity := mapper.MapDbNodeToTenantEntity(tenant)

	span.LogFields(
		tracingLog.Bool("tenant.found", true),
		tracingLog.String("tenant.name", tenantEntity.Name),
	)

	return tenantEntity.Name, nil
}

// GetTenantFromEmail extracts the tenant name from an email address
func (s *tenantService) GetTenantFromEmail(email string) (string, error) {
	// Regular expression to extract tenant name from email
	pattern := `@([^.]+)\.`
	tenantRegex, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid tenant regex pattern: %w", err)
	}

	matches := tenantRegex.FindStringSubmatch(email)
	if len(matches) < 2 {
		return "", nil
	}

	return strings.ToLower(matches[1]), nil
}

// IsTenantEmail checks if an email belongs to the specified tenant
func (s *tenantService) IsTenantEmail(email, tenant string) bool {
	expectedDomain := fmt.Sprintf("@%s.customeros.ai", strings.ToLower(tenant))
	return strings.HasSuffix(strings.ToLower(email), expectedDomain)
}

// GetTenantEmailDomain returns the email domain for a tenant
func (s *tenantService) GetTenantEmailDomain(tenant string) string {
	return fmt.Sprintf("%s.customeros.ai", strings.ToLower(tenant))
}

// FormatTenantEmail formats an email address for a tenant
func (s *tenantService) FormatTenantEmail(prefix, tenant string) string {
	return fmt.Sprintf("%s@%s", prefix, s.GetTenantEmailDomain(tenant))
}

// GetBccAddress returns the standard BCC address for a tenant
func (s *tenantService) GetBccAddress(tenant string) string {
	return s.FormatTenantEmail("bcc", tenant)
}
