package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type CountryWriteRepository interface {
	CreateCountry(ctx context.Context, id, name, codeA2, codeA3, phoneCode string, createdAt time.Time) error
	UpdateCountry(ctx context.Context, id, name, codeA2, codeA3, phoneCode string) error
}

type countryWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCountryWriteRepository(driver *neo4j.DriverWithContext, database string) CountryWriteRepository {
	return &countryWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *countryWriteRepository) prepareWriteSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *countryWriteRepository) CreateCountry(ctx context.Context, id, name, codeA2, codeA3, phoneCode string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetDefaultCountryCodeA3")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "id", id)
	tracing.LogObjectAsJson(span, "name", name)
	tracing.LogObjectAsJson(span, "codeA2", codeA2)
	tracing.LogObjectAsJson(span, "codeA3", codeA3)
	tracing.LogObjectAsJson(span, "phoneCode", phoneCode)
	tracing.LogObjectAsJson(span, "createdAt", createdAt)

	cypher := " MERGE (c:Country {id: $id})" +
		" ON CREATE SET c.name=$name, " +
		"				c.codeA2=$codeA2, " +
		"				c.codeA3=$codeA3, " +
		"				c.phoneCode=$phoneCode, " +
		" 				c.createdAt=$createdAt, " +
		" 				c.updatedAt=datetime() "
	params := map[string]any{
		"id":        id,
		"name":      name,
		"codeA2":    codeA2,
		"codeA3":    codeA3,
		"phoneCode": phoneCode,
		"createdAt": createdAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *countryWriteRepository) UpdateCountry(ctx context.Context, id, name, codeA2, codeA3, phoneCode string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryReadRepository.GetDefaultCountryCodeA3")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "id", id)
	tracing.LogObjectAsJson(span, "name", name)
	tracing.LogObjectAsJson(span, "codeA2", codeA2)
	tracing.LogObjectAsJson(span, "codeA3", codeA3)
	tracing.LogObjectAsJson(span, "phoneCode", phoneCode)

	cypher := `
			MATCH (c:Country {id:$id})
			SET c.name=$name, 
				c.codeA2=$codeA2,
				c.codeA3=$codeA3,
				c.phoneCode=$phoneCode,
				c.updatedAt=datetime()`
	params := map[string]any{
		"id":        id,
		"name":      name,
		"codeA2":    codeA2,
		"codeA3":    codeA3,
		"phoneCode": phoneCode,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
