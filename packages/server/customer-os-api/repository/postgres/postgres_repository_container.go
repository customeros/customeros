package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository/postgres/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositoryContainer struct {
	AppKeyRepo AppKeyRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositoryContainer {
	p := &PostgresRepositoryContainer{
		AppKeyRepo: NewAppKeyRepo(db),
	}

	err := db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
