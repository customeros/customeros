package neo4j_repository

import (
	"fmt"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type PhoneNumberReadRepository interface {
	GetPhoneNumberIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
	GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error)
	GetById(ctx context.Context, tenant, phoneNumberId string) (*dbtype.Node, error)
	GetAllForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	Exists(ctx context.Context, tenant string, e164 string) (bool, error)
	GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error)
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.GetPhoneNumberIdIfExists")
	defer spans.Finish()

	spans.LogKV("phoneNumber", phoneNumber)

	cypher := fmt.Sprintf(`MATCH (p:PhoneNumber_%s) WHERE p.e164 = $phoneNumber OR p.rawPhoneNumber = $phoneNumber RETURN p.id LIMIT 1`, tenant)
	params := map[string]any{
		"phoneNumber": phoneNumber,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return "", err
	}
	if len(result.([]*db.Record)) == 0 {
		spans.LogKV("result", "")
		return "", nil
	}
	spans.LogKV("result", result.([]*db.Record)[0].Values[0].(string))
	return result.([]*db.Record)[0].Values[0].(string), err
}

func (r *phoneNumberReadRepository) GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.GetCountryCodeA2ForPhoneNumber")
	defer spans.Finish()

	spans.TagEntity(phoneNumberId)

	cypher := `MATCH (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
				OPTIONAL MATCH (p)-[:LINKED_TO]->(c:Country)
				OPTIONAL MATCH (tenant)-[:DEFAULT_COUNTRY]->(dc:Country)
				RETURN COALESCE(c.codeA2, dc.codeA2, '') AS countryCodeA2 LIMIT 1`
	params := map[string]any{
		"tenant":        tenant,
		"phoneNumberId": phoneNumberId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	spans.LogKV("result", result.(string))
	return result.(string), nil
}

func (r *phoneNumberReadRepository) GetById(ctx context.Context, tenant, phoneNumberId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("phoneNumberId", phoneNumberId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId}) return p`
	params := map[string]any{
		"tenant":        tenant,
		"phoneNumberId": phoneNumberId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *phoneNumberReadRepository) GetAllForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.GetAllForLinkedEntityIds")
	defer spans.Finish()

	cypher := ""
	switch entityType {
	case neo4jenum.CONTACT:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(entity:Contact)`
	case neo4jenum.USER:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(entity:User)`
	case neo4jenum.ORGANIZATION:
		cypher = `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(entity:Organization)`
	}
	cypher = cypher + `, (entity)-[rel:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
					WHERE entity.id IN $entityIds
					RETURN p, rel, entity.id ORDER BY p.e164, p.rawPhoneNumber`
	params := map[string]any{
		"tenant":    tenant,
		"entityIds": entityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *phoneNumberReadRepository) Exists(ctx context.Context, tenant string, e164 string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.Exists")
	defer spans.Finish()

	cypher := fmt.Sprintf("MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1", tenant)
	params := map[string]any{
		"e164": e164,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil

		}
	})
	if err != nil {
		return false, err
	}
	return result.(bool), err
}

func (r *phoneNumberReadRepository) GetByPhoneNumber(ctx context.Context, tenant, e164 string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberReadRepository.GetByPhoneNumber")
	defer spans.Finish()

	cypher := fmt.Sprintf("MATCH (p:PhoneNumber_%s) WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 RETURN p LIMIT 1", tenant)
	params := map[string]any{
		"e164": e164,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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
