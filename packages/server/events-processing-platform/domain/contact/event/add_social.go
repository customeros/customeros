package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/pkg/errors"
	"time"
)

type AddSocialEvent struct {
	Tenant    string             `json:"tenant" validate:"required"`
	SocialId  string             `json:"socialId" validate:"required"`
	Url       string             `json:"url" validate:"required"`
	Source    commonmodel.Source `json:"source"`
	CreatedAt time.Time          `json:"createdAt"`
}

func NewAddSocialEvent(aggregate eventstore.Aggregate, request *contactpb.ContactAddSocialGrpcRequest, socialId string, source commonmodel.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := AddSocialEvent{
		Tenant:    aggregate.GetTenant(),
		CreatedAt: createdAt,
		Source:    source,
		SocialId:  socialId,
		Url:       request.Url,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate AddSocialEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactAddSocialV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for AddSocialEvent")
	}
	return event, nil
}
