package service

import (
	"encoding/xml"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"time"
)

type CurrencyService interface {
	GetCurrencyRatesECB()
}

type currencyService struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCurrencyService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories) CurrencyService {
	return &currencyService{
		cfg:          cfg,
		log:          log,
		repositories: repositories,
	}
}

func (c currencyService) GetCurrencyRatesECB() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "CurrencyService.GetCurrencyRatesECB")
	defer span.Finish()
	// Make HTTP GET request to ECB API endpoint
	resp, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		tracing.TraceErr(span, err)
		c.log.Errorf("Error making HTTP request to ECB API: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		c.log.Errorf("Error reading response body: %s", err.Error())
		return
	}

	// Define struct for XML unmarshalling
	type Envelope struct {
		Cube struct {
			Cube struct {
				Time       string `xml:"time,attr"`
				Currencies []struct {
					Currency string  `xml:"currency,attr"`
					Rate     float64 `xml:"rate,attr"`
				} `xml:"Cube"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	}

	// Unmarshal XML response
	var envelope Envelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		tracing.TraceErr(span, err)
		c.log.Errorf("Error unmarshalling XML response: %s", err.Error())
		return
	}

	// Extract date from response
	date, err := time.Parse("2006-01-02", envelope.Cube.Cube.Time)
	if err != nil {
		tracing.TraceErr(span, err)
		c.log.Errorf("Error parsing date: %s", err.Error())
		return
	}

	// Step 1, iterate over currencies and find USD rate
	usdToEurRate := float64(1)
	for _, currency := range envelope.Cube.Cube.Currencies {
		if currency.Currency == "USD" {
			usdToEurRate = utils.TruncateFloat64(float64(1)/currency.Rate, 5)
			err := c.repositories.PostgresRepositories.CurrencyRateRepository.SaveCurrencyRate("EUR", usdToEurRate, date, "European Central Bank")
			if err != nil {
				tracing.TraceErr(span, err)
				c.log.Errorf("Error saving currency rate: %s", err.Error())
				return
			}
			break
		}
	}
	// Step 2, iterate over currencies and convert rates to USD
	for _, currency := range envelope.Cube.Cube.Currencies {
		if currency.Currency != "USD" {
			usdToCurrency := utils.TruncateFloat64(usdToEurRate*currency.Rate, 5)
			err := c.repositories.PostgresRepositories.CurrencyRateRepository.SaveCurrencyRate(currency.Currency, usdToCurrency, date, "European Central Bank")
			if err != nil {
				tracing.TraceErr(span, err)
				c.log.Errorf("Error saving currency rate: %s", err.Error())
			}
		}
	}
}
