package postgres_repository

import (
	"context"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type UserCalendarAvailabilityRepository interface {
	GetByTenantAndEmail(ctx context.Context, tenant string, email string) (*postgres_entity.UserCalendarAvailability, error)
	SaveOrUpdate(ctx context.Context, availability *postgres_entity.UserCalendarAvailability) (*postgres_entity.UserCalendarAvailability, error)
}

type userCalendarAvailabilityRepository struct {
	db *gorm.DB
}

func NewUserCalendarAvailabilityRepository(db *gorm.DB) UserCalendarAvailabilityRepository {
	return &userCalendarAvailabilityRepository{
		db: db,
	}
}

func (r *userCalendarAvailabilityRepository) GetByTenantAndEmail(ctx context.Context, tenant string, email string) (*postgres_entity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "UserCalendarAvailabilityRepository.GetByTenantAndEmail")
	defer spans.Finish()
	spans.LogKV("email", email)

	var availability postgres_entity.UserCalendarAvailability
	result := r.db.Where("tenant = ? AND LOWER(email) = LOWER(?)", tenant, email).First(&availability)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		return nil, result.Error
	}

	spans.LogKV("result.found", true)
	return &availability, nil
}

func (r *userCalendarAvailabilityRepository) SaveOrUpdate(ctx context.Context, availability *postgres_entity.UserCalendarAvailability) (*postgres_entity.UserCalendarAvailability, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "UserCalendarAvailabilityRepository.SaveOrUpdate")
	defer spans.Finish()
	spans.LogKV("email", availability.Email)

	// Check if record exists
	existing, err := r.GetByTenantAndEmail(ctx, availability.Tenant, availability.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing record
		availability.ID = existing.ID
		availability.CreatedAt = existing.CreatedAt
		result := r.db.Save(availability)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		// Create new record
		availability.Email = strings.ToLower(availability.Email)
		result := r.db.Create(availability)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return availability, nil
}
