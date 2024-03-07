package entity

import "time"

// All rates are relative to USD, i.e., 1 USD = Rate Currency (1 USD = 0.93 EUR)
type CurrencyRate struct {
	ID        uint64    `gorm:"primary_key;autoIncrement:true"`
	Currency  string    `gorm:"type:varchar(3);not null"` // Currency code (e.g., EUR, GBP, etc.)
	Rate      float64   `gorm:"not null"`                 // Exchange rate relative to USD
	Date      time.Time `gorm:"not null"`                 // Date of the exchange rate
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp"`
	Source    string    `gorm:"type:varchar(255);not null"` // Source of the exchange rate (e.g., ECB, etc.)
}

func (CurrencyRate) TableName() string {
	return "currency_rates"
}

func (CurrencyRate) UniqueIndex() [][]string {
	return [][]string{
		{"Date", "Currency", "Source"},
	}
}
