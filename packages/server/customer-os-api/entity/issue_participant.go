package entity

type IssueParticipantDetails struct {
	Type string
}

type IssueParticipant interface {
	IsIssueParticipant()
	EntityLabel() string
	GetDataloaderKey() string
}

type IssueParticipants []IssueParticipant
