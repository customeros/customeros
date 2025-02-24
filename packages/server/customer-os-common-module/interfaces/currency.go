package interfaces

import "context"

type CurrencyService interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
	GetAmountInCurrency(ctx context.Context, amount float64, fromCurrency, toCurrency string) (float64, error)
}
