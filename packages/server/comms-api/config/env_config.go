package config

type Config struct {
	Service struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
		ServerAddress    string `env:"COMMS_API_SERVER_ADDRESS,required"`
		CorsUrl          string `env:"COMMS_API_CORS_URL,required"`
	}
	Mail struct {
		ApiKey string `env:"COMMS_API_MAIL_API_KEY,required"`
	}
	GMail struct {
		ClientId     string `env:"GMAIL_CLIENT_ID,unset"`
		ClientSecret string `env:"GMAIL_CLIENT_SECRET,unset"`
		RedirectUris string `env:"GMAIL_REDIRECT_URIS"`
	}
	WebChat struct {
		PingInterval int `env:"WEBSOCKET_PING_INTERVAL"`
	}
	VCon struct {
		ApiKey          string `env:"COMMS_API_VCON_API_KEY,required"`
		AwsAccessKey    string `env:"AWS_ACCESS_KEY"`
		AwsAccessSecret string `env:"AWS_ACCESS_SECRET"`
		AwsRegion       string `env:"AWS_REGION"`
		AwsBucket       string `env:"AWS_BUCKET"`
	}
	Postgres struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
	}
	WebRTC struct {
		AuthSecret   string `env:"WEBRTC_AUTH_SECRET,required"`
		TTL          int    `env:"WEBRTC_AUTH_TTL,required"`
		PingInterval int    `env:"WEBSOCKET_PING_INTERVAL"`
	}
}
