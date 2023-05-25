package entity

import (
	"fmt"
	"time"
)

type PlayerRelation string

const (
	IDENTIFIES PlayerRelation = "IDENTIFIES"
)

type PlayerEntity struct {
	Id            string
	IdentityId    *string
	AuthId        string
	Provider      string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	DataloaderKey string
}

func (player PlayerEntity) ToString() string {
	return fmt.Sprintf("id: %s\nAuthID: %s\nidentityId: %s", player.Id, player.AuthId, *player.IdentityId)
}

type PersonEntities []PlayerEntity

func (PlayerEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Player,
	}
}
