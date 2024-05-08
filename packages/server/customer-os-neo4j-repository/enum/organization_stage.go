package enum

type OrganizationStage string

const (
	Lead       OrganizationStage = "LEAD"
	Target     OrganizationStage = "TARGET"
	Interested OrganizationStage = "INTERESTED"
	Engaged    OrganizationStage = "ENGAGED"
	Nurture    OrganizationStage = "NURTURE"
	ClosedLost OrganizationStage = "CLOSED_LOST"
	ClosedWon  OrganizationStage = "CLOSED_WON"
	NotAFit    OrganizationStage = "NOT_A_FIT"
)

func (e OrganizationStage) String() string {
	return string(e)
}

func DecodeOrganizationStage(str string) OrganizationStage {
	switch str {
	case Lead.String():
		return Lead
	case Target.String():
		return Target
	case Interested.String():
		return Interested
	case Engaged.String():
		return Engaged
	case ClosedLost.String():
		return ClosedLost
	case ClosedWon.String():
		return ClosedWon
	case NotAFit.String():
		return NotAFit
	case Nurture.String():
		return Nurture
	default:
		return ""
	}
}

func (e OrganizationStage) IsValid() bool {
	switch e {
	case Lead, Target, Interested, Engaged, ClosedWon, ClosedLost, NotAFit, Nurture:
		return true
	}
	return false
}
