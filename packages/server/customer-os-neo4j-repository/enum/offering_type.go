package enum

type OfferingType string

const (
	OfferingTypeNone    OfferingType = ""
	OfferingTypeProduct OfferingType = "PRODUCT"
	OfferingTypeService OfferingType = "SERVICE"
)

var AllOfferingTypes = []OfferingType{
	OfferingTypeNone,
	OfferingTypeProduct,
	OfferingTypeService,
}

func DecodeOfferingType(s string) OfferingType {
	if IsValidOfferingType(s) {
		return OfferingType(s)
	}
	return OfferingTypeNone
}

func IsValidOfferingType(s string) bool {
	for _, ms := range AllOfferingTypes {
		if ms == OfferingType(s) {
			return true
		}
	}
	return false
}

func (e OfferingType) String() string {
	return string(e)
}
