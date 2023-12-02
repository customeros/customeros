package entity

type RenewalCycle string

const (
	RenewalCycleNone             RenewalCycle = ""
	RenewalCycleMonthlyRenewal   RenewalCycle = "MONTHLY"
	RenewalCycleQuarterlyRenewal RenewalCycle = "QUARTERLY"
	RenewalCycleAnnualRenewal    RenewalCycle = "ANNUALLY"
)

var AllRenewalCycles = []RenewalCycle{
	RenewalCycleNone,
	RenewalCycleMonthlyRenewal,
	RenewalCycleQuarterlyRenewal,
	RenewalCycleAnnualRenewal,
}

func GetRenewalCycle(s string) RenewalCycle {
	if IsValidRenewalCycle(s) {
		return RenewalCycle(s)
	}
	return RenewalCycleNone
}

func IsValidRenewalCycle(s string) bool {
	for _, ms := range AllRenewalCycles {
		if ms == RenewalCycle(s) {
			return true
		}
	}
	return false
}
