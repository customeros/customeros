package entity

type ContractRenewalCycle string

const (
	ContractRenewalCycleNone           ContractRenewalCycle = ""
	ContractRenewalCycleMonthlyRenewal ContractRenewalCycle = "MONTHLY_RENEWAL"
	ContractRenewalCycleAnnualRenewal  ContractRenewalCycle = "ANNUAL_RENEWAL"
)

var AllContractRenewalCycles = []ContractRenewalCycle{
	ContractRenewalCycleNone,
	ContractRenewalCycleMonthlyRenewal,
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
