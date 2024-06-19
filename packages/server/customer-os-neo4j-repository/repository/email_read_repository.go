package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type EmailReadRepository interface {
	GetEmailIdIfExists(ctx context.Context, tenant, email string) (string, error)
	GetEmailForUser(ctx context.Context, tenant string, userId string) (*dbtype.Node, error)
	GetById(ctx context.Context, tenant, emailId string) (*dbtype.Node, error)
}

type emailReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewEmailReadRepository(driver *neo4j.DriverWithContext, database string) EmailReadRepository {
	return &emailReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *emailReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *emailReadRepository) GetEmailIdIfExists(ctx context.Context, tenant, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReadRepository.GetEmailIdIfExists")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("email", email))

	cypher := fmt.Sprintf(`MATCH (e:Email_%s) WHERE e.email = $email OR e.rawEmail = $email RETURN e.id LIMIT 1`, tenant)
	params := map[string]any{
		"email": email,
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
		tracing.TraceErr(span, err)
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		span.LogFields(log.String("result", ""))
		return "", nil
	}
	span.LogFields(log.String("result", result.([]*db.Record)[0].Values[0].(string)))
	return result.([]*db.Record)[0].Values[0].(string), err
}

func (r *emailReadRepository) GetEmailForUser(ctx context.Context, tenant string, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReadRepository.GetEmailForUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("userId", userId))

	cypher := fmt.Sprintf("MATCH (e:Email_%s)<-[:HAS]-(u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) WHERE u:User_%s return e", tenant, tenant)
	params := map[string]any{
		"userId": userId,
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

func (r *emailReadRepository) GetById(ctx context.Context, tenant, emailId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReadRepository.GetById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("emailId", emailId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(:Email {id:$emailId}) return e`
	params := map[string]any{
		"tenant":  tenant,
		"emailId": emailId,
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
