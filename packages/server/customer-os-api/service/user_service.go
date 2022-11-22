package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"reflect"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error)
	FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	FindContactOwner(ctx context.Context, contactId string) (*entity.UserEntity, error)
	getDriver() neo4j.Driver
}

type userService struct {
	repository *repository.RepositoryContainer
}

func NewUserService(repository *repository.RepositoryContainer) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *userService) Create(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {
				  id: randomUUID(),
				  firstName: $firstName,
				  lastName: $lastName,
				  email: $email,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:USER_BELONGS_TO_TENANT]->(t)
			RETURN u`,
			map[string]interface{}{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"tenant":    common.GetContext(ctx).Tenant,
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
	return s.mapDbNodeToUserEntity(utils.NodePtr(queryResult.(dbtype.Node))), nil
}

func (s *userService) FindAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.UserEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(entity.UserEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repository.UserRepository.GetPaginatedUsers(
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

	users := entity.UserEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		users = append(users, *s.mapDbNodeToUserEntity(v))
	}
	paginatedResult.SetRows(&users)
	return &paginatedResult, nil
}

func (s *userService) FindContactOwner(ctx context.Context, contactId string) (*entity.UserEntity, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	ownerDbNode, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return s.repository.UserRepository.FindOwnerForContact(tx, common.GetContext(ctx).Tenant, contactId)
	})
	if err != nil {
		return nil, err
	} else if ownerDbNode.(*dbtype.Node) == nil {
		return nil, nil
	} else {
		return s.mapDbNodeToUserEntity(ownerDbNode.(*dbtype.Node)), nil
	}
}

func (s *userService) mapDbNodeToUserEntity(dbNode *dbtype.Node) *entity.UserEntity {
	props := utils.GetPropsFromNode(*dbNode)
	contact := entity.UserEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		FirstName: utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:  utils.GetStringPropOrEmpty(props, "lastName"),
		Email:     utils.GetStringPropOrEmpty(props, "email"),
		CreatedAt: props["createdAt"].(time.Time),
	}
	return &contact
}
