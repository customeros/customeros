package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	awsSes "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	GetFilePublicUrl(ctx context.Context, tenantName, id string) (string, error)
}

type fileService struct {
	cfg            *config.Config
	graphqlClient  *graphql.Client
	log            logger.Logger
	commonServices *commonService.Services
}

func NewFileService(cfg *config.Config, commonServices *commonService.Services, graphqlClient *graphql.Client, log logger.Logger) FileService {
	return &fileService{
		cfg:            cfg,
		graphqlClient:  graphqlClient,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *fileService) GetById(ctx context.Context, userEmail, tenantName, id string) (*model.File, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.GetById")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenantName)
	span.LogFields(log.String("fileId", id), log.String("userEmail", userEmail))

	attachment, err := s.commonServices.AttachmentService.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting attachment by id"))
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
		tracing.TraceErr(span, errors.Wrap(err, "Error storing multipart file to temp"))
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error opening file"))
		return nil, err
	}
	defer file.Close()

	headBytes, err := utils.GetFileTypeHeadFromMultipart(file)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting file type head"))
		return nil, err
	}

	fileType := filetype.Unknown
	// check if file type is csv from file name
	if strings.HasSuffix(strings.ToLower(multipartFileHeader.Filename), ".csv") {
		// Detect the MIME type using the file content
		mimeType := http.DetectContentType(headBytes)
		span.LogFields(log.String("mimeType", mimeType))

		// Validate if the detected MIME type is "text/csv"
		if mimeType != "text/csv" && mimeType != "application/octet-stream" {
			err = errors.New("Invalid mime type for CSV")
			tracing.TraceErr(span, errors.Wrap(err, "Unexpected file type"))
			s.log.Error("Unexpected file type")
			//return nil, err
		}
		fileType = types.NewType("csv", "text/csv")
	} else {
		fileType, err = utils.GetFileType(headBytes)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error getting file type"))
			return nil, err
		}

		if fileType == filetype.Unknown {
			err = errors.New("Unknown file type")
			tracing.TraceErr(span, errors.Wrap(err, "Unknown file type"))
			s.log.Error("Unknown multipartFile type")
			return nil, err
		}
	}

	attachmentEntity := neo4jentity.AttachmentEntity{
		Id:        fileId,
		FileName:  multipartFileHeader.Filename,
		MimeType:  multipartFileHeader.Header.Get(http.CanonicalHeaderKey("Content-Type")),
		Size:      multipartFileHeader.Size,
		AppSource: constants.AppSourceFileStoreApi,
	}

	if s.cfg.Service.CloudflareImageUploadApiKey != "" && s.cfg.Service.CloudflareImageUploadAccountId != "" && s.cfg.Service.CloudflareImageUploadSignKey != "" &&
		cdnUpload && (fileType.Extension == "gif" || fileType.Extension == "png" || fileType.Extension == "jpg" || fileType.Extension == "jpeg") {

		cloudflareApi, err := cloudflare.NewWithAPIToken(s.cfg.Service.CloudflareImageUploadApiKey)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error creating cloudflare api"))
			return nil, err
		}

		open, err := os.Open(fileName)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error opening file"))
			return nil, err
		}

		readCloser := io.NopCloser(open)

		uploadedFileToCdn, err := cloudflareApi.UploadImage(context.Background(), cloudflare.AccountIdentifier(s.cfg.Service.CloudflareImageUploadAccountId), cloudflare.UploadImageParams{
			File:              readCloser,
			Name:              fileId,
			RequireSignedURLs: true,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error uploading file to cdn"))
			return nil, err
		}

		attachmentEntity.CdnUrl = generateSignedURL(uploadedFileToCdn.Variants[0], s.cfg.Service.CloudflareImageUploadSignKey)
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating aws session"))
		s.log.Fatal(err)
	}

	if basePath == "" {
		basePath = "/GLOBAL"
	}

	extension := utils.FirstNotEmptyString(filepath.Ext(multipartFileHeader.Filename), fileType.Extension)
	// remove starting dot if exists
	if strings.HasPrefix(extension, ".") {
		extension = extension[1:]
	}
	err = uploadFileToS3(ctx, s.cfg, session, tenantName, basePath, fileId+"."+extension, multipartFileHeader)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error uploading file to s3"))
		s.log.Fatal(err)
	}
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error uploading file to s3"))
		return nil, err
	}

	attachmentEntity.BasePath = basePath

	created, err := s.commonServices.AttachmentService.Create(ctx, &attachmentEntity)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating attachment"))
		return nil, err
	}

	return mapper.MapAttachmentResponseToFileEntity(created), nil
}

