package neo4j_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
)

type TenantSettingsProperty string

const (
	TenantSettingsPropertyLogoRepositoryFileId     TenantSettingsProperty = "logoRepositoryFileId"
	TenantSettingsPropertyBaseCurrency             TenantSettingsProperty = "baseCurrency"
	TenantSettingsPropertyInvoicingPostpaid        TenantSettingsProperty = "invoicingPostpaid"
	TenantSettingsPropertyWorkspaceLogo            TenantSettingsProperty = "workspaceLogo"
	TenantSettingsPropertyWorkspaceLogoIdentifier  TenantSettingsProperty = "workspaceLogoIdentifier"
	TenantSettingsPropertyWorkspaceName            TenantSettingsProperty = "workspaceName"
	TenantSettingsPropertyEnrichContacts           TenantSettingsProperty = "enrichContacts"
	TenantSettingsPropertyStripeCustomerPortalLink TenantSettingsProperty = "stripeCustomerPortalLink"
	TenantSettingsPropertySlackChannelUrl          TenantSettingsProperty = "slackChannelUrl"
)

type TenantSettingsEntity struct {
	Id                       string
	BaseCurrency             enum.Currency
	InvoicingPostpaid        bool
	WorkspaceLogo            string // Deprecated
	WorkspaceLogoIdentifier  string
	WorkspaceName            string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	EnrichContacts           bool
	StripeCustomerPortalLink string
	SharedSlackChannelUrl    string
}

func (t *TenantSettingsEntity) GetWorkspaceLogoCdnUrl() string {
	if t.WorkspaceLogoIdentifier == "" {
		return ""
	}
	return "base" + "/" + t.WorkspaceLogoIdentifier // TODO alexb implement it
}
