package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"time"
)

type LinkedinConnectionRequestReadRepository interface {
	GetPendingRequestByUserForSocialUrl(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId, socialUrl string) (*dbtype.Node, error)
	GetLastScheduledForUser(ctx context.Context, tx *neo4j.ManagedTransaction, userId string) (*dbtype.Node, error)

	CountRequestsPerUserPerDay(ctx context.Context, tx *neo4j.ManagedTransaction, userId string, startDate, endDate time.Time) (int64, error)
}

type linkedinConnectionRequestReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLinkedinConnectionRequestReadRepository(driver *neo4j.DriverWithContext, database string) LinkedinConnectionRequestReadRepository {
	return &linkedinConnectionRequestReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *linkedinConnectionRequestReadRepository) GetPendingRequestByUserForSocialUrl(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId, socialUrl string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LinkedinConnectionRequestReadRepository.GetPendingRequestByUserForSocialUrl")
	defer spans.Finish()

	spans.LogKV("tenant", tenant)
	spans.LogKV("userId", userId)
	spans.LogKV("socialUrl", socialUrl)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.status = 'PENDING' and l.userId = $userId and l.socialUrl = $socialUrl return l`, tenant)

	params := map[string]any{
		"tenant":    tenant,
		"userId":    userId,
		"socialUrl": socialUrl,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *linkedinConnectionRequestReadRepository) GetLastScheduledForUser(ctx context.Context, tx *neo4j.ManagedTransaction, userId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionExecutionReadRepository.GetLastScheduledForUser")
	defer spans.Finish()

	spans.LogKV("userId", userId)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.userId = $userId RETURN l order by l.scheduledAt DESC limit 1`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"userId": userId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *linkedinConnectionRequestReadRepository) CountRequestsPerUserPerDay(ctx context.Context, tx *neo4j.ManagedTransaction, userId string, startDate, endDate time.Time) (int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionExecutionReadRepository.CountRequestsPerUserPerDay")
	defer spans.Finish()

	spans.LogKV("userId", userId)
	spans.LogObjectAsJson("startDate", startDate)
	spans.LogObjectAsJson("endDate", endDate)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.scheduledAt >= $startDate and l.scheduledAt <= $endDate and l.userId = $userId RETURN count(l)`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"userId":    userId,
		"startDate": startDate,
		"endDate":   endDate,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	queryResult, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})

	if err != nil {
		spans.TraceError(err)
		return 0, err
	}

	count := queryResult.(*db.Record).Values[0].(int64)
	spans.LogKV("result", count)
	return count, nil
}
