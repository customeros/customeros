package entity

type ContractRenewalCycle string

const (
	ContractRenewalCycleNone             ContractRenewalCycle = ""
	ContractRenewalCycleMonthlyRenewal   ContractRenewalCycle = "MONTHLY"
	ContractRenewalCycleQuarterlyRenewal ContractRenewalCycle = "QUARTERLY"
	ContractRenewalCycleAnnualRenewal    ContractRenewalCycle = "ANNUALLY"
)

var AllContractRenewalCycles = []ContractRenewalCycle{
	ContractRenewalCycleNone,
	ContractRenewalCycleMonthlyRenewal,
	ContractRenewalCycleQuarterlyRenewal,
	ContractRenewalCycleAnnualRenewal,
}

func GetContractRenewalCycle(s string) ContractRenewalCycle {
	if IsValidContractRenewalCycle(s) {
		return ContractRenewalCycle(s)
	}
	return ContractRenewalCycleNone
}

func IsValidContractRenewalCycle(s string) bool {
	for _, ms := range AllContractRenewalCycles {
		if ms == ContractRenewalCycle(s) {
			return true
		}
	}
	return false
}
