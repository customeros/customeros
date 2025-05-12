package config

type NATSConfig struct {
	Node1 string `env:"NATS_NODE_1_URL"`
	Node2 string `env:"NATS_NODE_2_URL"`
	Node3 string `env:"NATS_NODE_3_URL"`
}
