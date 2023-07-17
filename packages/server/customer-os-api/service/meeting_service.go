package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
	"reflect"
	"time"
)

type MeetingService interface {
	mapDbNodeToMeetingEntity(node dbtype.Node) *entity.MeetingEntity

	Update(ctx context.Context, input *MeetingUpdateData) (*entity.MeetingEntity, error)
	Create(ctx context.Context, newMeeting *MeetingCreateData) (*entity.MeetingEntity, error)

	LinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error
	UnlinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error

	LinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)
	UnlinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)

	LinkRecordingAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)
	UnlinkRecordingAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)

	GetMeetingById(ctx context.Context, meetingId string) (*entity.MeetingEntity, error)
	GetMeetingForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.MeetingEntity, error)
	GetMeetingsForInteractionEvents(ctx context.Context, ids []string) (*entity.MeetingEntities, error)
	GetParticipantsForMeetings(ctx context.Context, ids []string, relation entity.MeetingRelation) (*entity.MeetingParticipants, error)

	FindAll(ctx context.Context, externalSystemID string, externalID *string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
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
	ExternalReference *entity.ExternalSystemEntity
}

type MeetingUpdateData struct {
	MeetingEntity *entity.MeetingEntity
	NoteEntity    *entity.NoteEntity
	Meeting       *string
}

type meetingService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewMeetingService(log logger.Logger, repositories *repository.Repositories, services *Services) MeetingService {
	return &meetingService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *meetingService) Create(ctx context.Context, newMeeting *MeetingCreateData) (*entity.MeetingEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, s.createMeetingInDBTxWork(ctx, newMeeting))
	if err != nil {
		return nil, err
	}

	for _, participant := range newMeeting.CreatedBy {
		if participant.ContactId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *participant.ContactId)
		}
		if participant.OrganizationId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *participant.OrganizationId)
		}
	}
	for _, participant := range newMeeting.AttendedBy {
		if participant.ContactId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *participant.ContactId)
		}
		if participant.OrganizationId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *participant.OrganizationId)
		}
	}

	return s.mapDbNodeToMeetingEntity(*queryResult.(*dbtype.Node)), nil
}

func (s *meetingService) Update(ctx context.Context, input *MeetingUpdateData) (*entity.MeetingEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		meetingDbNode, err := s.repositories.MeetingRepository.Update(ctx, tx, common.GetTenantFromContext(ctx), input.MeetingEntity)
		if err != nil {
			return nil, err
		}
		if input.NoteEntity != nil {
			_, err := s.services.NoteService.UpdateNote(ctx, input.NoteEntity)
			if err != nil {
				return nil, err
			}
		}
		return meetingDbNode, nil
	})

	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToMeetingEntity(*queryResult.(*dbtype.Node)), nil
}

func (s *meetingService) LinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error) {
	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, nil, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
}

func (s *meetingService) UnlinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error) {
	node, err := s.services.AttachmentService.UnlinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, nil, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
}

func (s *meetingService) LinkRecordingAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error) {
	recording := repository.INCLUDE_NATURE_RECORDING
	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, &recording, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
}

func (s *meetingService) UnlinkRecordingAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error) {
	recording := repository.INCLUDE_NATURE_RECORDING
	node, err := s.services.AttachmentService.UnlinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, &recording, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
}

func (s *meetingService) LinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		tenant := common.GetContext(ctx).Tenant
		err := s.linkAttendedByTxWork(ctx, tx, tenant, meetingID, participant, entity.ATTENDED_BY)
		return nil, err
	})

	if participant.ContactId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *participant.ContactId)
	}
	if participant.OrganizationId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *participant.OrganizationId)
	}

	return err
}

func (s *meetingService) UnlinkAttendedBy(ctx context.Context, meetingID string, participant MeetingParticipant) error {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		tenant := common.GetContext(ctx).Tenant
		err := s.unlinkAttendedByTxWork(ctx, tx, tenant, meetingID, participant, entity.ATTENDED_BY)
		return nil, err
	})

	if participant.ContactId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *participant.ContactId)
	}
	if participant.OrganizationId != nil {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, *participant.OrganizationId)
	}

	return err
}

