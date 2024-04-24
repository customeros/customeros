package enum

// Deprecated
type BillingCycle string

const (
	// Deprecated
	BillingCycleNone BillingCycle = ""
	// Deprecated
	BillingCycleMonthlyBilling BillingCycle = "MONTHLY"
	// Deprecated
	BillingCycleQuarterlyBilling BillingCycle = "QUARTERLY"
	// Deprecated
	BillingCycleAnnuallyBilling BillingCycle = "ANNUALLY"
)

func (e BillingCycle) String() string {
	return string(e)
}
