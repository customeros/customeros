package enum

type PriceCalculationType string

const (
	PriceCalculationTypeNone         PriceCalculationType = ""
	PriceCalculationTypeRevenueShare PriceCalculationType = "REVENUE_SHARE"
)

var AllPriceCalculationTypes = []PriceCalculationType{
	PriceCalculationTypeNone,
	PriceCalculationTypeRevenueShare,
}

func DecodePriceCalculationType(s string) PriceCalculationType {
	if IsValidPriceCalculationType(s) {
		return PriceCalculationType(s)
	}
	return PriceCalculationTypeNone
}

func IsValidPriceCalculationType(s string) bool {
	for _, ms := range AllPriceCalculationTypes {
		if ms == PriceCalculationType(s) {
			return true
		}
	}
	return false
}

func (e PriceCalculationType) String() string {
	return string(e)
}
