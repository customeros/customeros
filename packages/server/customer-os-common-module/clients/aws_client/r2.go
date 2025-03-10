package aws_client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// R2Config holds configuration specific to Cloudflare R2
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// NewR2Client creates an S3Client configured for Cloudflare R2
func NewR2Client(config R2Config) (S3Client, error) {
	// Create custom AWS config for R2
	awsCfg := &aws.Config{
		// Use the R2 endpoint format
		Endpoint:    aws.String("https://" + config.AccountID + ".r2.cloudflarestorage.com"),
		Region:      aws.String("auto"), // R2 uses "auto" region
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKeySecret, ""),
		// This is important for R2 compatibility
		S3ForcePathStyle: aws.Bool(true),
	}

	// Create session with the R2 config
	s := session.Must(session.NewSession(awsCfg))

	// Return S3Client with R2 configuration
	return &s3Client{
		Uploader:   s3manager.NewUploader(s),
		Downloader: s3manager.NewDownloader(s),
		Config:     awsCfg,
		Session:    s,
	}, nil
}
