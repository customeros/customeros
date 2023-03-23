package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type IssueRepository interface {
	GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error)
}

type issueRepository struct {
	driver *neo4j.DriverWithContext
}

func NewIssueRepository(driver *neo4j.DriverWithContext) IssueRepository {
	return &issueRepository{
		driver: driver,
	}
}

func (r *issueRepository) GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})<-[:REPORTED_BY]-(i:Issue)
			WITH DISTINCT i
			RETURN i.status AS status, COUNT(i) AS count`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		status := ""
		if v.Values[0] != nil {
			status = v.Values[0].(string)
		}
		output[status] = v.Values[1].(int64)
	}
	return output, err
}
