package config

type NamecheapConfig struct {
	Url                   string  `env:"NAMECHEAP_URL" envDefault:"https://api.namecheap.com/xml.response" validate:"required"`
	ApiKey                string  `env:"NAMECHEAP_API_KEY" `
	ApiUser               string  `env:"NAMECHEAP_API_USER" `
	ApiUsername           string  `env:"NAMECHEAP_API_USERNAME"`
	ApiClientIp           string  `env:"NAMECHEAP_API_CLIENT_IP"`
	MaxPrice              float64 `env:"NAMECHEAP_MAX_PRICE" envDefault:"20.0" `
	Years                 int     `env:"NAMECHEAP_YEARS" envDefault:"1" `
	RegistrantFirstName   string  `env:"NAMECHEAP_REGISTRANT_FIRST_NAME" `
	RegistrantLastName    string  `env:"NAMECHEAP_REGISTRANT_LAST_NAME" `
	RegistrantCompanyName string  `env:"NAMECHEAP_REGISTRANT_COMPANY_NAME" `
	RegistrantJobTitle    string  `env:"NAMECHEAP_REGISTRANT_JOB_TITLE" `
	RegistrantAddress1    string  `env:"NAMECHEAP_REGISTRANT_ADDRESS1" `
	RegistrantCity        string  `env:"NAMECHEAP_REGISTRANT_CITY" `
	RegistrantState       string  `env:"NAMECHEAP_REGISTRANT_STATE" `
	RegistrantZIP         string  `env:"NAMECHEAP_REGISTRANT_ZIP" `
	RegistrantCountry     string  `env:"NAMECHEAP_REGISTRANT_COUNTRY" `
	RegistrantPhoneNumber string  `env:"NAMECHEAP_REGISTRANT_PHONE_NUMBER" `
	RegistrantEmail       string  `env:"NAMECHEAP_REGISTRANT_EMAIL" `
}
