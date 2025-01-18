package enum

import "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

type OnboardingStatus string

const (
	OnboardingStatusNotApplicable OnboardingStatus = "NOT_APPLICABLE"
	OnboardingStatusNotStarted    OnboardingStatus = "NOT_STARTED"
	OnboardingStatusOnTrack       OnboardingStatus = "ON_TRACK"
	OnboardingStatusLate          OnboardingStatus = "LATE"
	OnboardingStatusStuck         OnboardingStatus = "STUCK"
	OnboardingStatusDone          OnboardingStatus = "DONE"
	OnboardingStatusSuccessful    OnboardingStatus = "SUCCESSFUL"
)
const (
	OnboardingStatus_Order_NotStarted = 10
	OnboardingStatus_Order_Stuck      = 20
	OnboardingStatus_Order_Late       = 30
	OnboardingStatus_Order_OnTrack    = 40
	OnboardingStatus_Order_Done       = 50
	OnboardingStatus_Order_Successful = 60
)

func (s OnboardingStatus) String() string {
	return string(s)
}

func (s OnboardingStatus) GetOrder() *int64 {
	switch s {
	case OnboardingStatusNotStarted:
		return utils.Int64Ptr(OnboardingStatus_Order_NotStarted)
	case OnboardingStatusOnTrack:
		return utils.Int64Ptr(OnboardingStatus_Order_OnTrack)
	case OnboardingStatusLate:
		return utils.Int64Ptr(OnboardingStatus_Order_Late)
	case OnboardingStatusStuck:
		return utils.Int64Ptr(OnboardingStatus_Order_Stuck)
	case OnboardingStatusDone:
		return utils.Int64Ptr(OnboardingStatus_Order_Done)
	case OnboardingStatusSuccessful:
		return utils.Int64Ptr(OnboardingStatus_Order_Successful)
	default:
		return nil
	}
}

func (s OnboardingStatus) ReadableStringForActionMessage() string {
	switch s {
	case OnboardingStatusNotApplicable:
		return "Not applicable"
	case OnboardingStatusNotStarted:
		return "Not started"
	case OnboardingStatusOnTrack:
		return "On track"
	case OnboardingStatusLate:
		return "Late"
	case OnboardingStatusStuck:
		return "Stuck"
	case OnboardingStatusDone:
		return "Done"
	case OnboardingStatusSuccessful:
		return "Successful"
	default:
		return string(s)
	}
}

func DecodeOnboardingStatus(s string) OnboardingStatus {
	switch s {
	case "NOT_APPLICABLE":
		return OnboardingStatusNotApplicable
	case "NOT_STARTED":
		return OnboardingStatusNotStarted
	case "ON_TRACK":
		return OnboardingStatusOnTrack
	case "LATE":
		return OnboardingStatusLate
	case "STUCK":
		return OnboardingStatusStuck
	case "DONE":
		return OnboardingStatusDone
	case "SUCCESSFUL":
		return OnboardingStatusSuccessful
	default:
		return OnboardingStatusNotApplicable
	}
}
