package aws_client

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type S3Client interface {
	Upload(ctx context.Context, uploadContainer s3manager.UploadInput) error
	Download(ctx context.Context, bucket, key string) (string, error)
	ListFiles(ctx context.Context, bucket string) ([]string, error)
	ChangeRegion(ctx context.Context, region string)
	Delete(ctx context.Context, bucket, key string) error
	GetEndpoint() string
}

type s3Client struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
	Config     *aws.Config
	Session    *session.Session
}

func NewS3Client(config *aws.Config) S3Client {
	s := session.Must(session.NewSession(config))
	return &s3Client{
		Uploader:   s3manager.NewUploader(s),
		Downloader: s3manager.NewDownloader(s),
		Config:     config,
		Session:    s,
	}
}

func (s *s3Client) Upload(ctx context.Context, uploadContainer s3manager.UploadInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3Client.Upload")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	_, err := s.Uploader.Upload(&uploadContainer)
	return err
}

func (s *s3Client) Download(ctx context.Context, bucket, key string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3Client.Download")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	buffer := &aws.WriteAtBuffer{}
	_, err := s.Downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return "", err
	}

	return string(buffer.Bytes()), nil
}

func (s *s3Client) ListFiles(ctx context.Context, bucket string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3Client.ListFiles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := s3.New(s.Session)

	var files []string
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	// Use the ListObjectsV2 API to get all objects in the bucket
	err := session.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			if obj.Key != nil {
				files = append(files, *obj.Key)
			}
		}
		// Return true to continue paginating
		return true
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *s3Client) ChangeRegion(ctx context.Context, region string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3Client.ChangeRegion")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	s.Config.Region = aws.String(region)
}

func (s *s3Client) Delete(ctx context.Context, bucket, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "s3Client.Delete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	svc := s3.New(s.Session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *s3Client) GetEndpoint() string {
	return utils.IfNotNilString(*s.Config.Endpoint)
}
