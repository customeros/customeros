package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const (
	marketB2B         = "B2B"
	marketB2C         = "B2C"
	marketMarketplace = "Marketplace"
)

var marketByModel = map[model.Market]string{
	model.MarketB2b:         marketB2B,
	model.MarketB2c:         marketB2C,
	model.MarketMarketplace: marketMarketplace,
}

var marketByValue = utils.ReverseMap(marketByModel)

func MapMarketFromModel(input *model.Market) string {
	if input == nil {
		return ""
	}
	if v, exists := marketByModel[*input]; exists {
		return v
	} else {
		return ""
	}
}

func MapMarketToModel(input string) *model.Market {
	if v, exists := marketByValue[input]; exists {
		return &v
	} else {
		return nil
	}
}
