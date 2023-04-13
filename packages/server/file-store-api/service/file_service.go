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
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"golang.org/x/net/context"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type FileService interface {
	GetById(userEmail, tenantName string, id string) (*model.File, error)
	UploadSingleFile(userEmail, tenantName string, multipartFileHeader *multipart.FileHeader) (*model.File, error)
	DownloadSingleFile(userEmail, tenantName string, id string) (*model.File, []byte, error)
	Base64Image(userEmail, tenantName string, id string) (*string, error)
}

type fileService struct {
	cfg           *config.Config
	graphqlClient *graphql.Client
}

func NewFileService(cfg *config.Config, graphqlClient *graphql.Client) FileService {
	return &fileService{
		cfg:           cfg,
		graphqlClient: graphqlClient,
	}
}

func (s *fileService) GetById(userEmail, tenantName string, id string) (*model.File, error) {
	attachment, err := s.getCosAttachmentById(userEmail, tenantName, id)
	if err != nil {
		return nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(attachment), nil
}

func (s *fileService) UploadSingleFile(userEmail, tenantName string, multipartFileHeader *multipart.FileHeader) (*model.File, error) {
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

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}

	graphqlRequest := graphql.NewRequest(
		`mutation AttachmentCreate($mimeType: String!, $name: String!, $size: Int64!, $extension: String!, $appSource: String!) {
			attachment_Create(input: {
				mimeType: $mimeType	
				name: $name
				size: $size
				extension: $extension
				appSource: $appSource
			}) {
				id
				createdAt
				mimeType
				name
				size
				extension
			}
		}`)
	graphqlRequest.Var("mimeType", kind.MIME.Value)
	graphqlRequest.Var("name", multipartFileHeader.Filename)
	graphqlRequest.Var("size", multipartFileHeader.Size)
	graphqlRequest.Var("extension", kind.Extension)
	graphqlRequest.Var("appSource", "file-store-api")

	err = s.addHeadersToGraphRequest(graphqlRequest, tenantName, userEmail)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.AttachmentCreateResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}

	err = uploadFileToS3(s.cfg, session, graphqlResponse.Attachment.Id, multipartFileHeader)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(&graphqlResponse.Attachment), nil
}

func (s *fileService) DownloadSingleFile(userEmail, tenantName string, id string) (*model.File, []byte, error) {
	attachment, err := s.getCosAttachmentById(userEmail, tenantName, id)
	if err != nil {
		return nil, nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}
	downloader := s3manager.NewDownloader(session)

	fileBytes := make([]byte, attachment.Size)
	_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(attachment.Id),
		})
	if err != nil {
		return nil, nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(attachment), fileBytes, nil
}

func (s *fileService) Base64Image(userEmail, tenantName string, id string) (*string, error) {
	attachment, err := s.getCosAttachmentById(userEmail, tenantName, id)
	if err != nil {
		return nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}
	downloader := s3manager.NewDownloader(session)

	fileBytes := make([]byte, attachment.Size)
	_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(attachment.Id),
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

func (s *fileService) getCosAttachmentById(userEmail, tenantName string, id string) (*model.Attachment, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetAttachment($id: ID!) {
			attachment(id: $id) {
				id
				createdAt
				mimeType
				name
				size
				extension
			}
		}`)
	graphqlRequest.Var("id", id)

	err := s.addHeadersToGraphRequest(graphqlRequest, tenantName, userEmail)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.AttachmentResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return &graphqlResponse.Attachment, nil
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

func (s *fileService) addHeadersToGraphRequest(req *graphql.Request, tenant string, userEmail string) error {
	req.Header.Add("X-Openline-API-KEY", s.cfg.Service.CustomerOsAPIKey)
	if userEmail != "" {
		req.Header.Add("X-Openline-USERNAME", userEmail)
	}
	if tenant != "" {
		req.Header.Add("X-Openline-TENANT", tenant)
	}

	return nil
}

func (s *fileService) contextWithTimeout() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	return ctx, cancel, nil
}
