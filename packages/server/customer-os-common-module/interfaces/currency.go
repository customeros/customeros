package interfaces

import "context"

type CurrencyService interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
}
