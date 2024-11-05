package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type LinkedinConnectionRequestReadRepository interface {
	GetPendingRequestByUserForSocialUrl(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId, socialUrl string) (*dbtype.Node, error)
	GetFastestUser(ctx context.Context, tx *neo4j.ManagedTransaction, userIds []string) (string, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "LinkedinConnectionRequestReadRepository.GetPendingRequestByUserForSocialUrl")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(log.String("tenant", tenant), log.String("userId", userId), log.String("socialUrl", socialUrl))

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.status = 'PENDING' and l.userId = $userId and l.socialUrl = $socialUrl return l`, tenant)

	params := map[string]any{
		"tenant":    tenant,
		"userId":    userId,
		"socialUrl": socialUrl,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil {
		return nil, nil
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *linkedinConnectionRequestReadRepository) GetFastestUser(ctx context.Context, tx *neo4j.ManagedTransaction, userIds []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetFastestUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.Object("userIds", userIds))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.status = 'PENDING' and l.userId in $userIds RETURN l.userId order by l.scheduledAt ASC limit 1`, tenant)
	params := map[string]any{
		"tenant":  tenant,
		"userIds": userIds,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (r *linkedinConnectionRequestReadRepository) GetLastScheduledForUser(ctx context.Context, tx *neo4j.ManagedTransaction, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetLastScheduledForUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.Object("userId", userId))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.status = 'PENDING' and l.userId in $userIds RETURN l.userId order by l.scheduledAt DESC limit 1`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"userId": userId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.CountRequestsPerUserPerDay")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("userId", userId), log.Object("startDate", startDate), log.Object("endDate", endDate))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest_%s) where l.scheduledAt >= $startDate and l.scheduledAt <= $endDate and l.userId = $userId RETURN count(l)`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"userId":    userId,
		"startDate": startDate,
		"endDate":   endDate,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	queryResult, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	count := queryResult.(*db.Record).Values[0].(int64)
	span.LogFields(log.Int64("result", count))
	return count, nil
}
