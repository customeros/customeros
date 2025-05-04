package neo4j_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
)

type TenantSettingsProperty string

const (
	TenantSettingsPropertyBaseCurrency             TenantSettingsProperty = "baseCurrency"
	TenantSettingsPropertyInvoicingPostpaid        TenantSettingsProperty = "invoicingPostpaid"
	TenantSettingsPropertyWorkspaceLogoKey         TenantSettingsProperty = "workspaceLogoKey"
	TenantSettingsPropertyWorkspaceLogoUrl         TenantSettingsProperty = "workspaceLogoUrl"
	TenantSettingsPropertyWorkspaceName            TenantSettingsProperty = "workspaceName"
	TenantSettingsPropertyEnrichContacts           TenantSettingsProperty = "enrichContacts"
	TenantSettingsPropertyStripeCustomerPortalLink TenantSettingsProperty = "stripeCustomerPortalLink"
	TenantSettingsPropertySlackChannelUrl          TenantSettingsProperty = "slackChannelUrl"
	// Deprecated
	TenantSettingsPropertyWorkspaceLogo TenantSettingsProperty = "workspaceLogo"
)

type TenantSettingsEntity struct {
	Id                       string
	BaseCurrency             enum.Currency
	InvoicingPostpaid        bool
	WorkspaceLogoUrl         string
	WorkspaceLogoKey         string
	WorkspaceName            string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	EnrichContacts           bool
	StripeCustomerPortalLink string
	SharedSlackChannelUrl    string

	// Deprecated
	WorkspaceLogo string
}
