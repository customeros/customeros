package entity

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

var AllOnboardingStatuses = []OnboardingStatus{
	OnboardingStatusNotApplicable,
	OnboardingStatusNotStarted,
	OnboardingStatusOnTrack,
	OnboardingStatusLate,
	OnboardingStatusStuck,
	OnboardingStatusDone,
	OnboardingStatusSuccessful,
}

func GetOnboardingStatus(s string) OnboardingStatus {
	if IsValidOnboardingStatus(s) {
		return OnboardingStatus(s)
	}
	return OnboardingStatusNotApplicable
}

func IsValidOnboardingStatus(s string) bool {
	for _, ms := range AllOnboardingStatuses {
		if ms == OnboardingStatus(s) {
			return true
		}
	}
	return false
}
