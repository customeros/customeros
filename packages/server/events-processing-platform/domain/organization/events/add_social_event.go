package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationAddSocialEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	SocialId       string    `json:"socialId" validate:"required"`
	Url            string    `json:"url"`
	Alias          string    `json:"alias"`
	ExternalId     string    `json:"externalId"`
	FollowersCount int64     `json:"followersCount"`
	Source         string    `json:"source"`
	SourceOfTruth  string    `json:"sourceOfTruth"`
	AppSource      string    `json:"appSource"`
	CreatedAt      time.Time `json:"createdAt"`
}

func NewOrganizationAddSocialEvent(aggregate eventstore.Aggregate, socialId, url, alias, externalId string, followersCount int64, sourceFields cmnmod.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationAddSocialEvent{
		Tenant:         aggregate.GetTenant(),
		SocialId:       socialId,
		Url:            url,
		Alias:          alias,
		ExternalId:     externalId,
		FollowersCount: followersCount,
		Source:         sourceFields.Source,
		SourceOfTruth:  sourceFields.SourceOfTruth,
		AppSource:      sourceFields.AppSource,
		CreatedAt:      createdAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationAddSocialEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationAddSocialV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationAddSocialEvent")
	}
	return event, nil
}
