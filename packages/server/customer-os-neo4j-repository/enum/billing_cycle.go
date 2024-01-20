package enum

type BillingCycle string

const (
	BillingCycleNone             BillingCycle = ""
	BillingCycleMonthlyBilling   BillingCycle = "MONTHLY"
	BillingCycleQuarterlyBilling BillingCycle = "QUARTERLY"
	BillingCycleAnnuallyBilling  BillingCycle = "ANNUALLY"
)

var AllBillingCycles = []BillingCycle{
	BillingCycleNone,
	BillingCycleMonthlyBilling,
	BillingCycleQuarterlyBilling,
	BillingCycleAnnuallyBilling,
}

func DecodeBillingCycle(s string) BillingCycle {
	if IsValidBillingCycle(s) {
		return BillingCycle(s)
	}
	return BillingCycleNone
}

func IsValidBillingCycle(s string) bool {
	for _, ms := range AllBillingCycles {
		if ms == BillingCycle(s) {
			return true
		}
	}
	return false
}

func (e BillingCycle) String() string {
	return string(e)
}
