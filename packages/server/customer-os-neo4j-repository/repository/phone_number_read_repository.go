package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type PhoneNumberReadRepository interface {
	GetPhoneNumberIdIfExists(ctx context.Context, tenant, phoneNumber string) (string, error)
	GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error)
	GetById(ctx context.Context, tenant, phoneNumberId string) (*dbtype.Node, error)
	GetAllForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *phoneNumberReadRepository) GetCountryCodeA2ForPhoneNumber(ctx context.Context, tenant, phoneNumberId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberReadRepository.GetCountryCodeA2ForPhoneNumber")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, phoneNumberId)

	cypher := `MATCH (p:PhoneNumber {id:$phoneNumberId})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
				OPTIONAL MATCH (p)-[:LINKED_TO]->(c:Country)
				OPTIONAL MATCH (tenant)-[:DEFAULT_COUNTRY]->(dc:Country)
				RETURN COALESCE(c.codeA2, dc.codeA2, '') AS countryCodeA2 LIMIT 1`
	params := map[string]any{
		"tenant":        tenant,
		"phoneNumberId": phoneNumberId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return "", err
	}
	span.LogFields(log.String("result", result.(string)))
	return result.(string), nil
}

func (r *phoneNumberReadRepository) GetById(ctx context.Context, tenant, phoneNumberId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberReadRepository.GetById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId}) return p`
	params := map[string]any{
		"tenant":        tenant,
		"phoneNumberId": phoneNumberId,
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

func (r *phoneNumberReadRepository) GetAllForLinkedEntityIds(ctx context.Context, tenant string, entityType neo4jenum.EntityType, entityIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberReadRepository.GetAllForLinkedEntityIds")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
