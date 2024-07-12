package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type MeetingEntity struct {
	DataLoaderKey
	Id                 string
	Name               *string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	StartedAt          *time.Time
	EndedAt            *time.Time
	ConferenceUrl      *string
	MeetingExternalUrl *string
	AppSource          string
	Agenda             *string
	AgendaContentType  *string
	Source             DataSource
	SourceOfTruth      DataSource
	Recording          *string
	Status             *enum.MeetingStatus
}

func (MeetingEntity) IsTimelineEvent() {
}

func (MeetingEntity) TimelineEventLabel() string {
	return model.NodeLabelMeeting
}

func (e *MeetingEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *MeetingEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
