package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const (
	fundingRoundPreSeed = "Pre-Seed"
	fundingRoundSeed    = "Seed"
	fundingRoundSeriesA = "Series A"
	fundingRoundSeriesB = "Series B"
	fundingRoundSeriesC = "Series C"
	fundingRoundSeriesD = "Series D"
	fundingRoundSeriesE = "Series E"
	fundingRoundSeriesF = "Series F+"
	fundingRoundIpo     = "IPO"
	fundingRoundFF      = "Friends & Family"
	fundingRoundAngel   = "Angel"
	fundingRoundBridge  = "Bridge"
)

var fundingRoundByModel = map[model.FundingRound]string{
	model.FundingRoundPreSeed:          fundingRoundPreSeed,
	model.FundingRoundSeed:             fundingRoundSeed,
	model.FundingRoundSeriesA:          fundingRoundSeriesA,
	model.FundingRoundSeriesB:          fundingRoundSeriesB,
	model.FundingRoundSeriesC:          fundingRoundSeriesC,
	model.FundingRoundSeriesD:          fundingRoundSeriesD,
	model.FundingRoundSeriesE:          fundingRoundSeriesE,
	model.FundingRoundSeriesF:          fundingRoundSeriesF,
	model.FundingRoundIPO:              fundingRoundIpo,
	model.FundingRoundFriendsAndFamily: fundingRoundFF,
	model.FundingRoundAngel:            fundingRoundAngel,
	model.FundingRoundBridge:           fundingRoundBridge,
}

var fundingRoundByValue = utils.ReverseMap(fundingRoundByModel)

func MapFundingRoundFromModel(input *model.FundingRound) string {
	if input == nil {
		return ""
	}
	if v, exists := fundingRoundByModel[*input]; exists {
		return v
	} else {
		return ""
	}
}

func MapFundingRoundToModel(input string) *model.FundingRound {
	if v, exists := fundingRoundByValue[input]; exists {
		return &v
	} else {
		return nil
	}
}
