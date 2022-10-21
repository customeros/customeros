package mapper

import (
	"github.com.openline-ai.customer-os-analytics-api/graph/model"
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
)

func MapSessions(sessionEntities *entity.SessionEntities) []*model.AppSession {
	var sessions []*model.AppSession
	for _, sessionEntity := range *sessionEntities {
		sessions = append(sessions, MapSession(&sessionEntity))
	}
	return sessions
}

func MapSession(sessionEntity *entity.SessionEntity) *model.AppSession {
	return &model.AppSession{
		ID:             sessionEntity.ID,
		Country:        sessionEntity.Country,
		Region:         sessionEntity.Region,
		City:           sessionEntity.City,
		ReferrerSource: sessionEntity.ReferrerSource,

		UtmCampaign: sessionEntity.UtmCampaign,
		UtmContent:  sessionEntity.UtmContent,
		UtmMedium:   sessionEntity.UtmMedium,
		UtmNetwork:  sessionEntity.UtmNetwork,
		UtmSource:   sessionEntity.UtmSource,
		UtmTerm:     sessionEntity.UtmTerm,

		DeviceBrand: sessionEntity.DeviceBrand,
		DeviceName:  sessionEntity.DeviceName,
		DeviceClass: sessionEntity.DeviceClass,

		AgentName:    sessionEntity.AgentName,
		AgentVersion: sessionEntity.AgentVersion,

		StartedAt:   sessionEntity.Start,
		EndedAt:     sessionEntity.End,
		EngagedTime: sessionEntity.EngagedTime,
	}
}
