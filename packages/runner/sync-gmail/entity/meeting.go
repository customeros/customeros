package entity

import (
	"fmt"
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
	Source             string
	SourceOfTruth      string
	Recording          *string
	Status             *MeetingStatus
}

func (meeting MeetingEntity) ToString() string {
	return fmt.Sprintf("id: %s\n", meeting.Id)
}
