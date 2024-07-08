package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type TenantSettingsProperty string

const (
	TenantSettingsPropertyLogoRepositoryFileId TenantSettingsProperty = "logoRepositoryFileId"
	TenantSettingsPropertyBaseCurrency         TenantSettingsProperty = "baseCurrency"
	TenantSettingsPropertyInvoicingEnabled     TenantSettingsProperty = "invoicingEnabled"
	TenantSettingsPropertyInvoicingPostpaid    TenantSettingsProperty = "invoicingPostpaid"
	TenantSettingsPropertyOpportunityStages    TenantSettingsProperty = "opportunityStages"
	TenantSettingsPropertyEnrichContacts       TenantSettingsProperty = "enrichContacts"
)

type TenantSettingsEntity struct {
	Id                   string
	LogoRepositoryFileId string
	BaseCurrency         enum.Currency
	InvoicingEnabled     bool
	InvoicingPostpaid    bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	OpportunityStages    []string
	EnrichContacts       bool
}

func (t *TenantSettingsEntity) DefaultOpportunityStage() string {
	if len(t.OpportunityStages) == 0 {
		return ""
	}
	return t.OpportunityStages[0]
}
