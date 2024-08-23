package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type currencyRateRepo struct {
	db *gorm.DB
}

type CurrencyRateRepository interface {
	GetLatestCurrencyRate(ctx context.Context, currency string) (*entity.CurrencyRate, error)
	SaveCurrencyRate(ctx context.Context, currency string, rate float64, date time.Time, source string) error
}

func NewCurrencyRateRepository(db *gorm.DB) CurrencyRateRepository {
	return &currencyRateRepo{db: db}
}

func (r *currencyRateRepo) GetLatestCurrencyRate(ctx context.Context, currency string) (*entity.CurrencyRate, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CurrencyRateRepository.GetLatestCurrencyRate")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var rate entity.CurrencyRate
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
	span, _ := opentracing.StartSpanFromContext(ctx, "CurrencyRateRepository.SaveCurrencyRate")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Check if the currency rate already exists for the given currency and date
	var existingRate entity.CurrencyRate
	err := r.db.
		Where("currency = ?", currency).
		Where("date = ?", date).
		Where("source = ?", source).
		First(&existingRate).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// If the currency rate doesn't exist, create a new record
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newRate := entity.CurrencyRate{
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
