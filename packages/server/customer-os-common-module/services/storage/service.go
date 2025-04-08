package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/aws_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

// ObjectStorageService implements StorageService using S3Client
type ObjectStorageService struct {
	client     aws_client.S3Client
	bucketName string
	isPublic   bool
	cdnDomain  string // Optional CDN domain for public URLs
}

// StorageConfig holds configuration for object storage
type StorageConfig struct {
	BucketName string
	IsPublic   bool   // Whether objects should be publicly accessible
	CDNDomain  string // Optional CDN domain for public URLs
}

// NewStorageService creates a new object storage service
func NewStorageService(client aws_client.S3Client, config StorageConfig) interfaces.StorageService {
	return &ObjectStorageService{
		client:     client,
		bucketName: config.BucketName,
		isPublic:   config.IsPublic,
		cdnDomain:  config.CDNDomain,
	}
}

// Upload stores data in object storage
func (s *ObjectStorageService) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ObjectStorageService.Upload")
	defer spans.Finish()

	uploadInput := s3manager.UploadInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	// Set ACL if public
	if s.isPublic {
		uploadInput.ACL = aws.String("public-read")
	}

	return s.client.Upload(ctx, uploadInput)
}

// Download retrieves data from object storage
func (s *ObjectStorageService) Download(ctx context.Context, key string) ([]byte, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ObjectStorageService.Download")
	defer spans.Finish()

	content, err := s.client.Download(ctx, s.bucketName, key)
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}

// Delete removes an object from storage
func (s *ObjectStorageService) Delete(ctx context.Context, key string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ObjectStorageService.Delete")
	defer spans.Finish()

	if client, ok := interface{}(s.client).(interface {
		Delete(ctx context.Context, bucket, key string) error
	}); ok {
		return client.Delete(ctx, s.bucketName, key)
	}

	return nil
}

// GetPublicURL returns a public URL for the object
func (s *ObjectStorageService) GetPublicURL(key string) string {
	// Use CDN domain if provided
	if s.cdnDomain != "" {
		return "https://" + s.cdnDomain + "/" + key
	}

	return ""
}