func (s *fileService) DownloadSingleFile(ctx context.Context, userEmail, tenantName, id string, ginContext *gin.Context, inline bool) (*model.File, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.DownloadSingleFile")
	defer span.Finish()
	tracing.TagTenant(span, tenantName)
	span.LogFields(log.String("userEmail", userEmail), log.String("fileId", id), log.Bool("inline", inline))

	attachment, err := s.commonServices.AttachmentService.GetById(ctx, id)
	byId := mapper.MapAttachmentResponseToFileEntity(attachment)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting attachment by id"))
		return nil, err
	}
	tracing.LogObjectAsJson(span, "attachment", attachment)

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating aws session"))
		log.Error(err)
		ginContext.AbortWithError(http.StatusInternalServerError, err)
	}

	ginContext.Header("Accept-Ranges", "bytes")

	svc := s3.New(session)

	extension := filepath.Ext(attachment.FileName)
	if extension == "" {
		tracing.TraceErr(span, errors.New("No file extension found"))
		fmt.Println("No file extension found.")
	} else {
		extension = extension[1:]
		fmt.Println("File Extension:", extension)
	}

	// Get the object metadata to determine the file size and ETag
	bucket := s.cfg.AWS.Bucket
	key := tenantName + byId.BasePath + "/" + attachment.Id + "." + extension
	span.LogFields(log.String("bucket", bucket), log.String("key", key))
	respHead, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting object metadata"))
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
		Key:    aws.String(tenantName + byId.BasePath + "/" + attachment.Id + "." + extension),
		Range:  aws.String("bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)),
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting object"))
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

	attachment, err := s.commonServices.AttachmentService.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting attachment by id"))
		return nil, err
	}

	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating aws session"))
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
			Key:    aws.String(attachment.Id),
		})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error downloading file"))
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
		tracing.TraceErr(span, errors.Wrap(err, "Error opening file"))
		return fmt.Errorf("uploadFileToS3: %w", err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(cfg.AWS.Bucket),
		Key:           aws.String(tenantName + basePath + "/" + fileId),
		ACL:           aws.String("private"),
		ContentLength: aws.Int64(0),
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error putting object"))
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
		tracing.TraceErr(span, errors.Wrap(err, "Error creating temp file"))
		return "", err
	}
	src, err := multipartFileHeader.Open()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error opening multipart file"))
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(file, src)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error copying file"))
		return "", err
	}

	err = file.Close()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error closing file"))
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

func (s *fileService) GetFilePublicUrl(ctx context.Context, tenant, fileId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileService.GetFilePublicUrl")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("fileId", fileId)

	attachmentDbNode, err := s.commonServices.Neo4jRepositories.AttachmentReadRepository.GetById(ctx, tenant, fileId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting attachment by id"))
		return "", err
	}
	attachmentEntity := neo4jmapper.MapDbNodeToAttachmentEntity(attachmentDbNode)
	if attachmentEntity.PublicUrl != "" && attachmentEntity.PublicUrlExpiresAt != nil && attachmentEntity.PublicUrlExpiresAt.After(utils.Now()) {
		return attachmentEntity.PublicUrl, nil
	}

	// generate public url and store it
	extension := filepath.Ext(attachmentEntity.FileName)
	if extension == "" {
		span.LogKV("message", "No file extension found.")
	} else {
		span.LogKV("extension", extension)
		extension = extension[1:]
	}
	awsBucket := aws.String(s.cfg.AWS.Bucket)
	awsS3ObjectKey := aws.String(tenant + attachmentEntity.BasePath + "/" + attachmentEntity.Id + "." + extension)

	// Initialize the S3 service client
	session, err := awsSes.NewSession(&aws.Config{Region: aws.String(s.cfg.AWS.Region)})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating aws session"))
		s.log.Fatal(err)
	}
	svc := s3.New(session)

	// Create a request to get a presigned URL for the object
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: awsBucket,
		Key:    awsS3ObjectKey,
	})

	// Presign the URL with the expiration time
	expiration := 7 * 24 * time.Hour // 7 days
	publicUrl, err := req.Presign(expiration)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error presigning request"))
		return "", fmt.Errorf("failed to presign request: %v", err)
	}

	// set public url and expiration time
	err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateStringProperty(ctx, nil, tenant, commonmodel.ATTACHMENT.Neo4jLabel(), fileId, string(neo4jentity.AttachmentPropertyPublicUrl), publicUrl)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error updating attachment public url"))
	}

	// set expiration time to 6 days and 23 hours
	expiredAt := utils.Now().Add(7 * 24 * time.Hour).Add(-1 * time.Hour)
	err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, commonmodel.ATTACHMENT.Neo4jLabel(), fileId, string(neo4jentity.AttachmentPropertyPublicUrlExpiresAt), &expiredAt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error updating attachment public url expiration time"))
	}

	return publicUrl, nil
}
