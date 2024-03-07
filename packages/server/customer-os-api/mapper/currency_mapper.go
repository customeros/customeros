package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var currencyByModel = map[model.Currency]neo4jenum.Currency{
	model.CurrencyAud: neo4jenum.CurrencyAUD,
	model.CurrencyBrl: neo4jenum.CurrencyBRL,
	model.CurrencyCad: neo4jenum.CurrencyCAD,
	model.CurrencyChf: neo4jenum.CurrencyCHF,
	model.CurrencyCny: neo4jenum.CurrencyCNY,
	model.CurrencyEur: neo4jenum.CurrencyEUR,
	model.CurrencyGbp: neo4jenum.CurrencyGBP,
	model.CurrencyHkd: neo4jenum.CurrencyHKD,
	model.CurrencyInr: neo4jenum.CurrencyINR,
	model.CurrencyJpy: neo4jenum.CurrencyJPY,
	model.CurrencyKrw: neo4jenum.CurrencyKRW,
	model.CurrencyMxn: neo4jenum.CurrencyMXN,
	model.CurrencyNok: neo4jenum.CurrencyNOK,
	model.CurrencyNzd: neo4jenum.CurrencyNZD,
	model.CurrencyRon: neo4jenum.CurrencyRON,
	model.CurrencySek: neo4jenum.CurrencySEK,
	model.CurrencySgd: neo4jenum.CurrencySGD,
	model.CurrencyTry: neo4jenum.CurrencyTRY,
	model.CurrencyUsd: neo4jenum.CurrencyUSD,
	model.CurrencyZar: neo4jenum.CurrencyZAR,
}

var currencyByValue = utils.ReverseMap(currencyByModel)

func MapCurrencyFromModel(input model.Currency) neo4jenum.Currency {
	return currencyByModel[input]
}

func MapCurrencyToModel(input neo4jenum.Currency) model.Currency {
	return currencyByValue[input]
}
