package entity

type InteractionSessionParticipantDetails struct {
	Type string
}

type InteractionSessionParticipant interface {
	IsInteractionSessionParticipant()
	EntityLabel() string
	GetDataloaderKey() string
	//GetDataLoaderKey() string
}

type InteractionSessionParticipants []InteractionSessionParticipant
