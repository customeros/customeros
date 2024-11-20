package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type UserReadRepository interface {
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error)
	FindFirstUserWithRolesByEmail(ctx context.Context, email string) (string, string, []string, error)
	FindTestUser(ctx context.Context) (*dbtype.Node, error)
	GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error)
	GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error)
	GetOwnerForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	IsOwner(parentCtx context.Context, tenant, userId string) (bool, error)
	GetOwnerForContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetCreatorForNote(ctx context.Context, tenant, noteId string) (*dbtype.Node, error)
	GetPaginatedCustomerUsers(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	GetAllOwnersForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForServiceLineItems(ctx context.Context, tenant string, serviceLineItemIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
	GetAllAuthorsForLogEntries(ctx context.Context, tenant string, logEntryIDs []string) ([]*utils.DbNodeAndId, error)
	GetAllAuthorsForComments(ctx context.Context, tenant string, commentIds []string) ([]*utils.DbNodeAndId, error)
	GetAllSendersForFlowSenders(ctx context.Context, tenant string, flowSenderIds []string) ([]*utils.DbNodeAndId, error)
	GetUsersConnectedForContacts(ctx context.Context, tenant string, contactsIds []string) ([]*utils.DbNodeAndId, error)
	GetDistinctOrganizationOwners(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetOwnerForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOwnerForReminder(ctx context.Context, tenant, reminderId string) (*dbtype.Node, error)
}

type userReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewUserReadRepository(driver *neo4j.DriverWithContext, database string) UserReadRepository {
	return &userReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *userReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *userReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetAllForTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := r.prepareReadSession(ctx)
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

func (r *userReadRepository) GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetByIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodes := make([]*dbtype.Node, 0)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"tenant": tenant,
			"ids":    ids,
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) where u.id in $ids RETURN u`),
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

func (r *userReadRepository) GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetUserById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("userId", userId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$id}) RETURN u`
	params := map[string]any{
		"tenant": tenant,
		"id":     userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
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

func (u *userReadRepository) FindFirstUserWithRolesByEmail(ctx context.Context, email string) (string, string, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.FindFirstUserWithRolesByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("email", email))

	session := utils.NewNeo4jReadSession(ctx, *u.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (e:Email)<-[:HAS]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN t.name, u.id, u.roles ORDER BY u.createdAt ASC LIMIT 1`,
			map[string]interface{}{
				"email": email,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", "", []string{}, err
	}
	if len(records.([]*neo4j.Record)) > 0 {
		tenant := records.([]*neo4j.Record)[0].Values[0].(string)
		userId := records.([]*neo4j.Record)[0].Values[1].(string)
		roleList, ok := records.([]*neo4j.Record)[0].Values[2].([]interface{})
		var roles []string
		if !ok {
			roles = []string{}
		} else {
			roles = u.toStringList(roleList)
		}
		return userId, tenant, roles, nil
	} else {
		return "", "", []string{}, nil
	}
}

func (u *userReadRepository) toStringList(values []interface{}) []string {
	var result []string
	for _, value := range values {
		result = append(result, value.(string))
	}
	return result
}

func (r *userReadRepository) GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetFirstUserByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT(u) ORDER by u.createdAt ASC limit 1`
	params := map[string]any{
		"tenant": tenant,
		"email":  email,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) FindTestUser(ctx context.Context) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.FindTestUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)
			WHERE u.test = true and u.firstName = "Test" and u.lastName = "Sender"
			RETURN u limit 1`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetAllOwnersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:OWNS]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN u, o.id as orgId`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIDs,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetOwnerForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetOwnerForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})<-[:OWNS]-(u:User) RETURN u`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})

	if err != nil {
		return nil, err
	} else if len(result.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		return utils.NodePtr(result.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
	}
}

func (r *userReadRepository) IsOwner(parentCtx context.Context, tenant, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.IsOwner")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:OWNS]-(u:User{id:$userId})
			RETURN count(o)`
	params := map[string]any{
		"tenant": tenant,
		"userId": userId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	queryResult, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	count := queryResult.(*db.Record).Values[0].(int64)
	span.LogFields(log.Int64("result", count))
	return count > 0, nil
}

func (r *userReadRepository) GetOwnerForContact(parentCtx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetOwnerForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})<-[:OWNS]-(u:User)
			RETURN u`)
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (r *userReadRepository) GetCreatorForNote(parentCtx context.Context, tenant, noteId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetCreatorForNote")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:CREATED]->(n:Note {id:$noteId})
			RETURN u`)
	params := map[string]any{
		"tenant": tenant,
		"noteId": noteId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (r *userReadRepository) GetPaginatedCustomerUsers(parentCtx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetPaginatedCustomerUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("u")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) 
											WHERE (u.internal=false OR u.internal is null) AND (u.test=false OR u.test is null) AND (u.bot=false OR u.bot is null)
											WITH u
											%s RETURN count(u) as count`, filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) 
					WHERE (u.internal=false OR u.internal is null) AND (u.test=false OR u.test is null) AND (u.bot=false OR u.bot is null)
					WITH u
					%s
					RETURN u 
					%s 
					SKIP $skip LIMIT $limit`, filterCypherStr, sort.SortingCypherFragment("u")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *userReadRepository) GetAllForEmails(parentCtx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE e.id IN $emailIds
			RETURN u, e.id as emailId ORDER BY u.firstName, u.lastName`,
			map[string]any{
				"tenant":   tenant,
				"emailIds": emailIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllForPhoneNumbers(parentCtx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
			WHERE p.id IN $phoneNumberIds
			RETURN u, p.id as phoneNumberId ORDER BY u.firstName, u.lastName`,
			map[string]any{
				"tenant":         tenant,
				"phoneNumberIds": phoneNumberIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllOwnersForOpportunities(parentCtx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllOwnersForOpportunities")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:OWNS]->(op:Opportunity)
			WHERE op.id IN $opportunityIds
			RETURN u, op.id as opId`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllCreatorsForOpportunities(parentCtx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllCreatorsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(op:Opportunity)
			WHERE op.id IN $opportunityIds
			RETURN u, op.id as opId`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllCreatorsForServiceLineItems(parentCtx context.Context, tenant string, serviceLineItemIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllCreatorsForServiceLineItems")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(sli:ServiceLineItem)
			WHERE sli.id IN serviceLineItemIds
			RETURN u, sli.id as sliId`
	params := map[string]any{
		"tenant":             tenant,
		"serviceLineItemIds": serviceLineItemIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllCreatorsForContracts(parentCtx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllCreatorsForContracts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(c:Contract)
			WHERE c.id IN $contractIds
			RETURN u, c.id as cId`
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllAuthorsForLogEntries(parentCtx context.Context, tenant string, logEntryIDs []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllAuthorsForLogEntries")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("logEntryIDs", logEntryIDs))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(l:LogEntry_%s)
			WHERE l.id IN $logEntryIDs
			RETURN u, l.id as logEntryId`, tenant)
	span.LogFields(log.String("cypher", cypher))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
			map[string]any{
				"tenant":      tenant,
				"logEntryIDs": logEntryIDs,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllAuthorsForComments(parentCtx context.Context, tenant string, commentIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetAllAuthorsForComments")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("commentIds", commentIds))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(c:Comment_%s)
			WHERE c.id IN $commentIds
			RETURN u, c.id`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"commentIds": commentIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllSendersForFlowSenders(ctx context.Context, tenant string, flowSenderIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetAllSendersForFlowSenders")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("flowSenderIds", flowSenderIds))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:HAS]-(fs:FlowSender_%s)
			WHERE fs.id IN $flowSenderIds
			RETURN u, fs.id`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"flowSenderIds": flowSenderIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetUsersConnectedForContacts(ctx context.Context, tenant string, contactsIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetUsersConnectedForContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("contactsIds", contactsIds))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CONNECTED_WITH]-(c:Contact_%s)
			WHERE c.id IN $contactsIds
			RETURN u, c.id`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"contactsIds": contactsIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetDistinctOrganizationOwners(parentCtx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetDistinctOrganizationOwners")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:OWNS]->(:Organization)
			RETURN distinct(u) order by u.firstName, u.lastName, u.name`

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*dbtype.Node), err
}

func (r *userReadRepository) GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)
			WHERE u.id IN $ids
			RETURN u`
	span.LogFields(log.String("cypher", cypher))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)
	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*dbtype.Node), err
}

func (r *userReadRepository) GetOwnerForContract(parentCtx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetOwnerForContract")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})<-[:OWNS]-(u:User)
			RETURN u`)
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (r *userReadRepository) GetOwnerForReminder(parentCtx context.Context, tenant, reminderId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserReadRepository.GetOwnerForReminder")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$reminderId})-[:REMINDER_BELONGS_TO_USER]->(u:User)
			RETURN u`)
	params := map[string]any{
		"tenant":     tenant,
		"reminderId": reminderId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
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
