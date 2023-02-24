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
	ValidateEmail(ctx context.Context, email string) (bool, error)
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

func (s *emailValidationService) ValidateEmail(ctx context.Context, email string) (bool, error) {
	message := map[string]string{"to_email": email}
	bytesRepresentation, _ := json.Marshal(message)

	resp, _ := http.Post(s.config.ReacherApiPath, "application/json", bytes.NewBuffer(bytesRepresentation))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	d := new(dto.RancherEmailResponseDTO)

	err = json.Unmarshal([]byte(body), &d)
	if err != nil {
		return false, err
	}

	return true, nil
}
