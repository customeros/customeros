package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"time"
)

type InteractionEventBuilderOption func(*InteractionEventBuilderOptions)

type InteractionEventBuilderOptions struct {
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
	sentBy           []model.InteractionEventParticipantInput
	sentTo           []model.InteractionEventParticipantInput
	appSource        *string
	createdAt        *time.Time
}

func WithSessionId(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.sessionId = value
	}
}

func WithMeetingId(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.meetingId = value
	}
}

func WithRepliesTo(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.repliesTo = value
	}
}

func WithContent(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.content = value
	}
}

func WithCreatedAt(value *time.Time) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.createdAt = value
	}
}
func WithContentType(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.contentType = value
	}
}

func WithEventType(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.eventType = value
	}
}

func WithChannel(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.channel = value
	}
}

func WithSentBy(value []model.InteractionEventParticipantInput) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.sentBy = value
	}
}

func WithSentTo(value []model.InteractionEventParticipantInput) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.sentTo = value
	}
}

func WithAppSource(value *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.appSource = value
	}
}

func WithEventIdentifier(eventIdentifier string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.eventIdentifier = &eventIdentifier
	}
}

func WithExternalId(externalId string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.externalId = &externalId
	}
}

func WithExternalSystemId(externalSystemId string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.externalSystemId = &externalSystemId
	}
}

func WithChannelData(ChannelData *string) InteractionEventBuilderOption {
	return func(options *InteractionEventBuilderOptions) {
		options.channelData = ChannelData
	}
}
