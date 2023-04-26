package service

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
	"golang.org/x/net/context"
	"time"
)

type MeetingService interface {
	mapDbNodeToMeetingEntity(node dbtype.Node) *entity.MeetingEntity

	Update(ctx context.Context, input *MeetingUpdateData) (*entity.MeetingEntity, error)
	Create(ctx context.Context, newMeeting *MeetingCreateData) (*entity.MeetingEntity, error)

	LinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)
	UnlinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error)

	GetMeetingById(ctx context.Context, meetingId string) (*entity.MeetingEntity, error)
	GetParticipantsForMeetings(ctx context.Context, ids []string, relation entity.MeetingRelation) (*entity.MeetingParticipants, error)
}

type MeetingParticipantAddressData struct {
	ContactId *string
	UserId    *string
	Type      *string
}

type MeetingCreateData struct {
	MeetingEntity *entity.MeetingEntity
	CreatedBy     []MeetingParticipantAddressData
	AttendedBy    []MeetingParticipantAddressData
	NoteInput     *model.NoteInput
}

type MeetingUpdateData struct {
	MeetingEntity *entity.MeetingEntity
	NoteEntity    *entity.NoteEntity
}

type meetingService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewMeetingService(repositories *repository.Repositories, services *Services) MeetingService {
	return &meetingService{
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
	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
}

func (s *meetingService) UnlinkAttachment(ctx context.Context, meetingID string, attachmentID string) (*entity.MeetingEntity, error) {
	node, err := s.services.AttachmentService.UnlinkNodeWithAttachment(ctx, repository.INCLUDED_BY_MEETING, attachmentID, meetingID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToMeetingEntity(*node), nil
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
			if createdBy.ContactId != nil {
				err := s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenant, entity.CONTACT, meetingId, *createdBy.ContactId, createdBy.Type, entity.CREATED_BY)
				if err != nil {
					return nil, err
				}
			} else if createdBy.UserId != nil {
				err := s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenant, entity.USER, meetingId, *createdBy.UserId, createdBy.Type, entity.CREATED_BY)
				if err != nil {
					return nil, err
				}
			}

		}

		for _, attendedBy := range newMeeting.AttendedBy {
			if attendedBy.ContactId != nil {
				err := s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenant, entity.CONTACT, meetingId, *attendedBy.ContactId, attendedBy.Type, entity.ATTENDED_BY)
				if err != nil {
					return nil, err
				}
			} else if attendedBy.UserId != nil {
				err := s.repositories.MeetingRepository.LinkWithParticipantInTx(ctx, tx, tenant, entity.USER, meetingId, *attendedBy.UserId, attendedBy.Type, entity.ATTENDED_BY)
				if err != nil {
					return nil, err
				}
			}

		}
		if newMeeting.NoteInput != nil {
			toEntity := mapper.MapNoteInputToEntity(newMeeting.NoteInput)
			_, err := s.repositories.NoteRepository.CreateNoteForMeetingTx(ctx, tx, tenant, meetingId, toEntity)
			if err != nil {
				return nil, err
			}
		}
		return meetingDbNode, nil
	}
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
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrNil(props, "name"),
		Location:          utils.GetStringPropOrNil(props, "location"),
		Agenda:            utils.GetStringPropOrNil(props, "agenda"),
		AgendaContentType: utils.GetStringPropOrNil(props, "agendaContentType"),
		CreatedAt:         s.migrateStartedAt(props),
		UpdatedAt:         utils.GetTimePropOrNow(props, "updatedAt"),
		Start:             utils.GetTimePropOrNil(props, "start"),
		End:               utils.GetTimePropOrNil(props, "end"),
		Recording:         utils.GetStringPropOrNil(props, "recording"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &MeetingEntity
}

func (s *meetingService) mapDbRelationshipToParticipantDetails(relationship dbtype.Relationship) entity.MeetingParticipantDetails {
	props := utils.GetPropsFromRelationship(relationship)
	details := entity.MeetingParticipantDetails{
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return details
}

func MapMeetingParticipantInputToAddressData(input []*model.MeetingParticipantInput) []MeetingParticipantAddressData {
	var inputData []MeetingParticipantAddressData
	for _, participant := range input {
		inputData = append(inputData, MeetingParticipantAddressData{
			UserId:    participant.UserID,
			ContactId: participant.ContactID,
			Type:      participant.Type,
		})
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

func (s *meetingService) convertDbNodesToMeetingParticipants(records []*utils.DbNodeWithRelationAndId) entity.MeetingParticipants {
	meetingParticipants := entity.MeetingParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_User) {
			participant := s.services.UserService.mapDbNodeToUserEntity(*v.Node)
			participant.MeetingParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			meetingParticipants = append(meetingParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			participant := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			participant.MeetingParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			meetingParticipants = append(meetingParticipants, participant)
		}
	}
	return meetingParticipants
}

func (s *meetingService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
