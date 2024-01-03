package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
	"time"
)

type NoteService interface {
	GetNotesForContactPaginated(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error)
	GetNotesForContactTimeRange(ctx context.Context, contactId string, start, end time.Time) (*entity.NoteEntities, error)
	GetNotesForOrganization(ctx context.Context, organizationId string, page, limit int) (*utils.Pagination, error)
	GetNotesForMeetings(ctx context.Context, ids []string) (*entity.NoteEntities, error)

	CreateNoteForContact(ctx context.Context, contactId string, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	CreateNoteForOrganization(ctx context.Context, organizationId string, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	CreateNoteForMeeting(ctx context.Context, meetingId string, entity *entity.NoteEntity) (*entity.NoteEntity, error)

	UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	DeleteNote(ctx context.Context, noteId string) (bool, error)

	GetNotedEntities(ctx context.Context, ids []string) (*entity.NotedEntities, error)
	NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.NoteEntity, error)
	NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.NoteEntity, error)

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

func (s *noteService) NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.NoteLinkAttachment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("noteID", noteID), log.String("attachmentID", attachmentID))

	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.LINKED_WITH_NOTE, nil, attachmentID, noteID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToNoteEntity(*node), nil
}

func (s *noteService) NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.NoteUnlinkAttachment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("noteID", noteID), log.String("attachmentID", attachmentID))

	node, err := s.services.AttachmentService.UnlinkNodeWithAttachment(ctx, repository.LINKED_WITH_NOTE, nil, attachmentID, noteID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToNoteEntity(*node), nil
}

func (s *noteService) GetNotesForContactPaginated(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.GetNotesForContactPaginated")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.Int("page", page), log.Int("limit", limit))

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	noteDbNodesWithTotalCount, err := s.repositories.NoteRepository.GetPaginatedNotesForContact(
		ctx,
		common.GetContext(ctx).Tenant,
		contactId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(noteDbNodesWithTotalCount.Count)

	entities := make(entity.NoteEntities, 0, len(noteDbNodesWithTotalCount.Nodes))
	for _, v := range noteDbNodesWithTotalCount.Nodes {
		noteEntity := *s.mapDbNodeToNoteEntity(*v.Node)
		entities = append(entities, noteEntity)
	}
	paginatedResult.SetRows(&entities)
	return &paginatedResult, nil
}

func (s *noteService) GetNotesForContactTimeRange(ctx context.Context, contactId string, start time.Time, end time.Time) (*entity.NoteEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.GetNotesForContactTimeRange")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.Object("start", start.String()), log.Object("end", end.String()))

	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	nodes, err := s.repositories.NoteRepository.GetTimeRangeNotesForContact(
		ctx,
		session,
		common.GetContext(ctx).Tenant,
		contactId,
		start,
		end)
	if err != nil {
		return nil, err
	}
	result := make(entity.NoteEntities, len(nodes))

	for i, v := range nodes {
		noteEntity := s.mapDbNodeToNoteEntity(*v)
		result[i] = *noteEntity
	}
	return &result, nil
}

func (s *noteService) GetNotesForOrganization(ctx context.Context, organizationId string, page, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.GetNotesForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.Int("page", page), log.Int("limit", limit))

	session := utils.NewNeo4jReadSession(ctx, *s.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	noteDbNodesWithTotalCount, err := s.repositories.NoteRepository.GetPaginatedNotesForOrganization(
		ctx,
		session,
		common.GetContext(ctx).Tenant,
		organizationId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(noteDbNodesWithTotalCount.Count)

	entities := entity.NoteEntities{}

	for _, v := range noteDbNodesWithTotalCount.Nodes {
		noteEntity := *s.mapDbNodeToNoteEntity(*v.Node)
		entities = append(entities, noteEntity)
	}
	paginatedResult.SetRows(&entities)
	return &paginatedResult, nil
}

func (s *noteService) CreateNoteForContact(ctx context.Context, contactId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.CreateNoteForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId))

	dbNodePtr, err := s.repositories.NoteRepository.CreateNoteForContact(ctx, common.GetContext(ctx).Tenant, contactId, *entity)
	if err != nil {
		return nil, err
	}
	// set note creator
	if len(common.GetUserIdFromContext(ctx)) > 0 {
		props := utils.GetPropsFromNode(*dbNodePtr)
		noteId := utils.GetStringPropOrEmpty(props, "id")
		_ = s.repositories.NoteRepository.SetNoteCreator(ctx, common.GetTenantFromContext(ctx), common.GetUserIdFromContext(ctx), noteId)
	}
	return s.mapDbNodeToNoteEntity(*dbNodePtr), nil
}

func (s *noteService) CreateNoteForOrganization(ctx context.Context, organization string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.CreateNoteForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organization", organization))

	dbNodePtr, err := s.repositories.NoteRepository.CreateNoteForOrganization(ctx, common.GetContext(ctx).Tenant, organization, *entity)
	if err != nil {
		return nil, err
	}
	// set note creator
	if len(common.GetUserIdFromContext(ctx)) > 0 {
		props := utils.GetPropsFromNode(*dbNodePtr)
		noteId := utils.GetStringPropOrEmpty(props, "id")
		_ = s.repositories.NoteRepository.SetNoteCreator(ctx, common.GetTenantFromContext(ctx), common.GetUserIdFromContext(ctx), noteId)
	}
	return s.mapDbNodeToNoteEntity(*dbNodePtr), nil
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

func (s *noteService) GetNotedEntities(ctx context.Context, ids []string) (*entity.NotedEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.GetNotedEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.repositories.NoteRepository.GetNotedEntitiesForNotes(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	notedEntities := entity.NotedEntities{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_Organization) {
			notedEntity := s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v.Node)
			notedEntity.DataloaderKey = v.LinkedNodeId
			notedEntities = append(notedEntities, notedEntity)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			notedEntity := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			notedEntity.DataloaderKey = v.LinkedNodeId
			notedEntities = append(notedEntities, notedEntity)
		}
	}

	return &notedEntities, nil
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
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
