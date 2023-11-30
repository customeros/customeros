package entity

type InternalType string

const (
	InternalTypeNbo       InternalType = "NBO"
	InternalTypeUpsell    InternalType = "UPSELL"
	InternalTypeCrossSell InternalType = "CROSS_SELL"
	InternalTypeRenewal   InternalType = "RENEWAL"
)

var AllInternalTypes = []InternalType{
	InternalTypeNbo,
	InternalTypeUpsell,
	InternalTypeCrossSell,
	InternalTypeRenewal,
}

func GetInternalType(s string) InternalType {
	if IsValidInternalType(s) {
		return InternalType(s)
	}
	return InternalTypeNbo
}

func IsValidInternalType(s string) bool {
	for _, ms := range AllInternalTypes {
		if ms == InternalType(s) {
			return true
		}
	}
	return false
}
