package config

type PdfConverterConfig struct {
	PdfConverterUrl string `env:"PDF_CONVERTER_URL" envDefault:"http://localhost:11006"`
}
