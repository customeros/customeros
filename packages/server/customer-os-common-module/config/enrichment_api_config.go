package config

type EnrichmentAPIConfig struct {
	Url    string `env:"ENRICHMENT_API_URL"`
	ApiKey string `env:"ENRICHMENT_API_KEY"`
}
