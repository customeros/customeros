package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type ContractEntity struct {
	Id                              string
	CreatedAt                       time.Time
	UpdatedAt                       time.Time
	Source                          neo4jentity.DataSource
	SourceOfTruth                   neo4jentity.DataSource
	AppSource                       string
	Name                            string
	ContractUrl                     string
	Status                          string
	RenewalCycle                    string
	RenewalPeriods                  *int64
	SignedAt                        *time.Time
	ServiceStartedAt                *time.Time
	EndedAt                         *time.Time
	TriggeredOnboardingStatusChange bool
}

func (c ContractEntity) IsEnded() bool {
	return c.EndedAt != nil && c.EndedAt.Before(utils.Now())
}

func (c ContractEntity) IsSigned() bool {
	return c.SignedAt != nil && c.SignedAt.Before(utils.Now())
}

func (c ContractEntity) IsServiceStarted() bool {
	return c.ServiceStartedAt != nil && c.ServiceStartedAt.Before(utils.Now())
}

func (c ContractEntity) IsEligibleToStartOnboarding() bool {
	return !c.TriggeredOnboardingStatusChange && (c.IsSigned() || c.IsServiceStarted()) && !c.IsEnded()
}
