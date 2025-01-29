package postgres_repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type CommonRepository interface {
	UpdateProperty(ctx context.Context, tenant string, postgres_entityType interface{}, id any, propertyName string, newValue interface{}) error
	PermanentlyDelete(ctx context.Context, tenant string) error
}

type commonRepository struct {
	postgresDB *config.PostgresDB
}

func NewCommonRepository(postgresDB *config.PostgresDB) CommonRepository {
	return &commonRepository{postgresDB: postgresDB}
}

func (r *commonRepository) UpdateProperty(ctx context.Context, tenant string, postgres_entityType interface{}, id any, propertyName string, newValue interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.UpdateProperty")
	defer span.Finish()

	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("postgres_entityType", reflect.TypeOf(postgres_entityType).String()))
	span.LogFields(log.String("id", fmt.Sprintf("%v", id)))
	span.LogFields(log.String("propertyName", propertyName))
	span.LogFields(log.String("newValue", fmt.Sprintf("%v", newValue)))

	// Create a new instance of the postgres_entity type
	postgres_entity := reflect.New(reflect.TypeOf(postgres_entityType)).Interface()

	// Fetch the postgres_entity by ID and tenant using context
	query := r.postgresDB.GormDB.WithContext(ctx).Where("tenant = ? and id = ?", tenant, id).First(postgres_entity)
	if err := query.Error; err != nil {
		tracing.TraceErr(span, err)
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("postgres_entity not found for tenant %s with id %v", tenant, id)
		}
		return fmt.Errorf("failed to find postgres_entity: %w", err)
	}

	// Use reflection to update the property
	v := reflect.ValueOf(postgres_entity).Elem()
	field := v.FieldByName(propertyName)

	if !field.IsValid() {
		err := fmt.Errorf("property %s does not exist on postgres_entity", propertyName)
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

	// Save the updated postgres_entity with context
	if err := r.postgresDB.GormDB.WithContext(ctx).Save(postgres_entity).Error; err != nil {
		err := fmt.Errorf("failed to save updated postgres_entity: %w", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *commonRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.PermanentlyDelete")
	defer span.Finish()

	asyncTablesWithTenantNameColumn := []string{
		postgres_entity.GoogleServiceAccountKey{}.TableName(),
		postgres_entity.OAuthTokenEntity{}.TableName(),
	}

	asyncTablesWithTenantColumn := []string{
		postgres_entity.RawEmail{}.TableName(),
		postgres_entity.UserEmailImportState{}.TableName(),
		postgres_entity.UserEmailImportStateHistory{}.TableName(),
	}

	tableNamesWithTenantNameColumn := []string{
		postgres_entity.PersonalIntegration{}.TableName(),
		postgres_entity.PostmarkApiKey{}.TableName(),
		postgres_entity.SlackChannel{}.TableName(),
		postgres_entity.SlackSettingsEntity{}.TableName(),
		postgres_entity.TenantSettings{}.TableName(),
		postgres_entity.TenantWebhook{}.TableName(),
		postgres_entity.TenantWebhookApiKey{}.TableName(),
	}

	tableNamesWithTenantColumn := []string{
		postgres_entity.AiPromptLog{}.TableName(),
		postgres_entity.ApiBillableEvent{}.TableName(),
		postgres_entity.BrowserAutomationsRun{}.TableName(),
		postgres_entity.BrowserConfig{}.TableName(),
		postgres_entity.CosApiEnrichPersonTempResult{}.TableName(),
		postgres_entity.CustomerOsIds{}.TableName(),
		postgres_entity.EmailLookup{}.TableName(),
		postgres_entity.EmailMessage{}.TableName(),
		postgres_entity.EmailTracking{}.TableName(),
		postgres_entity.EmailValidationRecord{}.TableName(),
		postgres_entity.EmailValidationRequestBulk{}.TableName(),
		postgres_entity.MailStackDomain{}.TableName(),
		postgres_entity.SlackChannelNotification{}.TableName(),
		postgres_entity.StatsApiCalls{}.TableName(),
		postgres_entity.TableViewDefinition{}.TableName(),
		postgres_entity.TenantSettingsEmailExclusion{}.TableName(),
		postgres_entity.TenantSettingsMailbox{}.TableName(),
		postgres_entity.TenantSettingsOpportunityStage{}.TableName(),
		postgres_entity.UserWorkingSchedule{}.TableName(),
		postgres_entity.WebSession{}.TableName(),
		postgres_entity.WebTrackerEvents{}.TableName(),
	}

	for _, tableName := range asyncTablesWithTenantNameColumn {
		if err := r.postgresDB.AsyncGormDB.Exec("DELETE FROM "+tableName+" WHERE tenant_name = ?", tenant).Error; err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	for _, tableName := range asyncTablesWithTenantColumn {
		if err := r.postgresDB.AsyncGormDB.Exec("DELETE FROM "+tableName+" WHERE tenant = ?", tenant).Error; err != nil {
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
