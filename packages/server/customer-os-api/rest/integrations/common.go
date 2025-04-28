package integrations

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func (h *IntegrationHandler) GetCustomerOSUser(ctx context.Context, emails []string, agent enum.AgentType) (string, error) {
	spans, ctx := telemetry.StartRestSpan(ctx, "IntegrationHandler.getCustomerOSUser")
	defer spans.Finish()

	if len(emails) == 0 {
		return "", coserrors.ErrCannotIdentifyUser
	}

	for _, owner := range emails {
		validation := mailvalidate.ValidateEmailSyntax(owner)
		email := validation.CleanEmail
		if email == "" {
			continue
		}
		user, err := h.services.CommonServices.UserService.FindUserByEmail(ctx, email)
		if err != nil {
			spans.TraceError(err)
		}
		if user != nil && user.Id != "" {
			ctx = common.SetUserIdInContext(ctx, user.Id)
			activeAgent, err := h.services.Repositories.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByUserAndType(ctx, []enum.AgentType{agent})
			if activeAgent != nil {
				return user.Id, nil
			}
			if err != nil {
				spans.TraceError(err)
			}
		}
	}
	return "", coserrors.ErrCannotIdentifyUser
}

func getParticipantOrganizationIds(c *gin.Context, s *cosapi_services.Services, domains []string) ([]string, error) {
	spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.getParticipantOrganizationIds")
	defer spans.Finish()

	var results []string
	tenantDomains, err := s.CommonServices.WorkspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		return results, err
	}
	for _, domain := range domains {
		if utils.Contains(tenantDomains, domain) {
			continue
		}
		dataFields := data_fields.OrganizationFields{
			Domains: []string{
				domain,
			},
		}
		orgId, err := s.CommonServices.OrganizationService.Save(ctx, nil, nil, dataFields)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "Error saving organization by domain"))
		}

		results = append(results, orgId)
	}
	return results, nil
}
