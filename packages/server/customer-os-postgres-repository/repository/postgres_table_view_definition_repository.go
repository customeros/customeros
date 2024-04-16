package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type tableViewDefinitionRepository struct {
	gormDb *gorm.DB
}

func (t tableViewDefinitionRepository) GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "AppKeyRepo.FindByKey")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.LogFields(log.String("tenant", tenant))

	var tableViewDefinitions []entity.TableViewDefinition
	err := t.gormDb.
		Where("tenant = ?", tenant).
		Where("user_id = ?", userId).
		Find(&tableViewDefinitions).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: tableViewDefinitions}
}

func (t tableViewDefinitionRepository) CreateTableViewDefinition(viewDefinition entity.TableViewDefinition) helper.QueryResult {
	//TODO implement me
	panic("implement me")
}

func (t tableViewDefinitionRepository) UpdateTableViewDefinition(viewDefinition entity.TableViewDefinition) helper.QueryResult {
	//TODO implement me
	panic("implement me")
}

type TableViewDefinitionRepository interface {
	GetTableViewDefinitions(ctx context.Context, tenant, userId string) helper.QueryResult
	CreateTableViewDefinition(viewDefinition entity.TableViewDefinition) helper.QueryResult
	UpdateTableViewDefinition(viewDefinition entity.TableViewDefinition) helper.QueryResult
}

func NewTableViewDefinitionRepository(gormDb *gorm.DB) TableViewDefinitionRepository {
	return &tableViewDefinitionRepository{gormDb: gormDb}
}
