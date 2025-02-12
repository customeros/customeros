package config

type IntegrationAppConfig struct {
	WorkspaceKey                    string `env:"INTEGRATION_APP_WORKSPACE_KEY"`
	WorkspaceSecret                 string `env:"INTEGRATION_APP_WORKSPACE_SECRET"`
	ApiTriggerUrlCreatePaymentLinks string `env:"INTEGRATION_APP_API_TRIGGER_URL_CREATE_PAYMENT_LINKS"`

	IntegrationAppEventWebhookUrls struct {
		GeneratePaymentLinkUrl string `env:"INVOICE_GENERATE_PAYMENT_LINK_URL" envDefault:"" required:"true"`
		InvoiceFinalizedUrl    string `env:"INVOICE_READY_URL" envDefault:"" required:"true"`
	}
}
