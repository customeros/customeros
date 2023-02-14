package entity

import (
	"time"
)

type ExternalReferenceRelationship struct {
	Id               string
	SyncDate         time.Time
	ExternalSystemId string
}
