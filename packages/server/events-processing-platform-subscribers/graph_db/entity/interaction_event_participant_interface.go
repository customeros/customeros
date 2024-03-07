package entity

type InteractionEventParticipantDetails struct {
	Type string
}

type InteractionEventParticipant interface {
	IsInteractionEventParticipant()
	ParticipantLabel() string
}

type InteractionEventParticipants []InteractionEventParticipant
