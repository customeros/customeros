package enum

import "fmt"

type CustomerJourneyStage string

const (
	CustomerJourneyProblemRecognition  CustomerJourneyStage = "Problem Recognition"
	CustomerJourneySolutionEvaluation  CustomerJourneyStage = "Solution Evaluation"
	CustomerJourneyDecisionPreparation CustomerJourneyStage = "Decision Preparation"
	CustomerJourneyOnboarding          CustomerJourneyStage = "Onboarding"
	CustomerJourneyOutcomeAttainment   CustomerJourneyStage = "Outcome Attainment"
	CustomerJourneySustainedSuccess    CustomerJourneyStage = "Sustained Success"
)

func (a CustomerJourneyStage) String() string {
	return string(a)
}

func GetCustomerJourneyStage(s string) (CustomerJourneyStage, error) {
	switch CustomerJourneyStage(s) {
	case
		CustomerJourneyProblemRecognition,
		CustomerJourneySolutionEvaluation,
		CustomerJourneyDecisionPreparation,
		CustomerJourneyOnboarding,
		CustomerJourneyOutcomeAttainment,
		CustomerJourneySustainedSuccess:
		return CustomerJourneyStage(s), nil

	default:
		return "", fmt.Errorf("invalid CustomerJourneyStage: %s", s)
	}
}
