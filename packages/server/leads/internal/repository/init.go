package repository

import (
	"github.com/customeros/customeros/packages/server/leads/internal/database"
)

type Repositories struct {
	APICallLogRepository APICallLogRepository
	IPIntelligence       IPIntelligenceRepository
	Outbox               OutboxRepository
	WebEvent             WebTrackerEventRepository
	WebTracker           WebTrackerRepository
}

func InitRepositories(leadsDB, warehouseDB *database.DbConnections) *Repositories {
	return &Repositories{
		APICallLogRepository: NewAPICallLogRepository(warehouseDB),
		IPIntelligence:       NewIPIntelligenceRepository(leadsDB),
		Outbox:               NewOutboxRepository(leadsDB),
		WebEvent:             NewWebTrackerEventRepository(warehouseDB),
		WebTracker:           NewWebTrackerRepository(leadsDB),
	}
}
