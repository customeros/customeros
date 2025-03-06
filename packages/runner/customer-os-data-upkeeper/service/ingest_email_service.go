package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	commonservice "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_listeners"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type IngestEmailService interface {
	SyncEmailsInState(state postgresEntity.EmailImportState)
	SendIngestedEmailsToAgents()
}

type ingestEmailService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.CommonServices
}

func NewIngestEmailService(cfg *config.Config, log logger.Logger, commonServices *commonservice.CommonServices) IngestEmailService {
	return &ingestEmailService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *ingestEmailService) SyncEmailsInState(state postgresEntity.EmailImportState) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit
	span, ctx := tracing.StartTracerSpan(ctx, "IngestEmailService.SyncEmailsInState")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	span.LogKV("state", state)

	runImportFor := []map[string]interface{}{}

	//TODO base logic on email keeper agents to get oauth tokens

	agents, err := s.commonServices.PostgresRepositories.AgentRepository.GetAllAgentsByTypesCrossTenant(ctx, []commonenum.AgentType{commonenum.AgentEmailKeeper})
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	for _, agent := range agents {
		if agent.IsActive == false {
			continue
		}

		for _, listener := range agent.Listeners {
			var listenerConfig agent_listeners.NewEmailConfig
			err := json.Unmarshal(listener.Config, &listenerConfig)
			if err != nil {
				tracing.TraceErr(span, err)
				continue
			}

			agentError := ""

			for _, email := range listenerConfig.Emails.Value {
				emailAddress := email.(map[string]interface{})["email"].(string)
				provider := email.(map[string]interface{})["provider"].(string)

				oAuthTokenEntities, err := s.commonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, agent.Tenant, provider, emailAddress)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				if oAuthTokenEntities == nil {
					//todo mark error on agent
					agentError += fmt.Sprintf("No OAuth token found for provider: %s and email: %s", provider, emailAddress)
				} else {
					runImportFor = append(runImportFor, map[string]interface{}{
						"agentId":  agent.ID,
						"tenant":   agent.Tenant,
						"provider": provider,
						"email":    emailAddress,
					})
				}
			}

			if agentError != "" {
				// TODO set error on agent / listener ??
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(runImportFor))

	for _, importFor := range runImportFor {
		agentId := importFor["agentId"].(string)
		tenant := importFor["tenant"].(string)
		provider := importFor["provider"].(string)
		email := importFor["email"].(string)

		go func(agentId, tenant, provider, email string) {
			defer wg.Done()

			oAuthTokenEntities, err := s.commonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant, provider, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			s.syncEmailsForEmailAddress(ctx, oAuthTokenEntities, state)
		}(agentId, tenant, provider, email)
	}

	wg.Wait()
}

func (s *ingestEmailService) SendIngestedEmailsToAgents() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit
	span, ctx := tracing.StartTracerSpan(ctx, "IngestEmailService.SendIngestedEmailsToAgents")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	distinctUsersForImport, err := s.commonServices.PostgresRepositories.IngestEmailMessageRepository.GetDistinctUsersForPendingMessages(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if len(distinctUsersForImport) == 0 {
		span.LogKV("result", "no distinct users for import with pending messages")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(distinctUsersForImport))

	for _, dt := range distinctUsersForImport {
		go func(distinctUser postgresEntity.IngestEmailMessage) {
			defer wg.Done()
			localCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant: distinctUser.Tenant,
			})

			// Check if agent is enabled to process message
			agentListener, err := s.commonServices.AgentService.GetListener(dto.NewEmail{}.ListenerEvent())
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
			agentTypes := agentListener.ExecutingAgents()
			agents, err := s.commonServices.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(localCtx, agentTypes)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}
			if len(agents) == 0 {
				span.LogKV("message", "no active agents found for tenant: %s", distinctUser.Tenant)
				return
			}

			ingestEmailMessages, err := s.commonServices.PostgresRepositories.IngestEmailMessageRepository.GetEmailsForUserForSync(ctx, distinctUser.Tenant, distinctUser.Username)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			for _, ingestEmailMessage := range ingestEmailMessages {
				err = s.commonServices.Events.Publisher.PublishFanoutEvent(localCtx, ingestEmailMessage.Id, model.INGEST_EMAIL_MESSAGE, dto.NewEmail{})
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				err := s.commonServices.PostgresRepositories.IngestEmailMessageRepository.UpdateState(localCtx, ingestEmailMessage.Id, postgresEntity.IngestEmailMessageStateSentToAgent)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

			}

			span.LogKV("message", "sent emails to agents for user: %s in tenant: %s", distinctUser.Tenant, distinctUser.Username)
		}(dt)
	}

	wg.Wait()
	span.LogKV("message", "sent emails to agents for all users")
}

