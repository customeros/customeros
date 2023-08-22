package entity

type InteractionSessionParticipantDetails struct {
	Type string
}

type InteractionSessionParticipant interface {
	IsInteractionSessionParticipant()
	InteractionSessionParticipantLabel() string
	GetDataloaderKey() string
}

type InteractionSessionParticipants []InteractionSessionParticipant
