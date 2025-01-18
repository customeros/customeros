package config

type GrpcClientConfig struct {
	EventsProcessingPlatformEnabled     bool   `env:"EVENTS_PROCESSING_PLATFORM_ENABLED" envDefault:"true"`
	EventsProcessingPlatformUrl         string `env:"EVENTS_PROCESSING_PLATFORM_URL" `
	EventsProcessingPlatformApiKey      string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" `
	EventsProcessingPlatformServername  string `env:"EVENTS_PROCESSING_PLATFORM_SERVER_NAME" envDefault:""`
	EventsProcessingPlatformCertificate string `env:"EVENTS_PROCESSING_PLATFORM_CERTIFICATE" envDefault:""`
}
