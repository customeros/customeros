package integrations

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func (h *IntegrationHandler) GetCustomerOSUser(ctx context.Context, emails []string, agent enum.AgentType) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Integrations.getCustomerOSUser")
	defer span.Finish()
	tracing.TagComponentRest(span)

	if len(emails) == 0 {
		return "", coserrors.ErrCannotIdentifyUser
	}

	for _, owner := range emails {
		validation := mailvalidate.ValidateEmailSyntax(owner)
		email := validation.CleanEmail
		if email == "" {
			continue
		}
		user, _ := h.services.CommonServices.UserService.FindUserByEmail(ctx, email)
		if user != nil && user.Id != "" {
			ctx = common.SetUserIdInContext(ctx, user.Id)
			agent, err := h.services.Repositories.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByUserAndType(ctx, []enum.AgentType{agent})
			if agent != nil {
				return user.Id, nil
			}
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
	}
	return "", coserrors.ErrCannotIdentifyUser
}

func getParticipantOrganizationIds(c *gin.Context, s *cosapi_services.Services, domains []string) ([]string, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "flows.getParticipantOrganizationIds")
	defer span.Finish()
	tracing.TagComponentRest(span)

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
			tracing.TraceErr(span, errors.Wrap(err, "Error saving organization by domain"))
		}

		results = append(results, orgId)
	}
	return results, nil
}
