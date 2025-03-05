package config

type OpensearchConfig struct {
	Tracing *OpensearchTracingConfig
	Events  *OpensearchEventsConfig
	AI      *OpensearchAIConfig
}

type OpensearchTracingConfig struct {
	Url      string `env:"OPENSEARCH_TRACING_URL"`
	Username string `env:"OPENSEARCH_TRACING_USERNAME"`
	Password string `env:"OPENSEARCH_TRACING_PASSWORD"`
}

type OpensearchEventsConfig struct {
	Url      string `env:"OPENSEARCH_EVENTS_URL"`
	Username string `env:"OPENSEARCH_EVENTS_USERNAME"`
	Password string `env:"OPENSEARCH_EVENTS_PASSWORD"`
}

type OpensearchAIConfig struct {
	Url      string `env:"OPENSEARCH_AI_URL"`
	Username string `env:"OPENSEARCH_AI_USERNAME"`
	Password string `env:"OPENSEARCH_AI_PASSWORD"`
}
