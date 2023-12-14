package model

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
