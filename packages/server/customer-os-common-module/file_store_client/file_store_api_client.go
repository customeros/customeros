package fsc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

const (
	TenantHeader = "X-Openline-TENANT"
	ApiKeyHeader = "X-Openline-API-KEY"
)

type FileStoreApiService interface {
	GetFile(tenantName, fileId string, span opentracing.Span) (*FileDTO, *[]byte, error)
	GetFileMetadata(tenantName, fileId string, span opentracing.Span) (*FileDTO, error)
	GetFileBytes(tenantName, fileId string, span opentracing.Span) (*[]byte, error)
	GetFilePublicUrl(ctx context.Context, tenant, fileId string) (string, error)

	UploadSingleMultipartFile(tenantName, basePath string, multipartFileHeader *multipart.FileHeader, span opentracing.Span) (*FileDTO, error)
	UploadSingleFileBytes(tenantName, basePath, fileId, fileName string, fileBytes []byte, span opentracing.Span) (*FileDTO, error)
}

type fileStoreApiService struct {
	conf *FileStoreApiConfig
}

func (fsas *fileStoreApiService) GetFile(tenantName, fileId string, span opentracing.Span) (*FileDTO, *[]byte, error) {
	fileMetadata, err := fsas.GetFileMetadata(tenantName, fileId, span)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetFileMetadata")
	}

	fileBytes, err := fsas.GetFileBytes(tenantName, fileId, span)
	if err != nil {
		return nil, nil, errors.Wrap(err, "GetFileBytes")
	}

	return fileMetadata, fileBytes, nil
}

func (fsas *fileStoreApiService) GetFileMetadata(tenantName, fileId string, span opentracing.Span) (*FileDTO, error) {

	url := fmt.Sprintf("%s/file/%s", fsas.conf.ApiPath, fileId)
	log.Printf("DownloadFile: url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetFile: failed to create new request: %w", err)
	}

	req.Header.Add(ApiKeyHeader, fsas.conf.ApiKey)
	req.Header.Add(TenantHeader, tenantName)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ext.LogError(span, err)
		return nil, fmt.Errorf("GetFile: failed to perform request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var fileResponse FileDTO
		if err := json.NewDecoder(resp.Body).Decode(&fileResponse); err != nil {
			ext.LogError(span, err)
			return nil, fmt.Errorf("GetFile: failed to decode response: %w", err)
		}
		return &fileResponse, nil
	} else {
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			ext.LogError(span, err)
			return nil, err
		}

		err = fmt.Errorf("Got error from File Store API: Status: %d Response: %s", resp.StatusCode, responseBody.String())
		ext.LogError(span, err)
		return nil, err
	}
}

func (fsas *fileStoreApiService) GetFileBytes(tenantName, fileId string, span opentracing.Span) (*[]byte, error) {

	url := fmt.Sprintf("%s/file/%s/download", fsas.conf.ApiPath, fileId)
	log.Printf("DownloadFile: url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("DownloadFile: failed to create new request: %w", err)
	}

	req.Header.Add(ApiKeyHeader, fsas.conf.ApiKey)
	req.Header.Add(TenantHeader, tenantName)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ext.LogError(span, err)
		return nil, fmt.Errorf("DownloadFile: failed to perform request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var fileResponse []byte
		if fileResponse, err = io.ReadAll(resp.Body); err != nil {
			ext.LogError(span, err)
			return nil, fmt.Errorf("DownloadFile: failed to read response: %w", err)
		}
		return &fileResponse, nil
	} else {
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			ext.LogError(span, err)
			return nil, err
		}

		err = fmt.Errorf("Got error from File Store API: Status: %d Response: %s", resp.StatusCode, responseBody.String())
		ext.LogError(span, err)
		return nil, err
	}
}

func (fsas *fileStoreApiService) UploadSingleMultipartFile(tenantName, basePath string, multipartFileHeader *multipart.FileHeader, span opentracing.Span) (*FileDTO, error) {
	file, err := multipartFileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("UploadSingleMultipartFile: failed to open multipart file: %w", err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("UploadSingleMultipartFile: failed to read multipart file: %w", err)
	}

	return sendRequest(fsas.conf, tenantName, basePath, "", multipartFileHeader.Filename, fileBytes, span)
}

