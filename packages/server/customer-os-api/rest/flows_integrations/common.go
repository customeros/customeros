package integrations

import (
	"github.com/gin-gonic/gin"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
)

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
