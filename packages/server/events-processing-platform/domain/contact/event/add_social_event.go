package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactAddSocialEvent struct {
	Tenant         string             `json:"tenant" validate:"required"`
	SocialId       string             `json:"socialId" validate:"required"`
	Url            string             `json:"url" `
	Alias          string             `json:"alias"`
	FollowersCount int64              `json:"followersCount"`
	Source         commonmodel.Source `json:"source"`
	CreatedAt      time.Time          `json:"createdAt"`
}

func NewContactAddSocialEvent(aggregate eventstore.Aggregate, socialId, url, alias string, followersCount int64, sourceFields commonmodel.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := ContactAddSocialEvent{
		Tenant:         aggregate.GetTenant(),
		SocialId:       socialId,
		Url:            url,
		Alias:          alias,
		FollowersCount: followersCount,
		Source:         sourceFields,
		CreatedAt:      createdAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactAddSocialEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactAddSocialV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactAddSocialEvent")
	}
	return event, nil
}
