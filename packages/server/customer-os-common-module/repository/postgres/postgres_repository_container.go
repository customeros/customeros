package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresCommonRepositoryContainer struct {
	AppKeyRepo       AppKeyRepository
	UserToTenantRepo UserToTenantRepository
}

func InitCommonRepositories(db *gorm.DB) *PostgresCommonRepositoryContainer {
	p := &PostgresCommonRepositoryContainer{
		AppKeyRepo:       NewAppKeyRepo(db),
		UserToTenantRepo: NewUserToTenantRepo(db),
	}

	var err error

	err = db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.UserToTenant{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
