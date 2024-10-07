package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type EmailRepository interface {
	FindEmailsForUser(ctx context.Context, tenant string, userId string) ([]*dbtype.Node, error)
}

type emailRepository struct {
	driver *neo4j.DriverWithContext
}

func NewEmailRepository(driver *neo4j.DriverWithContext) EmailRepository {
	return &emailRepository{
		driver: driver,
	}
}

func (r *emailRepository) FindEmailsForUser(ctx context.Context, tenant string, userId string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User{id:$userId})-[:HAS]->(e:Email) 
			RETURN DISTINCT e`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}
