package event

const (
	TenantAddBillingProfileV1    = "V1_TENANT_BILLING_PROFILE_NEW"
	TenantUpdateBillingProfileV1 = "V1_TENANT_BILLING_PROFILE_UPDATE"
	TenantUpdateSettingsV1       = "V1_TENANT_SETTINGS_UPDATE"
	TenantAddBankAccountV1       = "V1_TENANT_BANK_ACCOUNT_CREATE"
	TenantUpdateBankAccountV1    = "V1_TENANT_BANK_ACCOUNT_UPDATE"
	TenantDeleteBankAccountV1    = "V1_TENANT_BANK_ACCOUNT_DELETE"
)

const (
	FieldMaskPhone                  = "phone"
	FieldMaskAddressLine1           = "addressLine1"
	FieldMaskAddressLine2           = "addressLine2"
	FieldMaskAddressLine3           = "addressLine3"
	FieldMaskZip                    = "zip"
	FieldMaskCountry                = "country"
	FieldMaskRegion                 = "region"
	FieldMaskLocality               = "locality"
	FieldMaskLegalName              = "legalName"
	FieldMaskVatNumber              = "vatNumber"
	FieldMaskSendInvoicesFrom       = "sendInvoicesFrom"
	FieldMaskSendInvoicesBcc        = "sendInvoicesBcc"
	FieldMaskCanPayWithPigeon       = "canPayWithPigeon"
	FieldMaskCanPayWithBankTransfer = "canPayWithBankTransfer"
	FieldMaskCheck                  = "check"

	FieldMaskLogoRepositoryFileId = "logoRepositoryFileId"
	FieldMaskBaseCurrency         = "baseCurrency"
	FieldMaskInvoicingEnabled     = "invoicingEnabled"
	FieldMaskInvoicingPostpaid    = "invoicingPostpaid"
	FieldMaskWorkspaceLogo        = "workspaceLogo"
	FieldMaskWorkspaceName        = "workspaceName"

	FieldMaskBankAccountBankName            = "bankAccountBankName"
	FieldMaskBankAccountBankTransferEnabled = "bankAccountBankTransferEnabled"
	FieldMaskBankAccountAllowInternational  = "bankAccountAllowInternational"
	FieldMaskBankAccountCurrency            = "bankAccountCurrency"
	FieldMaskBankAccountIban                = "bankAccountIban"
	FieldMaskBankAccountBic                 = "bankAccountBic"
	FieldMaskBankAccountSortCode            = "bankAccountSortCode"
	FieldMaskBankAccountAccountNumber       = "bankAccountAccountNumber"
	FieldMaskBankAccountRoutingNumber       = "bankAccountRoutingNumber"
	FieldMaskBankAccountOtherDetails        = "bankAccountOtherDetails"
)
