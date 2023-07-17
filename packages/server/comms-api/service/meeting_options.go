package service

import (
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

type MeetingOption func(*MeetingOptions)

type MeetingOptions struct {
	name           *string
	appSource      *string
	tenant         *string
	username       *string
	startedAt      *time.Time
	endedAt        *time.Time
	attendedBy     []cosModel.MeetingParticipantInput
	createdBy      []cosModel.MeetingParticipantInput
	noteInput      *cosModel.NoteInput
	externalSystem *cosModel.ExternalSystemReferenceInput
}

func WithMeetingName(value *string) MeetingOption {
	return func(options *MeetingOptions) {
		options.name = value
	}
}

func WithMeetingAppSource(value *string) MeetingOption {
	return func(options *MeetingOptions) {
		options.appSource = value
	}
}

func WithMeetingTenant(value *string) MeetingOption {
	return func(options *MeetingOptions) {
		options.tenant = value
	}
}

func WithMeetingUsername(value *string) MeetingOption {
	return func(options *MeetingOptions) {
		options.username = value
	}
}

func WithMeetingStartedAt(value *time.Time) MeetingOption {
	return func(options *MeetingOptions) {
		options.startedAt = value
	}
}

func WithMeetingEndedAt(value *time.Time) MeetingOption {
	return func(options *MeetingOptions) {
		options.endedAt = value
	}
}

func WithMeetingAttendedBy(value []cosModel.MeetingParticipantInput) MeetingOption {
	return func(options *MeetingOptions) {
		options.attendedBy = value
	}
}

func WithMeetingCreatedBy(value []cosModel.MeetingParticipantInput) MeetingOption {
	return func(options *MeetingOptions) {
		options.createdBy = value
	}
}

func WithMeetingNote(value *cosModel.NoteInput) MeetingOption {
	return func(options *MeetingOptions) {
		options.noteInput = value
	}
}

func WithExternalSystem(value *cosModel.ExternalSystemReferenceInput) MeetingOption {
	return func(options *MeetingOptions) {
		options.externalSystem = value
	}
}
