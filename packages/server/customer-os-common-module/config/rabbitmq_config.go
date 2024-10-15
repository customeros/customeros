package config

type RabbitMQConfig struct {
	Url string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@127.0.0.1:5672/"`
}