func (s *meetingService) GetMeetingById(ctx context.Context, id string) (*entity.MeetingEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (m:Meeting_%s {id:$id}) RETURN m`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToMeetingEntity(queryResult.(dbtype.Node)), nil
}

func (s *meetingService) createMeetingInDBTxWork(ctx context.Context, newMeeting *MeetingCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		meetingDbNode, err := s.repositories.MeetingRepository.Create(ctx, tx, common.GetTenantFromContext(ctx), newMeeting.MeetingEntity)
		if err != nil {
			return nil, err
		}
		var meetingId = utils.GetPropsFromNode(*meetingDbNode)["id"].(string)

		for _, createdBy := range newMeeting.CreatedBy {
			err := s.linkAttendedByTxWork(ctx, tx, tenant, meetingId, createdBy, entity.CREATED_BY)
			if err != nil {
				return nil, err
			}
		}

		for _, attendedBy := range newMeeting.AttendedBy {
			err := s.linkAttendedByTxWork(ctx, tx, tenant, meetingId, attendedBy, entity.ATTENDED_BY)
			if err != nil {
				return nil, err
			}
		}

		if newMeeting.NoteInput != nil {
			toEntity := mapper.MapNoteInputToEntity(newMeeting.NoteInput)
			_, err := s.repositories.NoteRepository.CreateNoteForMeetingTx(ctx, tx, tenant, meetingId, toEntity)
			if err != nil {
				return nil, err
			}
		}

		if newMeeting.ExternalReference != nil {
			err := s.repositories.ExternalSystemRepository.LinkNodeWithExternalSystemInTx(ctx, tx, tenant, meetingId, entity.ExternalNodeMeeting, *newMeeting.ExternalReference)
			if err != nil {
				return nil, err
			}
		}

		return meetingDbNode, nil
	}
}

func (s *meetingService) linkAttendedByTxWork(ctx context.Context, tx neo4j.ManagedTransaction, tenantName, meetingId string, participant MeetingParticipant, relationType entity.MeetingRelation) error {
	var err error
	if participant.ContactId != nil {
		err = s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenantName, meetingId, *participant.ContactId, entity.CONTACT, relationType)
	} else if participant.UserId != nil {
		err = s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenantName, meetingId, *participant.UserId, entity.USER, relationType)
	} else if participant.OrganizationId != nil {
		err = s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenantName, meetingId, *participant.OrganizationId, entity.ORGANIZATION, relationType)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *meetingService) unlinkAttendedByTxWork(ctx context.Context, tx neo4j.ManagedTransaction, tenantName, meetingId string, participant MeetingParticipant, relationType entity.MeetingRelation) error {
	var err error
	if participant.ContactId != nil {
		err = s.repositories.MeetingRepository.UnlinkParticipantInTx(ctx, tx, tenantName, meetingId, *participant.ContactId, entity.CONTACT, relationType)
	} else if participant.UserId != nil {
		err = s.repositories.MeetingRepository.UnlinkParticipantInTx(ctx, tx, tenantName, meetingId, *participant.UserId, entity.USER, relationType)
	} else if participant.OrganizationId != nil {
		err = s.repositories.MeetingRepository.UnlinkParticipantInTx(ctx, tx, tenantName, meetingId, *participant.OrganizationId, entity.ORGANIZATION, relationType)
	}
	if err != nil {
		return err
	}
	return nil
}

// createdAt takes priority over startedAt
func (s *meetingService) migrateStartedAt(props map[string]any) time.Time {
	if props["createdAt"] != nil {
		return utils.GetTimePropOrNow(props, "createdAt")
	}
	if props["startedAt"] != nil {
		return utils.GetTimePropOrNow(props, "startedAt")
	}
	return time.Now()
}

func (s *meetingService) mapDbNodeToMeetingEntity(node dbtype.Node) *entity.MeetingEntity {
	props := utils.GetPropsFromNode(node)
	MeetingEntity := entity.MeetingEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		Name:               utils.GetStringPropOrNil(props, "name"),
		ConferenceUrl:      utils.GetStringPropOrNil(props, "conferenceUrl"),
		MeetingExternalUrl: utils.GetStringPropOrNil(props, "meetingExternalUrl"),
		Agenda:             utils.GetStringPropOrNil(props, "agenda"),
		AgendaContentType:  utils.GetStringPropOrNil(props, "agendaContentType"),
		CreatedAt:          s.migrateStartedAt(props),
		UpdatedAt:          utils.GetTimePropOrNow(props, "updatedAt"),
		StartedAt:          utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:            utils.GetTimePropOrNil(props, "endedAt"),
		Recording:          utils.GetStringPropOrNil(props, "recording"),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		Source:             entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &MeetingEntity
}

func MapMeetingParticipantInputToParticipant(participant *model.MeetingParticipantInput) MeetingParticipant {
	meetingParticipant := MeetingParticipant{
		UserId:         participant.UserID,
		ContactId:      participant.ContactID,
		OrganizationId: participant.OrganizationID,
	}
	return meetingParticipant
}

func MapMeetingParticipantInputListToParticipant(input []*model.MeetingParticipantInput) []MeetingParticipant {
	var inputData []MeetingParticipant
	for _, participant := range input {
		inputData = append(inputData, MapMeetingParticipantInputToParticipant(participant))
	}
	return inputData
}

func (s *meetingService) GetParticipantsForMeetings(ctx context.Context, ids []string, relation entity.MeetingRelation) (*entity.MeetingParticipants, error) {
	records, err := s.repositories.MeetingRepository.GetParticipantsForMeetings(ctx, common.GetTenantFromContext(ctx), ids, relation)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToMeetingParticipants(records)

	return &interactionEventParticipants, nil
}

func (s *meetingService) GetMeetingForInteractionEvent(ctx context.Context, interactionEventId string) (*entity.MeetingEntity, error) {
	record, err := s.repositories.MeetingRepository.GetMeetingForInteractionEvent(ctx, common.GetTenantFromContext(ctx), interactionEventId)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	return s.mapDbNodeToMeetingEntity(*record), nil
}

func (s *meetingService) GetMeetingsForInteractionEvents(ctx context.Context, ids []string) (*entity.MeetingEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingService.GetMeetingsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	issues, err := s.repositories.MeetingRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	meetingEntities := make(entity.MeetingEntities, 0, len(issues))
	for _, v := range issues {
		meetingEntity := s.mapDbNodeToMeetingEntity(*v.Node)
		meetingEntity.DataloaderKey = v.LinkedNodeId
		meetingEntities = append(meetingEntities, *meetingEntity)
	}
	return &meetingEntities, nil
}

func (s *meetingService) FindAll(ctx context.Context, externalSystemID string, externalID *string, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.MeetingEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.MeetingEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.MeetingRepository.GetPaginatedMeetings(
		ctx, session,
		externalSystemID,
		externalID,
		common.GetContext(ctx).Tenant,
		common.GetContext(ctx).UserEmail,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	meetings := make(entity.MeetingEntities, 0, len(dbNodesWithTotalCount.Nodes))

	for _, v := range dbNodesWithTotalCount.Nodes {
		meetings = append(meetings, *s.mapDbNodeToMeetingEntity(*v))
	}
	paginatedResult.SetRows(&meetings)
	return &paginatedResult, nil
}

func (s *meetingService) convertDbNodesToMeetingParticipants(records []*utils.DbNodeWithRelationAndId) entity.MeetingParticipants {
	meetingParticipants := entity.MeetingParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_User) {
			participant := s.services.UserService.mapDbNodeToUserEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			meetingParticipants = append(meetingParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			participant := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			meetingParticipants = append(meetingParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Organization) {
			participant := s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			meetingParticipants = append(meetingParticipants, participant)
		}
	}
	return meetingParticipants
}

func (s *meetingService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
