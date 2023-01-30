package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type NoteService interface {
	GetNotesForContact(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error)
	MergeNoteToContact(ctx context.Context, contactId string, toEntity *entity.NoteEntity) (*entity.NoteEntity, error)
	UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	DeleteNote(ctx context.Context, noteId string) (bool, error)
}

type noteService struct {
	repositories *repository.Repositories
}

func NewNoteService(repositories *repository.Repositories) NoteService {
	return &noteService{
		repositories: repositories,
	}
}

func (s *noteService) getNeo4jDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *noteService) GetNotesForContact(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.repositories.Drivers.Neo4jDriver)
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	noteDbNodesWithTotalCount, err := s.repositories.NoteRepository.GetPaginatedNotesForContact(
		session,
		common.GetContext(ctx).Tenant,
		contactId,
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

func (s *noteService) MergeNoteToContact(ctx context.Context, contactId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	dbNodePtr, err := s.repositories.NoteRepository.CreateNoteForContact(session, common.GetContext(ctx).Tenant, contactId, *entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToNoteEntity(*dbNodePtr), nil
}

func (s *noteService) UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	dbNodePtr, err := s.repositories.NoteRepository.UpdateNote(session, common.GetTenantFromContext(ctx), *entity)
	if err != nil {
		return nil, err
	}

	var emailEntity = s.mapDbNodeToNoteEntity(*dbNodePtr)
	return emailEntity, nil
}

func (s *noteService) DeleteNote(ctx context.Context, noteId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getNeo4jDriver())
	defer session.Close()

	err := s.repositories.NoteRepository.Delete(session, common.GetTenantFromContext(ctx), noteId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *noteService) mapDbNodeToNoteEntity(node dbtype.Node) *entity.NoteEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.NoteEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Html:          utils.GetStringPropOrEmpty(props, "html"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrNow(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}
