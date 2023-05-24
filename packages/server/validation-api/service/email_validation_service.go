package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/dto"
	"io/ioutil"
	"net/http"
)

type EmailValidationService interface {
	ValidateEmail(ctx context.Context, email string) (*dto.RancherEmailResponseDTO, error)
}

type emailValidationService struct {
	config   *config.Config
	Services *Services
}

func NewEmailValidationService(config *config.Config, services *Services) EmailValidationService {
	return &emailValidationService{
		config:   config,
		Services: services,
	}
}

func (s *emailValidationService) ValidateEmail(ctx context.Context, email string) (*dto.RancherEmailResponseDTO, error) {
	message := map[string]string{"to_email": email}
	bytesRepresentation, _ := json.Marshal(message)

	client := http.Client{}
	// Create the request
	req, err := http.NewRequest("POST", s.config.ReacherApiPath, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-reacher-secret", s.config.ReacherSecret)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Process the response
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := new(dto.RancherEmailResponseDTO)

	err = json.Unmarshal([]byte(body), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
