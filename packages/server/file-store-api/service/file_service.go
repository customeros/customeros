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
	graph_model "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
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
	GetById(ctx context.Context, userEmail, tenantName string, id string) (*model.File, error)
	UploadSingleFile(ctx context.Context, userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader, cdnUpload bool) (*model.File, error)
	DownloadSingleFile(ctx context.Context, userEmail, tenantName, id string, context *gin.Context, inline bool) (*model.File, error)
	Base64Image(ctx context.Context, userEmail, tenantName string, id string) (*string, error)
}

type fileService struct {
	cfg           *config.Config
	graphqlClient *graphql.Client
	log           logger.Logger
}

func NewFileService(cfg *config.Config, graphqlClient *graphql.Client, log logger.Logger) FileService {
	return &fileService{
		cfg:           cfg,
		graphqlClient: graphqlClient,
		log:           log,
	}
}

func (s *fileService) GetById(ctx context.Context, userEmail, tenantName, id string) (*model.File, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.GetById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("fileId", id), log.String("userEmail", userEmail))

	attachment, err := s.getCosAttachmentById(ctx, userEmail, tenantName, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(attachment), nil
}

func (s *fileService) UploadSingleFile(ctx context.Context, userEmail, tenantName, basePath, fileId string, multipartFileHeader *multipart.FileHeader, cdnUpload bool) (*model.File, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.UploadSingleFile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("userEmail", userEmail), log.String("basePath", basePath), log.String("fileId", fileId))
	if multipartFileHeader != nil {
		span.LogFields(log.String("fileName", multipartFileHeader.Filename), log.Int64("size", multipartFileHeader.Size))
	}

	if fileId == "" {
		fileId = uuid.New().String()
	}

	fileName, err := storeMultipartFileToTemp(ctx, fileId, multipartFileHeader)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer file.Close()

	headBytes, err := utils.GetFileTypeHeadFromMultipart(file)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	fileType, err := utils.GetFileType(headBytes)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if fileType == filetype.Unknown {
		err = errors.New("Unknown file type")
		tracing.TraceErr(span, err)
		s.log.Error("Unknown multipartFile type")
		return nil, err
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
			tracing.TraceErr(span, err)
			return nil, err
		}

		open, err := os.Open(fileName)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		readCloser := io.NopCloser(open)

		uploadedFileToCdn, err := cloudflareApi.UploadImage(context.Background(), cloudflare.AccountIdentifier(s.cfg.Service.CloudflareImageUploadAccountId), cloudflare.UploadImageParams{
			File:              readCloser,
			Name:              fileId,
			RequireSignedURLs: true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		graphqlRequest.Var("cdnUrl", generateSignedURL(uploadedFileToCdn.Variants[0], s.cfg.Service.CloudflareImageUploadSignKey))
	} else {
		graphqlRequest.Var("cdnUrl", "")
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Fatal(err)
	}

	if basePath == "" {
		basePath = "/GLOBAL"
	}

	err = uploadFileToS3(ctx, s.cfg, session, tenantName, basePath, fileId+"."+fileType.Extension, multipartFileHeader)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Fatal(err)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	graphqlRequest.Var("basePath", basePath)

	err = s.addHeadersToGraphRequest(graphqlRequest, tenantName, userEmail)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.AttachmentCreateResponse
	tracing.InjectSpanContextIntoGraphQLRequest(graphqlRequest, span)
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(&graphqlResponse.Attachment), nil
}

