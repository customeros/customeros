package enum

type CustomerJourneyStage string

const (
	JourneyStageTarget        CustomerJourneyStage = "target"
	JourneyStageAwareness     CustomerJourneyStage = "awareness"
	JourneyStageConsideration CustomerJourneyStage = "consideration"
	JourneyStageOpportunity   CustomerJourneyStage = "opportunity"
	JourneyStageCustomer      CustomerJourneyStage = "customer"
	JourneyStageNotAFit       CustomerJourneyStage = "not a fit"
)

func (e CustomerJourneyStage) String() string {
	return string(e)
}
