package media

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/aws_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/storage"
)

const (
	AWS_REGION            = "eu-west-1"
	R2_IMAGES_BUCKET_NAME = "images"
)

type mediaService struct {
	postgresRepository   *postgres_repository.Repositories
	s3client             aws_client.S3Client
	r2ImageStoragePublic interfaces.StorageService
}

func NewMediaService(postgres *postgres_repository.Repositories, cfg *config.R2StorageConfig) interfaces.MediaService {
	s3 := aws_client.NewS3Client(&aws.Config{Region: aws.String(AWS_REGION)})
	r2ImageStorage := storage.NewR2StorageService(
		cfg.AccountID,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
		R2_IMAGES_BUCKET_NAME,
		true,
	)

	return &mediaService{
		postgresRepository:   postgres,
		s3client:             s3,
		r2ImageStoragePublic: r2ImageStorage,
	}
}

// UploadImageToS3 downloads an image from a URL directly to an S3 bucket
// using the existing S3Client implementation
func (s *mediaService) UploadImageToS3(ctx context.Context, imageURL, bucketName, s3FilePath string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MediaService.UploadImageToS3")
	defer spans.Finish()
	spans.LogKV("url", imageURL)
	spans.LogKV("bucket", bucketName)
	spans.LogKV("key", s3FilePath)

	if imageURL == "" {
		err := errors.New("imageURL cannot be empty")
		spans.TraceError(err)
		return "", err
	}
	if s3FilePath == "" {
		err := errors.New("s3FilePath cannot be empty")
		spans.TraceError(err)
		return "", err
	}

	// Get the image from the URL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	// Add User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.linkedin.com/")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		spans.LogKV("response.statusCode", resp.StatusCode)
		switch {
		case resp.StatusCode == http.StatusNotFound:
			return "", coserrors.ErrResourceNotFound
		case resp.StatusCode == http.StatusForbidden:
			return "", coserrors.ErrResourceForbidden
		default:
			err := errors.New("failed to download image")
			spans.TraceError(err)
			return "", err
		}
	}

	// Get content type from response headers
	contentType := resp.Header.Get("Content-Type")
	spans.LogKV("content_type", contentType)

	// Check if s3FilePath already has an extension
	extension := filepath.Ext(s3FilePath)
	if extension == "" {
		// No extension found, let's add one based on content type
		newExt := detectExtension(contentType, imageURL)
		if newExt != "" {
			s3FilePath = s3FilePath + newExt
			spans.LogKV("detected_extension", newExt)
			spans.LogKV("updated_key", s3FilePath)
		}
	}

	if contentType == "" {
		contentType = getContentTypeFromExtension(s3FilePath)
	}

	// Upload directly to S3 using the provided S3Client
	err = s.s3client.Upload(ctx, s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(s3FilePath),
		Body:        resp.Body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return s3FilePath, nil
}

func (s *mediaService) UploadImageDataToR2(ctx context.Context, data []byte, r2FilePath, fileName string, generateNewFileName bool) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MediaService.UploadImageDataToR2")
	defer spans.Finish()
	spans.LogKV("r2FilePath", r2FilePath)
	spans.LogKV("fileName", fileName)
	spans.LogKV("generateNewFileName", generateNewFileName)

	if r2FilePath == "" {
		err := errors.New("r2FilePath cannot be empty")
		spans.TraceError(err)
		return "", err
	}

	if fileName == "" {
		err := errors.New("fileName cannot be empty")
		spans.TraceError(err)
		return "", err
	}

	// prepare file name for r2 storage
	storedFileName := fileName
	if generateNewFileName {
		fileExt := filepath.Ext(fileName)
		storedFileName = utils.GenerateLowerAlphaNumeric(20)
		if fileExt != "" {
			storedFileName = storedFileName + fileExt
		}
	}

	// prepare content type
	contentType := getContentTypeFromExtension(fileName)
	spans.LogKV("contentType", contentType)

	// Construct storage key without bucket name since it's handled by the storage service
	storageKey := fmt.Sprintf("%s/%s", r2FilePath, storedFileName)
	spans.LogKV("storageKey", storageKey)

	// Store the file in the storage service
	if err := s.r2ImageStoragePublic.Upload(ctx, storageKey, data, contentType); err != nil {
		spans.TraceError(err)
		return "", errors.Wrap(err, "failed to upload file to r2")
	}

	return storageKey, nil
}

// detectExtension determines the appropriate file extension based on content type
// and falls back to extracting from URL if content type is not recognized
func detectExtension(contentType, url string) string {
	// Map of content types to file extensions
	contentTypeMap := map[string]string{
		"image/jpeg":      ".jpg",
		"image/jpg":       ".jpg",
		"image/png":       ".png",
		"image/gif":       ".gif",
		"image/webp":      ".webp",
		"image/tiff":      ".tiff",
		"image/bmp":       ".bmp",
		"image/x-icon":    ".ico",
		"image/svg+xml":   ".svg",
		"image/svg":       ".svg",
		"application/pdf": ".pdf",
	}

	// Check if we have a mapping for this content type
	if ext, ok := contentTypeMap[contentType]; ok {
		return ext
	}

	// If content type doesn't match or is empty, try to extract from URL
	urlExt := filepath.Ext(url)
	if urlExt != "" {
		return urlExt
	}

	// If we still can't determine the extension, default to .jpg
	// since it's a common image format
	return ".jpg"
}

// getContentTypeFromExtension determines the content type based on file extension
func getContentTypeFromExtension(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg" // fallback
	}
}

func (s *mediaService) GetR2PublicURL(storageKey string) string {
	return s.r2ImageStoragePublic.GetPublicURL(storageKey)
}
