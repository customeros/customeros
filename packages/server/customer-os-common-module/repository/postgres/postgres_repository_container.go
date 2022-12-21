package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresCommonRepositoryContainer struct {
	AppKeyRepo AppKeyRepository
}

func InitCommonRepositories(db *gorm.DB) *PostgresCommonRepositoryContainer {
	p := &PostgresCommonRepositoryContainer{
		AppKeyRepo: NewAppKeyRepo(db),
	}

	err := db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
