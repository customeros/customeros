package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositories struct {
	FileRepository FileRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		FileRepository: NewFileRepo(db),
	}

	err := db.AutoMigrate(&entity.File{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
