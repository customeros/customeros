package entity

import (
	"time"
)

type MeetingEntity struct {
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
	DataloaderKey      string
	Status             *MeetingStatus
}

func (MeetingEntity) IsTimelineEvent() {
}

func (MeetingEntity) TimelineEventLabel() string {
	return NodeLabel_Meeting
}
