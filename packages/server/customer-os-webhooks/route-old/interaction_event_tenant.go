package route

import (
	"context"
	"errors"
	"regexp"
)

func (h *EmailWebhookHandler) resolveTenant(ctx context.Context, span opentracing.Span, data *model.PostmarkEmailWebhookData) (string, error) {
	pattern := regexp.MustCompile(`@([^.]+)\.`)

	for _, email := range data.BccFull {
		if matches := pattern.FindStringSubmatch(email.Email); len(matches) >= 2 {
			return h.validateTenant(ctx, matches[1])
		}
	}
	return "", errors.New("tenant not found")
}

func (h *EmailWebhookHandler) validateTenant(ctx context.Context, name string) (string, error) {
	tenant, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantByNameIgnoreCase(ctx, name)
	if err != nil || tenant == nil {
		return "", errors.New("invalid tenant")
	}
	return mapper.MapDbNodeToTenantEntity(tenant).Name, nil
}
