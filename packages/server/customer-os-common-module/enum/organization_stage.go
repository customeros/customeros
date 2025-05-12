package enum

type OrganizationStage string

const (
	Target      OrganizationStage = "TARGET"
	Education   OrganizationStage = "EDUCATION"
	Solution    OrganizationStage = "SOLUTION"
	Evaluation  OrganizationStage = "EVALUATION"
	ReadyToBuy  OrganizationStage = "READY_TO_BUY"
	Opportunity OrganizationStage = "OPPORTUNITY"
	Customer    OrganizationStage = "CUSTOMER"
	NotAFit     OrganizationStage = "NOT_A_FIT"

	// Deprecated
	Lead OrganizationStage = "LEAD"
	// Deprecated
	Engaged OrganizationStage = "ENGAGED"
	// Deprecated
	Unqualified OrganizationStage = "UNQUALIFIED"
	// Deprecated
	Onboarding OrganizationStage = "ONBOARDING"
	// Deprecated
	InitialValue OrganizationStage = "INITIAL_VALUE"
	// Deprecated
	RecurringValue OrganizationStage = "RECURRING_VALUE"
	// Deprecated
	MaxValue OrganizationStage = "MAX_VALUE"
	// Deprecated
	PendingChurn OrganizationStage = "PENDING_CHURN"
	// Deprecated
	Trial OrganizationStage = "TRIAL"
)

func (e OrganizationStage) String() string {
	return string(e)
}

func DecodeOrganizationStage(str string) OrganizationStage {
	switch str {
	case Target.String():
		return Target
	case ReadyToBuy.String():
		return ReadyToBuy
	case Education.String():
		return Education
	case Solution.String():
		return Solution
	case Evaluation.String():
		return Evaluation
	case Opportunity.String():
		return Opportunity
	case Customer.String():
		return Customer
	case NotAFit.String():
		return NotAFit

	case Lead.String():
		return Lead
	case Engaged.String():
		return Engaged
	case Unqualified.String():
		return Unqualified
	case Onboarding.String():
		return Onboarding
	case InitialValue.String():
		return InitialValue
	case RecurringValue.String():
		return RecurringValue
	case MaxValue.String():
		return MaxValue
	case PendingChurn.String():
		return PendingChurn
	case Trial.String():
		return Trial
	default:
		return ""
	}
}
