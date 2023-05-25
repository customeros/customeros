package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/dto"
	"github.com/sirupsen/logrus"
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
		logrus.Printf("Error on creating request: %v", err.Error())
		return nil, err
	}
	req.Header.Set("x-reacher-secret", s.config.ReacherSecret)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logrus.Printf("Error on sending request: %v", err.Error())
		return nil, err
	}
	// Process the response
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Printf("Error on reading response: %v", err.Error())
		return nil, err
	}
	if resp.StatusCode == 200 {
		d := new(dto.RancherEmailResponseDTO)

		err = json.Unmarshal(body, &d)
		if err != nil {
			logrus.Printf("Error on unmarshal body: %v", err.Error())
			return nil, err
		}
		return d, nil
	} else {
		return nil, errors.New(fmt.Sprintf("validation error: %s", body))
	}
}
