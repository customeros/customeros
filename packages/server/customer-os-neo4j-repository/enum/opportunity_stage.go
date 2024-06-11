package enum

type OpportunityInternalStage string

const (
	OpportunityInternalStageOpen       OpportunityInternalStage = "OPEN"
	OpportunityInternalStageClosedWon  OpportunityInternalStage = "CLOSED_WON"
	OpportunityInternalStageClosedLost OpportunityInternalStage = "CLOSED_LOST"
	OpportunityInternalStageSuspended  OpportunityInternalStage = "SUSPENDED"
)

var AllOpportunityInternalStage = []OpportunityInternalStage{
	OpportunityInternalStageOpen,
	OpportunityInternalStageClosedWon,
	OpportunityInternalStageClosedLost,
	OpportunityInternalStageSuspended,
}

func (e OpportunityInternalStage) IsValid() bool {
	switch e {
	case OpportunityInternalStageOpen, OpportunityInternalStageClosedWon,
		OpportunityInternalStageClosedLost, OpportunityInternalStageSuspended:
		return true
	}
	return false
}

func (e OpportunityInternalStage) String() string {
	return string(e)
}

func DecodeOpportunityInternalStage(input string) OpportunityInternalStage {
	switch input {
	case "OPEN":
		return OpportunityInternalStageOpen
	case "CLOSED_WON":
		return OpportunityInternalStageClosedWon
	case "CLOSED_LOST":
		return OpportunityInternalStageClosedLost
	case "SUSPENDED":
		return OpportunityInternalStageSuspended
	}
	return ""
}
