package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

func MapEntityToMeeting(entity *entity.MeetingEntity) *model.Meeting {
	if entity == nil {
		return nil
	}

	meeting := model.Meeting{
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
	if entity.Status != nil {
		meeting.Status = MapMeetingStatusToModel(*entity.Status)
	} else {
		meeting.Status = model.MeetingStatusUndefined
	}
	return &meeting
}

func MapMeetingInputToEntity(model *model.MeetingUpdateInput) *entity.MeetingEntity {
	if model == nil {
		return nil
	}

	meetingEntity := entity.MeetingEntity{
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

	if model.Status != nil {
		status := MapMeetingStatusFromModel(*model.Status)
		meetingEntity.Status = &status
	} else {
		status := entity.MeetingStatusUndefined
		meetingEntity.Status = &status
	}

	return &meetingEntity
}

func MapMeetingToEntity(model *model.MeetingInput) *entity.MeetingEntity {
	if model == nil {
		return nil
	}
	var createdAt time.Time
	if model.CreatedAt != nil {
		createdAt = model.CreatedAt.UTC()
	} else {
		createdAt = utils.Now()
	}

	meetingEntity := entity.MeetingEntity{
		CreatedAt:          createdAt,
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
	if model.Status != nil {
		status := MapMeetingStatusFromModel(*model.Status)
		meetingEntity.Status = &status
	} else {
		status := entity.MeetingStatusUndefined
		meetingEntity.Status = &status
	}

	return &meetingEntity
}

func MapEntitiesToMeetings(entities *entity.MeetingEntities) []*model.Meeting {
	var meetings []*model.Meeting
	for _, meetingEntity := range *entities {
		meetings = append(meetings, MapEntityToMeeting(&meetingEntity))
	}
	return meetings
}