func (s *ingestEmailService) syncEmailsForEmailAddress(ctx context.Context, authTokenEntity *postgresEntity.OAuthTokenEntity, state postgresEntity.EmailImportState) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailService.syncEmailsForEmailAddress - "+authTokenEntity.TenantName+" - "+authTokenEntity.EmailAddress)
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	if authTokenEntity == nil {
		span.LogKV("message", "no oauth token found for tenant: %s and username: %s", authTokenEntity.TenantName, authTokenEntity.EmailAddress)
		return
	}

	tenant := authTokenEntity.TenantName
	provider := authTokenEntity.Provider
	email := authTokenEntity.EmailAddress

	var emailImportState *postgresEntity.UserEmailImportState
	if state == postgresEntity.REAL_TIME {

		//activate real time sync only if there are more than Batch Size emails imported, to not overlap with the history sync
		//if there are no emails imported, activate real time sync only if there are no other active syncs

		emailImportStateLastWeek, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_WEEK)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if emailImportStateLastWeek == nil {
			span.LogKV("message", "no gmail import state found for tenant: %s and username: %s", tenant, email)
			return
		}

		if emailImportStateLastWeek.Active == true {
			span.LogKV("message", "gmail import state for tenant: %s and username: %s is active for last week. skipping real time import", tenant, email)
			return
		}

		emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if emailImportState.Active == false {
			importedEmails, err := s.commonServices.PostgresRepositories.RawEmailRepository.CountForUsername(ctx, "gmail", tenant, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return
			}

			//batch size of 100
			if importedEmails > 100 {
				err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.ActivateEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}
			} else {

				lastWeek, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_WEEK)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				last3Months, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_3_MONTHS)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				lastYear, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_YEAR)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				olderThanOneYear, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.OLDER_THAN_ONE_YEAR)
				if err != nil {
					tracing.TraceErr(span, err)
					return
				}

				if lastWeek.Active == false && last3Months.Active == false && lastYear.Active == false && olderThanOneYear.Active == false {
					err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.ActivateEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
					if err != nil {
						tracing.TraceErr(span, err)
						return
					}

					emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
					if err != nil {
						tracing.TraceErr(span, err)
						return
					}
				}

				span.LogKV("message", "gmail import state for tenant: %s and username: %s is not active for real time. skipping real time import", tenant, email)
				return
			}
		}
	} else {

		err := s.initializeEmailImportState(ctx, authTokenEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		emailImportState, err = s.getHistoryImportState(ctx, authTokenEntity, postgresEntity.LAST_WEEK)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		if emailImportState == nil {
			span.LogKV("message", "no gmail import state found for tenant: %s and username: %s", tenant, email)
			return
		}
	}

	emailImportState, err := s.syncEmailsForState(ctx, emailImportState)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	if state == postgresEntity.HISTORY && emailImportState.Cursor == "" {
		err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.DeactivateEmailImportState(ctx, tenant, provider, email, emailImportState.State)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}

}

func (s *ingestEmailService) getHistoryImportState(ctx context.Context, authTokenEntity *postgresEntity.OAuthTokenEntity, state postgresEntity.EmailImportState) (*postgresEntity.UserEmailImportState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailService.getHistoryImportState")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	tenant := authTokenEntity.TenantName
	provider := authTokenEntity.Provider
	username := authTokenEntity.EmailAddress

	emailImportState, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, username, state)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if emailImportState == nil {
		err := fmt.Errorf("failed to get gmail import state for tenant: %s and username: %s and week: %s", tenant, username, state)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if emailImportState.Active {
		return emailImportState, nil
	}

	if state == postgresEntity.OLDER_THAN_ONE_YEAR {
		return nil, nil
	}

	importState, err := s.getNextEmailImportState(state)
	if err != nil {
		return nil, err
	}

	return s.getHistoryImportState(ctx, authTokenEntity, importState)
}

func (s *ingestEmailService) getNextEmailImportState(state postgresEntity.EmailImportState) (postgresEntity.EmailImportState, error) {
	switch state {
	case postgresEntity.LAST_WEEK:
		return postgresEntity.LAST_3_MONTHS, nil
	case postgresEntity.LAST_3_MONTHS:
		return postgresEntity.LAST_YEAR, nil
	case postgresEntity.LAST_YEAR:
		return postgresEntity.OLDER_THAN_ONE_YEAR, nil
	case postgresEntity.OLDER_THAN_ONE_YEAR:
		return postgresEntity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	default:
		return postgresEntity.OLDER_THAN_ONE_YEAR, fmt.Errorf("invalid state: %s", state)
	}
}

