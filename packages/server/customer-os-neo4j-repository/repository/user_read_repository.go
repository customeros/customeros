package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type UserReadRepository interface {
	GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error)
	FindFirstUserWithRolesByEmail(ctx context.Context, email string) (string, string, []string, error)
	GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error)
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

func (r *userReadRepository) GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserReadRepository.GetUserById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, userId)

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
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (u *userReadRepository) FindFirstUserWithRolesByEmail(ctx context.Context, email string) (string, string, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindFirstUserWithRolesByEmail")
	defer span.Finish()
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

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

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}
