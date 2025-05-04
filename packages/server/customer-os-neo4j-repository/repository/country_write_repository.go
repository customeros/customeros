package neo4j_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CountryReadRepository.GetDefaultCountryCodeA3")
	defer spans.Finish()
	spans.LogObjectAsJson("id", id)
	spans.LogObjectAsJson("name", name)
	spans.LogObjectAsJson("codeA2", codeA2)
	spans.LogObjectAsJson("codeA3", codeA3)
	spans.LogObjectAsJson("phoneCode", phoneCode)
	spans.LogObjectAsJson("createdAt", createdAt)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *countryWriteRepository) UpdateCountry(ctx context.Context, id, name, codeA2, codeA3, phoneCode string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CountryReadRepository.GetDefaultCountryCodeA3")
	defer spans.Finish()
	spans.LogObjectAsJson("id", id)
	spans.LogObjectAsJson("name", name)
	spans.LogObjectAsJson("codeA2", codeA2)
	spans.LogObjectAsJson("codeA3", codeA3)
	spans.LogObjectAsJson("phoneCode", phoneCode)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}

	return err
}
