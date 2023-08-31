package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/caches"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository"
	rawrepo "github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository/postgres_raw"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"sync"
)

var postgresTablesAutomigrated = false

type MessageOpenlineFields struct {
	ChannelUserIds   []string          `json:"channel_user_ids"`
	ChannelUserNames map[string]string `json:"channel_user_names"`
	ChannelId        string            `json:"channel_id"`
	ChannelName      string            `json:"channel_name"`
	OrganizationId   string            `json:"organization_id"`
}

type UserOpenlineFields struct {
	OrganizationId string `json:"organization_id"`
	TenantDomain   string `json:"tenant_domain"`
	TenantTeamId   string `json:"tenant_team_id"`
}

type SyncFromSourceService interface {
	SyncSlackRawData()
}

type syncFromSourceService struct {
	cfg            *config.Config
	log            logger.Logger
	repositories   *repository.Repositories
	slackService   SlackService
	rawDataStoreDb *config.RawDataStoreDB
	cache          caches.Cache
}

func NewSlackRawService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, rawDataStoreDb *config.RawDataStoreDB) SyncFromSourceService {
	return &syncFromSourceService{
		cfg:            cfg,
		log:            log,
		repositories:   repositories,
		slackService:   NewSlackService(cfg, log, repositories),
		rawDataStoreDb: rawDataStoreDb,
		cache:          caches.InitCaches(),
	}
}

func (s *syncFromSourceService) SyncSlackRawData() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel() // Cancel context on exit

	if !postgresTablesAutomigrated {
		err := s.repositories.SlackSyncSettingsRepository.AutoMigrate()
		if err != nil {
			s.log.Errorf("Failed to auto migrate %s table: %v", entity.SlackSyncSettings{}.TableName(), err)
			return
		}
		err = s.repositories.SlackSyncRunRepository.AutoMigrate()
		if err != nil {
			s.log.Errorf("Failed to auto migrate %s table: %v", entity.SlackSyncRunStatus{}.TableName(), err)
			return
		}
		postgresTablesAutomigrated = true
	}

	slackSyncSettings, err := s.repositories.SlackSyncSettingsRepository.GetChannelsToSync(ctx)
	if err != nil {
		s.log.Errorf("Failed to get tenants for slack sync: %v", err)
		return
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	maxWorkers := 5
	workerLimit := make(chan struct{}, maxWorkers)

	// Sync all Slack channels concurrently
	for _, slackDtls := range slackSyncSettings {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
			// Continue with Slack sync
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(slackDtls entity.SlackSyncSettings) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			s.prepareAndSyncSlackChannel(ctx, slackDtls)
		}(slackDtls)
	}
	// Wait for all workers to finish
	wg.Wait()
}

func (s *syncFromSourceService) prepareAndSyncSlackChannel(ctx context.Context, slackStngs entity.SlackSyncSettings) {
	span, ctx := tracing.StartTracerSpan(ctx, "SyncFromSourceService.prepareAndSyncSlackChannel")
	defer span.Finish()
	span.SetTag("tenant", slackStngs.Tenant)
	span.LogFields(log.Object("slackSettings", slackStngs))

	err := s.autoMigrateRawTables(slackStngs.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to auto migrate raw tables for tenant %s: %v", slackStngs.Tenant, err)
		return
	}

	slackToken, err := s.getSlackToken(ctx, slackStngs.Tenant)
	if err != nil || slackToken == "" {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get slack token for tenant %s: %v", slackStngs.Tenant, err)
		return
	}

	tenantDomain, err := s.repositories.TenantRepository.GetTenantDomain(ctx, slackStngs.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get tenant domain for tenant %s: %v", slackStngs.Tenant, err)
		return
	}

	organization, err := s.repositories.OrganizationRepository.GetOrganization(ctx, slackStngs.Tenant, slackStngs.OrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get organization %s for tenant %s: %v", slackStngs.OrganizationId, slackStngs.Tenant, err)
		return
	}
	orgName := utils.GetStringPropOrEmpty(organization.Props, "name")

	runId, err := uuid.NewUUID()
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to generate sync id for organization %s: %v", organization.Id, err)
		return
	}
	runStatus := &entity.SlackSyncRunStatus{
		Tenant:         slackStngs.Tenant,
		OrganizationId: slackStngs.OrganizationId,
		RunId:          runId.String(),
		StartAt:        utils.Now(),
	}

	err = s.syncSlackChannelForOrganization(ctx, slackStngs, tenantDomain, orgName, runId.String(), slackToken, runStatus)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to sync slack channels for organization %s: %v", slackStngs.OrganizationId, err)
		runStatus.Failed = true
	}

	runStatus.EndAt = utils.Now()
	s.repositories.SlackSyncRunRepository.Save(ctx, *runStatus)
}

