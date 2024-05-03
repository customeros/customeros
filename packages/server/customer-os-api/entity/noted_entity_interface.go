package entity

type NotedEntity interface {
	IsNotedEntity()
	EntityLabel() string
	GetDataloaderKey() string
}

type NotedEntities []NotedEntity
