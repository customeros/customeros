package config

type FileStoreApiConfig struct {
	ApiPath string `env:"FILE_STORE_API,required" envDefault:""`
	ApiKey  string `env:"FILE_STORE_API_KEY,required" envDefault:""`
}
