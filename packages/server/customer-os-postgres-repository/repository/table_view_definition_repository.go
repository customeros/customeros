package repository

import (
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
	UpdateTableViewPresetDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult
}

func NewTableViewDefinitionRepository(gormDb *gorm.DB) TableViewDefinitionRepository {
	return &tableViewDefinitionRepository{gormDb: gormDb}
}

func (t tableViewDefinitionRepository) GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.GetTableViewDefinitions")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagUserId, userId)

	var tableViewDefinitions []entity.TableViewDefinition
	var tableViewPresetDefinitions []entity.TableViewDefinition

	defsErr := t.gormDb.
		Where("tenant = ?", tenant).
		Where("user_id = ?", userId).
		Order("position asc").
		Find(&tableViewDefinitions).Error

	presetsErr := t.gormDb.
		Where("tenant = ?", tenant).
		Where("is_preset = ?", true).
		Order("position asc").
		Find(&tableViewPresetDefinitions).Error

	if defsErr != nil {
		return helper.QueryResult{Error: defsErr}
	}
	if presetsErr != nil {
		return helper.QueryResult{Error: presetsErr}
	}

	allTableViewDefinitions := append(tableViewDefinitions, tableViewPresetDefinitions...)

	return helper.QueryResult{Result: allTableViewDefinitions}
}

func (t tableViewDefinitionRepository) CreateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.CreateTableViewDefinition")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, viewDefinition.Tenant)
	span.SetTag(tracing.SpanTagUserId, viewDefinition.UserId)

	// if the view is a preset, set the UserId to empty string
	if viewDefinition.IsPreset {
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
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, viewDefinition.Tenant)
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

	// Verify that the tableViewDef is not a preset
	if existing.IsPreset {
		return helper.QueryResult{Error: errors.New("record is a preset. use tableViewDef_UpdatePreset instead")}
	}
	// Verify that UserId and TenantName are unchanged
	if existing.UserId != viewDefinition.UserId || existing.Tenant != viewDefinition.Tenant {
		return helper.QueryResult{Error: errors.New("user ID or tenant name mismatch")}
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

func (t tableViewDefinitionRepository) UpdateTableViewPresetDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TableViewDefinitionRepository.UpdateTableViewPresetDefinition")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, viewDefinition.Tenant)
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

	// Verify that the tableViewDef is a preset
	if !existing.IsPreset {
		return helper.QueryResult{Error: errors.New("record is not a preset. use tableViewDef_Update instead")}
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
