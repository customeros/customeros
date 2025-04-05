package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type currencyRateRepo struct {
	db *gorm.DB
}

type CurrencyRateRepository interface {
	GetLatestCurrencyRate(ctx context.Context, currency string) (*postgres_entity.CurrencyRate, error)
	SaveCurrencyRate(ctx context.Context, currency string, rate float64, date time.Time, source string) error
}

func NewCurrencyRateRepository(db *gorm.DB) CurrencyRateRepository {
	return &currencyRateRepo{db: db}
}

func (r *currencyRateRepo) GetLatestCurrencyRate(ctx context.Context, currency string) (*postgres_entity.CurrencyRate, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CurrencyRateRepository.GetLatestCurrencyRate")
	defer spans.Finish()

	var rate postgres_entity.CurrencyRate
	err := r.db.
		Where("currency = ?", currency).
		Order("date desc").
		First(&rate).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rate, nil
}

func (r *currencyRateRepo) SaveCurrencyRate(ctx context.Context, currency string, rate float64, date time.Time, source string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CurrencyRateRepository.SaveCurrencyRate")
	defer spans.Finish()

	// Check if the currency rate already exists for the given currency and date
	var existingRate postgres_entity.CurrencyRate
	err := r.db.
		Where("currency = ?", currency).
		Where("date = ?", date).
		Where("source = ?", source).
		First(&existingRate).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return err
	}

	// If the currency rate doesn't exist, create a new record
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newRate := postgres_entity.CurrencyRate{
			Currency: currency,
			Rate:     rate,
			Date:     date,
			Source:   source,
		}
		return r.db.Create(&newRate).Error
	}

	// If the currency rate exists, update the rate
	existingRate.Rate = rate
	existingRate.UpdatedAt = utils.Now()
	return r.db.Save(&existingRate).Error
}
