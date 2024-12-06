package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"time"
)

type emailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type EmailService interface {
	SyncEmailsForState(ctx context.Context, importState *postgresEntity.UserEmailImportState) (*postgresEntity.UserEmailImportState, error)
}

func (s *emailService) SyncEmailsForState(ctx context.Context, importState *postgresEntity.UserEmailImportState) (*postgresEntity.UserEmailImportState, error) {
	countEmailsExists := int64(0)

	var externalSystem string
	var rawEmails []*postgresEntity.EmailRawData
	var next string
	var err error

	if importState.Provider == "google" {
		externalSystem = neo4jenum.GMail.String()
		rawEmails, next, err = s.services.CommonServices.GoogleService.ReadEmails(ctx, s.cfg.SyncData.BatchSize, importState)
		if err != nil {
			return nil, fmt.Errorf("unable to read emails from google: %v", err)
		}
	} else if importState.Provider == "azure-ad" {
		externalSystem = neo4jenum.Outlook.String()
		rawEmails, next, err = s.services.CommonServices.AzureService.ReadEmailsFromAzureAd(ctx, importState)
		if err != nil {
			return nil, fmt.Errorf("unable to read emails from azure ad: %v", err)
		}
	}

	for _, emailRawData := range rawEmails {

		emailExists, err := s.services.CommonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(ctx, externalSystem, importState.Tenant, importState.Username, emailRawData.MessageId)
		if err != nil {
			return nil, fmt.Errorf("unable to check if email exists: %v", err)
		}

		//counting emails that are already imported based on the batch size
		//if the job is stopped in the middle of execution and we haven't saved the latest token
		//we are going to lose the history
		if emailExists {

			if importState.State == postgresEntity.REAL_TIME {
				countEmailsExists = countEmailsExists + 1

				if countEmailsExists >= s.cfg.SyncData.BatchSize {
					importState, err = s.services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, "")
					if err != nil {
						return nil, fmt.Errorf("unable to update the gmail page token for username: %v", err)
					}
					return importState, nil
				}
			}

			continue
		} else {
			zeroTime := time.Time{}
			if emailRawData.Sent != zeroTime && importState.StopDate != nil && emailRawData.Sent.Before(*importState.StopDate) {
				importState, err = s.services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, "")
				if err != nil {
					return nil, fmt.Errorf("unable to update the gmail page token for username: %v", err)
				}
				return importState, nil
			}

		}

		jsonContent, err := JSONMarshal(emailRawData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal email content: %v", err)
		}

		err = s.services.CommonServices.PostgresRepositories.RawEmailRepository.Store(ctx, externalSystem, importState.Tenant, importState.Username, emailRawData.ProviderMessageId, emailRawData.MessageId, string(jsonContent), emailRawData.Sent, importState.State)
		if err != nil {
			return nil, fmt.Errorf("failed to store email content: %v", err)
		}
	}

	importState, err = s.services.CommonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, next)
	if err != nil {
		return nil, fmt.Errorf("unable to update the email page token for username: %v", err)
	}

	return importState, nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func NewEmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) EmailService {
	return &emailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
