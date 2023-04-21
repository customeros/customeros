package entity

type MeetingParticipantDetails struct {
	Type string
}

type MeetingParticipant interface {
	IsMeetingParticipant()
	MeetingParticipantLabel() string
	GetDataloaderKey() string
}

type MeetingParticipants []MeetingParticipant
