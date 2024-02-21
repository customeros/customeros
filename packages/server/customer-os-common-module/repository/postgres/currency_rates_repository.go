package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
	"time"
)

type currencyRateRepo struct {
	db *gorm.DB
}

type CurrencyRateRepository interface {
	GetLatestCurrencyRate(currency string) (*entity.CurrencyRate, error)
	SaveCurrencyRate(currency string, rate float64, date time.Time, source string) error
}

func NewCurrencyRateRepository(db *gorm.DB) CurrencyRateRepository {
	return &currencyRateRepo{db: db}
}

func (r *currencyRateRepo) GetLatestCurrencyRate(currency string) (*entity.CurrencyRate, error) {
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

func (r *currencyRateRepo) SaveCurrencyRate(currency string, rate float64, date time.Time, source string) error {
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
	return r.db.Save(&existingRate).Error
}
