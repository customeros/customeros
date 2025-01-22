package config

type CloudflareConfig struct {
	Url    string `env:"CLOUDFLARE_URL" envDefault:"https://api.cloudflare.com/client/v4" validate:"required"`
	ApiKey string `env:"CLOUDFLARE_API_KEY" `
	Email  string `env:"CLOUDFLARE_API_EMAIL"`
}

type OpenSRSConfig struct {
	Url      string `env:"OPENSRS_URL" envDefault:"https://admin.a.hostedemail.com"`
	ApiKey   string `env:"OPENSRS_API_KEY"`
	Username string `env:"OPENSRS_API_USERNAME"`
}

type PostmarkConfig struct {
	Url                         string `env:"POSTMARK_URL" envDefault:"https://api.postmarkapp.com"`
	AccountApiKey               string `env:"POSTMARK_ACCOUNT_API_KEY"`
	DefaultInboundStreamWebhook string `env:"POSTMARK_DEFAULT_INBOUND_STREAM_WEBHOOK"`
}

type IntegrationAppConfig struct {
	WorkspaceKey                    string `env:"INTEGRATION_APP_WORKSPACE_KEY"`
	WorkspaceSecret                 string `env:"INTEGRATION_APP_WORKSPACE_SECRET"`
	ApiTriggerUrlCreatePaymentLinks string `env:"INTEGRATION_APP_API_TRIGGER_URL_CREATE_PAYMENT_LINKS"`
}

type SlackConfig struct {
	ClientID                        string `env:"SLACK_CLIENT_ID"`
	ClientSecret                    string `env:"SLACK_CLIENT_SECRET"`
	NotifyNewTenantRegisteredHook   string `env:"SLACK_NOTIFY_NEW_TENANT_REGISTERED_WEBHOOK"`
	InternalAlertsRegisteredWebhook string `env:"SLACK_INTERNAL_ALERTS_REGISTERED_WEBHOOK" envDefault:""`
	NotifyPostmarkEmail             string `env:"SLACK_NOTIFY_POSTMARK_EMAIL"`
	NotifyFlowGoalAchieved          string `env:"SLACK_NOTIFY_FLOW_GOAL_ACHIEVED"`
}

type StripeConfig struct {
	ApiKey string `env:"STRIPE_API_KEY" envDefault:"N/A"`
}

type NovuCofig struct {
	ApiKey      string `env:"NOVU_API_KEY"`
	FronteraUrl string `env:"NOVU_FRONTERA_URL"`
}

type QuickbooksCofig struct {
	ClientId     string `env:"QUICKBOOKS_CLIENT_ID"`
	ClientSecret string `env:"QUICKBOOKS_CLIENT_SECRET"`
	RedirectUrl  string `env:"QUICKBOOKS_REDIRECT_URL"`
}
