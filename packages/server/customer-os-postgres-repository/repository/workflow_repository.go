package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type workflowRepository struct {
	gormDb *gorm.DB
}

type WorkflowRepository interface {
	GetWorkflowByTypeIfExists(ctx context.Context, tenant string, workflowType entity.WorkflowType) (*entity.Workflow, error)
	CreateWorkflow(ctx context.Context, workflow *entity.Workflow) (entity.Workflow, error)
	UpdateWorkflow(ctx context.Context, id uint64, name, condition, actionParam1 *string, live *bool) error
	GetWorkflowByTenantAndId(ctx context.Context, tenant string, id uint64) (entity.Workflow, error)
	GetWorkflows(ctx context.Context, tenant string) ([]entity.Workflow, error)
	GetAllTenantsLiveWorkflows(ctx context.Context) ([]entity.Workflow, error)
}

func NewWorkflowRepository(gormDb *gorm.DB) WorkflowRepository {
	return &workflowRepository{gormDb: gormDb}
}

func (t workflowRepository) GetWorkflowByTypeIfExists(ctx context.Context, tenant string, workflowType entity.WorkflowType) (*entity.Workflow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.GetWorkflowByTypeIfExists")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var workflow entity.Workflow
	err := t.gormDb.
		Where("tenant = ?", tenant).
		Where("workflow_type = ?", workflowType).
		First(&workflow).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &workflow, nil
}

func (t workflowRepository) CreateWorkflow(ctx context.Context, workflow *entity.Workflow) (entity.Workflow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.CreateWorkflow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := t.gormDb.Create(workflow).Error
	if err != nil {
		return entity.Workflow{}, err
	}

	return *workflow, nil
}

func (t workflowRepository) UpdateWorkflow(ctx context.Context, id uint64, name, condition, actionParam1 *string, live *bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.UpdateWorkflow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	updateMap := make(map[string]interface{})
	if name != nil {
		updateMap["name"] = *name
	}
	if condition != nil {
		updateMap["condition"] = *condition
	}
	if actionParam1 != nil {
		updateMap["action_param1"] = *actionParam1
	}
	if live != nil {
		updateMap["live"] = *live
	}
	updateMap["updated_at"] = utils.Now()

	return t.gormDb.Model(&entity.Workflow{}).Where("id = ?", id).Updates(updateMap).Error
}

func (t workflowRepository) GetWorkflowByTenantAndId(ctx context.Context, tenant string, id uint64) (entity.Workflow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.GetWorkflowByTenantAndId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var workflow entity.Workflow
	err := t.gormDb.
		Where("id = ?", id).
		Where("tenant = ?", tenant).
		First(&workflow).Error
	if err != nil {
		return workflow, err
	}

	return workflow, nil
}

func (t workflowRepository) GetAllTenantsLiveWorkflows(ctx context.Context) ([]entity.Workflow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.GetAllTenantsLiveWorkflows")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var workflows []entity.Workflow
	err := t.gormDb.
		Where("live = ?", true).
		Order("created_at").
		Find(&workflows).Error
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func (t workflowRepository) GetWorkflows(ctx context.Context, tenant string) ([]entity.Workflow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "WorkflowRepository.GetWorkflows")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var workflows []entity.Workflow
	err := t.gormDb.
		Where("tenant = ?", tenant).
		Order("created_at").
		Find(&workflows).Error
	if err != nil {
		return nil, err
	}

	return workflows, nil
}
