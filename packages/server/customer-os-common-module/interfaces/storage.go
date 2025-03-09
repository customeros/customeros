package interfaces

import "context"

type StorageService interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) error
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetPublicURL(key string) string
}
