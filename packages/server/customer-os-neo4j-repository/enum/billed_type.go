package enum

type BilledType string

const (
	BilledTypeNone      BilledType = ""
	BilledTypeMonthly   BilledType = "MONTHLY"
	BilledTypeQuarterly BilledType = "QUARTERLY"
	BilledTypeAnnually  BilledType = "ANNUALLY"
	BilledTypeOnce      BilledType = "ONCE"
	BilledTypeUsage     BilledType = "USAGE"
)

var AllBilledTypes = []BilledType{
	BilledTypeNone,
	BilledTypeMonthly,
	BilledTypeQuarterly,
	BilledTypeAnnually,
	BilledTypeOnce,
	BilledTypeUsage,
}

func DecodeBilledType(s string) BilledType {
	if IsValidBilledType(s) {
		return BilledType(s)
	}
	return BilledTypeNone
}

func IsValidBilledType(s string) bool {
	for _, ms := range AllBilledTypes {
		if ms == BilledType(s) {
			return true
		}
	}
	return false
}

func (bt BilledType) String() string {
	return string(bt)
}
