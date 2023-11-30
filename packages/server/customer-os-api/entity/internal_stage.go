package entity

type InternalStage string

const (
	InternalStageOpen       InternalStage = "OPEN"
	InternalStageEvaluating InternalStage = "EVALUATING"
	InternalStageClosedWon  InternalStage = "CLOSED_WON"
	InternalStageClosedLost InternalStage = "CLOSED_LOST"
)

var AllInternalStages = []InternalStage{
	InternalStageOpen,
	InternalStageEvaluating,
	InternalStageClosedWon,
	InternalStageClosedLost,
}

func GetInternalStage(s string) InternalStage {
	if IsValidInternalStage(s) {
		return InternalStage(s)
	}
	return InternalStageOpen
}

func IsValidInternalStage(s string) bool {
	for _, ms := range AllInternalStages {
		if ms == InternalStage(s) {
			return true
		}
	}
	return false
}
