package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
	"time"
)

type TenantUserService interface {
	Create(ctx context.Context, tenantUser *entity.TenantUserEntity) (*entity.TenantUserEntity, error)
}

type tenantUserService struct {
	driver *neo4j.Driver
}

func NewTenantUserService(driver *neo4j.Driver) TenantUserService {
	return &tenantUserService{
		driver: driver,
	}
}

func (s *tenantUserService) Create(ctx context.Context, tenantUser *entity.TenantUserEntity) (*entity.TenantUserEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (u:TenantUser {
				  id: randomUUID(),
				  firstName: $firstName,
				  lastName: $lastName,
				  email: $email,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:USER_BELONGS_TO_TENANT]->(t)
			RETURN u`,
			map[string]interface{}{
				"firstName": tenantUser.FirstName,
				"lastName":  tenantUser.LastName,
				"email":     tenantUser.Email,
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
	return s.mapDbNodeToTenantUserEntity(queryResult.(dbtype.Node)), nil
}

func (s *tenantUserService) mapDbNodeToTenantUserEntity(dbNode dbtype.Node) *entity.TenantUserEntity {
	props := utils.GetPropsFromNode(dbNode)
	contact := entity.TenantUserEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		FirstName: utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:  utils.GetStringPropOrEmpty(props, "lastName"),
		Email:     utils.GetStringPropOrEmpty(props, "email"),
		CreatedAt: props["createdAt"].(time.Time),
	}
	return &contact
}
