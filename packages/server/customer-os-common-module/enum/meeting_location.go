package enum

type MeetingLocation string

const (
	MeetingLocationGoogleMeet MeetingLocation = "Google Meet"
)

var AllMeetingLocations = []MeetingLocation{
	MeetingLocationGoogleMeet,
}

func IsValidMeetingLocation(s string) bool {
	for _, loc := range AllMeetingLocations {
		if loc == MeetingLocation(s) {
			return true
		}
	}
	return false
}

func (m MeetingLocation) String() string {
	return string(m)
}

func IsGoogleMeet(location string) bool {
	return location == MeetingLocationGoogleMeet.String()
}
