package config

type RabbitMQConfig struct {
	Url string `env:"RABBITMQ_URL"`
}
