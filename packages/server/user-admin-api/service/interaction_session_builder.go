package service

import "github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"

type InteractionSessionBuilderOption func(*InteractionSessionBuilderOptions)

type InteractionSessionBuilderOptions struct {
	channel           *string
	name              *string
	status            *string
	appSource         *string
	sessionIdentifier *string
	sessionType       *string
	attendedBy        []model.InteractionSessionParticipantInput
}

func WithSessionIdentifier(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.sessionIdentifier = value
	}
}

func WithSessionChannel(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.channel = value
	}
}

func WithSessionName(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.name = value
	}
}

func WithSessionStatus(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.status = value
	}
}

func WithSessionAttendedBy(value []model.InteractionSessionParticipantInput) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.attendedBy = value
	}
}

func WithSessionAppSource(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.appSource = value
	}
}

func WithSessionType(value *string) InteractionSessionBuilderOption {
	return func(options *InteractionSessionBuilderOptions) {
		options.sessionType = value
	}
}
