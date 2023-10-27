package model

type InteractionEventDataFields struct {
	Content            string
	ContentType        string
	Identifier         string
	EventType          string
	Channel            string
	ChannelData        string
	BelongsToIssueId   *string
	BelongsToSessionId *string
	Hide               bool
}
