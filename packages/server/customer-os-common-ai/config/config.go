package config

type Config struct {
	OpenAi struct {
		ApiKey       string `env:"OPENAI_API_KEY,required" envDefault:"N/A"`
		Organization string `env:"OPENAI_ORGANIZATION,required" envDefault:""`
		Model        string `env:"OPENAI_MODEL,required" envDefault:"gpt-3.5-turbo-1106"` // 1106 has an extra parameter available that locks response as JSON)
	}
	Anthropic struct {
		ApiPath string `env:"ANTHROPIC_API,required" envDefault:"N/A"`
		ApiKey  string `env:"ANTHROPIC_API_KEY,required" envDefault:"N/A"`
	}
}
