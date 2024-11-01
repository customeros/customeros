package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type NoteService interface {
	GetById(ctx context.Context, id string) (*entity.NoteEntity, error)
	GetNotesForMeetings(ctx context.Context, ids []string) (*entity.NoteEntities, error)

	CreateNoteForMeeting(ctx context.Context, meetingId string, entity *entity.NoteEntity) (*entity.NoteEntity, error)

	UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	DeleteNote(ctx context.Context, noteId string) (bool, error)

	NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) error
	NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) error

	mapDbNodeToNoteEntity(node dbtype.Node) *entity.NoteEntity
}

type noteService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewNoteService(log logger.Logger, repositories *repository.Repositories, services *Services) NoteService {
	return &noteService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *noteService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *noteService) GetById(ctx context.Context, id string) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.GetById")
	defer span.Finish()

	byId, err := s.services.CommonServices.Neo4jRepositories.CommonReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id, commonModel.NodeLabelNote)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if byId == nil {
		return nil, nil
	}

	return s.mapDbNodeToNoteEntity(*byId), nil
}

func (s *noteService) NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.NoteLinkAttachment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("noteID", noteID), log.String("attachmentID", attachmentID))

	tenant := common.GetTenantFromContext(ctx)

	err := s.services.CommonServices.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, tenant, neo4jrepository.LinkDetails{
		FromEntityId:   noteID,
		FromEntityType: commonModel.NOTE,
		Relationship:   commonModel.INCLUDES,
		ToEntityId:     attachmentID,
		ToEntityType:   commonModel.ATTACHMENT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *noteService) NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.NoteUnlinkAttachment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("noteID", noteID), log.String("attachmentID", attachmentID))

	tenant := common.GetTenantFromContext(ctx)

	err := s.services.CommonServices.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, neo4jrepository.LinkDetails{
		FromEntityId:   noteID,
		FromEntityType: commonModel.NOTE,
		Relationship:   commonModel.INCLUDES,
		ToEntityId:     attachmentID,
		ToEntityType:   commonModel.ATTACHMENT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *noteService) UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.UpdateNote")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	dbNodePtr, err := s.repositories.NoteRepository.UpdateNote(ctx, session, common.GetTenantFromContext(ctx), *entity)
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToNoteEntity(*dbNodePtr)
	return emailEntity, nil
}

func (s *noteService) DeleteNote(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.DeleteNote")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("noteId", noteId))

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

func (s *noteService) CreateNoteForMeeting(ctx context.Context, meetingId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.CreateNoteForMeeting")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("meetingId", meetingId))

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
