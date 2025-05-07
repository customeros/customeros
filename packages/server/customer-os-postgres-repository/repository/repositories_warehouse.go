package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/database"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type WarehouseRepositories struct {
	ReadDb  *gorm.DB
	WriteDb *gorm.DB
}

func InitWarehouseRepositories(dbConns *database.DbConnections) *WarehouseRepositories {
	repositories := &WarehouseRepositories{
		ReadDb:  dbConns.ReadDB,
		WriteDb: dbConns.WriteDB,
	}

	return repositories
}

func (r *WarehouseRepositories) MigrateDataWarehouse() error {
	if r.WriteDb == nil {
		return nil
	}
	err := r.WriteDb.AutoMigrate(
		&postgres_entity.APICallLog{},
	)
	return err
}
