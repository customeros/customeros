package currency

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type currencyService struct {
	postgres *postgresRepository.Repositories
}

func NewCurrencyService(postgres *postgresRepository.Repositories) interfaces.CurrencyService {
	return &currencyService{
		postgres: postgres,
	}
}

func (c *currencyService) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CurrencyService.GetRate")
	defer spans.Finish()
	spans.LogKV("fromCurrency", fromCurrency)
	spans.LogKV("toCurrency", toCurrency)

	// Note, currency rates keep currencies from USD to other currencies

	finalRate := 1.0

	// step 1 convert from fromCurrency into USD
	if fromCurrency != "" && fromCurrency != neo4jenum.CurrencyUSD.String() {
		rate, err := c.postgres.CurrencyRateRepository.GetLatestCurrencyRate(ctx, fromCurrency)
		if err != nil {
			spans.TraceError(err)
			return 1, err
		}
		if rate != nil {
			finalRate = 1 / rate.Rate
		}
	}

	// step 2 covert from USD to toCurrency
	if toCurrency != "" && toCurrency != neo4jenum.CurrencyUSD.String() {
		rate, err := c.postgres.CurrencyRateRepository.GetLatestCurrencyRate(ctx, toCurrency)
		if err != nil {
			spans.TraceError(err)
			return 1, err
		}
		if rate != nil {
			finalRate = finalRate * rate.Rate
		}
	}

	return finalRate, nil
}

func (c *currencyService) GetAmountInCurrency(ctx context.Context, amount float64, fromCurrency, toCurrency string) (float64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CurrencyService.GetAmountInCurrency")
	defer spans.Finish()
	spans.LogKV("amount", amount)
	spans.LogKV("fromCurrency", fromCurrency)
	spans.LogKV("toCurrency", toCurrency)

	if amount == 0 {
		return 0, nil
	}

	rate, err := c.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		spans.TraceError(err)
		return 0, err
	}

	return amount * rate, nil
}
