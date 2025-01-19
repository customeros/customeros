package config

type AwsConfig struct {
	Region string `env:"AWS_S3_REGION"`
	Bucket string `env:"AWS_S3_BUCKET"`
}
