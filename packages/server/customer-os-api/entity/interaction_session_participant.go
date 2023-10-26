package entity

type InteractionSessionParticipantDetails struct {
	Type string
}

type InteractionSessionParticipant interface {
	IsInteractionSessionParticipant()
	ParticipantLabel() string
	GetDataloaderKey() string
}

type InteractionSessionParticipants []InteractionSessionParticipant
