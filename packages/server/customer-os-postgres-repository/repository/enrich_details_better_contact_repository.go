package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsBetterContactRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsBetterContactRepository interface {
	RegisterRequest(ctx context.Context, request entity.EnrichDetailsBetterContact) helper.QueryResult
	AddResponse(ctx context.Context, requestId, response string) helper.QueryResult
	GetLatestByRequestId(ctx context.Context, requestId string) helper.QueryResult
}

func NewEnrichDetailsBetterContactRepository(gormDb *gorm.DB) EnrichDetailsBetterContactRepository {
	return &enrichDetailsBetterContactRepository{gormDb: gormDb}
}

func (r enrichDetailsBetterContactRepository) RegisterRequest(ctx context.Context, request entity.EnrichDetailsBetterContact) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.RegisterRequest")
	defer span.Finish()

	err := r.gormDb.Create(&request).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: request}
}

func (e enrichDetailsBetterContactRepository) AddResponse(ctx context.Context, requestId, response string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.AddResponse")
	defer span.Finish()
	span.SetTag("requestId", requestId)

	// Add response to the request with the given requestId, empty response and latest by created_at
	err := e.gormDb.
		Model(&entity.EnrichDetailsBetterContact{}).
		Where("request_id = ?", requestId).
		Where("response = ?", "").
		Order("created_at desc").
		Limit(1).
		UpdateColumn("response", response).
		UpdateColumn("updated_at", gorm.Expr("current_timestamp")).
		Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: response}
}

func (r enrichDetailsBetterContactRepository) GetLatestByRequestId(ctx context.Context, requestId string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetLatestByRequestId")
	defer span.Finish()
	span.SetTag("requestId", requestId)

	var request entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("request_id = ?", requestId).
		Order("created_at desc").
		Limit(1).
		First(&request).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: request}
}
