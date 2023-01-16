package entity

type ConversationInitiator struct {
	Id             string
	ExternalId     string
	ExternalSystem string
	FirstName      string
	LastName       string
	Email          string
	InitiatorType  SenderType
}
