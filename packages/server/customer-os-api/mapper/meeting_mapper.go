package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToMeeting(entity *entity.MeetingEntity) *model.Meeting {
	if entity == nil {
		return nil
	}
	return &model.Meeting{
		ID:                 entity.Id,
		Name:               entity.Name,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		StartedAt:          entity.StartedAt,
		EndedAt:            entity.EndedAt,
		ConferenceURL:      entity.ConferenceUrl,
		MeetingExternalURL: entity.MeetingExternalUrl,
		Agenda:             entity.Agenda,
		AgendaContentType:  entity.AgendaContentType,
		AppSource:          entity.AppSource,
		Source:             MapDataSourceToModel(entity.Source),
		SourceOfTruth:      MapDataSourceToModel(entity.SourceOfTruth),
	}
}

func MapMeetingInputToEntity(model *model.MeetingUpdateInput) *entity.MeetingEntity {
	if model == nil {
		return nil
	}
	return &entity.MeetingEntity{
		CreatedAt:          utils.Now(),
		Name:               model.Name,
		AppSource:          model.AppSource,
		ConferenceUrl:      model.ConferenceURL,
		MeetingExternalUrl: model.MeetingExternalURL,
		StartedAt:          model.StartedAt,
		EndedAt:            model.EndedAt,
		Agenda:             model.Agenda,
		AgendaContentType:  model.AgendaContentType,
		Source:             entity.DataSourceOpenline,
		SourceOfTruth:      entity.DataSourceOpenline,
	}
}

func MapMeetingToEntity(model *model.MeetingInput) *entity.MeetingEntity {
	if model == nil {
		return nil
	}
	return &entity.MeetingEntity{
		CreatedAt:          utils.Now(),
		Name:               model.Name,
		AppSource:          model.AppSource,
		ConferenceUrl:      model.ConferenceURL,
		MeetingExternalUrl: model.MeetingExternalURL,
		StartedAt:          model.StartedAt,
		EndedAt:            model.EndedAt,
		Agenda:             model.Agenda,
		AgendaContentType:  model.AgendaContentType,
		Source:             entity.DataSourceOpenline,
		SourceOfTruth:      entity.DataSourceOpenline,
	}
}

func MapEntitiesToMeetings(entities *entity.MeetingEntities) []*model.Meeting {
	var meetings []*model.Meeting
	for _, meetingEntity := range *entities {
		meetings = append(meetings, MapEntityToMeeting(&meetingEntity))
	}
	return meetings
}
