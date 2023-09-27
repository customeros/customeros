package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type UserRepository interface {
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	FindUserByEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, email string) (*dbtype.Node, error)
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbNodes := make([]*dbtype.Node, 0)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) RETURN u `),
			params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodes, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, email string) (*dbtype.Node, error) {
	dbResult, err := tx.Run(ctx,
		`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email{rawEmail:$email}) 
				RETURN DISTINCT u limit 1`,
		map[string]interface{}{
			"tenant": tenant,
			"email":  email,
		})
	if err != nil {
		return nil, err
	}
	records, err := dbResult.Collect(ctx)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return utils.NodePtr(records[0].Values[0].(neo4j.Node)), nil
}
