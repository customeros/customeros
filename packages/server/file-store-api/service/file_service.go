package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SIGN_TOKEN_EXPIRATION = 60 * 60 * 24 * 365 * 99 // 99 years
)

type FileService interface {
	GetById(userEmail, tenantName string, id string) (*model.File, error)
	UploadSingleFile(userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader, cdnUpload bool) (*model.File, error)
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

func (s *fileService) UploadSingleFile(userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader, cdnUpload bool) (*model.File, error) {
	if fileId == "" {
		fileId = uuid.New().String()
	}

	fileName, err := storeMultipartFileToTemp(fileId, multipartFileHeader)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	headBytes, err := utils.GetFileTypeHeadFromMultipart(file)
	if err != nil {
		return nil, err
	}

	fileType, err := utils.GetFileType(headBytes)
	if err != nil {
		return nil, err
	}

	if fileType == filetype.Unknown {
		fmt.Println("Unknown multipartFile type")
		return nil, errors.New("Unknown multipartFile type")
	}

	graphqlRequest := graphql.NewRequest(
		`mutation AttachmentCreate($id: ID, $cdnUrl: String!, $basePath: String!, $mimeType: String!, $size: Int64!, $fileName: String!, $appSource: String!) {
			attachment_Create(input: {
				id: $id	
				cdnUrl: $cdnUrl	
				basePath: $basePath	
				mimeType: $mimeType	
				fileName: $fileName
				size: $size
				appSource: $appSource
			}) {
				id
				cdnUrl
				basePath
				mimeType
				fileName
				size
			}
		}`)

	graphqlRequest.Var("id", fileId)
	graphqlRequest.Var("mimeType", multipartFileHeader.Header.Get(http.CanonicalHeaderKey("Content-Type")))
	graphqlRequest.Var("fileName", multipartFileHeader.Filename)
	graphqlRequest.Var("size", multipartFileHeader.Size)
	graphqlRequest.Var("appSource", "file-store-api")

	if s.cfg.Service.CloudflareImageUploadApiKey != "" && s.cfg.Service.CloudflareImageUploadAccountId != "" && s.cfg.Service.CloudflareImageUploadSignKey != "" &&
		cdnUpload && (fileType.Extension == "gif" || fileType.Extension == "png" || fileType.Extension == "jpg" || fileType.Extension == "jpeg") {

		cloudflareApi, err := cloudflare.NewWithAPIToken(s.cfg.Service.CloudflareImageUploadApiKey)
		if err != nil {
			return nil, err
		}

		open, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}

		readCloser := io.NopCloser(open)

		uploadedFileToCdn, err := cloudflareApi.UploadImage(context.Background(), cloudflare.AccountIdentifier(s.cfg.Service.CloudflareImageUploadAccountId), cloudflare.UploadImageParams{
			File:              readCloser,
			Name:              fileId,
			RequireSignedURLs: true,
		})
		if err != nil {
			return nil, err
		}

		graphqlRequest.Var("cdnUrl", generateSignedURL(uploadedFileToCdn.Variants[0], s.cfg.Service.CloudflareImageUploadSignKey))
	} else {
		graphqlRequest.Var("cdnUrl", "")
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		log.Fatal(err)
	}

	if basePath == "" {
		basePath = "/GLOBAL"
	}

	err = uploadFileToS3(s.cfg, session, tenantName, basePath, fileId+"."+fileType.Extension, multipartFileHeader)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	graphqlRequest.Var("basePath", basePath)

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
				cdnUrl
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

func storeMultipartFileToTemp(fileId string, multipartFileHeader *multipart.FileHeader) (string, error) {
	file, err := os.CreateTemp("", fileId)
	if err != nil {
		return "", err
	}
	src, err := multipartFileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(file, src)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func generateSignedURL(imageDeliveryURL, key string) string {
	// Parse the URL
	parsedURL, err := url.Parse(imageDeliveryURL)
	if err != nil {
		return fmt.Sprintf("Error parsing URL: %v", err)
	}

	// Attach the expiration value to the URL
	expiry := time.Now().Unix() + SIGN_TOKEN_EXPIRATION
	q := parsedURL.Query()
	q.Set("exp", fmt.Sprintf("%d", expiry))
	parsedURL.RawQuery = q.Encode()

	// Extract path and query from the URL
	stringToSign := parsedURL.Path + "?" + parsedURL.RawQuery

	// Generate the signature
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(stringToSign))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Attach the signature to the URL
	q.Set("sig", signature)
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String()
}
