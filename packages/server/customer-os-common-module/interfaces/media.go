package interfaces

import (
	"context"
)

const (
	WorkspacePath   = "workspace"
	UserProfilePath = "user-profile"
)

type MediaService interface {
	UploadImageToS3(ctx context.Context, imageURL, bucketName, s3FilePath string) (string, error)

	UploadImageDataToR2(ctx context.Context, data []byte, r2FilePath, fileName string, generateNewFileName bool) (string, error)

	GetR2ImagePublicURL(storageKey string) string
}
