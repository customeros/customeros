package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type EmailReadRepository interface {
	GetEmailIdIfExists(ctx context.Context, tenant, email string) (string, error)
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
