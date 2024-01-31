package fsc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type FileStoreApiService interface {
	UploadSingleMultipartFile(tenantName, basePath string, multipartFileHeader *multipart.FileHeader, span opentracing.Span) (*FileDTO, error)
	UploadSingleFileBytes(tenantName, basePath, fileId, fileName string, fileBytes []byte, span opentracing.Span) (*FileDTO, error)
	DownloadFile(tenantName, fileId string, span opentracing.Span) (*[]byte, error)
}

type fileStoreApiService struct {
	conf *FileStoreApiConfig
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

func (fsas *fileStoreApiService) DownloadFile(tenantName, fileId string, span opentracing.Span) (*[]byte, error) {

	url := fmt.Sprintf("%s/file/%s/download", fsas.conf.ApiPath, fileId)
	log.Printf("DownloadFile: url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("DownloadFile: failed to create new request: %w", err)
	}

	req.Header.Add("X-Openline-API-KEY", fsas.conf.ApiKey)
	req.Header.Add("X-Openline-TENANT", tenantName)

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
	req.Header.Add("X-Openline-API-KEY", conf.ApiKey)
	req.Header.Add("X-Openline-TENANT", tenantName)

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
