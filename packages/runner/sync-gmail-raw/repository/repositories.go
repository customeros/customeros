package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	TenantRepository                   TenantRepository
	UserRepository                     UserRepository
	EmailRepository                    EmailRepository
	ApiKeyRepository                   ApiKeyRepository
	UserGmailImportPageTokenRepository UserGmailImportPageTokenRepository
	RawEmailRepository                 RawEmailRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{

		TenantRepository:                   NewTenantRepository(driver),
		UserRepository:                     NewUserRepository(driver),
		EmailRepository:                    NewEmailRepository(driver),
		ApiKeyRepository:                   NewApiKeyRepository(gormDb),
		UserGmailImportPageTokenRepository: NewUserGmailImportPageTokenRepository(gormDb),
		RawEmailRepository:                 NewRawEmailRepository(gormDb),
	}

	var err error

	err = gormDb.AutoMigrate(&entity.UserGmailImportPageToken{})
	if err != nil {
		panic(err)
	}

	err = gormDb.AutoMigrate(&entity.RawEmail{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
