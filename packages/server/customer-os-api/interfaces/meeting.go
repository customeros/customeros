package cosapi_interfaces

import (
	"context"

	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type MeetingService interface {
	Update(ctx context.Context, input *MeetingUpdateData) (*entity.MeetingEntity, error)
	Create(ctx context.Context, newMeeting *MeetingCreateData) (*entity.MeetingEntity, error)

	LinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error
	UnlinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error

	GetMeetingById(ctx context.Context, meetingId string) (*entity.MeetingEntity, error)
	GetMeetingForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.MeetingEntity, error)
	GetMeetingsForInteractionEvents(ctx context.Context, ids []string) (*entity.MeetingEntities, error)
	GetParticipantsForMeetings(ctx context.Context, ids []string, relation entity.MeetingRelation) (*neo4jentity.MeetingParticipants, error)

	FindAll(ctx context.Context, externalSystemID string, externalID *string, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
}

type MeetingParticipant struct {
	ContactId      *string
	UserId         *string
	OrganizationId *string
}

type MeetingCreateData struct {
	MeetingEntity     *entity.MeetingEntity
	CreatedBy         []MeetingParticipant
	AttendedBy        []MeetingParticipant
	NoteInput         *model.NoteInput
	ExternalReference *neo4jentity.ExternalSystemEntity
}

type MeetingUpdateData struct {
	MeetingEntity     *entity.MeetingEntity
	NoteEntity        *entity.NoteEntity
	Meeting           *string
	ExternalReference *neo4jentity.ExternalSystemEntity
}
