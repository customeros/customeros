package config

type RabbitMQConfig struct {
	Url string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672"`
}
