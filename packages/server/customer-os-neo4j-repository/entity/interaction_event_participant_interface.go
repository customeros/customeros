package entity

type InteractionEventParticipantDetails struct {
	Type string
}

type InteractionEventParticipant interface {
	IsInteractionEventParticipant()
	EntityLabel() string
	GetDataloaderKey() string
}

type InteractionEventParticipants []InteractionEventParticipant
