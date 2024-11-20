package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type Repositories struct {
	TenantRepository              TenantRepository
	UserRepository                UserRepository
	EmailRepository               EmailRepository
	UserGCalImportStateRepository UserGCalImportStateRepository
	RawCalendarEventRepository    RawCalendarEventRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		TenantRepository:              NewTenantRepository(driver),
		UserRepository:                NewUserRepository(driver),
		EmailRepository:               NewEmailRepository(driver),
		UserGCalImportStateRepository: NewUserGCalImportStateRepository(gormDb),
		RawCalendarEventRepository:    NewRawCalendarEventRepository(gormDb),
	}

	return &repositories
}
