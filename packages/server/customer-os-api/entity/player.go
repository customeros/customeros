package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
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
