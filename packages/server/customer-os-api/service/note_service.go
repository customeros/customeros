package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type NoteService interface {
	GetNotesForContact(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error)
	MergeNoteToContact(ctx context.Context, contactId string, toEntity *entity.NoteEntity) (*entity.NoteEntity, error)
	UpdateNoteInContact(ctx context.Context, contactId string, toEntity *entity.NoteEntity) (*entity.NoteEntity, error)
	DeleteFromContact(ctx context.Context, contactId string, noteId string) (bool, error)
	getDriver() neo4j.Driver
}

type noteService struct {
	repository *repository.Repositories
}

func NewNoteService(repository *repository.Repositories) NoteService {
	return &noteService{
		repository: repository,
	}
}

func (s *noteService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *noteService) GetNotesForContact(ctx context.Context, contactId string, page, limit int) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.repository.Drivers.Neo4jDriver)
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	noteDbNodesWithTotalCount, err := s.repository.NoteRepository.GetPaginatedNotesForContact(
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
		noteEntity := *s.mapDbNodeTNoteEntity(v.Node)
		entities = append(entities, noteEntity)
	}
	paginatedResult.SetRows(&entities)
	return &paginatedResult, nil
}

func (s *noteService) MergeNoteToContact(ctx context.Context, contactId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MERGE (c)-[:NOTED]->(n:Note {id:randomUUID(), html: $html})
			RETURN n`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"html":      entity.Html,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Single()
	})
	if err != nil {
		return nil, err
	}

	node := queryResult.(*db.Record).Values[0].(dbtype.Node)
	return s.mapDbNodeTNoteEntity(&node), nil
}

func (s *noteService) UpdateNoteInContact(ctx context.Context, contactId string, entity *entity.NoteEntity) (*entity.NoteEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(c)-[r:NOTED]->(n:Note {id:$noteId})
			SET n.html=$html
			RETURN n`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"noteId":    entity.Id,
				"html":      entity.Html,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return nil, err
	}

	node := queryResult.(*db.Record).Values[0].(dbtype.Node)
	var emailEntity = s.mapDbNodeTNoteEntity(&node)
	return emailEntity, nil
}

func (s *noteService) DeleteFromContact(ctx context.Context, contactId string, noteId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact {id:$contactId})-[:NOTED]->(n:Note {id:$noteId})
            DETACH DELETE n
			`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"noteId":    noteId,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *noteService) mapDbNodeTNoteEntity(node *dbtype.Node) *entity.NoteEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.NoteEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Html: utils.GetStringPropOrEmpty(props, "html"),
	}
	return &result
}
