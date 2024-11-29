package flows

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
)

func getParticipantOrganizationIds(ctx rest.HTTPContext, domains []string) ([]string, error) {
	var results []string
	tenantDomains, err := ctx.Services.CommonServices.WorkspaceService.GetWorkspaceDomainsForTenant(*ctx.ServiceContext)
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
		orgId, err := ctx.Services.CommonServices.OrganizationService.Save(*ctx.ServiceContext, nil, nil, dataFields)
		if err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "Error saving organization by domain"))
		}

		results = append(results, orgId)
	}
	return results, nil
}
