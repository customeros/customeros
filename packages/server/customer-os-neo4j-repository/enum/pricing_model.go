package enum

type PricingModel string

const (
	PricingModelNone         PricingModel = ""
	PricingModelSubscription PricingModel = "SUBSCRIPTION"
	PricingModelOneTime      PricingModel = "ONE_TIME"
	PricingModelUsage        PricingModel = "USAGE"
)

var AllPricingModels = []PricingModel{
	PricingModelNone,
	PricingModelSubscription,
	PricingModelOneTime,
	PricingModelUsage,
}

func DecodePricingModel(s string) PricingModel {
	if IsValidPricingModel(s) {
		return PricingModel(s)
	}
	return PricingModelNone
}

func IsValidPricingModel(s string) bool {
	for _, ms := range AllPricingModels {
		if ms == PricingModel(s) {
			return true
		}
	}
	return false
}

func (e PricingModel) String() string {
	return string(e)
}
