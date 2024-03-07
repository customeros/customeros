package entity

type NotedEntity interface {
	IsNotedEntity()
	NotedEntityLabel() string
	GetDataloaderKey() string
}

type NotedEntities []NotedEntity
