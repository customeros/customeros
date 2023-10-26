package model

type InteractionEventDataFields struct {
	Content         string
	ContentType     string
	Identifier      string
	EventType       string
	Channel         string
	ChannelData     string
	PartOfIssueId   *string
	PartOfSessionId *string
	Hide            bool
}
