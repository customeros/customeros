package entity

type MeetingParticipant interface {
	IsMeetingParticipant()
	EntityLabel() string
	GetDataloaderKey() string
}

type MeetingParticipants []MeetingParticipant
