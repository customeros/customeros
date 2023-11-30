package entity

type BilledType string

const (
	BilledTypeNone     BilledType = ""
	BilledTypeMonthly  BilledType = "MONTHLY"
	BilledTypeAnnually BilledType = "ANNUALLY"
	BilledTypeOnce     BilledType = "ONCE"
	BilledTypeUsage    BilledType = "USAGE"
)

var AllBilledTypes = []BilledType{
	BilledTypeNone,
	BilledTypeMonthly,
	BilledTypeAnnually,
	BilledTypeOnce,
	BilledTypeUsage,
}

func GetBilledType(s string) BilledType {
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
