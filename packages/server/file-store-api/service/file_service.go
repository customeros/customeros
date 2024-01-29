package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileService interface {
	GetById(userEmail, tenantName string, id string) (*model.File, error)
	UploadSingleFile(userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader) (*model.File, error)
	DownloadSingleFile(userEmail, tenantName, id string, context *gin.Context, inline bool) (*model.File, error)
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

func (s *fileService) UploadSingleFile(userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader) (*model.File, error) {
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

	if basePath == "" {
		basePath = "/GLOBAL"
	}

	if fileId == "" {
		fileId = uuid.New().String()
	}

	err = uploadFileToS3(s.cfg, session, tenantName, basePath, fileId+"."+kind.Extension, multipartFileHeader)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	graphqlRequest := graphql.NewRequest(
		`mutation AttachmentCreate($id: ID, $basePath: String!, $mimeType: String!, $size: Int64!, $fileName: String!, $appSource: String!) {
			attachment_Create(input: {
				id: $id	
				basePath: $basePath	
				mimeType: $mimeType	
				fileName: $fileName
				size: $size
				appSource: $appSource
			}) {
				id
			}
		}`)
	graphqlRequest.Var("id", fileId)
	graphqlRequest.Var("basePath", basePath)
	graphqlRequest.Var("mimeType", multipartFileHeader.Header.Get(http.CanonicalHeaderKey("Content-Type")))
	graphqlRequest.Var("fileName", multipartFileHeader.Filename)
	graphqlRequest.Var("size", multipartFileHeader.Size)
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

	return mapper.MapAttachmentResponseToFileEntity(&graphqlResponse.Attachment), nil
}

func (s *fileService) DownloadSingleFile(userEmail, tenantName, id string, context *gin.Context, inline bool) (*model.File, error) {
	attachment, err := s.getCosAttachmentById(userEmail, tenantName, id)
	byId := mapper.MapAttachmentResponseToFileEntity(attachment)
	if err != nil {
		return nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Printf("Error creating session: %v", err)
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	context.Header("Accept-Ranges", "bytes")

	svc := s3.New(session)

	extension := filepath.Ext(attachment.FileName)
	if extension == "" {
		fmt.Println("No file extension found.")
	} else {
		extension = extension[1:]
		fmt.Println("File Extension:", extension)
	}

	// Get the object metadata to determine the file size and ETag
	respHead, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(tenantName + byId.BasePath + "/" + attachment.Id + "." + extension),
	})
	if err != nil {
		// Handle error
		context.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}

	// Get the ETag header value
	eTag := aws.StringValue(respHead.ETag)

	// Parse the range header
	rangeHeader := context.GetHeader("Range")
	var start, end int64
	if rangeHeader != "" {
		log.Printf("Range header: %s", rangeHeader)
		log.Printf("Content Length: %d", *respHead.ContentLength)

		rangeParts := strings.Split(rangeHeader, "=")[1]
		rangeBytes := strings.Split(rangeParts, "-")
		start, _ = strconv.ParseInt(rangeBytes[0], 10, 64)
		log.Printf("rangeBytes %v", rangeBytes)
		if len(rangeBytes) > 1 && rangeBytes[1] != "" {
			end, _ = strconv.ParseInt(rangeBytes[1], 10, 64)
		} else {
			end = *respHead.ContentLength - 1
		}
	} else {
		start = 0
		end = *respHead.ContentLength - 1
	}

	// Set the content length header to the file size
	context.Header("Content-Length", strconv.FormatInt(end-start+1, 10))

	// Set the content range header to indicate the range of bytes being served
	context.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(*respHead.ContentLength, 10))

	// If the ETag matches, send a 304 Not Modified response and exit early
	if match := context.GetHeader("If-Range"); match != "" && match != eTag {
		context.Status(http.StatusRequestedRangeNotSatisfiable)
		return byId, nil
	}

	if !inline {
		context.Header("Content-Disposition", "attachment; filename="+byId.FileName)
	} else {
		context.Header("Content-Disposition", "inline; filename="+byId.FileName)
	}
	context.Header("Content-Type", fmt.Sprintf("%s", byId.MimeType))
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(tenantName + byId.BasePath + "/" + attachment.Id + "." + extension),
		Range:  aws.String("bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)),
	})
	if err != nil {
		// Handle error
		log.Printf("Error getting object: %v", err)
		context.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Serve the file contents
	io.Copy(context.Writer, resp.Body)
	return byId, nil
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

	if attachment.Size > s.cfg.MaxFileSizeMB*1024*1024 {
		return nil, errors.New("file is too big for base64 encoding")
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
				fileName
				basePath
				size
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

func uploadFileToS3(cfg *config.Config, session *awsSes.Session, tenantName, basePath, fileId string, multipartFile *multipart.FileHeader) error {
	fileStream, err := multipartFile.Open()
	if err != nil {
		return fmt.Errorf("uploadFileToS3: %w", err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(cfg.AWS.Bucket),
		Key:           aws.String(tenantName + basePath + "/" + fileId),
		ACL:           aws.String("private"),
		ContentLength: aws.Int64(0),
	})
	if err != nil {
		return fmt.Errorf("uploadFileToS3: %w", err)
	}

	_, err2 := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(cfg.AWS.Bucket),
		Key:                  aws.String(tenantName + basePath + "/" + fileId),
		ACL:                  aws.String("private"),
		Body:                 fileStream,
		ContentLength:        aws.Int64(multipartFile.Size),
		ContentType:          aws.String(multipartFile.Header.Get("Content-Type")),
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
