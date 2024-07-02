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

type workflowRepository struct {
	gormDb *gorm.DB
}

type WorkflowRepository interface {
	GetWorkflowByTypeIfExists(ctx context.Context, tenant, wo string) (*entity.Workflow, error)
	CreateWorkflow(ctx context.Context, workflow entity.Workflow) (entity.Workflow, error)
	UpdateWorkflow(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult
}

func NewWorkflowRepository(gormDb *gorm.DB) WorkflowRepository {
	return &workflowRepository{gormDb: gormDb}
}

func (t workflowRepository) GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.GetTableViewDefinitions")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagUserId, userId)
	span.LogFields(log.String("tenant", tenant))

	var tableViewDefinitions []entity.TableViewDefinition
	err := t.gormDb.
		Where("tenant = ?", tenant).
		Where("user_id = ?", userId).
		Order("position asc").
		Find(&tableViewDefinitions).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: tableViewDefinitions}
}

func (t workflowRepository) CreateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.CreateTableViewDefinition")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.SetTag(tracing.SpanTagTenant, viewDefinition.Tenant)
	span.SetTag(tracing.SpanTagUserId, viewDefinition.UserId)

	err := t.gormDb.Create(&viewDefinition).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return helper.QueryResult{Error: err}
	}
	span.LogFields(log.Uint64("result.createdID", viewDefinition.ID))
	return helper.QueryResult{Result: viewDefinition}
}

func (t workflowRepository) UpdateTableViewDefinition(ctx context.Context, viewDefinition entity.TableViewDefinition) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.CreateTableViewDefinition")
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
