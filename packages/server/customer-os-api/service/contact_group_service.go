package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"reflect"
)

type ContactGroupService interface {
	Create(ctx context.Context, contactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error)
	Update(ctx context.Context, contactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error)
	Delete(ctx context.Context, id string) (bool, error)

	FindContactGroupById(ctx context.Context, id string) (*entity.ContactGroupEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupEntities, error)

	AddContactToGroup(ctx context.Context, contactId, groupId string) (bool, error)
	RemoveContactFromGroup(ctx context.Context, contactId, groupId string) (bool, error)
}

type contactGroupService struct {
	repository *repository.RepositoryContainer
}

func NewContactGroupService(repository *repository.RepositoryContainer) ContactGroupService {
	return &contactGroupService{
		repository: repository,
	}
}

func (s *contactGroupService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *contactGroupService) Create(ctx context.Context, newContactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (g:ContactGroup {
				  id: randomUUID(),
				  name: $name})-[:GROUP_BELONGS_TO_TENANT]->(t)
			RETURN g`,
			map[string]any{
				"name":   newContactGroup.Name,
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
	return s.mapDbNodeToContactGroup(utils.NodePtr(queryResult.(dbtype.Node))), nil
}

func (s *contactGroupService) Update(ctx context.Context, contactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
			MATCH (g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			SET g.name=$name
			RETURN g`,
			map[string]any{
				"tenant":  common.GetContext(ctx).Tenant,
				"groupId": contactGroup.Id,
				"name":    contactGroup.Name,
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
	return s.mapDbNodeToContactGroup(utils.NodePtr(queryResult.(dbtype.Node))), nil
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

	dbNodesWithTotalCount, err := s.repository.ContactGroupRepository.GetPaginatedContactGroups(
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
		contactGroups = append(contactGroups, *s.mapDbNodeToContactGroup(v))
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

	return s.mapDbNodeToContactGroup(utils.NodePtr(queryResult.(dbtype.Node))), nil
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
		contactGroup := s.mapDbNodeToContactGroup(utils.NodePtr(dbRecord.Values[0].(dbtype.Node)))
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

func (s *contactGroupService) mapDbNodeToContactGroup(dbContactGroupNode *dbtype.Node) *entity.ContactGroupEntity {
	props := utils.GetPropsFromNode(*dbContactGroupNode)
	contactGroup := entity.ContactGroupEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &contactGroup
}