func (fsas *fileStoreApiService) UploadSingleFileBytes(tenantName, basePath, fileId, fileName string, fileBytes []byte, span opentracing.Span) (*FileDTO, error) {
	return sendRequest(fsas.conf, tenantName, basePath, fileId, fileName, fileBytes, span)
}

func sendRequest(conf *FileStoreApiConfig, tenantName, basePath, fileId, fileName string, fileBytes []byte, span opentracing.Span) (*FileDTO, error) {
	// Create a new buffer to store the request body
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field for the file
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}

	// Copy the file content (bytes) to the form file field
	_, err = fileWriter.Write(fileBytes)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}

	err = addMultipartValue(writer, basePath, "basePath")
	if err != nil {
		ext.LogError(span, err)
		return nil, errors.Wrap(err, "addMultipartValue basePath")
	}
	err = addMultipartValue(writer, fileId, "fileId")
	if err != nil {
		ext.LogError(span, err)
		return nil, errors.Wrap(err, "addMultipartValue fileId")
	}

	// Close the multipart writer to finalize the request body
	writer.Close()

	url := fmt.Sprintf("%s/file", conf.ApiPath)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		ext.LogError(span, err)
		return nil, fmt.Errorf("UploadSingleMultipartFile: failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("accept", "application/json")
	req.Header.Add(ApiKeyHeader, conf.ApiKey)
	req.Header.Add(TenantHeader, tenantName)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ext.LogError(span, err)
		return nil, fmt.Errorf("UploadSingleMultipartFile: failed to perform request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var fileResponse FileDTO
		if err := json.NewDecoder(resp.Body).Decode(&fileResponse); err != nil {
			ext.LogError(span, err)
			return nil, fmt.Errorf("UploadSingleMultipartFile: failed to decode response: %w", err)
		}
		return &fileResponse, nil
	} else {
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			ext.LogError(span, err)
			return nil, err
		}

		err = fmt.Errorf("Got error from File Store API: Status: %d Response: %s", resp.StatusCode, responseBody.String())
		ext.LogError(span, err)
		return nil, err
	}
}

func addMultipartValue(writer *multipart.Writer, value string, partName string) error {
	part, err := writer.CreateFormField(partName)
	if err != nil {
		return errors.Wrap(err, "writer.CreateFormFile")
	}
	_, err = part.Write([]byte(value))
	if err != nil {
		return errors.Wrap(err, "part.Write")
	}
	return nil
}

func NewFileStoreApiService(conf *FileStoreApiConfig) *fileStoreApiService {
	return &fileStoreApiService{
		conf: conf,
	}
}

func (fsas *fileStoreApiService) GetFilePublicUrl(ctx context.Context, tenant, fileId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FileStoreApiService.GetFilePublicUrl")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogKV("fileId", fileId)

	fileStoreUrl := fmt.Sprintf("%s/file/%s/public-url", fsas.conf.ApiPath, fileId)
	req, err := http.NewRequest("GET", fileStoreUrl, nil)
	if err != nil {
		return "", fmt.Errorf("GetFilePublicUrl: failed to create new request: %s", err.Error())
	}

	req.Header.Add(ApiKeyHeader, fsas.conf.ApiKey)
	req.Header.Add(TenantHeader, tenant)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ext.LogError(span, err)
		return "", fmt.Errorf("GetFilePublicUrl: failed to perform request: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Create a struct to capture the JSON response
		var response struct {
			PublicUrl string `json:"publicUrl"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			ext.LogError(span, err)
			return "", fmt.Errorf("GetFilePublicUrl: failed to decode response: %w", err)
		}

		// Return the public URL
		return response.PublicUrl, nil
	} else {
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			ext.LogError(span, err)
			return "", err
		}

		err = fmt.Errorf("Got error from File Store API: Status: %d Response: %s", resp.StatusCode, responseBody.String())
		ext.LogError(span, err)
		return "", err
	}
}
