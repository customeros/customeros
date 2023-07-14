package service

import cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

type SessionOption func(*SessionOptions)

type SessionOptions struct {
	channel           *string
	name              *string
	status            *string
	appSource         *string
	tenant            *string
	username          *string
	sessionIdentifier *string
	sessionType       *string
	attendedBy        []cosModel.InteractionSessionParticipantInput
}

func WithSessionIdentifier(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.sessionIdentifier = value
	}
}

func WithSessionChannel(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.channel = value
	}
}

func WithSessionName(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.name = value
	}
}

func WithSessionStatus(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.status = value
	}
}

func WithSessionAttendedBy(value []cosModel.InteractionSessionParticipantInput) SessionOption {
	return func(options *SessionOptions) {
		options.attendedBy = value
	}
}

func WithSessionAppSource(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.appSource = value
	}
}

func WithSessionTenant(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.tenant = value
	}
}

func WithSessionUsername(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.username = value
	}
}

func WithSessionType(value *string) SessionOption {
	return func(options *SessionOptions) {
		options.sessionType = value
	}
}
