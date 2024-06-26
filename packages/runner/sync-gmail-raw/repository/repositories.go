package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	TenantRepository TenantRepository
	UserRepository   UserRepository
	EmailRepository  EmailRepository

	UserGCalImportStateRepository UserGCalImportStateRepository

	RawCalendarEventRepository RawCalendarEventRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		TenantRepository: NewTenantRepository(driver),
		UserRepository:   NewUserRepository(driver),
		EmailRepository:  NewEmailRepository(driver),

		UserGCalImportStateRepository: NewUserGCalImportStateRepository(gormDb),

		RawCalendarEventRepository: NewRawCalendarEventRepository(gormDb),
	}

	var err error

	err = gormDb.AutoMigrate(&entity.UserGCalImportState{})
	if err != nil {
		panic(err)
	}
	err = gormDb.AutoMigrate(&entity.RawCalendarEvent{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
