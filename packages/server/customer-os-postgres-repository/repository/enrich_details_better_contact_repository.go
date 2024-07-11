package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type enrichDetailsBetterContactRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsBetterContactRepository interface {
	RegisterRequest(ctx context.Context, request entity.EnrichDetailsBetterContact) error
	AddResponse(ctx context.Context, requestId, response string) error
	GetByLinkedInUrl(ctx context.Context, linkedInUrl string) (*entity.EnrichDetailsBetterContact, error)
	GetByRequestId(ctx context.Context, requestId string) (*entity.EnrichDetailsBetterContact, error)
	GetBy(ctx context.Context, firstName, lastName, companyName, companyDomain string) ([]*entity.EnrichDetailsBetterContact, error)
	GetWithoutResponses(ctx context.Context) ([]*entity.EnrichDetailsBetterContact, error)
}

func NewEnrichDetailsBetterContactRepository(gormDb *gorm.DB) EnrichDetailsBetterContactRepository {
	return &enrichDetailsBetterContactRepository{gormDb: gormDb}
}

func (r enrichDetailsBetterContactRepository) RegisterRequest(ctx context.Context, request entity.EnrichDetailsBetterContact) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.RegisterRequest")
	defer span.Finish()

	err := r.gormDb.Create(&request).Error
	if err != nil {
		return err
	}

	return nil
}

func (e enrichDetailsBetterContactRepository) AddResponse(ctx context.Context, requestId, response string) error {
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
		return err
	}

	return nil
}

func (r enrichDetailsBetterContactRepository) GetByLinkedInUrl(ctx context.Context, linkedInUrl string) (*entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetLatestByRequestId")
	defer span.Finish()
	span.LogFields(tracingLog.String("linkedInUrl", linkedInUrl))

	var entity *entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("contact_linkedin_url = ?", linkedInUrl).
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return entity, err
}

func (r enrichDetailsBetterContactRepository) GetByRequestId(ctx context.Context, requestId string) (*entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetByRequestId")
	defer span.Finish()
	span.LogFields(tracingLog.String("requestId", requestId))

	var entity *entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("request_id = ?", requestId).
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return entity, err
}

func (r enrichDetailsBetterContactRepository) GetBy(ctx context.Context, firstName, lastName, companyName, companyDomain string) ([]*entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetLatestByRequestId")
	defer span.Finish()
	span.LogFields(tracingLog.String("firstName", firstName), tracingLog.String("lastName", lastName), tracingLog.String("companyName", companyName), tracingLog.String("companyDomain", companyDomain))

	var entity []*entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("contact_first_name = ?", firstName).
		Where("contact_last_name = ?", lastName).
		Where("company_name = ?", companyName).
		Where("company_domain = ?", companyDomain).
		Find(&entity).Error

	return entity, err
}

func (r enrichDetailsBetterContactRepository) GetWithoutResponses(ctx context.Context) ([]*entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetWithoutResponses")
	defer span.Finish()

	var entity []*entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("response = ?", "").
		Where("created_at < ?", time.Now().Add(-10*time.Minute)).
		Limit(50).
		Find(&entity).Error

	return entity, err
}
