package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type tableViewDefinitionRepository struct {
	gormDb *gorm.DB
}

type TableViewDefinitionRepository interface {
	GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult
	CreateTableViewDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult
	UpdateTableViewDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult
	UpdateTableViewSharedDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult
	ArchiveTableViewDefinition(ctx context.Context, viewDefinitionId uint64) error
	GetTableViewDefinition(ctx context.Context, tenant string, id uint64) (postgres_entity.TableViewDefinition, error)
}

func NewTableViewDefinitionRepository(gormDb *gorm.DB) TableViewDefinitionRepository {
	return &tableViewDefinitionRepository{gormDb: gormDb}
}

func (t tableViewDefinitionRepository) GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.GetTableViewDefinitions")
	defer spans.Finish()

	var tableViewDefinitions []postgres_entity.TableViewDefinition
	var tableViewSharedDefinitions []postgres_entity.TableViewDefinition

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

	spans.LogKV("result.count", len(allTableViewDefinitions))
	return helper.QueryResult{Result: allTableViewDefinitions}
}

func (t tableViewDefinitionRepository) CreateTableViewDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.CreateTableViewDefinition")
	defer spans.Finish()

	spans.LogKV("userId", viewDefinition.UserId)
	spans.LogObjectAsJson("viewDefinition", viewDefinition)

	// if the view is a preset, set the UserId to empty string
	if viewDefinition.IsShared {
		viewDefinition.UserId = ""
	}

	err := t.gormDb.Create(&viewDefinition).Error
	if err != nil {
		spans.TraceError(err)
		return helper.QueryResult{Error: err}
	}
	spans.LogKV("result.createdID", viewDefinition.ID)
	return helper.QueryResult{Result: viewDefinition}
}

func (t tableViewDefinitionRepository) UpdateTableViewDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.UpdateTableViewDefinition")
	defer spans.Finish()

	spans.LogKV("userId", viewDefinition.UserId)

	// Retrieve the existing record by ID
	var existing postgres_entity.TableViewDefinition
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

func (t tableViewDefinitionRepository) UpdateTableViewSharedDefinition(ctx context.Context, viewDefinition postgres_entity.TableViewDefinition) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.UpdateTableViewSharedDefinition")
	defer spans.Finish()

	spans.LogKV("userId", viewDefinition.UserId)

	// Retrieve the existing record by ID
	var existing postgres_entity.TableViewDefinition
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

func (t tableViewDefinitionRepository) ArchiveTableViewDefinition(ctx context.Context, viewDefinitionId uint64) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.ArchiveTableViewDefinition")
	defer spans.Finish()

	result := t.gormDb.
		Where("tenant = ?", common.GetTenantFromContext(ctx)).
		Where("user_id = ? OR is_shared = ?", common.GetUserIdFromContext(ctx), true).
		Where("id = ?", viewDefinitionId).
		Delete(&postgres_entity.TableViewDefinition{})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		err := errors.New("TableViewDef not found")
		spans.TraceError(err)
		return err
	}

	return nil
}

func (t tableViewDefinitionRepository) GetTableViewDefinition(ctx context.Context, tenant string, id uint64) (postgres_entity.TableViewDefinition, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TableViewDefinitionRepository.GetTableViewDefinition")
	defer spans.Finish()

	var tableViewDefinition postgres_entity.TableViewDefinition
	err := t.gormDb.
		Where("tenant = ?", tenant).
		Where("id = ?", id).
		First(&tableViewDefinition).Error
	if err != nil {
		spans.TraceError(err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return tableViewDefinition, nil
}
