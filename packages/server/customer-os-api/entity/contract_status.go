package entity

type ContractStatus string

const (
	ContractStatusUndefined ContractStatus = ""
	ContractStatusDraft     ContractStatus = "DRAFT"
	ContractStatusLive      ContractStatus = "LIVE"
	ContractStatusEnded     ContractStatus = "ENDED"
)

var AllContractStatuses = []ContractStatus{
	ContractStatusDraft,
	ContractStatusLive,
	ContractStatusEnded,
}

func GetContractStatus(s string) ContractStatus {
	if IsValidContractStatus(s) {
		return ContractStatus(s)
	}
	return ContractStatusUndefined
}

func IsValidContractStatus(s string) bool {
	for _, ms := range AllContractStatuses {
		if ms == ContractStatus(s) {
			return true
		}
	}
	return false
}
