package entity

type IssueParticipantDetails struct {
	Type string
}

type IssueParticipant interface {
	IsIssueParticipant()
	ParticipantLabel() string
	GetDataloaderKey() string
}

type IssueParticipants []IssueParticipant
