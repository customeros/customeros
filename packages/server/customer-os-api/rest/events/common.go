package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"
)

func getParticipantOrganizationIds(ctx rest.HTTPContext, domains []string) ([]string, error) {
	var results []string
	tenantDomains, err := ctx.Services.CommonServices.WorkspaceService.GetWorkspaceDomainsForTenant(*ctx.ServiceContext)
	if err != nil {
		return results, err
	}
	for _, domain := range domains {
		if isDomainTenantDomain(domain, tenantDomains) {
			continue
		}
		dataFields := data_fields.OrganizationFields{
			Domains: []string{
				domain,
			},
		}
		id, err := ctx.Services.CommonServices.OrganizationService.Save(*ctx.ServiceContext, nil, nil, dataFields)
		if err != nil {
			tracing.TraceErr(ctx.Span, errors.Wrap(err, "Error saving organization by domain"))
		}

		results = append(results, id)
	}
	return results, nil
}

func isDomainTenantDomain(domain string, tenantDomains []string) bool {
	for _, tenantDomain := range tenantDomains {
		if domain == tenantDomain {
			return true
		}
	}
	return false
}
