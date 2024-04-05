package enum

type RenewalCycle string

// Deprecated: RenewalCycle is deprecated
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

func DecodeRenewalCycle(s string) RenewalCycle {
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

func (e RenewalCycle) String() string {
	return string(e)
}

func IsFrequencyBasedRenewalCycle(renewalCycle RenewalCycle) bool {
	return renewalCycle == RenewalCycleMonthlyRenewal ||
		renewalCycle == RenewalCycleAnnualRenewal ||
		renewalCycle == RenewalCycleQuarterlyRenewal
}
