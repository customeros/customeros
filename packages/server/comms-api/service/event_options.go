package service

import (
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

type EventOption func(*EventOptions)

type EventOptions struct {
	tenant           *string
	username         *string
	sessionId        *string
	meetingId        *string
	eventIdentifier  *string
	externalId       *string
	externalSystemId *string
	repliesTo        *string
	content          *string
	contentType      *string
	channel          *string
	channelData      *string
	eventType        *string
	sentBy           []cosModel.InteractionEventParticipantInput
	sentTo           []cosModel.InteractionEventParticipantInput
	appSource        *string
	createdAt        *time.Time
}

func WithTenant(value *string) EventOption {
	return func(options *EventOptions) {
		options.tenant = value
	}
}

func WithUsername(value *string) EventOption {
	return func(options *EventOptions) {
		options.username = value
	}
}

func WithSessionId(value *string) EventOption {
	return func(options *EventOptions) {
		options.sessionId = value
	}
}

func WithMeetingId(value *string) EventOption {
	return func(options *EventOptions) {
		options.meetingId = value
	}
}

func WithRepliesTo(value *string) EventOption {
	return func(options *EventOptions) {
		options.repliesTo = value
	}
}

func WithContent(value *string) EventOption {
	return func(options *EventOptions) {
		options.content = value
	}
}

func WithCreatedAt(value *time.Time) EventOption {
	return func(options *EventOptions) {
		options.createdAt = value
	}
}
func WithContentType(value *string) EventOption {
	return func(options *EventOptions) {
		options.contentType = value
	}
}

func WithEventType(value *string) EventOption {
	return func(options *EventOptions) {
		options.eventType = value
	}
}

func WithChannel(value *string) EventOption {
	return func(options *EventOptions) {
		options.channel = value
	}
}

func WithSentBy(value []cosModel.InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentBy = value
	}
}

func WithSentTo(value []cosModel.InteractionEventParticipantInput) EventOption {
	return func(options *EventOptions) {
		options.sentTo = value
	}
}

func WithAppSource(value *string) EventOption {
	return func(options *EventOptions) {
		options.appSource = value
	}
}

func WithEventIdentifier(eventIdentifier string) EventOption {
	return func(options *EventOptions) {
		options.eventIdentifier = &eventIdentifier
	}
}

func WithExternalId(externalId string) EventOption {
	return func(options *EventOptions) {
		options.externalId = &externalId
	}
}

func WithExternalSystemId(externalSystemId string) EventOption {
	return func(options *EventOptions) {
		options.externalSystemId = &externalSystemId
	}
}

func WithChannelData(ChannelData *string) EventOption {
	return func(options *EventOptions) {
		options.channelData = ChannelData
	}
}
