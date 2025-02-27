package neo4j_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"time"
)

type AuthenticationEntity struct {
	Id         string
	CreatedAt  time.Time
	IdentityId string
	AuthId     string
	Provider   string
}

func (AuthenticationEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelAuthentication,
	}
}

type AuthenticationUserEntity struct {
	Id            string
	FirstName     string
	LastName      string
	DefaultTenant string
	CurrentTenant string
}

func (AuthenticationUserEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelAuthenticationUser,
	}
}
