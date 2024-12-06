package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
)

type Repositories struct {
	UserGCalImportStateRepository UserGCalImportStateRepository
	RawCalendarEventRepository    RawCalendarEventRepository
}

func InitRepos(driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		UserGCalImportStateRepository: NewUserGCalImportStateRepository(postgresDB.AsyncGormDB),
		RawCalendarEventRepository:    NewRawCalendarEventRepository(postgresDB.AsyncGormDB),
	}

	return &repositories
}
