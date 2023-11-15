package entity

type TimelineEvent interface {
	IsTimelineEvent()
	TimelineEventLabel() string
}

type TimelineEventEntities []TimelineEvent
