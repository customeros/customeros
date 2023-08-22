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
	"strings"
)

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

	err := s.repositories.SlackSyncReposiotry.AutoMigrate()
	if err != nil {
		s.log.Errorf("Failed to auto migrate slack_sync table: %v", err)
		return
	}
	err = s.repositories.SlackSyncRunRepository.AutoMigrate()
	if err != nil {
		s.log.Errorf("Failed to auto migrate slack_sync_run_status table: %v", err)
		return
	}

	slackSyncSettings, err := s.repositories.SlackSyncReposiotry.GetChannelsToSync(ctx)
	if err != nil {
		s.log.Errorf("Failed to get tenants for slack sync: %v", err)
		return
	}

	// Long-running process
	for _, slackDtls := range slackSyncSettings {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
			// Continue processing tenants
		}
		s.prepareAndSyncSlackChannel(ctx, slackDtls)
	}
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
		s.log.Errorf("Failed to get organizations for tenant %s: %v", slackStngs.Tenant, err)
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
		err = s.repositories.SlackSyncReposiotry.SaveSyncRun(ctx, slackStngs.Tenant, slackStngs.ChannelId, currentSyncTime)
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
		err = s.syncUser(ctx, slackStngs.Tenant, tenantDomain, userId, slackStngs.OrganizationId, token)
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
			ChannelUserIds         []string          `json:"channel_user_ids"`
			ChannelUserNames       map[string]string `json:"channel_user_names"`
			ChannelId              string            `json:"channel_id"`
			ChannelName            string            `json:"channel_name"`
			OpenlineOrganizationId string            `json:"openline_organization_id"`
		}{
			Message:                message,
			ChannelUserIds:         channelRealUserIds,
			ChannelUserNames:       channelUserNames,
			ChannelId:              slackStngs.ChannelId,
			ChannelName:            slackStngs.ChannelName,
			OpenlineOrganizationId: slackStngs.OrganizationId,
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
			ChannelUserIds         []string          `json:"channel_user_ids"`
			ChannelUserNames       map[string]string `json:"channel_user_names"`
			ChannelId              string            `json:"channel_id"`
			ChannelName            string            `json:"channel_name"`
			OpenlineOrganizationId string            `json:"openline_organization_id"`
		}{
			Message:                message,
			ChannelUserIds:         channelRealUserIds,
			ChannelUserNames:       channelUserNames,
			ChannelId:              slackStngs.ChannelId,
			ChannelName:            slackStngs.ChannelName,
			OpenlineOrganizationId: slackStngs.OrganizationId,
		}
		messageJson, _ := json.Marshal(output)
		err = rawrepo.RawThreadMessages_Save(ctx, s.getDb(slackStngs.Tenant), string(messageJson))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	err = s.repositories.SlackSyncReposiotry.SaveSyncRun(ctx, slackStngs.Tenant, slackStngs.ChannelId, currentSyncTime)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *syncFromSourceService) syncUser(ctx context.Context, tenant, tenantDomain, userId, orgId, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.syncUser")
	defer span.Finish()
	span.SetTag("tenant", tenant)
	span.LogFields(log.String("tenantDomain", tenantDomain), log.String("slackUserId", userId), log.String("organizationId", orgId))

	cachedSlackUser, ok := s.cache.GetSlackUser(tenant, userId)
	var okContactCheck = true
	if ok && cachedSlackUser.UserType == caches.UserType_Contact {
		_, okContactCheck = s.cache.GetSlackUserAsContactForOrg(orgId, userId)
	}
	if !ok || !okContactCheck {
		slackUser, err := s.slackService.FetchUserInfo(ctx, token, userId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to fetch user info for user %s: %v", userId, err)
			return err
		}
		if slackUser == nil {
			span.LogFields(log.String("output", "slack user not found"))
			return nil
		}
		if slackUser.Deleted || slackUser.IsBot || slackUser.IsAppUser {
			// save as non-user
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				UserType: caches.UserType_NonUser,
				Name:     slackUser.Name,
			})
			span.LogFields(log.String("output", "slack user is not real user"))
			return nil
		}
		if strings.HasSuffix(slackUser.Profile.Email, tenantDomain) {
			// save as user
			userJson, err := json.Marshal(slackUser)
			err = rawrepo.RawUsers_Save(ctx, s.getDb(tenant), string(userJson))
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				UserType: caches.UserType_User,
				Name:     slackUser.Profile.RealNameNormalized,
			})
			span.LogFields(log.String("output", "slack user is user"))
		} else {
			// save as contact
			output := struct {
				slack.User
				OpenlineOrganizationId string `json:"openline_organization_id"`
			}{
				User:                   *slackUser,
				OpenlineOrganizationId: orgId,
			}
			userJson, err := json.Marshal(output)

			err = rawrepo.RawContacts_Save(ctx, s.getDb(tenant), string(userJson))
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			s.cache.SetSlackUser(tenant, userId, caches.SlackUser{
				UserType: caches.UserType_Contact,
				Name:     slackUser.Profile.RealNameNormalized,
			})
			s.cache.SetSlackUserAsContactForOrg(orgId, userId, string(caches.UserType_Contact))
			span.LogFields(log.String("output", "slack user is contact"))
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
	err = rawrepo.RawContacts_AutoMigrate(s.getDb(tenant))
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
	return user.UserType == caches.UserType_User || user.UserType == caches.UserType_Contact
}

func (s *syncFromSourceService) getUserName(tenant, userId string) string {
	user, _ := s.cache.GetSlackUser(tenant, userId)
	return user.Name
}