func (s *ingestEmailService) initializeEmailImportState(ctx context.Context, authTokenEntity *postgresEntity.OAuthTokenEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailService.initializeEmailImportState")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	now := time.Now()
	tenant := authTokenEntity.TenantName
	provider := authTokenEntity.Provider
	email := authTokenEntity.EmailAddress

	emailImportState, err := s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailImportState == nil {
		_, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(ctx, tenant, provider, email, postgresEntity.REAL_TIME, nil, nil, false, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_WEEK)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(0, 0, -7)
		_, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_WEEK, &now, &stop, true, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_3_MONTHS)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(0, -3, 0)
		_, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_3_MONTHS, &now, &stop, true, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_YEAR)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(-1, 0, 0)
		_, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(ctx, tenant, provider, email, postgresEntity.LAST_YEAR, &now, &stop, true, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	emailImportState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.GetEmailImportState(ctx, tenant, provider, email, postgresEntity.OLDER_THAN_ONE_YEAR)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if emailImportState == nil {
		stop := now.AddDate(-50, 0, 0)
		_, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.CreateEmailImportState(ctx, tenant, provider, email, postgresEntity.OLDER_THAN_ONE_YEAR, &now, &stop, true, "")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *ingestEmailService) syncEmailsForState(ctx context.Context, importState *postgresEntity.UserEmailImportState) (*postgresEntity.UserEmailImportState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailService.syncEmailsForState")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	batchSize := int64(100)
	countEmailsExists := int64(0)

	var externalSystem string
	var rawEmails []*postgresEntity.EmailRawData
	var next string
	var err error

	if importState.Provider == commonenum.WorkspaceProviderGoogle.String() {
		externalSystem = commonenum.SourceGmail.String()
		rawEmails, next, err = s.commonServices.GoogleService.ReadEmails(ctx, batchSize, importState)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else if importState.Provider == commonenum.WorkspaceProviderAzure.String() {
		externalSystem = commonenum.SourceOutlook.String()
		rawEmails, next, err = s.commonServices.AzureService.ReadEmailsFromAzureAd(ctx, importState)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	for _, emailRawData := range rawEmails {

		emailExists, err := s.commonServices.PostgresRepositories.RawEmailRepository.EmailExistsByMessageId(ctx, externalSystem, importState.Tenant, importState.Username, emailRawData.MessageId)
		if err != nil {
			return nil, fmt.Errorf("unable to check if email exists: %v", err)
		}

		// counting emails that are already imported based on the batch size
		// if the job is stopped in the middle of execution and we haven't saved the latest token
		// we are going to lose the history
		if emailExists {

			if importState.State == postgresEntity.REAL_TIME {
				countEmailsExists = countEmailsExists + 1

				if countEmailsExists >= batchSize {
					importState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, "")
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}
					return importState, nil
				}
			}

			continue
		} else {
			zeroTime := time.Time{}
			if emailRawData.Sent != zeroTime && importState.StopDate != nil && emailRawData.Sent.Before(*importState.StopDate) {
				importState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, "")
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				return importState, nil
			}

		}

		jsonContent, err := JSONMarshal(emailRawData)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.commonServices.PostgresRepositories.RawEmailRepository.Store(ctx, externalSystem, importState.Tenant, importState.Username, emailRawData.ProviderMessageId, emailRawData.MessageId, string(jsonContent), emailRawData.Sent, importState.State)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		headersString, err := JSONMarshal(emailRawData.Headers)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		ingestEmailMessage := postgresEntity.IngestEmailMessage{
			Tenant:   importState.Tenant,
			Username: importState.Username,
			Provider: importState.Provider,
			State:    postgresEntity.IngestEmailMessageStatePending,

			Subject:     emailRawData.Subject,
			TextContent: emailRawData.Text,
			HtmlContent: emailRawData.Html,

			SentAt: emailRawData.Sent,

			From: emailRawData.From,
			To:   emailRawData.To,
			Cc:   emailRawData.Cc,
			Bcc:  emailRawData.Bcc,

			ProviderMessageId:  emailRawData.ProviderMessageId,
			ProviderThreadId:   emailRawData.ThreadId,
			ProviderInReplyTo:  emailRawData.InReplyTo,
			ProviderReferences: emailRawData.Reference,

			Headers: string(headersString),
		}

		err = s.commonServices.PostgresRepositories.IngestEmailMessageRepository.Store(ctx, externalSystem, importState.Tenant, importState.Username, emailRawData.ProviderMessageId, &ingestEmailMessage)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	importState, err = s.commonServices.PostgresRepositories.UserEmailImportPageTokenRepository.UpdateEmailImportState(ctx, importState.Tenant, importState.Provider, importState.Username, importState.State, next)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
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
