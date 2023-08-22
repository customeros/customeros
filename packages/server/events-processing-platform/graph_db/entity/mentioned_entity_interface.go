package entity

type MentionedEntity interface {
	IsMentionedEntity()
	MentionedEntityLabel() string
	GetDataloaderKey() string
}

type MentionedEntities []MentionedEntity
