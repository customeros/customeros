package config

type FileStoreApiConfig struct {
	ApiPath string `env:"FILE_STORE_API,required"`
	ApiKey  string `env:"FILE_STORE_API_KEY,required"`
}
