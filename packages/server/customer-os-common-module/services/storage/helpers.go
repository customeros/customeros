package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/aws_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

// NewS3StorageService creates a StorageService configured for AWS S3
func NewS3StorageService(awsRegion, accessKeyID, accessKeySecret, bucketName string, isPublic bool) interfaces.StorageService {
	s3Client := aws_client.NewS3Client(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(accessKeyID, accessKeySecret, ""),
	})

	return NewStorageService(s3Client, StorageConfig{
		BucketName: bucketName,
		IsPublic:   isPublic,
	})
}

// NewR2StorageService creates a StorageService configured for Cloudflare R2
func NewR2StorageService(accountID, accessKeyID, accessKeySecret, bucketName string, isPublic bool) interfaces.StorageService {
	r2Client := aws_client.NewS3Client(&aws.Config{
		Endpoint:         aws.String("https://" + accountID + ".r2.cloudflarestorage.com"),
		Region:           aws.String("auto"),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, accessKeySecret, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	return NewStorageService(r2Client, StorageConfig{
		BucketName: bucketName,
		IsPublic:   isPublic,
	})
}
