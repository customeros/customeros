package config

type GrpcClientConfig struct {
	EventsProcessingPlatformEnabled     bool   `env:"EVENTS_PROCESSING_PLATFORM_ENABLED" envDefault:"true"`
	EventsProcessingPlatformUrl         string `env:"EVENTS_PROCESSING_PLATFORM_URL" validate:"required"`
	EventsProcessingPlatformApiKey      string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" validate:"required"`
	EventsProcessingPlatformServername  string `env:"EVENTS_PROCESSING_PLATFORM_SERVER_NAME" envDefault:""`
	EventsProcessingPlatformCertificate string `env:"EVENTS_PROCESSING_PLATFORM_CERTIFICATE" envDefault:""`
}
