package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"reflect"
)

type ContactGroupService interface {
	Create(ctx context.Context, entity *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error)
	Update(ctx context.Context, contactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error)
	Delete(ctx context.Context, id string) (bool, error)
	FindContactGroupById(ctx context.Context, id string) (*entity.ContactGroupEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupEntities, error)
	AddContactToGroup(ctx context.Context, contactId, groupId string) (bool, error)
	RemoveContactFromGroup(ctx context.Context, contactId, groupId string) (bool, error)
}

type contactGroupService struct {
	repositories *repository.Repositories
}

func NewContactGroupService(repositories *repository.Repositories) ContactGroupService {
	return &contactGroupService{
		repositories: repositories,
	}
}

func (s *contactGroupService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *contactGroupService) Create(ctx context.Context, entity *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	dbNode, err := s.repositories.ContactGroupRepository.Create(session, common.GetContext(ctx).Tenant, *entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactGroupEntity(*dbNode), nil
}

func (s *contactGroupService) Update(ctx context.Context, entity *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	dbNode, err := s.repositories.ContactGroupRepository.Update(session, common.GetContext(ctx).Tenant, *entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToContactGroupEntity(*dbNode), nil
}

func (s *contactGroupService) Delete(ctx context.Context, id string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
            DETACH DELETE g
			`,
			map[string]any{
				"groupId": id,
				"tenant":  common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactGroupService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ContactGroupEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.ContactGroupEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.ContactGroupRepository.GetPaginatedContactGroups(
		session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	contactGroups := entity.ContactGroupEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		contactGroups = append(contactGroups, *s.mapDbNodeToContactGroupEntity(*v))
	}
	paginatedResult.SetRows(&contactGroups)
	return &paginatedResult, nil
}

func (s *contactGroupService) FindContactGroupById(ctx context.Context, id string) (*entity.ContactGroupEntity, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
			MATCH (c:ContactGroup {id:$id})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) RETURN c`,
			map[string]any{
				"id":     id,
				"tenant": common.GetContext(ctx).Tenant,
			})
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToContactGroupEntity(queryResult.(dbtype.Node)), nil
}

func (s *contactGroupService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$id})-[:BELONGS_TO_GROUP]->(g:ContactGroup)-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
				RETURN g 
				ORDER BY g.name`,
			map[string]any{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	contactGroups := entity.ContactGroupEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := s.mapDbNodeToContactGroupEntity(dbRecord.Values[0].(dbtype.Node))
		contactGroups = append(contactGroups, *contactGroup)
	}

	return &contactGroups, nil
}

func (s *contactGroupService) AddContactToGroup(ctx context.Context, contactId, groupId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH 	(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
					(g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (c)-[:BELONGS_TO_GROUP]->(g)
			`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"groupId":   groupId,
			})
		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactGroupService) RemoveContactFromGroup(ctx context.Context, contactId, groupId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH 	(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
					(g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MATCH (c)-[r:BELONGS_TO_GROUP]->(g)
            DELETE r`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"groupId":   groupId,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactGroupService) mapDbNodeToContactGroupEntity(node dbtype.Node) *entity.ContactGroupEntity {
	props := utils.GetPropsFromNode(node)
	contactGroup := entity.ContactGroupEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &contactGroup
}
