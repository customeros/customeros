package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
)

type OpportunityService interface {
	UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood enum.RenewalLikelihood, amount *float64, comments *string, ownerUserId *string, adjustedRate *int64, appSource string) error
	UpdateRenewalsForOrganization(ctx context.Context, organizationId string, renewalLikelihood enum.RenewalLikelihood, renewalAdjustedRate *int64) error
}
