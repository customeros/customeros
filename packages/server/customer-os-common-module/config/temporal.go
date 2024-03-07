package config

type TemporalConfig struct {
	HostPort  string `env:"TEMPORAL_HOSTPORT" envDefault:"openline-dev-temporal-frontend.openline.svc.cluster.local:7233"`
	Namespace string `env:"TEMPORAL_NAMESPACE" envDefault:"default"`
	RunWorker bool   `env:"TEMPORAL_RUN_WORKER" envDefault:"false"`
}
