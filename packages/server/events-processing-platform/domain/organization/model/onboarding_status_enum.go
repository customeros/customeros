package model

type OnboardingStatusString string

const (
	OnboardingStatusNotApplicable OnboardingStatusString = "NOT_APPLICABLE"
	OnboardingStatusNotStarted    OnboardingStatusString = "NOT_STARTED"
	OnboardingStatusOnTrack       OnboardingStatusString = "ON_TRACK"
	OnboardingStatusLate          OnboardingStatusString = "LATE"
	OnboardingStatusStuck         OnboardingStatusString = "STUCK"
	OnboardingStatusDone          OnboardingStatusString = "DONE"
	OnboardingStatusSuccessful    OnboardingStatusString = "SUCCESSFUL"
)

type OnboardingStatus int32

const (
	NotApplicable OnboardingStatus = iota
	NotStarted
	OnTrack
	Late
	Stuck
	Done
	Successful
)

func (os OnboardingStatus) String() string {
	switch os {
	case NotApplicable:
		return string(OnboardingStatusNotApplicable)
	case NotStarted:
		return string(OnboardingStatusNotStarted)
	case OnTrack:
		return string(OnboardingStatusOnTrack)
	case Late:
		return string(OnboardingStatusLate)
	case Stuck:
		return string(OnboardingStatusStuck)
	case Done:
		return string(OnboardingStatusDone)
	case Successful:
		return string(OnboardingStatusSuccessful)
	default:
		return string(OnboardingStatusNotApplicable)
	}
}
