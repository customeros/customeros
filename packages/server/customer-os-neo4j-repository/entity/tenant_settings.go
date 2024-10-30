package entity

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type TenantSettingsProperty string

const (
	TenantSettingsPropertyLogoRepositoryFileId     TenantSettingsProperty = "logoRepositoryFileId"
	TenantSettingsPropertyBaseCurrency             TenantSettingsProperty = "baseCurrency"
	TenantSettingsPropertyInvoicingEnabled         TenantSettingsProperty = "invoicingEnabled"
	TenantSettingsPropertyInvoicingPostpaid        TenantSettingsProperty = "invoicingPostpaid"
	TenantSettingsPropertyWorkspaceLogo            TenantSettingsProperty = "workspaceLogo"
	TenantSettingsPropertyWorkspaceName            TenantSettingsProperty = "workspaceName"
	TenantSettingsPropertyEnrichContacts           TenantSettingsProperty = "enrichContacts"
	TenantSettingsPropertyStripeCustomerPortalLink TenantSettingsProperty = "stripeCustomerPortalLink"
	TenantSettingsPropertySlackChannelUrl          TenantSettingsProperty = "slackChannelUrl"
)

type TenantSettingsEntity struct {
	Id                       string
	LogoRepositoryFileId     string
	BaseCurrency             enum.Currency
	InvoicingEnabled         bool
	InvoicingPostpaid        bool
	WorkspaceLogo            string
	WorkspaceName            string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	EnrichContacts           bool
	StripeCustomerPortalLink string
	SharedSlackChannelUrl    string
}
