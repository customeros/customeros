package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error)
	FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error)
}

type userService struct {
	driver *neo4j.Driver
}

func NewUserService(driver *neo4j.Driver) UserService {
	return &userService{
		driver: driver,
	}
}

func (s *userService) Create(ctx context.Context, user *entity.UserEntity) (*entity.UserEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
	return s.mapDbNodeToUserEntity(queryResult.(dbtype.Node)), nil
}

func (s *userService) FindAll(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	dataResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) RETURN count(u) as count`,
			map[string]interface{}{
				"tenant": common.GetContext(ctx).Tenant,
			})
		count, _ := result.Single()
		paginatedResult.SetTotalRows(count.Values[0].(int64))

		result, err = tx.Run(`
				MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) RETURN u SKIP $skip LIMIT $limit`,
			map[string]interface{}{
				"tenant": common.GetContext(ctx).Tenant,
				"skip":   paginatedResult.GetSkip(),
				"limit":  paginatedResult.GetLimit(),
			})
		data, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	users := entity.UserEntities{}

	for _, dbRecord := range dataResult.([]*db.Record) {
		user := s.mapDbNodeToUserEntity(dbRecord.Values[0].(dbtype.Node))
		users = append(users, *user)
	}
	paginatedResult.SetRows(&users)
	return &paginatedResult, nil
}

func (s *userService) mapDbNodeToUserEntity(dbNode dbtype.Node) *entity.UserEntity {
	props := utils.GetPropsFromNode(dbNode)
	contact := entity.UserEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		FirstName: utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:  utils.GetStringPropOrEmpty(props, "lastName"),
		Email:     utils.GetStringPropOrEmpty(props, "email"),
		CreatedAt: props["createdAt"].(time.Time),
	}
	return &contact
}
