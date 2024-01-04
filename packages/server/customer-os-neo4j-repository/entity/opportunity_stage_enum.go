package entity

type OpportunityInternalStage string

const (
	OpportunityInternalStageOpen       OpportunityInternalStage = "OPEN"
	OpportunityInternalStageEvaluating OpportunityInternalStage = "EVALUATING"
	OpportunityInternalStageClosedWon  OpportunityInternalStage = "CLOSED_WON"
	OpportunityInternalStageClosedLost OpportunityInternalStage = "CLOSED_LOST"
)

var AllOpportunityInternalStage = []OpportunityInternalStage{
	OpportunityInternalStageOpen,
	OpportunityInternalStageEvaluating,
	OpportunityInternalStageClosedWon,
	OpportunityInternalStageClosedLost,
}

func (e OpportunityInternalStage) IsValid() bool {
	switch e {
	case OpportunityInternalStageOpen, OpportunityInternalStageEvaluating, OpportunityInternalStageClosedWon, OpportunityInternalStageClosedLost:
		return true
	}
	return false
}

func (e OpportunityInternalStage) String() string {
	return string(e)
}
