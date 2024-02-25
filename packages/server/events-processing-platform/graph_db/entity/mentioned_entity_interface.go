package entity

// Deprecated
type MentionedEntity interface {
	IsMentionedEntity()
	MentionedEntityLabel() string
	GetDataloaderKey() string
}

type MentionedEntities []MentionedEntity