func (s *fileService) DownloadSingleFile(ctx context.Context, userEmail, tenantName, id string, ginContext *gin.Context, inline bool) (*model.File, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.DownloadSingleFile")
	defer span.Finish()
	tracing.TagTenant(span, tenantName)
	span.LogFields(log.String("userEmail", userEmail), log.String("fileId", id), log.Bool("inline", inline))

	attachment, err := s.getCosAttachmentById(ctx, userEmail, tenantName, id)
	byId := mapper.MapAttachmentResponseToFileEntity(attachment)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, err)
		log.Error(err)
		ginContext.AbortWithError(http.StatusInternalServerError, err)
	}

	ginContext.Header("Accept-Ranges", "bytes")

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
		Key:    aws.String(tenantName + byId.BasePath + "/" + attachment.ID + "." + extension),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		ginContext.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}

	// Get the ETag header value
	eTag := aws.StringValue(respHead.ETag)

	// Parse the range header
	rangeHeader := ginContext.GetHeader("Range")
	var start, end int64
	if rangeHeader != "" {
		s.log.Infof("Range header: %s", rangeHeader)
		s.log.Infof("Content Length: %d", *respHead.ContentLength)

		rangeParts := strings.Split(rangeHeader, "=")[1]
		rangeBytes := strings.Split(rangeParts, "-")
		start, _ = strconv.ParseInt(rangeBytes[0], 10, 64)
		s.log.Infof("rangeBytes %v", rangeBytes)
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
	ginContext.Header("Content-Length", strconv.FormatInt(end-start+1, 10))

	// Set the content range header to indicate the range of bytes being served
	ginContext.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(*respHead.ContentLength, 10))

	// If the ETag matches, send a 304 Not Modified response and exit early
	if match := ginContext.GetHeader("If-Range"); match != "" && match != eTag {
		ginContext.Status(http.StatusRequestedRangeNotSatisfiable)
		return byId, nil
	}

	if !inline {
		ginContext.Header("Content-Disposition", "attachment; filename="+byId.FileName)
	} else {
		ginContext.Header("Content-Disposition", "inline; filename="+byId.FileName)
	}
	ginContext.Header("Content-Type", fmt.Sprintf("%s", byId.MimeType))
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AWS.Bucket),
		Key:    aws.String(tenantName + byId.BasePath + "/" + attachment.ID + "." + extension),
		Range:  aws.String("bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		// Handle error
		s.log.Errorf("Error getting object: %v", err)
		ginContext.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Serve the file contents
	io.Copy(ginContext.Writer, resp.Body)
	return byId, nil
}

func (s *fileService) Base64Image(ctx context.Context, userEmail, tenantName string, id string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.Base64Image")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("userEmail", userEmail), log.String("fileId", id))

	attachment, err := s.getCosAttachmentById(ctx, userEmail, tenantName, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err)
	}

	if attachment.Size > s.cfg.MaxFileSizeMB*1024*1024 {
		return nil, errors.New("file is too big for base64 encoding")
	}

	downloader := s3manager.NewDownloader(session)

	fileBytes := make([]byte, attachment.Size)
	_, err = downloader.Download(aws.NewWriteAtBuffer(fileBytes),
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AWS.Bucket),
			Key:    aws.String(attachment.ID),
		})
	if err != nil {
		tracing.TraceErr(span, err)
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

func (s *fileService) getCosAttachmentById(ctx context.Context, userEmail, tenantName, id string) (*graph_model.Attachment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.getCosAttachmentById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("userEmail", userEmail), log.String("fileId", id))

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.AttachmentResponse
	tracing.InjectSpanContextIntoGraphQLRequest(graphqlRequest, span)
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &graphqlResponse.Attachment, nil
}

func uploadFileToS3(ctx context.Context, cfg *config.Config, session *awsSes.Session, tenantName, basePath, fileId string, multipartFile *multipart.FileHeader) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.uploadFileToS3")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("basePath", basePath), log.String("fileId", fileId))
	if multipartFile != nil {
		span.LogFields(log.String("fileName", multipartFile.Filename), log.Int64("size", multipartFile.Size))
	}

	fileStream, err := multipartFile.Open()
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("uploadFileToS3: %w", err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(cfg.AWS.Bucket),
		Key:           aws.String(tenantName + basePath + "/" + fileId),
		ACL:           aws.String("private"),
		ContentLength: aws.Int64(0),
	})
	if err != nil {
		tracing.TraceErr(span, err)
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

func (s *fileService) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Second)
	return ctx, cancel, nil
}

func storeMultipartFileToTemp(ctx context.Context, fileId string, multipartFileHeader *multipart.FileHeader) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.storeMultipartFileToTemp")
	defer span.Finish()
	span.LogFields(log.String("fileId", fileId))

	file, err := os.CreateTemp("", fileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	src, err := multipartFileHeader.Open()
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(file, src)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	err = file.Close()
	if err != nil {
		tracing.TraceErr(span, err)
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
