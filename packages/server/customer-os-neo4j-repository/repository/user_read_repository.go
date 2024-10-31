package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error)
	GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error)
	GetOwnerForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindFirstUserWithRolesByEmail")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetFirstUserByEmail")
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
