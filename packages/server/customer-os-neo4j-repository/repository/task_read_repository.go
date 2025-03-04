package neo4j_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TaskReadRepository interface {
	GetById(ctx context.Context, tenant, taskId string) (*dbtype.Node, error)
	GetAllByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	CountByTenant(ctx context.Context, tenant string) (int64, error)
	SearchTasks(ctx context.Context, tenant string, limit int, where *model.Filter, sort *model.SortBy) (*utils.StringsWithTotalCount, error)
}

type taskReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *taskReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func NewTaskReadRepository(driver *neo4j.DriverWithContext, database string) TaskReadRepository {
	return &taskReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *taskReadRepository) GetById(ctx context.Context, tenant, taskId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.TagEntity(span, taskId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task {id:$taskId})
				RETURN tsk`
	params := map[string]any{
		"tenant": tenant,
		"taskId": taskId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *taskReadRepository) GetAllByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.GetAllByIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "ids", ids)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task)
				WHERE tsk.id IN $ids
				RETURN tsk`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(dbRecords.([]*dbtype.Node))))
	return dbRecords.([]*dbtype.Node), err
}

func (r *taskReadRepository) CountByTenant(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.CountByTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task)
				RETURN count(tsk) as count`
	params := map[string]any{
		"tenant": tenant,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0].(int64), nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}
	span.LogFields(log.Int64("result", count.(int64)))
	return count.(int64), nil
}

func (r *taskReadRepository) SearchTasks(ctx context.Context, tenant string, limit int, where *model.Filter, sort *model.SortBy) (*utils.StringsWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.SearchTasks")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit))
	tracing.LogObjectAsJson(span, "where", where)
	tracing.LogObjectAsJson(span, "sort", sort)

	taskFilterCypher, taskFilterParams := "", make(map[string]interface{})
	userAuthorFilterCypher, userAuthorFilterParams := "", make(map[string]interface{})
	userAssigneeFilterCypher, userAssigneeFilterParams := "", make(map[string]interface{})

	if where != nil {
		taskFilter := new(utils.CypherFilter)
		taskFilter.Negate = false
		taskFilter.LogicalOperator = utils.AND
		taskFilter.Filters = make([]*utils.CypherFilter, 0)

		userAuthorFilterCypher := new(utils.CypherFilter)
		userAuthorFilterCypher.Negate = false
		userAuthorFilterCypher.LogicalOperator = utils.AND
		userAuthorFilterCypher.Filters = make([]*utils.CypherFilter, 0)

		userAssigneeFilterCypher := new(utils.CypherFilter)
		userAssigneeFilterCypher.Negate = false
		userAssigneeFilterCypher.LogicalOperator = utils.AND
		userAssigneeFilterCypher.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {

		}
	}

	params := make(map[string]any)
	params["tenant"] = tenant
	params["limit"] = limit

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task)
				OPTIONAL MATCH tsk<-[:TASK_ASSIGNED_TO]-(assignee:User)
				OPTIONAL MATCH tsk<-[:TASK_CREATED_BY]-(author:User)
				OPTIONAL MATCH tsk<-[:LINKED_TO]-(opp:Opportunity)
				WHERE 1=1`

	if where != nil && where.Filter != nil {
		switch where.Filter.Property {
		case "TASKS_SUBJECT":
			if where.Filter.Operation == model.ComparisonOperatorContains {
				cypher += ` AND toLower(tsk.subject) CONTAINS toLower($subject)`
				params["subject"] = where.Filter.Value.Str
			}
		case "TASKS_STATUS":
			if where.Filter.Operation == model.ComparisonOperatorEquals {
				cypher += ` AND tsk.status = $status`
				params["status"] = where.Filter.Value.Str
			}
		case "TASKS_DUE_DATE":
			if where.Filter.Operation == model.ComparisonOperatorLte {
				cypher += ` AND tsk.dueAt IS NOT NULL AND tsk.dueAt <= $dueDate`
				params["dueDate"] = where.Filter.Value.Time
			}
		case "TASKS_ASSIGNEES":
			if where.Filter.Operation == model.ComparisonOperatorIn {
				cypher += ` AND assignee.id IN $assigneeIds`
				params["assigneeIds"] = where.Filter.Value.ArrayStr
			}
		case "TASKS_OPPORTUNITIES":
			if where.Filter.Operation == model.ComparisonOperatorIn {
				cypher += ` AND opp.id IN $opportunityIds`
				params["opportunityIds"] = where.Filter.Value.ArrayStr
			}
		}
	}

	cypher += ` WITH tsk, count(*) as totalCount`

	if sort != nil {
		switch sort.By {
		case "TASKS_SUBJECT":
			cypher += ` ORDER BY tsk.subject`
		case "TASKS_STATUS":
			cypher += ` ORDER BY tsk.status`
		case "TASKS_DESCRIPTION":
			cypher += ` ORDER BY tsk.description`
		case "TASKS_AUTHOR":
			cypher += ` ORDER BY author.firstName, author.lastName`
		case "TASKS_DUE_DATE":
			cypher += ` ORDER BY tsk.dueAt`
		case "TASKS_ASSIGNEES":
			cypher += ` ORDER BY assignee.firstName, assignee.lastName`
		case "TASKS_OPPORTUNITIES":
			cypher += ` ORDER BY opp.name`
		default:
			cypher += ` ORDER BY tsk.createdAt`
		}

		if sort.Direction == "DESC" {
			cypher += ` DESC`
		}
	}

	cypher += ` RETURN tsk.id as id, totalCount
				SKIP 0 LIMIT $limit`

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var ids []string
		var totalCount int64
		records, err := queryResult.Collect(ctx)
		if err != nil {
			return nil, err
		}

		if len(records) > 0 {
			totalCount = records[0].Values[1].(int64)
			for _, record := range records {
				ids = append(ids, record.Values[0].(string))
			}
		}

		return &utils.StringsWithTotalCount{
			Strings: ids,
			Count:   totalCount,
		}, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result.(*utils.StringsWithTotalCount), nil
}
