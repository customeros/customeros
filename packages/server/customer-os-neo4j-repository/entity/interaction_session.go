package neo4j_entity

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"time"
)

type InteractionSessionEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Identifier    string
	Name          string
	Status        commonenum.InteractionSessionStatus
	Type          commonenum.InteractionSessionType
	Channel       commonenum.InteractionSessionChannel
	ChannelData   string
	AppSource     string
	Source        DataSource
	SourceOfTruth DataSource
}

type InteractionSessionEntities []InteractionSessionEntity

func (InteractionSessionEntity) IsTimelineEvent() {
}

func (InteractionSessionEntity) TimelineEventLabel() string {
	return model.NodeLabelInteractionSession
}

func (e *InteractionSessionEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *InteractionSessionEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
