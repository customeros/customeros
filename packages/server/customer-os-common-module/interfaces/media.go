package interfaces

import (
	"context"
)

type MediaService interface {
	DownloadImageToS3(ctx context.Context, imageURL, bucketName, s3FilePath string) (string, error)
}
