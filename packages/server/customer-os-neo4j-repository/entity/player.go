package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type PlayerRelation string

const (
	IDENTIFIES PlayerRelation = "IDENTIFIES"
)

type PlayerEntity struct {
	Id            string
	IdentityId    string
	AuthId        string
	Provider      string
	Source        string
	SourceOfTruth string
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	DataloaderKey string
}

type PersonEntities []PlayerEntity

func (PlayerEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelPlayer,
	}
}
