package enum

type ContractStatus string

const (
	ContractStatusUndefined     ContractStatus = ""
	ContractStatusDraft         ContractStatus = "DRAFT"
	ContractStatusScheduled     ContractStatus = "SCHEDULED"
	ContractStatusLive          ContractStatus = "LIVE"
	ContractStatusEnded         ContractStatus = "ENDED"
	ContractStatusOutOfContract ContractStatus = "OUT_OF_CONTRACT"
)

var AllContractStatuses = []ContractStatus{
	ContractStatusDraft,
	ContractStatusScheduled,
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
	for _, cs := range AllContractStatuses {
		if cs == ContractStatus(s) {
			return true
		}
	}
	return false
}

func (c ContractStatus) String() string {
	return string(c)
}
