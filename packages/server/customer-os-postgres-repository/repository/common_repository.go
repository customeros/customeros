package repository

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"reflect"
)

type CommonRepository interface {
	UpdateProperty(ctx context.Context, tenant string, entityType interface{}, id any, propertyName string, newValue interface{}) error
	PermanentlyDelete(ctx context.Context, tenant string) error
}

type commonRepository struct {
	postgresDB *config.PostgresDB
}

func NewCommonRepository(postgresDB *config.PostgresDB) CommonRepository {
	return &commonRepository{postgresDB: postgresDB}
}

func (r *commonRepository) UpdateProperty(ctx context.Context, tenant string, entityType interface{}, id any, propertyName string, newValue interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.UpdateProperty")
	defer span.Finish()

	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("entityType", reflect.TypeOf(entityType).String()))
	span.LogFields(log.String("id", fmt.Sprintf("%v", id)))
	span.LogFields(log.String("propertyName", propertyName))
	span.LogFields(log.String("newValue", fmt.Sprintf("%v", newValue)))

	// Create a new instance of the entity type
	entity := reflect.New(reflect.TypeOf(entityType)).Interface()

	// Fetch the entity by ID and tenant using context
	query := r.postgresDB.GormDB.WithContext(ctx).Where("tenant = ? and id = ?", tenant, id).First(entity)
	if err := query.Error; err != nil {
		tracing.TraceErr(span, err)
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("entity not found for tenant %s with id %v", tenant, id)
		}
		return fmt.Errorf("failed to find entity: %w", err)
	}

	// Use reflection to update the property
	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName(propertyName)

	if !field.IsValid() {
		err := fmt.Errorf("property %s does not exist on entity", propertyName)
		tracing.TraceErr(span, err)
		return err
	}

	if !field.CanSet() {
		err := fmt.Errorf("property %s cannot be set", propertyName)
		tracing.TraceErr(span, err)
		return err
	}

	// Set the new value
	field.Set(reflect.ValueOf(newValue))

	// Save the updated entity with context
	if err := r.postgresDB.GormDB.WithContext(ctx).Save(entity).Error; err != nil {
		err := fmt.Errorf("failed to save updated entity: %w", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *commonRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.PermanentlyDelete")
	defer span.Finish()

	asyncTablesWithTenantNameColumn := []string{
		entity.GoogleServiceAccountKey{}.TableName(),
		entity.OAuthTokenEntity{}.TableName(),
	}

	asyncTablesWithTenantColumn := []string{
		entity.RawEmail{}.TableName(),
		entity.UserEmailImportState{}.TableName(),
		entity.UserEmailImportStateHistory{}.TableName(),
	}

	tableNamesWithTenantNameColumn := []string{
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
		entity.SlackChannelNotification{}.TableName(),
		entity.StatsApiCalls{}.TableName(),
		entity.TableViewDefinition{}.TableName(),
		entity.TenantSettingsEmailExclusion{}.TableName(),
		entity.TenantSettingsMailbox{}.TableName(),
		entity.TenantSettingsOpportunityStage{}.TableName(),
		entity.Tracking{}.TableName(),
		entity.TrackingAllowedOrigin{}.TableName(),
		entity.UserWorkingSchedule{}.TableName(),
	}

	for _, tableName := range asyncTablesWithTenantNameColumn {
		if err := r.postgresDB.AsyncGormDB.Exec("DELETE FROM "+tableName+" WHERE tenant_name = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, tableName := range asyncTablesWithTenantColumn {
		if err := r.postgresDB.AsyncGormDB.Exec("DELETE FROM "+tableName+" WHERE tenant_name = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, tableName := range tableNamesWithTenantNameColumn {
		if err := r.postgresDB.GormDB.Exec("DELETE FROM "+tableName+" WHERE tenant_name = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, tableName := range tableNamesWithTenantColumn {
		if err := r.postgresDB.GormDB.Exec("DELETE FROM "+tableName+" WHERE tenant = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
