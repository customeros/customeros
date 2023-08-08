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
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"strings"
	"time"
)

type SyncFromSourceService interface {
	FetchDataFromSlack()
}

type syncFromSourceService struct {
	cfg            *config.Config
	log            logger.Logger
	repositories   *repository.Repositories
	slackService   SlackService
	rawDataStoreDb *config.RawDataStoreDB
	cache          caches.Cache
}

type SlackWorkspaceDtls struct {
	token          string
	lookBackWindow string
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

func (s *syncFromSourceService) FetchDataFromSlack() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel() // Cancel context on exit

	tenants, err := s.repositories.TenantRepository.GetTenantsWithOrganizationsWithSlackChannels(ctx)
	if err != nil {
		s.log.Errorf("Failed to get tenants for slack sync: %v", err)
		return
	} else {
		s.log.Infof("Got %d tenants for slack sync", len(tenants))
	}

	err = s.repositories.SlackSyncReposiotry.AutoMigrate()
	if err != nil {
		s.log.Errorf("Failed to auto migrate slack sync table: %v", err)
		return
	}

	// Long-running process
	for _, tenant := range tenants {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
			// Continue processing tenants
		}
		s.syncSlackChannelsForTenant(ctx, tenant)
	}
}

func (s *syncFromSourceService) syncSlackChannelsForTenant(ctx context.Context, tenant string) {
	span, ctx := tracing.StartTracerSpan(ctx, "SyncFromSourceService.syncSlackChannelsForTenant")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	err := s.autoMigrateRawTables(tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to auto migrate raw tables for tenant %s: %v", tenant, err)
		return
	}

	slackDetails, err := s.getSlackConnectionDetailsForTenant(ctx, tenant)
	if err != nil || slackDetails == nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get slack details for tenant %s: %v", tenant, err)
		return
	}

	tenantDomain, err := s.repositories.TenantRepository.GetTenantDomain(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get tenant domain for tenant %s: %v", tenant, err)
		return
	}

	organizations, err := s.repositories.OrganizationRepository.GetOrganizationsWithSlackChannels(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed to get organizations for tenant %s: %v", tenant, err)
		return
	}
	for _, organization := range organizations {
		orgId := utils.GetStringPropOrEmpty(organization.Props, "id")
		orgName := utils.GetStringPropOrEmpty(organization.Props, "name")
		channelLink := utils.GetStringPropOrEmpty(organization.Props, "slackChannelLink")
		syncId, err := uuid.NewUUID()
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to generate sync id for organization %s: %v", organization.Id, err)
			continue
		}
		err = s.syncOrganizationSlackChannel(ctx, tenant, tenantDomain, orgId, orgName, channelLink, syncId.String(), slackDetails)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to sync slack channels for organization %s: %v", orgId, err)
		}
	}
}

