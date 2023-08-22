package entity

type InteractionEventParticipantDetails struct {
	Type string
}

type InteractionEventParticipant interface {
	IsInteractionEventParticipant()
	InteractionEventParticipantLabel() string
	GetDataloaderKey() string
}

type InteractionEventParticipants []InteractionEventParticipant
