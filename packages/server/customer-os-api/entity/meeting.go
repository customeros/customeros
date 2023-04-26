package entity

import (
	"fmt"
	"time"
)

type MeetingEntity struct {
	Id                string
	Name              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Start             *time.Time
	End               *time.Time
	ConferenceUrl     *string
	AppSource         string
	Agenda            *string
	AgendaContentType *string
	Source            DataSource
	SourceOfTruth     DataSource
	Recording         *string
	DataloaderKey     string
}

func (meeting MeetingEntity) ToString() string {
	return fmt.Sprintf("id: %s\n", meeting.Id)
}

type MeetingEntities []MeetingEntity

func (MeetingEntity) IsTimelineEvent() {
}

func (MeetingEntity) TimelineEventLabel() string {
	return NodeLabel_Meeting
}

func (MeetingEntity) IsAnalysisDescribe() {
}

func (MeetingEntity) AnalysisDescribeLabel() string {
	return NodeLabel_Meeting
}

func (meeting MeetingEntity) GetDataloaderKey() string {
	return meeting.DataloaderKey
}

func (meeting MeetingEntity) Labels(tenant string) []string {
	return []string{"Meeting", "TimelineEvent", "Meeting_" + tenant, "TimelineEvent_" + tenant}
}
