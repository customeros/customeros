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
	default:
		return ""
	}
}
