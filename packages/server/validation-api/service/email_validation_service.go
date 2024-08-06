package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
)

type EmailValidationService interface {
	ValidateEmail(ctx context.Context, email string) (*model.RancherEmailResponseDTO, error)
}

type emailValidationService struct {
	config   *config.Config
	Services *Services
	log      logger.Logger
}

func NewEmailValidationService(config *config.Config, services *Services, log logger.Logger) EmailValidationService {
	return &emailValidationService{
		config:   config,
		Services: services,
		log:      log,
	}
}

func (s *emailValidationService) ValidateEmail(ctx context.Context, email string) (*model.RancherEmailResponseDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidationService.OnEmailCreate")
	defer span.Finish()

	message := map[string]string{"to_email": email}
	bytesRepresentation, _ := json.Marshal(message)

	client := http.Client{}
	// Create the request
	req, err := http.NewRequest("POST", s.config.ReacherConfig.ApiPath, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on creating request: %v", err.Error())
		return nil, err
	}
	req.Header.Set("x-reacher-secret", s.config.ReacherConfig.Secret)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on sending request: %v", err.Error())
		return nil, err
	}
	// Process the response
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error on reading body: %v", err.Error())
		return nil, err
	}
	if resp.StatusCode == 200 {
		d := new(model.RancherEmailResponseDTO)

		err = json.Unmarshal(body, &d)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error on unmarshalling body: %v", err.Error())
			return nil, err
		}
		return d, nil
	} else {
		return nil, errors.New(fmt.Sprintf("validation error: %s", body))
	}
}
