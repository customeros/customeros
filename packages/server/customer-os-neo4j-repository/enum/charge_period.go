package enum

type ChargePeriod string

const (
	ChargePeriodNone      ChargePeriod = ""
	ChargePeriodMonthly   ChargePeriod = "MOTHLY"
	ChargePeriodQuarterly ChargePeriod = "QUARTERLY"
	ChargePeriodAnnually  ChargePeriod = "ANNUALLY"
)

var AllChargePeriods = []ChargePeriod{
	ChargePeriodNone,
	ChargePeriodMonthly,
	ChargePeriodQuarterly,
	ChargePeriodAnnually,
}

func DecodeChargePeriod(s string) ChargePeriod {
	if IsValidChargePeriod(s) {
		return ChargePeriod(s)
	}
	return ChargePeriodNone
}

func IsValidChargePeriod(s string) bool {
	for _, ms := range AllChargePeriods {
		if ms == ChargePeriod(s) {
			return true
		}
	}
	return false
}

func (e ChargePeriod) String() string {
	return string(e)
}
