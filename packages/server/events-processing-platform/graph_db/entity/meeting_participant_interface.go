package entity

type MeetingParticipant interface {
	IsMeetingParticipant()
	MeetingParticipantLabel() string
	GetDataloaderKey() string
}

type MeetingParticipants []MeetingParticipant
