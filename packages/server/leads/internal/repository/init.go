package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/config"
	"github.com/customeros/customeros/packages/server/leads/internal/database"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
)

type Repositories struct {
	APICallLogRepository APICallLogRepository
	IPIntelligence       IPIntelligenceRepository
	Outbox               OutboxRepository
	WebTrackerEvent      WebTrackerEventRepository
	WebTracker           WebTrackerRepository
}

func InitRepositories(leadsDB, warehouseDB *database.DbConnections) *Repositories {
	return &Repositories{
		APICallLogRepository: NewAPICallLogRepository(warehouseDB),
		IPIntelligence:       NewIPIntelligenceRepository(leadsDB),
		Outbox:               NewOutboxRepository(leadsDB),
		WebTrackerEvent:      NewWebTrackerEventRepository(warehouseDB),
		WebTracker:           NewWebTrackerRepository(leadsDB),
	}
}

func MigrateLeadsDB(dbConfig *config.LeadsDatabaseConfig, leadsDB *gorm.DB) error {
	db, err := leadsDB.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)

	err = leadsDB.AutoMigrate(
		&models.IPIntelligence{},
		&models.OutboxEvent{},
		&models.WebTracker{},
	)

	db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	db.SetMaxOpenConns(dbConfig.MaxConn)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Minute)

	return err
}

func MigrateDataWarehouse(dbConfig *config.DataWarehouseConfig, warehouseDB *gorm.DB) error {
	db, err := warehouseDB.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)

	err = warehouseDB.AutoMigrate(
		&models.APICallLog{},
		&models.WebTrackerEvent{},
	)

	db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	db.SetMaxOpenConns(dbConfig.MaxConn)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Minute)

	return err
}