func (s *syncFromSourceService) syncSlackChannelForOrganization(ctx context.Context, slackStngs entity.SlackSyncSettings, tenantDomain, orgName, runId, token string, runStatus *entity.SlackSyncRunStatus) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.syncSlackChannelForOrganization")
	defer span.Finish()
	span.SetTag("tenant", slackStngs.Tenant)
	span.LogFields(log.Object("slackSettings", slackStngs), log.String("orgName", orgName), log.String("runId", runId))

	runStatus.SlackChannelId = slackStngs.ChannelId

	currentSyncTime := utils.Now()

	if slackStngs.ChannelName == "" {
		channelInfo, err := s.slackService.FetchChannelInfo(ctx, token, slackStngs.ChannelId)
		if err != nil {
			s.log.Errorf("Failed to fetch channel info for channel %s: %s", slackStngs.ChannelId, err.Error())
		} else {
			if channelInfo.Name != "" {
				slackStngs.ChannelName = channelInfo.Name
			}
		}
	}
	if slackStngs.TeamId == "" {
		authTest, err := s.slackService.AuthTest(ctx, token)
		if err != nil {
			s.log.Errorf("Failed to test auth: %s", err.Error())
		} else {
			if authTest.TeamID != "" {
				slackStngs.TeamId = authTest.TeamID
			}
		}
	}

	// get new messages
	channelMessages, err := s.slackService.FetchNewMessagesFromSlackChannel(ctx, token, slackStngs.ChannelId, slackStngs.GetSyncStartDate(), currentSyncTime)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// get channel messages with thread replies
	channelMessagesWithReplies, err := s.slackService.FetchMessagesFromSlackChannelWithReplies(ctx, token, slackStngs.ChannelId, currentSyncTime, slackStngs.LookbackWindow)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// get thread messages
	threadMessages := make([]slack.Message, 0)
	for _, message := range channelMessagesWithReplies {
		messages, err := s.slackService.FetchNewThreadMessages(ctx, token, slackStngs.ChannelId, message.Timestamp, slackStngs.GetSyncStartDate(), currentSyncTime)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		for _, threadMessage := range messages {
			threadMessages = append(threadMessages, threadMessage)
		}
	}

	// if no new messages, save sync run and return
	if len(channelMessages) == 0 && len(threadMessages) == 0 {
		err = s.repositories.SlackSyncSettingsRepository.SaveSyncRun(ctx, slackStngs, currentSyncTime)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	// get current user ids from channel
	channelUserIds, err := s.slackService.FetchUserIdsFromSlackChannel(ctx, token, slackStngs.ChannelId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// get user details if not synced before
	for _, userId := range channelUserIds {
		err = s.syncUser(ctx, slackStngs.Tenant, tenantDomain, userId, slackStngs.OrganizationId, slackStngs.TeamId, token)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	channelRealUserIds := make([]string, 0)
	channelUserNames := make(map[string]string, 0)
	for _, userId := range channelUserIds {
		if s.isRealUser(slackStngs.Tenant, userId) {
			channelRealUserIds = append(channelRealUserIds, userId)
		}
		channelUserNames[userId] = s.getUserName(slackStngs.Tenant, userId)
	}

	for _, message := range channelMessages {
		output := struct {
			slack.Message
			MessageOpenlineFields `json:"openline_fields"`
		}{
			Message: message,
			MessageOpenlineFields: MessageOpenlineFields{
				ChannelUserIds:   channelRealUserIds,
				ChannelUserNames: channelUserNames,
				ChannelId:        slackStngs.ChannelId,
				ChannelName:      slackStngs.ChannelName,
				OrganizationId:   slackStngs.OrganizationId,
			},
		}
		messageJson, _ := json.Marshal(output)
		err = rawrepo.RawChannelMessages_Save(ctx, s.getDb(slackStngs.Tenant), string(messageJson))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, message := range threadMessages {
		output := struct {
			slack.Message
			MessageOpenlineFields `json:"openline_fields"`
		}{
			Message: message,
			MessageOpenlineFields: MessageOpenlineFields{
				ChannelUserIds:   channelRealUserIds,
				ChannelUserNames: channelUserNames,
				ChannelId:        slackStngs.ChannelId,
				ChannelName:      slackStngs.ChannelName,
				OrganizationId:   slackStngs.OrganizationId,
			},
		}
		messageJson, _ := json.Marshal(output)
		err = rawrepo.RawThreadMessages_Save(ctx, s.getDb(slackStngs.Tenant), string(messageJson))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	err = s.repositories.SlackSyncSettingsRepository.SaveSyncRun(ctx, slackStngs, currentSyncTime)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *syncFromSourceService) syncUser(ctx context.Context, tenant, tenantDomain, userId, orgId, teamId, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.syncUser")
	defer span.Finish()
	span.SetTag("tenant", tenant)
	span.LogFields(log.String("userId", userId), log.String("orgId", orgId), log.String("tenantDomain", tenantDomain))

	_, ok := s.cache.GetSlackUser(tenant, userId)
	if !ok {
		slackUser, err := s.slackService.FetchUserInfo(ctx, token, userId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to fetch slack user info for user id %s: %v", userId, err)
			return err
		}
		if slackUser == nil {
			span.LogFields(log.String("output", "slack user not found"))
			return nil
		}

		output := struct {
			slack.User
			UserOpenlineFields `json:"openline_fields"`
		}{
			User: *slackUser,
			UserOpenlineFields: UserOpenlineFields{
				TenantDomain:   tenantDomain,
				OrganizationId: orgId,
				TenantTeamId:   teamId,
			},
		}
		userJson, err := json.Marshal(output)

		err = rawrepo.RawUsers_Save(ctx, s.getDb(tenant), string(userJson))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if slackUser.Deleted {
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				Deleted: true,
				Bot:     false,
				Name:    slackUser.Name,
			})
		} else if slackUser.IsBot || slackUser.IsAppUser {
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				Deleted: false,
				Bot:     true,
				Name:    slackUser.RealName,
			})
		} else {
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				Deleted: false,
				Bot:     false,
				Name:    slackUser.Profile.RealNameNormalized,
			})
		}
	}
	return nil
}

