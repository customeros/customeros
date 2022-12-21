package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/file-storage-api/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositoryContainer struct {
	FileRepo FileRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositoryContainer {
	p := &PostgresRepositoryContainer{
		FileRepo: NewFileRepo(db),
	}

	err := db.AutoMigrate(&entity.File{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
