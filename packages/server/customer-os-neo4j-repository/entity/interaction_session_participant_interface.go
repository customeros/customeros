package entity

type InteractionSessionParticipantDetails struct {
	Type string
}

type InteractionSessionParticipant interface {
	IsInteractionSessionParticipant()
	EntityLabel() string
	GetDataloaderKey() string
}

type InteractionSessionParticipants []InteractionSessionParticipant
