package config

type OpensearchConfig struct {
	TracingUrl      string `env:"OPENSEARCH_TRACING_URL"`
	TracingUsername string `env:"OPENSEARCH_TRACING_USERNAME"`
	TracingPassword string `env:"OPENSEARCH_TRACING_PASSWORD"`

	EventsUrl      string `env:"OPENSEARCH_EVENTS_URL"`
	EventsUsername string `env:"OPENSEARCH_EVENTS_USERNAME"`
	EventsPassword string `env:"OPENSEARCH_EVENTS_PASSWORD"`

	AIUrl      string `env:"OPENSEARCH_AI_URL"`
	AIUsername string `env:"OPENSEARCH_AI_USERNAME"`
	AIPassword string `env:"OPENSEARCH_AI_PASSWORD"`
}
