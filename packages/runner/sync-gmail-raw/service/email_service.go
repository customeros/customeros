package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"time"
)

type emailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type EmailService interface {
	FindEmailsForUser(tenant, userId string) ([]*entity.EmailEntity, error)
	SyncEmailsForState(ctx context.Context, importState *postgresEntity.UserEmailImportState) (*postgresEntity.UserEmailImportState, error)
}

func (s *emailService) FindEmailsForUser(tenant, userId string) ([]*entity.EmailEntity, error) {
	ctx := context.Background()

	emails, err := s.repositories.EmailRepository.FindEmailsForUser(ctx, tenant, userId)
	if err != nil {
		return nil, fmt.Errorf("unable to find user by email: %v", err)
	}
	if emails == nil {
		return nil, nil
	}

	emailsEntities := make([]*entity.EmailEntity, len(emails))
	for i, email := range emails {
		emailsEntities[i] = s.mapDbNodeToEmailEntity(*email)
	}

	return emailsEntities, nil
}

func (s *emailService) SyncEmailsForState(ctx context.Context, importState *postgresEntity.UserEmailImportState) (*postgresEntity.UserEmailImportState, error) {
	countEmailsExists := int64(0)

	var externalSystem string
	var rawEmails []*postgresEntity.EmailRawData
	var next string
	var err error

	if importState.Provider == "google" {
		externalSystem = "gmail"
		rawEmails, next, err = s.services.CommonServices.GoogleService.ReadEmails(ctx, s.cfg.SyncData.BatchSize, importState)
		if err != nil {
			return nil, fmt.Errorf("unable to read emails from google: %v", err)
		}
	} else if importState.Provider == "azure-ad" {
		externalSystem = "outlook"
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

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Email:    utils.GetStringPropOrEmpty(props, "email"),
		RawEmail: utils.GetStringPropOrEmpty(props, "rawEmail"),
	}
	return &result
}

func NewEmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) EmailService {
	return &emailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
