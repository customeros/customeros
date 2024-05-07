package enum

type OrganizationStage string

const (
	Lead       OrganizationStage = "LEAD"
	Target     OrganizationStage = "TARGET"
	Interested OrganizationStage = "INTERESTED"
	Engaged    OrganizationStage = "ENGAGED"
	Contracted OrganizationStage = "CONTRACTED"
	Nurture    OrganizationStage = "NURTURE"
	Abandoned  OrganizationStage = "ABANDONED"
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
	case Contracted.String():
		return Contracted
	case Nurture.String():
		return Nurture
	case Abandoned.String():
		return Abandoned
	default:
		return ""
	}
}

func (e OrganizationStage) IsValid() bool {
	switch e {
	case Lead, Target, Interested, Engaged, Contracted, Nurture, Abandoned:
		return true
	}
	return false
}
