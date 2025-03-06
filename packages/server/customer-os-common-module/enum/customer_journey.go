package enum

import (
	"fmt"
	"strings"
)

type CustomerJourneyStage string

const (
	CustomerJourneyProblemRecognition  CustomerJourneyStage = "Problem Recognition"
	CustomerJourneySolutionEvaluation  CustomerJourneyStage = "Solution Evaluation"
	CustomerJourneyDecisionPreparation CustomerJourneyStage = "Decision Preparation"
	CustomerJourneyOnboarding          CustomerJourneyStage = "Onboarding"
	CustomerJourneyOutcomeAttainment   CustomerJourneyStage = "Outcome Attainment"
	CustomerJourneySustainedSuccess    CustomerJourneyStage = "Sustained Success"
	CustomerJourneyUnknown             CustomerJourneyStage = ""
)

func (a CustomerJourneyStage) String() string {
	return string(a)
}

func GetCustomerJourneyStage(s string) (CustomerJourneyStage, error) {
	cleanValue := strings.TrimSpace(s)
	switch CustomerJourneyStage(cleanValue) {
	case
		CustomerJourneyProblemRecognition,
		CustomerJourneySolutionEvaluation,
		CustomerJourneyDecisionPreparation,
		CustomerJourneyOnboarding,
		CustomerJourneyOutcomeAttainment,
		CustomerJourneySustainedSuccess:
		return CustomerJourneyStage(cleanValue), nil
	default:
		return CustomerJourneyUnknown, fmt.Errorf("invalid CustomerJourneyStage: %s", s)
	}
}

// CustomerJourneyStageValidator implements the EnumValidator interface
type CustomerJourneyStageValidator struct{}

func GetCustomerJourneyStageValidator() *CustomerJourneyStageValidator {
	return &CustomerJourneyStageValidator{}
}

func (v *CustomerJourneyStageValidator) IsValid(value string) bool {
	cleanValue := strings.TrimSpace(value)
	switch CustomerJourneyStage(cleanValue) {
	case
		CustomerJourneyProblemRecognition,
		CustomerJourneySolutionEvaluation,
		CustomerJourneyDecisionPreparation,
		CustomerJourneyOnboarding,
		CustomerJourneyOutcomeAttainment,
		CustomerJourneySustainedSuccess:
		return true
	default:
		return false
	}
}

func (v *CustomerJourneyStageValidator) ValidValues() []string {
	return []string{
		string(CustomerJourneyProblemRecognition),
		string(CustomerJourneySolutionEvaluation),
		string(CustomerJourneyDecisionPreparation),
		string(CustomerJourneyOnboarding),
		string(CustomerJourneyOutcomeAttainment),
		string(CustomerJourneySustainedSuccess),
	}
}

func (v *CustomerJourneyStageValidator) ParseJourneyStage(value string) (CustomerJourneyStage, error) {
	cleanValue := strings.TrimSpace(value)
	stage := CustomerJourneyStage(cleanValue)
	if !v.IsValid(cleanValue) {
		return CustomerJourneyUnknown, fmt.Errorf("invalid customer journey stage: %s", value)
	}
	return stage, nil
}
