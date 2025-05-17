package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsBetterContactRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsBetterContactRepository interface {
	RegisterRequest(ctx context.Context, data postgres_entity.EnrichDetailsBetterContact) (*postgres_entity.EnrichDetailsBetterContact, error)
	AddResponse(ctx context.Context, requestId, response string) error

	GetById(ctx context.Context, id string) (*postgres_entity.EnrichDetailsBetterContact, error)
	GetByRequestId(ctx context.Context, requestId string) (*postgres_entity.EnrichDetailsBetterContact, error)
	GetByRequestParams(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool, lookBack time.Duration) ([]*postgres_entity.EnrichDetailsBetterContact, error)
	GetWithoutResponses(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsBetterContact, error)
}

func NewEnrichDetailsBetterContactRepository(gormDb *gorm.DB) EnrichDetailsBetterContactRepository {
	return &enrichDetailsBetterContactRepository{gormDb: gormDb}
}

func (r enrichDetailsBetterContactRepository) RegisterRequest(ctx context.Context, data postgres_entity.EnrichDetailsBetterContact) (*postgres_entity.EnrichDetailsBetterContact, error) {
	spans, _ := telemetry.StartServiceSpan(ctx, "EnrichDetailsBetterContactRepository.RegisterRequest")
	defer spans.Finish()
	spans.LogObjectAsJson("data", data)

	data.CreatedAt = utils.Now()
	err := r.gormDb.Create(&data).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &data, nil
}

func (r enrichDetailsBetterContactRepository) AddResponse(ctx context.Context, requestId, response string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBetterContactRepository.AddResponse")
	defer spans.Finish()

	spans.LogKV("requestId", requestId)

	// Add response to the request with the given requestId
	err := r.gormDb.
		Model(&postgres_entity.EnrichDetailsBetterContact{}).
		Where("request_id = ?", requestId).
		Where("response = ?", "").
		Updates(
			map[string]interface{}{
				"response":   response,
				"updated_at": utils.Now(),
			},
		).
		Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r enrichDetailsBetterContactRepository) GetById(ctx context.Context, id string) (*postgres_entity.EnrichDetailsBetterContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBetterContactRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	var betterContactEntity *postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("id = ?", id).
		First(&betterContactEntity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return betterContactEntity, err
}

func (r enrichDetailsBetterContactRepository) GetByRequestId(ctx context.Context, requestId string) (*postgres_entity.EnrichDetailsBetterContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBetterContactRepository.GetByRequestId")
	defer spans.Finish()

	spans.LogKV("requestId", requestId)

	var betterContactEntity *postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("request_id = ?", requestId).
		First(&betterContactEntity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return betterContactEntity, err
}

func (r enrichDetailsBetterContactRepository) GetByRequestParams(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool, lookBack time.Duration) ([]*postgres_entity.EnrichDetailsBetterContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBetterContactRepository.GetByRequestParams")
	defer spans.Finish()

	spans.LogKV("linkedInUrl", linkedInUrl)
	spans.LogKV("firstName", firstName)
	spans.LogKV("lastName", lastName)
	spans.LogKV("companyName", companyName)
	spans.LogKV("companyDomain", companyDomain)
	spans.LogKV("enrichPhoneNumber", enrichPhoneNumber)

	var betterContactEntity []*postgres_entity.EnrichDetailsBetterContact
	tx := r.gormDb.
		Where("contact_linkedin_url = ?", linkedInUrl).
		Where("contact_first_name = ?", firstName).
		Where("contact_last_name = ?", lastName).
		Where("company_name = ?", companyName).
		Where("company_domain = ?", companyDomain).
		Where("created_at > ?", time.Now().Add(-lookBack))

	if enrichPhoneNumber {
		tx = tx.Where("enrich_phone_number = ?", enrichPhoneNumber)
	}
	err := tx.Find(&betterContactEntity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return betterContactEntity, err
}

func (r enrichDetailsBetterContactRepository) GetWithoutResponses(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsBetterContact, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBetterContactRepository.GetWithoutResponses")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	var betterContactEntities []*postgres_entity.EnrichDetailsBetterContact
	err := r.gormDb.
		Where("response = ?", "").
		Where("created_at < ?", utils.Now().Add(-10*time.Minute)).
		Limit(limit).
		Find(&betterContactEntities).Error

	return betterContactEntities, err
}
