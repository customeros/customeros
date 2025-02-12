package config

type OpensearchConfig struct {
	Url      string `env:"OPENSEARCH_URL"`
	Username string `env:"OPENSEARCH_USERNAME"`
	Password string `env:"OPENSEARCH_PASSWORD"`
}
