package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CurrencyService interface {
	GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error)
}

type currencyService struct {
	repositories *postgresRepository.Repositories
}

func NewCurrencyService(repositories *postgresRepository.Repositories) CurrencyService {
	return &currencyService{
		repositories: repositories,
	}
}

func (c *currencyService) GetRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CurrencyService.GetRate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("fromCurrency", fromCurrency))
	span.LogFields(log.String("toCurrency", toCurrency))

	// Note, currency rates keep currencies from USD to other currencies

	finalRate := 1.0

	// step 1 convert from fromCurrency into USD
	if fromCurrency != "" && fromCurrency != neo4jenum.CurrencyUSD.String() {
		rate, err := c.repositories.CurrencyRateRepository.GetLatestCurrencyRate(fromCurrency)
		if err != nil {
			tracing.TraceErr(span, err)
			return 1, err
		}
		if rate != nil {
			finalRate = 1 / rate.Rate
		}
	}

	// step 2 covert from USD to toCurrency
	if toCurrency != "" && toCurrency != neo4jenum.CurrencyUSD.String() {
		rate, err := c.repositories.CurrencyRateRepository.GetLatestCurrencyRate(toCurrency)
		if err != nil {
			tracing.TraceErr(span, err)
			return 1, err
		}
		if rate != nil {
			finalRate = finalRate * rate.Rate
		}
	}

	return finalRate, nil
}
