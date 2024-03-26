package enum

type ContractStatus string

const (
	ContractStatusUndefined     ContractStatus = ""
	ContractStatusDraft         ContractStatus = "DRAFT"
	ContractStatusLive          ContractStatus = "LIVE"
	ContractStatusEnded         ContractStatus = "ENDED"
	ContractStatusOutOfContract ContractStatus = "OUT_OF_CONTRACT"
)

var AllContractStatuses = []ContractStatus{
	ContractStatusDraft,
	ContractStatusLive,
	ContractStatusEnded,
	ContractStatusOutOfContract,
}

func DecodeContractStatus(s string) ContractStatus {
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

func (c ContractStatus) String() string {
	return string(c)
}
