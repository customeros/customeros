package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type tableViewDefinitionRepository struct {
	gormDb *gorm.DB
}

type TableViewDefinitionRepository interface {
	GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult
	CreateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult
	UpdateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult
	UpdateTableViewSharedDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult
	ArchiveTableViewDefinition(ctx context.Context, viewDefinitionId uint64) error
}

func NewTableViewDefinitionRepository(gormDb *gorm.DB) TableViewDefinitionRepository {
	return &tableViewDefinitionRepository{gormDb: gormDb}
}

func (t tableViewDefinitionRepository) GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.GetTableViewDefinitions")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagUserId, userId)

	var tableViewDefinitions []entity.TableViewDefinition
	var tableViewSharedDefinitions []entity.TableViewDefinition

	defsErr := t.gormDb.
		Where("tenant = ?", tenant).
		Where("user_id = ?", userId).
		Order("position asc").
		Find(&tableViewDefinitions).Error

	sharedErr := t.gormDb.
		Where("tenant = ?", tenant).
		Where("is_shared = ?", true).
		Order("position asc").
		Find(&tableViewSharedDefinitions).Error

	if defsErr != nil {
		return helper.QueryResult{Error: defsErr}
	}
	if sharedErr != nil {
		return helper.QueryResult{Error: sharedErr}
	}

	allTableViewDefinitions := append(tableViewDefinitions, tableViewSharedDefinitions...)

	return helper.QueryResult{Result: allTableViewDefinitions}
}

func (t tableViewDefinitionRepository) CreateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.CreateTableViewDefinition")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, viewDefinition.Tenant)
	span.SetTag(tracing.SpanTagUserId, viewDefinition.UserId)

	// if the view is a preset, set the UserId to empty string
	if viewDefinition.IsShared {
		viewDefinition.UserId = ""
	}

	err := t.gormDb.Create(&viewDefinition).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return helper.QueryResult{Error: err}
	}
	span.LogFields(log.Uint64("result.createdID", viewDefinition.ID))
	return helper.QueryResult{Result: viewDefinition}
}

func (t tableViewDefinitionRepository) UpdateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.UpdateTableViewDefinition")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, viewDefinition.Tenant)
	span.SetTag(tracing.SpanTagUserId, viewDefinition.UserId)

	// Retrieve the existing record by ID
	var existing entity.TableViewDefinition
	err := t.gormDb.First(&existing, viewDefinition.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.QueryResult{Error: errors.New("record not found")}
		}
		return helper.QueryResult{Error: err}
	}

	// Verify that the tableViewDef is not shared
	if existing.IsShared {
		return helper.QueryResult{Error: errors.New("record is shared. use tableViewDef_UpdateShared instead")}
	}
	// Verify that UserId and TenantName are unchanged
	if existing.UserId != viewDefinition.UserId || existing.Tenant != viewDefinition.Tenant {
		return helper.QueryResult{Error: errors.New("user ID or tenant name mismatch")}
	}

	// Update the record
	// Map the fields you want to allow updating, excluding UserId and TenantName
	updateData := map[string]interface{}{
		"table_name":      viewDefinition.Name,
		"position":        viewDefinition.Order,
		"icon":            viewDefinition.Icon,
		"filters":         viewDefinition.Filters,
		"default_filters": viewDefinition.DefaultFilters,
		"sorting":         viewDefinition.Sorting,
		"columns":         viewDefinition.ColumnsJson,
		"updated_at":      utils.Now(),
	}

	err = t.gormDb.Model(&existing).Updates(updateData).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: existing}
}

func (t tableViewDefinitionRepository) UpdateTableViewSharedDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.UpdateTableViewSharedDefinition")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, viewDefinition.Tenant)
	span.SetTag(tracing.SpanTagUserId, viewDefinition.UserId)

	// Retrieve the existing record by ID
	var existing entity.TableViewDefinition
	err := t.gormDb.First(&existing, viewDefinition.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.QueryResult{Error: errors.New("record not found")}
		}
		return helper.QueryResult{Error: err}
	}

	// Verify that the tableViewDef is shared
	if !existing.IsShared {
		return helper.QueryResult{Error: errors.New("record is not shared. use tableViewDef_Update instead")}
	}
	// Verify that tenantName is unchanged
	if existing.Tenant != viewDefinition.Tenant {
		return helper.QueryResult{Error: errors.New("tenant name mismatch")}
	}

	// Update the record
	// Map the fields you want to allow updating, excluding UserId and TenantName
	updateData := map[string]interface{}{
		"table_name": viewDefinition.Name,
		"position":   viewDefinition.Order,
		"icon":       viewDefinition.Icon,
		"filters":    viewDefinition.Filters,
		"sorting":    viewDefinition.Sorting,
		"columns":    viewDefinition.ColumnsJson,
		"updated_at": utils.Now(),
	}

	err = t.gormDb.Model(&existing).Updates(updateData).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: existing}
}

func (t tableViewDefinitionRepository) ArchiveTableViewDefinition(ctx context.Context, viewDefinitionId uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.ArchiveTableViewDefinition")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := t.gormDb.
		Where("tenant = ?", common.GetTenantFromContext(ctx)).
		Where("user_id = ? OR is_shared = ?", common.GetUserIdFromContext(ctx), true).
		Where("id = ?", viewDefinitionId).
		Delete(&entity.TableViewDefinition{})

	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		err := errors.New("TableViewDef not found")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
