package config

type NamecheapConfig struct {
	Url                   string  `env:"NAMECHEAP_URL" envDefault:"https://api.namecheap.com/xml.response" validate:"required"`
	ApiKey                string  `env:"NAMECHEAP_API_KEY" validate:"required"`
	ApiUser               string  `env:"NAMECHEAP_API_USER" validate:"required"`
	ApiUsername           string  `env:"NAMECHEAP_API_USERNAME" validate:"required"`
	ApiClientIp           string  `env:"NAMECHEAP_API_CLIENT_IP" validate:"required"`
	MaxPrice              float64 `env:"NAMECHEAP_MAX_PRICE" envDefault:"20.0" validate:"required"`
	Years                 int     `env:"NAMECHEAP_YEARS" envDefault:"1" validate:"required"`
	RegistrantFirstName   string  `env:"NAMECHEAP_REGISTRANT_FIRST_NAME" validate:"required"`
	RegistrantLastName    string  `env:"NAMECHEAP_REGISTRANT_LAST_NAME" validate:"required"`
	RegistrantCompanyName string  `env:"NAMECHEAP_REGISTRANT_COMPANY_NAME" validate:"required"`
	RegistrantJobTitle    string  `env:"NAMECHEAP_REGISTRANT_JOB_TITLE" validate:"required"`
	RegistrantAddress1    string  `env:"NAMECHEAP_REGISTRANT_ADDRESS1" validate:"required"`
	RegistrantCity        string  `env:"NAMECHEAP_REGISTRANT_CITY" validate:"required"`
	RegistrantState       string  `env:"NAMECHEAP_REGISTRANT_STATE" validate:"required"`
	RegistrantZIP         string  `env:"NAMECHEAP_REGISTRANT_ZIP" validate:"required"`
	RegistrantCountry     string  `env:"NAMECHEAP_REGISTRANT_COUNTRY" validate:"required"`
	RegistrantPhoneNumber string  `env:"NAMECHEAP_REGISTRANT_PHONE_NUMBER" validate:"required"`
	RegistrantEmail       string  `env:"NAMECHEAP_REGISTRANT_EMAIL" validate:"required"`
}