func (s *syncFromSourceService) syncOrganizationSlackChannel(ctx context.Context, tenant, tenantDomain, orgId, orgName, channelLink, syncId string, details *SlackWorkspaceDtls) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.syncOrganizationSlackChannel")
	defer span.Finish()
	span.LogFields(log.String("orgId", orgId), log.String("orgName", orgName), log.String("channelLink", channelLink), log.String("syncId", syncId))

	channelId := s.extractChannelIdFromLink(channelLink)
	if channelId == "" {
		err := fmt.Errorf("failed to extract channel id from link %s", channelLink)
		tracing.TraceErr(span, err)
		return err
	}

	syncRunAt := utils.Now()
	previousSyncRunAt, err := s.getPreviousSyncRunForChannel(ctx, tenant, channelId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// get new messages
	messages, err := s.slackService.FetchNewMessagesFromSlackChannel(ctx, channelId, previousSyncRunAt, syncRunAt, *details)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(messages) == 0 {
		err = s.repositories.SlackSyncReposiotry.SaveSyncRun(ctx, tenant, channelId, syncRunAt)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	// get current user ids from channel
	channelUserIds, err := s.slackService.FetchUserIdsFromSlackChannel(ctx, channelId, *details)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// sync user details if not sync before
	for _, userId := range channelUserIds {
		err = s.syncUser(ctx, tenant, tenantDomain, userId, *details)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	// filter for real users
	channelRealUserIds := make([]string, 0)
	for _, userId := range channelUserIds {
		if s.isRealUser(tenant, userId) {
			channelRealUserIds = append(channelRealUserIds, userId)
		}
	}

	for _, message := range messages {
		output := struct {
			slack.Message
			channelUserIds []string
		}{
			Message:        message,
			channelUserIds: channelRealUserIds,
		}
		messageJson, _ := json.Marshal(output)
		err = rawrepo.RawChannelMessages_Save(ctx, s.getDb(tenant), string(messageJson))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	err = s.repositories.SlackSyncReposiotry.SaveSyncRun(ctx, tenant, channelId, syncRunAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *syncFromSourceService) getSlackConnectionDetailsForTenant(ctx context.Context, tenant string) (*SlackWorkspaceDtls, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.getSlackConnectionDetailsForTenant")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	qr := s.repositories.TenantSettingsRepository.FindForTenantName(ctx, tenant)
	var settings entity.TenantSettings
	var ok bool
	if qr.Error != nil {
		return nil, qr.Error
	} else if qr.Result == nil {
		return nil, fmt.Errorf("GetForTenant: no settings found for tenant %s", tenant)
	} else {
		settings, ok = qr.Result.(entity.TenantSettings)
		if !ok {
			return nil, fmt.Errorf("GetForTenant: unexpected type %T", qr.Result)
		}
	}
	return &SlackWorkspaceDtls{
		token:          *settings.SlackApiToken,
		lookBackWindow: *settings.SlackLookbackWindow,
	}, nil
}

func (s *syncFromSourceService) extractChannelIdFromLink(channelLink string) string {
	split := strings.Split(channelLink, "/")

	// Channel link is usually in format
	// https://app.slack.com/client/xxx/C012345689
	channelID := split[len(split)-1]

	if strings.HasPrefix(channelID, "C") {
		return channelID
	} else {
		return ""
	}
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
	return nil
}

func (s *syncFromSourceService) getPreviousSyncRunForChannel(ctx context.Context, tenant string, channelId string) (time.Time, error) {
	queryResult := s.repositories.SlackSyncReposiotry.FindForTenantAndChannelId(ctx, tenant, channelId)
	if queryResult.Error != nil {
		return time.Time{}, queryResult.Error
	}
	if queryResult.Result != nil {
		return queryResult.Result.(entity.SlackSync).LastSyncAt, nil
	}
	oneYearAgo := utils.Now().AddDate(-1, 0, 0)
	return oneYearAgo, nil
}

func (s *syncFromSourceService) syncUser(ctx context.Context, tenant, tenantDomain, userId string, dtls SlackWorkspaceDtls) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncFromSourceService.syncUser")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.String("tenantDomain", tenantDomain), log.String("userId", userId))

	_, ok := s.cache.GetUser(tenant + "-" + userId)
	if !ok {
		slackUser, err := s.slackService.FetchUserInfo(ctx, userId, dtls)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Failed to fetch user info for user %s: %v", userId, err)
			return err
		}
		if slackUser == nil {
			return nil
		}
		if slackUser.Deleted || slackUser.IsBot || slackUser.IsAppUser || slackUser.IsRestricted || slackUser.IsUltraRestricted {
			// save as non-user
			s.cache.SetUser(tenant+"-"+userId, "non-user")
			return nil
		}
		userJson, err := json.Marshal(slackUser)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if strings.HasSuffix(slackUser.Profile.Email, tenantDomain) {
			// save as user
			err = rawrepo.RawUsers_Save(ctx, s.getDb(tenant), string(userJson))
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			s.cache.SetUser(tenant+"-"+userId, "user")
		} else {
			// save as contact
			err = rawrepo.RawContacts_Save(ctx, s.getDb(tenant), string(userJson))
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			s.cache.SetUser(tenant+"-"+userId, "contact")
		}
	}
	return nil
}

func (s *syncFromSourceService) isRealUser(tenant, userId string) bool {
	userType, _ := s.cache.GetUser(tenant + "-" + userId)
	return userType == "user" || userType == "contact"
}
