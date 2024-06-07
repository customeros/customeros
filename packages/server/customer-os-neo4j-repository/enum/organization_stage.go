package enum

type OrganizationStage string

const (
	Lead           OrganizationStage = "LEAD"
	Target         OrganizationStage = "TARGET"
	Engaged        OrganizationStage = "ENGAGED"
	Unqualified    OrganizationStage = "UNQUALIFIED"
	ReadyToBuy     OrganizationStage = "READY_TO_BUY"
	Onboarding     OrganizationStage = "ONBOARDING"
	InitialValue   OrganizationStage = "INITIAL_VALUE"
	RecurringValue OrganizationStage = "RECURRING_VALUE"
	MaxValue       OrganizationStage = "MAX_VALUE"
	PendingChurn   OrganizationStage = "PENDING_CHURN"
	Trial          OrganizationStage = "TRIAL"
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
	case Engaged.String():
		return Engaged
	case Unqualified.String():
		return Unqualified
	case ReadyToBuy.String():
		return ReadyToBuy
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

func (e OrganizationStage) IsValid() bool {
	switch e {
	case Lead, Target, Engaged, Unqualified, ReadyToBuy, Onboarding, InitialValue, RecurringValue, MaxValue, PendingChurn, Trial:
		return true
	}
	return false
}
