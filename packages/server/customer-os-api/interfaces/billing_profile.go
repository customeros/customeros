package cosapi_interfaces

import (
	"context"
	"time"
)

type BillingProfileService interface {
	CreateBillingProfile(ctx context.Context, organizationId, legalName, taxId string, createdAt *time.Time) (string, error)
	UpdateBillingProfile(ctx context.Context, organizationId, billingProfileId string, legalName, taxId *string, updatedAt *time.Time) error
	LinkEmailToBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string, primary bool) error
	UnlinkEmailFromBillingProfile(ctx context.Context, organizationId, billingProfileId, emailId string) error
	LinkLocationToBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error
	UnlinkLocationFromBillingProfile(ctx context.Context, organizationId, billingProfileId, locationId string) error
}
