package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	RegisterRequest(ctx context.Context, data postgres_entity.EnrichDetailsBetterContact) (*postgres_entity.EnrichDetailsBetterContact, error)
	AddResponse(ctx context.Context, requestId, response string) error
	GetByLinkedInUrl(ctx context.Context, linkedInUrl string, enrichPhoneNumber bool) (*postgres_entity.EnrichDetailsBetterContact, error)
	GetById(ctx context.Context, id string) (*postgres_entity.EnrichDetailsBetterContact, error)
	GetByRequestId(ctx context.Context, requestId string) (*postgres_entity.EnrichDetailsBetterContact, error)
	GetBy(ctx context.Context, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool) ([]*postgres_entity.EnrichDetailsBetterContact, error)
	GetWithoutResponses(ctx context.Context) ([]*postgres_entity.EnrichDetailsBetterContact, error)
}

func NewEnrichDetailsBetterContactRepository(gormDb *gorm.DB) EnrichDetailsBetterContactRepository {
	return &enrichDetailsBetterContactRepository{gormDb: gormDb}
}

func (r enrichDetailsBetterContactRepository) RegisterRequest(ctx context.Context, data postgres_entity.EnrichDetailsBetterContact) (*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.RegisterRequest")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "data", data)

	now := utils.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	err := r.gormDb.Create(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r enrichDetailsBetterContactRepository) AddResponse(ctx context.Context, requestId, response string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.AddResponse")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.SetTag("requestId", requestId)

	// Add response to the request with the given requestId, empty response and latest by created_at
	err := r.gormDb.
		Model(&postgres_entity.EnrichDetailsBetterContact{}).
		Where("request_id = ?", requestId).
		Where("response = ?", "").
		Order("created_at desc").
		Limit(1).
		UpdateColumn("response", response).
		UpdateColumn("updated_at", utils.Now()).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r enrichDetailsBetterContactRepository) GetByLinkedInUrl(ctx context.Context, linkedInUrl string, enrichPhoneNumber bool) (*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetLatestByRequestId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("linkedInUrl", linkedInUrl), tracingLog.Bool("enrichPhoneNumber", enrichPhoneNumber))

	var postgres_entity *postgres_entity.EnrichDetailsBetterContact
	tx := r.gormDb.
		Where("contact_linkedin_url = ?", linkedInUrl)
	if enrichPhoneNumber {
		tx = tx.Where("enrich_phone_number = ?", enrichPhoneNumber)
	}
	err := tx.First(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return postgres_entity, err
}

func (r enrichDetailsBetterContactRepository) GetById(ctx context.Context, id string) (*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("id", id))

	var postgres_entity *postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("id = ?", id).
		First(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return postgres_entity, err
}

func (r enrichDetailsBetterContactRepository) GetByRequestId(ctx context.Context, requestId string) (*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetByRequestId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("requestId", requestId))

	var postgres_entity *postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("request_id = ?", requestId).
		First(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return postgres_entity, err
}

func (r enrichDetailsBetterContactRepository) GetBy(ctx context.Context, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool) ([]*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetLatestByRequestId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(
		tracingLog.String("firstName", firstName),
		tracingLog.String("lastName", lastName),
		tracingLog.String("companyName", companyName),
		tracingLog.String("companyDomain", companyDomain),
		tracingLog.Bool("enrichPhoneNumber", enrichPhoneNumber))

	var postgres_entity []*postgres_entity.EnrichDetailsBetterContact
	tx := r.gormDb.
		Where("contact_first_name = ?", firstName).
		Where("contact_last_name = ?", lastName).
		Where("company_name = ?", companyName).
		Where("company_domain = ?", companyDomain)
	if enrichPhoneNumber {
		tx = tx.Where("enrich_phone_number = ?", enrichPhoneNumber)
	}
	err := tx.Find(&postgres_entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return postgres_entity, err
}

func (r enrichDetailsBetterContactRepository) GetWithoutResponses(ctx context.Context) ([]*postgres_entity.EnrichDetailsBetterContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsBetterContactRepository.GetWithoutResponses")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var postgres_entity []*postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("response = ?", "").
		Where("created_at < ?", utils.Now().Add(-10*time.Minute)).
		Limit(50).
		Find(&postgres_entity).Error

	return postgres_entity, err
}
