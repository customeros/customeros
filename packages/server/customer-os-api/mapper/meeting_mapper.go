package mapper

import (
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToMeeting(entity *neo4jentity.MeetingEntity) *model.Meeting {
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

func MapMeetingInputToEntity(model *model.MeetingUpdateInput) *neo4jentity.MeetingEntity {
	if model == nil {
		return nil
	}

	meetingEntity := neo4jentity.MeetingEntity{
		CreatedAt:          utils.Now(),
		Name:               model.Name,
		AppSource:          utils.IfNotNilStringWithDefault(model.AppSource, constants.AppSourceCustomerOsApi),
		ConferenceUrl:      model.ConferenceURL,
		MeetingExternalUrl: model.MeetingExternalURL,
		StartedAt:          model.StartedAt,
		EndedAt:            model.EndedAt,
		Agenda:             model.Agenda,
		AgendaContentType:  model.AgendaContentType,
		Source:             neo4jentity.DataSourceOpenline,
		SourceOfTruth:      neo4jentity.DataSourceOpenline,
	}

	if model.Status != nil {
		status := MapMeetingStatusFromModel(*model.Status)
		meetingEntity.Status = &status
	} else {
		status := neo4jenum.MeetingStatusUndefined
		meetingEntity.Status = &status
	}

	return &meetingEntity
}

func MapMeetingToEntity(model *model.MeetingInput) *neo4jentity.MeetingEntity {
	if model == nil {
		return nil
	}
	var createdAt time.Time
	if model.CreatedAt != nil {
		createdAt = model.CreatedAt.UTC()
	} else {
		createdAt = utils.Now()
	}

	meetingEntity := neo4jentity.MeetingEntity{
		CreatedAt:          createdAt,
		Name:               model.Name,
		AppSource:          utils.IfNotNilStringWithDefault(model.AppSource, constants.AppSourceCustomerOsApi),
		ConferenceUrl:      model.ConferenceURL,
		MeetingExternalUrl: model.MeetingExternalURL,
		StartedAt:          model.StartedAt,
		EndedAt:            model.EndedAt,
		Agenda:             model.Agenda,
		AgendaContentType:  model.AgendaContentType,
		Source:             neo4jentity.DataSourceOpenline,
		SourceOfTruth:      neo4jentity.DataSourceOpenline,
	}
	if model.Status != nil {
		status := MapMeetingStatusFromModel(*model.Status)
		meetingEntity.Status = &status
	} else {
		status := neo4jenum.MeetingStatusUndefined
		meetingEntity.Status = &status
	}

	return &meetingEntity
}

func MapEntitiesToMeetings(entities *neo4jentity.MeetingEntities) []*model.Meeting {
	var meetings []*model.Meeting
	for _, meetingEntity := range *entities {
		meetings = append(meetings, MapEntityToMeeting(&meetingEntity))
	}
	return meetings
}
