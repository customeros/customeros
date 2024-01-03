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

type PhoneNumberReadRepository interface {
	GetPhoneNumberIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
}

type phoneNumberReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewPhoneNumberReadRepository(driver *neo4j.DriverWithContext, database string) PhoneNumberReadRepository {
	return &phoneNumberReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *phoneNumberReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *phoneNumberReadRepository) GetPhoneNumberIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberReadRepository.GetPhoneNumberIdIfExists")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	cypher := fmt.Sprintf(`MATCH (p:PhoneNumber_%s) WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber RETURN p.id LIMIT 1`, tenant)
	params := map[string]any{
		"phoneNumber": phoneNumber,
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
