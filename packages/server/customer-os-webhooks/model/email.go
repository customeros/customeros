package model

type EmailData struct {
	BaseData
	Content     string    `json:"content,omitempty"`
	ContentType string    `json:"contentType,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	SentBy      string    `json:"sentBy,omitempty"`
	SentTo      []string  `json:"sentTo,omitempty"`
	Bcc         []string  `json:"bcc,omitempty"`
	Cc          []string  `json:"cc,omitempty"`
	Channel     string    `json:"channel,omitempty"`
	ChannelData string    `json:"channelData,omitempty"`
	Identifier  string    `json:"identifier,omitempty"`
	EventType   string    `json:"eventType,omitempty"`
	Hide        bool      `json:"hide,omitempty"`
	BelongsTo   BelongsTo `json:"belongsTo,omitempty"`
	// in sent to or sent by at least 1 contact should be available in the system
	ContactRequired bool                   `json:"contactRequired,omitempty"`
	ParentRequired  bool                   `json:"parentRequired,omitempty"`
	SessionDetails  InteractionSessionData `json:"sessionDetails,omitempty"`
}
