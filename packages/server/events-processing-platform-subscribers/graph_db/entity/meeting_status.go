package entity

type MeetingStatus string

const (
	MeetingStatusUndefined MeetingStatus = ""
	MeetingStatusAccepted  MeetingStatus = "ACCEPTED"
	MeetingStatusCanceled  MeetingStatus = "CANCELED"
)

var AllMeetingStatuses = []MeetingStatus{
	MeetingStatusAccepted,
	MeetingStatusCanceled,
}

func GetMeetingStatus(s string) MeetingStatus {
	if IsValidMeetingStatus(s) {
		return MeetingStatus(s)
	}
	return MeetingStatusUndefined
}

func IsValidMeetingStatus(s string) bool {
	for _, ms := range AllMeetingStatuses {
		if ms == MeetingStatus(s) {
			return true
		}
	}
	return false
}
