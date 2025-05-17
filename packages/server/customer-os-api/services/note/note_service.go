package api_note

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type noteService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewNoteService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.NoteService {
	return &noteService{
		log:          log,
		repositories: repositories,
	}
}

func (s *noteService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *noteService) GetById(ctx context.Context, id string) (*entity.NoteEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.GetById")
	defer spans.Finish()

	byId, err := s.repositories.Neo4jRepositories.CommonReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id, commonModel.NodeLabelNote)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if byId == nil {
		return nil, nil
	}

	return s.mapDbNodeToNoteEntity(*byId), nil
}

func (s *noteService) NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.NoteLinkAttachment")
	defer spans.Finish()
	spans.LogKV("noteID", noteID, "attachmentID", attachmentID)

	tenant := common.GetTenantFromContext(ctx)

	err := s.repositories.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, tenant, neo4jrepository.LinkDetails{
		FromEntityId:   noteID,
		FromEntityType: commonModel.NOTE,
		Relationship:   commonModel.INCLUDES,
		ToEntityId:     attachmentID,
		ToEntityType:   commonModel.ATTACHMENT,
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (s *noteService) NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.NoteUnlinkAttachment")
	defer spans.Finish()
	spans.LogKV("noteID", noteID, "attachmentID", attachmentID)

	tenant := common.GetTenantFromContext(ctx)

	err := s.repositories.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, neo4jrepository.LinkDetails{
		FromEntityId:   noteID,
		FromEntityType: commonModel.NOTE,
		Relationship:   commonModel.INCLUDES,
		ToEntityId:     attachmentID,
		ToEntityType:   commonModel.ATTACHMENT,
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (s *noteService) mapDbNodeToNoteEntity(node dbtype.Node) *entity.NoteEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.NoteEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}

func (s *noteService) convertDbNodesToNotes(records []*utils.DbNodeAndId) entity.NoteEntities {
	notes := entity.NoteEntities{}
	for _, v := range records {
		attachment := s.mapDbNodeToNoteEntity(*v.Node)
		attachment.DataloaderKey = v.LinkedNodeId
		notes = append(notes, *attachment)

	}
	return notes
}

func (s *noteService) UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.UpdateNote")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNodePtr, err := s.repositories.NoteRepository.UpdateNote(ctx, session, common.GetTenantFromContext(ctx), *entity)
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToNoteEntity(*dbNodePtr)
	return emailEntity, nil
}
func (s *noteService) CreateNoteForMeeting(ctx context.Context, meetingId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.CreateNoteForMeeting")
	defer spans.Finish()
	spans.LogKV("meetingId", meetingId)

	dbNodePtr, err := s.repositories.NoteRepository.CreateNoteForMeeting(ctx, common.GetContext(ctx).Tenant, meetingId, entity)
	if err != nil {
		return nil, err
	}
	// set note creator
	if len(common.GetUserIdFromContext(ctx)) > 0 {
		props := utils.GetPropsFromNode(*dbNodePtr)
		noteId := utils.GetStringPropOrEmpty(props, "id")
		s.repositories.NoteRepository.SetNoteCreator(ctx, common.GetTenantFromContext(ctx), common.GetUserIdFromContext(ctx), noteId)
	}
	return s.mapDbNodeToNoteEntity(*dbNodePtr), nil
}

func (s *noteService) DeleteNote(ctx context.Context, noteId string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NoteService.DeleteNote")
	defer spans.Finish()
	spans.LogKV("noteId", noteId)

	err := s.repositories.NoteRepository.Delete(ctx, common.GetTenantFromContext(ctx), noteId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *noteService) GetNotesForMeetings(ctx context.Context, ids []string) (*entity.NoteEntities, error) {

	records, err := s.repositories.NoteRepository.GetNotesForMeetings(ctx, common.GetContext(ctx).Tenant, ids)
	if err != nil {
		return nil, err
	}

	notes := s.convertDbNodesToNotes(records)

	return &notes, nil
}
