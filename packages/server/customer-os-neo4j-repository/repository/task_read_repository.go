package neo4j_repository

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	TasksAssignees     string = "TASKS_ASSIGNEES"
	TasksDueDate       string = "TASKS_DUE_DATE"
	TasksStatus        string = "TASKS_STATUS"
	TasksAuthor        string = "TASKS_AUTHOR"
	TasksCreatedAt     string = "TASKS_CREATED_AT"
	TasksUpdatedAt     string = "TASKS_UPDATED_AT"
	TasksOpportunities string = "TASKS_OPPORTUNITIES"
	TasksSubject       string = "TASKS_SUBJECT"
	TasksDescription   string = "TASKS_DESCRIPTION"
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

		userAuthorFilter := new(utils.CypherFilter)
		userAuthorFilter.Negate = false
		userAuthorFilter.LogicalOperator = utils.AND
		userAuthorFilter.Filters = make([]*utils.CypherFilter, 0)

		userAssigneeFilter := new(utils.CypherFilter)
		userAssigneeFilter.Negate = false
		userAssigneeFilter.LogicalOperator = utils.AND
		userAssigneeFilter.Filters = make([]*utils.CypherFilter, 0)

		for _, filter := range where.And {
			if filter.Filter.Property == TasksSubject {
				taskFilter.Filters = append(taskFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.TaskPropertySubject), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == TasksDescription {
				taskFilter.Filters = append(taskFilter.Filters, utils.CreateStringCypherFilter(string(neo4jentity.TaskPropertyDescription), filter.Filter.Value.Str, filter.Filter.Operation))
			}
			if filter.Filter.Property == TasksStatus {
				createInOrEmptyStringFilter(filter, taskFilter, string(neo4jentity.TaskPropertyStatus))
			}
			if filter.Filter.Property == TasksAssignees {
				// special case for not in tags
				if filter.Filter.Operation == model.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
					rawCypher := ""
					for _, v := range *filter.Filter.Value.ArrayStr {
						if rawCypher != "" {
							rawCypher += " AND "
						}
						rawCypher += fmt.Sprintf(` NOT (tsk)-[:ASSIGNED_TO]->(:User {id:"%s"}) `, v)
					}
					userAssigneeFilter.Filters = append(userAssigneeFilter.Filters, utils.CreateRawCypherFilter(rawCypher))
				} else {
					createInOrEmptyStringFilter(filter, userAssigneeFilter, "id")
				}
			}
			if filter.Filter.Property == TasksAuthor {
				createInOrEmptyStringFilter(filter, userAuthorFilter, "id")
			}
			if filter.Filter.Property == TasksDueDate {
				createTimeFilter(filter, taskFilter, string(neo4jentity.TaskPropertyDueAt))
			}
		}

		if len(taskFilter.Filters) > 0 {
			taskFilterCypher, taskFilterParams = taskFilter.BuildCypherFilterFragmentWithParamName("tsk", "tsk_param_")
		}
		if len(userAuthorFilter.Filters) > 0 {
			userAuthorFilterCypher, userAuthorFilterParams = userAuthorFilter.BuildCypherFilterFragmentWithParamName("ua", "ua_param_")
		}
		if len(userAssigneeFilter.Filters) > 0 {
			userAssigneeFilterCypher, userAssigneeFilterParams = userAssigneeFilter.BuildCypherFilterFragmentWithParamName("uas", "uas_param_")
		}
	}

	params := make(map[string]any)
	params["tenant"] = tenant
	params["limit"] = limit

	utils.MergeMapToMap(taskFilterParams, params)
	utils.MergeMapToMap(userAuthorFilterParams, params)
	utils.MergeMapToMap(userAssigneeFilterParams, params)

	//region count selectQuery
	countQuery := ""
	{
		countQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task_%s) `, tenant)
		if userAuthorFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (tsk)-[:CREATED_BY]->(ua:User) WITH *`
		}
		if userAssigneeFilterCypher != "" {
			countQuery += ` OPTIONAL MATCH (tsk)-[:ASSIGNED_TO]->(uas:User) WITH *`
		}

		if taskFilterCypher != "" || userAuthorFilterCypher != "" || userAssigneeFilterCypher != "" {
			countQuery += " WHERE "
		}

		countQueryParts := []string{}
		if taskFilterCypher != "" {
			countQueryParts = append(countQueryParts, taskFilterCypher)
		}
		if userAuthorFilterCypher != "" {
			countQueryParts = append(countQueryParts, userAuthorFilterCypher)
		}
		if userAssigneeFilterCypher != "" {
			countQueryParts = append(countQueryParts, userAssigneeFilterCypher)
		}
		countQuery += strings.Join(countQueryParts, " AND ") + ` RETURN count(distinct(tsk))`
	}
	//end count region

	//region selectQuery
	selectQuery := ""
	{
		selectQuery += fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task_%s) `, tenant)
		if userAuthorFilterCypher != "" || (sort != nil && (sort.By == TasksAuthor)) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (tsk)-[:CREATED_BY]->(ua:User) WITH *`)
		}
		if userAssigneeFilterCypher != "" || (sort != nil && (sort.By == TasksAssignees)) {
			selectQuery += fmt.Sprintf(` OPTIONAL MATCH (tsk)-[:ASSIGNED_TO]->(uas:User) WITH *`)
		}
		selectQuery += ` WHERE 1=1 `
		if taskFilterCypher != "" {
			selectQuery += fmt.Sprintf(` AND %s`, taskFilterCypher)
		}
		if userAuthorFilterCypher != "" {
			selectQuery += fmt.Sprintf(` AND %s`, userAuthorFilterCypher)
		}
		if userAssigneeFilterCypher != "" {
			selectQuery += fmt.Sprintf(` AND %s`, userAssigneeFilterCypher)
		}

		aliases := ""
		if sort != nil {
			if sort.By == TasksSubject {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN trim(tsk.subject) <> '' and not tsk.subject is null THEN toLower(trim(tsk.subject)) ELSE '𠀀' END as SORT_BY `
				} else {
					aliases += `CASE WHEN trim(tsk.subject) <> '' and not tsk.subject is null THEN toLower(trim(tsk.subject)) ELSE '' END as SORT_BY `
				}
			}
			if sort.By == TasksDescription {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN trim(tsk.description) <> '' and not tsk.description is null THEN toLower(trim(tsk.description)) ELSE '𠀀' END as SORT_BY `
				} else {
					aliases += `CASE WHEN trim(tsk.description) <> '' and not tsk.description is null THEN toLower(trim(tsk.description)) ELSE '' END as SORT_BY `
				}
			}
			if sort.By == TasksDueDate {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN tsk.dueAt IS NOT NULL THEN tsk.dueAt ELSE datetime({year:2100}) END as SORT_BY `
				} else {
					aliases += `CASE WHEN tsk.dueAt IS NOT NULL THEN tsk.dueAt ELSE datetime({year:1900}) END as SORT_BY `
				}
			}
			if sort.By == TasksCreatedAt {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN tsk.createdAt IS NOT NULL THEN tsk.createdAt ELSE datetime({year:2100}) END as SORT_BY `
				} else {
					aliases += `CASE WHEN tsk.createdAt IS NOT NULL THEN tsk.createdAt ELSE datetime({year:1900}) END as SORT_BY `
				}
			}
			if sort.By == TasksUpdatedAt {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN tsk.updatedAt IS NOT NULL THEN tsk.updatedAt ELSE datetime({year:2100}) END as SORT_BY `
				} else {
					aliases += `CASE WHEN tsk.updatedAt IS NOT NULL THEN tsk.updatedAt ELSE datetime({year:1900}) END as SORT_BY `
				}
			}
			if sort.By == TasksStatus {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN trim(tsk.status) <> '' and not tsk.status is null THEN toLower(trim(tsk.status)) ELSE '' END as SORT_BY `
				} else {
					aliases += `CASE WHEN trim(tsk.status) <> '' and not tsk.status is null THEN toLower(trim(tsk.status)) ELSE '𠀀' END as SORT_BY `
				}
			}
			if sort.By == TasksAssignees {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN (COALESCE(uas.firstName, '') + COALESCE(uas.lastName, '')) <> '' THEN toLower(trim(COALESCE(uas.firstName, '') + COALESCE(uas.lastName, ''))) ELSE '𠀀' END as SORT_BY `
				} else {
					aliases += `CASE WHEN (COALESCE(uas.firstName, '') + COALESCE(uas.lastName, '')) <> '' THEN toLower(trim(COALESCE(uas.firstName, '') + COALESCE(uas.lastName, ''))) ELSE '' END as SORT_BY `
				}
			}
			if sort.By == TasksAuthor {
				if sort.Direction == model.SortingDirectionAsc {
					aliases += `CASE WHEN (COALESCE(ua.firstName, '') + COALESCE(ua.lastName, '')) <> '' THEN toLower(trim(COALESCE(ua.firstName, '') + COALESCE(ua.lastName, ''))) ELSE '𠀀' END as SORT_BY `
				} else {
					aliases += `CASE WHEN (COALESCE(ua.firstName, '') + COALESCE(ua.lastName, '')) <> '' THEN toLower(trim(COALESCE(ua.firstName, '') + COALESCE(ua.lastName, ''))) ELSE '' END as SORT_BY `
				}
			}
		}
		// Always add task ID sorting at the end
		selectQuery += "WITH tsk "
		if aliases != "" {
			selectQuery += `, ` + aliases
		}
		if sort != nil && aliases != "" {
			selectQuery += " ORDER BY SORT_BY " + string(sort.Direction) + " , tsk.id ASC "
		} else {
			selectQuery += " ORDER BY tsk.id ASC"
		}
		selectQuery += " RETURN distinct(tsk.id) LIMIT $limit"
	}

	//end selectQuery

	stringsWithTotalCount := &utils.StringsWithTotalCount{}

	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	setError := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.SearchTasks.CountQuery")
		defer span.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		tracing.LogObjectAsJson(span, "params", params)
		span.LogFields(log.String("countQuery", countQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		countRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			countQueryResult, err := tx.Run(ctx, countQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, countQueryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(span, err)
			setError(err)
			return
		}

		result.Count = countRecord.(int64)
	}(ctx, stringsWithTotalCount)

	wg.Add(1)
	go func(ctx context.Context, result *utils.StringsWithTotalCount) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "TaskReadRepository.SearchTasks.SelectQuery")
		defer span.Finish()
		defer wg.Done()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		tracing.LogObjectAsJson(span, "params", params)
		span.LogFields(log.String("selectQuery", selectQuery))

		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			queryResult, err := tx.Run(ctx, selectQuery, params)
			if err != nil {
				return nil, err
			} else {
				return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
			}
		})

		if err != nil {
			tracing.TraceErr(span, err)
			setError(err)
			return
		}

		result.Strings = dbRecords.([]string)
	}(ctx, stringsWithTotalCount)

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	span.LogFields(log.Int("result.count", len(stringsWithTotalCount.Strings)), log.Int64("result.totalCount", stringsWithTotalCount.Count))

	return stringsWithTotalCount, nil
}

func createInOrEmptyStringFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Operation == model.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, model.ComparisonOperatorIsEmpty))
	} else if filter.Filter.Operation == model.ComparisonOperatorIsNotEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, model.ComparisonOperatorIsNotEmpty))
	} else if filter.Filter.Operation == model.ComparisonOperatorIn && filter.Filter.Value.ArrayStr != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.ArrayStr, model.ComparisonOperatorIn))
	} else if filter.Filter.Operation == model.ComparisonOperatorNotIn && filter.Filter.Value.ArrayStr != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, filter.Filter.Value.ArrayStr, model.ComparisonOperatorNotIn))
	}
}

func createTimeFilter(filter *model.Filter, cypherFilter *utils.CypherFilter, neo4jProperty string) {
	if filter.Filter.Operation == model.ComparisonOperatorBetween && filter.Filter.Value.ArrayTime != nil && len(*filter.Filter.Value.ArrayTime) == 2 {
		times := *filter.Filter.Value.ArrayTime
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[0], model.ComparisonOperatorGte))
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, times[1], model.ComparisonOperatorLte))
	} else if filter.Filter.Operation == model.ComparisonOperatorIsEmpty {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, nil, model.ComparisonOperatorIsEmpty))
	} else if filter.Filter.Operation == model.ComparisonOperatorGte && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, model.ComparisonOperatorGte))
	} else if filter.Filter.Operation == model.ComparisonOperatorGt && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, model.ComparisonOperatorGt))
	} else if filter.Filter.Operation == model.ComparisonOperatorLte && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, model.ComparisonOperatorLte))
	} else if filter.Filter.Operation == model.ComparisonOperatorLt && filter.Filter.Value.Time != nil {
		cypherFilter.Filters = append(cypherFilter.Filters, utils.CreateCypherFilter(neo4jProperty, *filter.Filter.Value.Time, model.ComparisonOperatorLt))
	}
}
