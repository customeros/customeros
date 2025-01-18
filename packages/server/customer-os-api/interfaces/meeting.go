package cosapi_interfaces

import (
	"context"

	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

type MeetingService interface {
	Update(ctx context.Context, input *MeetingUpdateData) (*neo4jentity.MeetingEntity, error)
	Create(ctx context.Context, newMeeting *MeetingCreateData) (*neo4jentity.MeetingEntity, error)

	LinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error
	UnlinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error

	GetMeetingById(ctx context.Context, meetingId string) (*neo4jentity.MeetingEntity, error)
	GetMeetingForInteractionEvent(ctx context.Context, interactionEventId string) (*neo4jentity.MeetingEntity, error)
	GetMeetingsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.MeetingEntities, error)
	GetParticipantsForMeetings(ctx context.Context, ids []string, relation entity.MeetingRelation) (*neo4jentity.MeetingParticipants, error)

	FindAll(ctx context.Context, externalSystemID string, externalID *string, page, limit int, filter *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
}

type MeetingParticipant struct {
	ContactId      *string
	UserId         *string
	OrganizationId *string
}

type MeetingCreateData struct {
	MeetingEntity     *neo4jentity.MeetingEntity
	CreatedBy         []MeetingParticipant
	AttendedBy        []MeetingParticipant
	NoteInput         *model.NoteInput
	ExternalReference *neo4jentity.ExternalSystemEntity
}

type MeetingUpdateData struct {
	MeetingEntity     *neo4jentity.MeetingEntity
	NoteEntity        *entity.NoteEntity
	Meeting           *string
	ExternalReference *neo4jentity.ExternalSystemEntity
}
