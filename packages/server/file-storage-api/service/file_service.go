package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/h2non/filetype"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/entity"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
)

type FileService interface {
	GetById(tenantName string, id string) (*entity.File, error)
	UploadSingleFile(tenantName string, multipartFileHeader *multipart.FileHeader) (*entity.File, error)
	DownloadSingleFile(tenantName string, id string) (*entity.File, []byte, error)
	Base64Image(tenantName string, id string) (*string, error)
}

type fileService struct {
	cfg          *config.Config
	db           *gorm.DB
	repositories *repository.PostgresRepositories
}

func NewFileService(cfg *config.Config, db *gorm.DB, repositories *repository.PostgresRepositories) FileService {
	return &fileService{
		cfg:          cfg,
		db:           db,
		repositories: repositories,
	}
}

func (s *fileService) GetById(tenantName string, id string) (*entity.File, error) {
	byId := s.repositories.FileRepository.FindById(tenantName, id)
	if byId.Error != nil {
		return nil, byId.Error
	}

	return byId.Result.(*entity.File), nil
}

func (s *fileService) UploadSingleFile(tenantName string, multipartFileHeader *multipart.FileHeader) (*entity.File, error) {
	multipartFile, err := multipartFileHeader.Open()
	if err != nil {
		return nil, err
	}

	head := make([]byte, 1024)
	_, err = multipartFile.Read(head)
	if err != nil {
		return nil, err
	}

	//TODO docx is not recognized
	//https://github.com/h2non/filetype/issues/121
	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		fmt.Println("Unknown multipartFile type")
		return nil, errors.New("Unknown multipartFile type")
	}

	fileEntity := entity.File{
		TenantName: tenantName,
		Name:       multipartFileHeader.Filename,
		Extension:  kind.Extension,
		MIME:       kind.MIME.Value,
		Length:     multipartFileHeader.Size,
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		result := s.repositories.FileRepository.Save(&fileEntity)

		if result.Error != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = uploadFileToS3(s.cfg, session, fileEntity.ID, multipartFileHeader)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	return &fileEntity, nil
}

func (s *fileService) DownloadSingleFile(tenantName string, id string) (*entity.File, []byte, error) {
	byId := s.repositories.FileRepository.FindById(tenantName, id)
	if byId.Error != nil {
		return nil, nil, byId.Error
	}

	fileEntity := byId.Result.(*entity.File)

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}
	downloader := s3manager.NewDownloader(session)

	fileBytes := make([]byte, fileEntity.Length)
	_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(fileEntity.ID),
		})
	if err != nil {
		return nil, nil, err
	}

	return fileEntity, fileBytes, nil
}

func (s *fileService) Base64Image(tenantName string, id string) (*string, error) {
	byId := s.repositories.FileRepository.FindById(tenantName, id)
	if byId.Error != nil {
		return nil, byId.Error
	}

	fileEntity := byId.Result.(*entity.File)

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}
	downloader := s3manager.NewDownloader(session)

	fileBytes := make([]byte, fileEntity.Length)
	_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(fileEntity.ID),
		})
	if err != nil {
		return nil, err
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(fileBytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
		break
	case "image/png":
		base64Encoding += "data:image/png;base64,"
		break
	default:
		return nil, err // TODO say that the file can not be preview
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(fileBytes)
	return &base64Encoding, nil
}

func uploadFileToS3(cfg *config.Config, session *awsSes.Session, fileId string, multipartFile *multipart.FileHeader) error {
	open, err := multipartFile.Open()
	if err != nil {
		log.Fatal(err)
	}

	fileBytes := make([]byte, multipartFile.Size)
	_, err = open.Read(fileBytes)
	if err != nil {
		return err
	}

	_, err2 := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(cfg.AWS.Bucket),
		Key:                  aws.String(fileId),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBytes),
		ContentLength:        aws.Int64(int64(len(fileBytes))),
		ContentType:          aws.String(http.DetectContentType(fileBytes)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err2
}
