package config

type AppConfig struct {
	APIPort     string `env:"API_PORT,required" envDefault:"8080"`
	Environment string `env:"ENVIRONMENT,required" envDefault:"dev"`
}

type NATSConfig struct {
	Node1 string `env:"NATS_NODE_1_URL,required"`
	Node2 string `env:"NATS_NODE_2_URL"`
	Node3 string `env:"NATS_NODE_3_URL"`
}

type AnthropicConfig struct {
	ApiPath string `env:"ANTHROPIC_API_PATH" envDefault:"https://api.anthropic.com/v1/messages"`
	ApiKey  string `env:"ANTHROPIC_API_KEY"`
}

type DeepseekConfig struct {
	Url    string `env:"DEEPSEEK_URL" envDefault:"https://api.deepseek.com"`
	ApiKey string `env:"DEEPSEEK_API_KEY" envDefault:"N/A"`
}

type GroqConfig struct {
	Url    string `env:"GROQ_URL" envDefault:"https://api.groq.com/openai/v1/chat/completions"`
	ApiKey string `env:"GROQ_API_KEY" envDefault:"N/A"`
}

type GeminiConfig struct {
	ApiKey string `env:"GEMINI_API_KEY" envDefault:"N/A"`
}
