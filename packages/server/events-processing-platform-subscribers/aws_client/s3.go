package aws_client

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3ClientI interface {
	Upload(bucket, key string, content io.Reader) error
	Download(bucket, key string) (string, error)
	ChangeRegion(region string)
}

type S3Client struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
	Config     *aws.Config
}

func NewS3Client(config *aws.Config) *S3Client {
	s := session.Must(session.NewSession(config))
	return &S3Client{
		Uploader:   s3manager.NewUploader(s),
		Downloader: s3manager.NewDownloader(s),
		Config:     config,
	}
}

func (s *S3Client) Upload(bucket, key string, content io.Reader) error {
	_, err := s.Uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   content,
	})
	return err
}

func (s *S3Client) Download(bucket, key string) (string, error) {
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

func (s *S3Client) ChangeRegion(region string) {
	s.Config.Region = aws.String(region)
}
