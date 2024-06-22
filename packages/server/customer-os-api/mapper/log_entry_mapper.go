package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToLogEntry(entity *neo4jentity.LogEntryEntity) *model.LogEntry {
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

func MapEntitiesToLogEntries(entities *neo4jentity.LogEntryEntities) []*model.LogEntry {
	var logEntries []*model.LogEntry
	for _, logEntryEntity := range *entities {
		logEntries = append(logEntries, MapEntityToLogEntry(&logEntryEntity))
	}
	return logEntries
}
