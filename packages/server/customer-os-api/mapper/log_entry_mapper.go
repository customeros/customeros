package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToLogEntry(entity *entity.LogEntryEntity) *model.LogEntry {
	logEntry := model.LogEntry{
		ID:            entity.Id,
		Content:       utils.StringPtr(entity.Content),
		ContentType:   utils.StringPtr(entity.ContentType),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		StartedAt:     entity.StartedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
	return &logEntry
}

func MapEntitiesToLogEntries(entities *entity.LogEntryEntities) []*model.LogEntry {
	var logEntries []*model.LogEntry
	for _, logEntryEntity := range *entities {
		logEntries = append(logEntries, MapEntityToLogEntry(&logEntryEntity))
	}
	return logEntries
}
