package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Tenant struct {
	Id              string                  `json:"id"`
	Name            string                  `json:"name"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
	SourceFields    commonmodel.Source      `json:"source"`
	BillingProfiles []*TenantBillingProfile `json:"billingProfiles"`
}

type TenantBillingProfile struct {
	Id           string             `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	SourceFields commonmodel.Source `json:"source"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	AddressLine1 string             `json:"addressLine1"`
	AddressLine2 string             `json:"addressLine2"`
	AddressLine3 string             `json:"addressLine3"`
	LegalName    string             `json:"legalName"`
}

func (t Tenant) HasBillingProfile(id string) bool {
	for _, bp := range t.BillingProfiles {
		if bp.Id == id {
			return true
		}
	}
	return false
}

func (t Tenant) GetBillingProfile(id string) *TenantBillingProfile {
	for _, bp := range t.BillingProfiles {
		if bp.Id == id {
			return bp
		}
	}
	return nil
}
