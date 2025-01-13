package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type FileService interface {
	GetById(ctx context.Context, id string) (*File, error)
	UploadSingleFile(ctx context.Context, basePath, fileId string, multipartFileHeader *multipart.FileHeader, cdnUpload bool) (*File, error)
	DownloadSingleFile(ctx context.Context, id string, context *gin.Context, inline bool) (*File, error)
	Base64Image(ctx context.Context, id string) (*string, error)
	GetFilePublicUrl(ctx context.Context, id string) (string, error)
	GetFileBytes(ctx context.Context, fileURL string) (*[]byte, error)
	UploadSingleFileBytes(ctx context.Context, basePath, fileID, fileName string, content *[]byte, cdn bool) (*File, error)
}

type File struct {
	ID        string
	FileName  string
	MimeType  string
	BasePath  string
	Size      int64
	CdnUrl    string
	PublicUrl string
}

type FileDTO struct {
	Id          string `json:"id"`
	FileName    string `json:"fileName"`
	MimeType    string `json:"mimeType"`
	Size        int64  `json:"size"`
	MetadataUrl string `json:"previewUrl"`
	DownloadUrl string `json:"downloadUrl"`
	CdnUrl      string `json:"cdnUrl"`
}