func (s *syncFromSourceService) getSlackToken(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.getSlackToken")
	defer span.Finish()
	span.SetTag("tenant", tenant)

	queryResult := s.repositories.TenantSettingsRepository.FindForTenantName(ctx, tenant)
	var settings entity.TenantSettings
	var ok bool
	if queryResult.Error != nil {
		return "", queryResult.Error
	} else if queryResult.Result == nil {
		return "", fmt.Errorf("GetForTenant: no settings found for tenant %s", tenant)
	} else {
		settings, ok = queryResult.Result.(entity.TenantSettings)
		if !ok {
			return "", fmt.Errorf("GetForTenant: unexpected type %T", queryResult.Result)
		}
	}
	if settings.SlackApiToken == nil {
		return "", errors.New("GetForTenant: no slack api token found")
	}
	return utils.IfNotNilStringWithDefault(settings.SlackApiToken, ""), nil
}

func (s *syncFromSourceService) getDb(tenant string) *gorm.DB {
	return s.rawDataStoreDb.GetDBHandler(&config.Context{
		Schema: "slack_" + tenant,
	})
}

func (s *syncFromSourceService) autoMigrateRawTables(tenant string) error {
	s.getDb(tenant).Exec("CREATE SCHEMA IF NOT EXISTS " + "slack_" + tenant)

	err := rawrepo.RawUsers_AutoMigrate(s.getDb(tenant))
	if err != nil {
		return err
	}
	err = rawrepo.RawChannelMessages_AutoMigrate(s.getDb(tenant))
	if err != nil {
		return err
	}
	err = rawrepo.RawThreadMessages_AutoMigrate(s.getDb(tenant))
	if err != nil {
		return err
	}
	return nil
}

func (s *syncFromSourceService) isRealUser(tenant, userId string) bool {
	user, _ := s.cache.GetSlackUser(tenant, userId)
	return !user.Bot && !user.Deleted
}

func (s *syncFromSourceService) getUserName(tenant, userId string) string {
	user, _ := s.cache.GetSlackUser(tenant, userId)
	return user.Name
}
