package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type CommonRepository interface {
	PermanentlyDelete(ctx context.Context, tenant string) error
}

type commonRepository struct {
	gormDb *gorm.DB
}

func NewCommonRepository(gormDb *gorm.DB) CommonRepository {
	return &commonRepository{gormDb: gormDb}
}

func (r *commonRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.PermanentlyDelete")
	defer span.Finish()

	tableNamesWithTenantNameColumn := []string{
		entity.GoogleServiceAccountKey{}.TableName(),
		entity.OAuthTokenEntity{}.TableName(),
		entity.PersonalIntegration{}.TableName(),
		entity.PostmarkApiKey{}.TableName(),
		entity.SlackChannel{}.TableName(),
		entity.SlackSettingsEntity{}.TableName(),
		entity.TenantSettings{}.TableName(),
		entity.TenantWebhook{}.TableName(),
		entity.TenantWebhookApiKey{}.TableName(),
	}

	tableNamesWithTenantColumn := []string{
		entity.AiPromptLog{}.TableName(),
		entity.ApiBillableEvent{}.TableName(),
		entity.BrowserAutomationsRun{}.TableName(),
		entity.BrowserConfig{}.TableName(),
		entity.CosApiEnrichPersonTempResult{}.TableName(),
		entity.CustomerOsIds{}.TableName(),
		entity.EmailLookup{}.TableName(),
		entity.EmailMessage{}.TableName(),
		entity.EmailTracking{}.TableName(),
		entity.EmailValidationRecord{}.TableName(),
		entity.EmailValidationRequestBulk{}.TableName(),
		entity.EventBuffer{}.TableName(),
		entity.MailStackDomain{}.TableName(),
		entity.RawEmail{}.TableName(),
		entity.SlackChannelNotification{}.TableName(),
		entity.StatsApiCalls{}.TableName(),
		entity.TableViewDefinition{}.TableName(),
		entity.TenantSettingsEmailExclusion{}.TableName(),
		entity.TenantSettingsMailbox{}.TableName(),
		entity.TenantSettingsOpportunityStage{}.TableName(),
		entity.Tracking{}.TableName(),
		entity.TrackingAllowedOrigin{}.TableName(),
		entity.UserEmailImportState{}.TableName(),
		entity.UserEmailImportStateHistory{}.TableName(),
		entity.UserWorkingSchedule{}.TableName(),
		entity.Workflow{}.TableName(),
	}

	for _, tableName := range tableNamesWithTenantNameColumn {
		if err := r.gormDb.Exec("DELETE FROM "+tableName+" WHERE tenant_name = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, tableName := range tableNamesWithTenantColumn {
		if err := r.gormDb.Exec("DELETE FROM "+tableName+" WHERE tenant = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
