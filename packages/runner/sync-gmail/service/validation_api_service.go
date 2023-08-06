package service

import (
	"bytes"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"net/http"
	"strings"
)

type EmailValidationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailValidationResponse struct {
	Error       string `json:"error"`
	IsReachable string `json:"isReachable"`
}

type ValidationApiService interface {
	ValidateEmail(email string) (EmailValidationResponse, error)
}

type validationApiService struct {
	cfg *config.Config
}

func (s *validationApiService) ValidateEmail(email string) (EmailValidationResponse, error) {
	request := EmailValidationRequest{
		Email: strings.TrimSpace(email),
	}

	evJSON, err := json.Marshal(request)
	if err != nil {
		return EmailValidationResponse{}, err
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", s.cfg.ValidationApi.Path+"/validateEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		return EmailValidationResponse{}, err
	}
	// Set the request headers
	req.Header.Set(common_module.ApiKeyHeader, s.cfg.ValidationApi.Key)
	req.Header.Set(common_module.TenantHeader, "openline")

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return EmailValidationResponse{}, err
	}
	defer response.Body.Close()
	var result EmailValidationResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return EmailValidationResponse{}, err
	}

	return result, nil
}

func NewValidationApiService(cfg *config.Config) ValidationApiService {
	return &validationApiService{
		cfg: cfg,
	}
}
