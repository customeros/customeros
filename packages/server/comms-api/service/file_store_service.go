package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	fileDto "github.com/openline-ai/openline-customer-os/packages/server/file-store-api/dto"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type fileStoreApiService struct {
	conf *c.Config
}

type FileStoreApiService interface {
	UploadSingleFile(tenantName string, multipartFileHeader *multipart.FileHeader) (*fileDto.File, error)
}

func (fsas *fileStoreApiService) UploadSingleFile(tenantName string, multipartFileHeader *multipart.FileHeader) (*fileDto.File, error) {

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", multipartFileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("UploadSingleFile: failed to create form file: %w", err)
	}

	fileStream, err := multipartFileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("UploadSingleFile: failed to open multipart file header: %w", err)
	}
	defer fileStream.Close()

	if _, err = io.Copy(fw, fileStream); err != nil {
		return nil, fmt.Errorf("UploadSingleFile: failed to copy multipart file header: %w", err)
	}
	w.Close()

	url := fmt.Sprintf("%s/file", fsas.conf.Service.FileStoreAPI)
	log.Printf("UploadSingleFile: url: %s", url)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("UploadSingleFile: failed to create new request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Openline-API-KEY", fsas.conf.Service.FileStoreAPIKey)
	req.Header.Add("X-Openline-Tenant", tenantName)

	req.Header.Set("Content-Type", w.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("UploadSingleFile: failed to perform request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var fileResponse fileDto.File
		if err := json.NewDecoder(resp.Body).Decode(&fileResponse); err != nil {
			return nil, fmt.Errorf("UploadSingleFile: failed to decode response: %w", err)
		}
		return &fileResponse, nil
	} else {
		var responseBody bytes.Buffer
		_, err = io.Copy(&responseBody, resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return nil, err
		}

		err = fmt.Errorf("Got error from File Store API: Status: %d Response: %s", resp.StatusCode, responseBody.String())
		return nil, err
	}

}

func NewFileStoreApiService(conf *c.Config) *fileStoreApiService {
	return &fileStoreApiService{
		conf: conf,
	}
}
