package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type TenantAndContactDetails struct {
	Tenant         string
	ContractId     string
	OrganizationId string
}

type TenantAndContactId struct {
	Tenant     string
	ContractId string
}

type ContactReadRepository interface {
	GetContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetContactInOrganizationByEmail(ctx context.Context, tenant, organizationId, email string) (*neo4j.Node, error)
	GetContactCountByOrganizations(ctx context.Context, tenant string, ids []string) (map[string]int64, error)
	GetContactsToFindEmail(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]TenantAndContactDetails, error)
	GetContactsToEnrichByEmail(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, limit int) ([]TenantAndContactId, error)
	GetLinkedOrgDomains(ctx context.Context, tenant, contactId string) ([]string, error)
}

type contactReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *contactReadRepository) GetLinkedOrgDomains(ctx context.Context, tenant, contactId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetLinkedOrgDomains")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id})--(:JobRole)--(o:Organization)--(d:Domain)
				RETURN DISTINCT d.domain`
	params := map[string]any{
		"tenant": tenant,
		"id":     contactId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range result.([]*neo4j.Record) {
		output = append(output, v.Values[0].(string))
	}
	span.LogFields(log.Object("result", output))
	return output, nil
}

func NewContactReadRepository(driver *neo4j.DriverWithContext, database string) ContactReadRepository {
	return &contactReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *contactReadRepository) GetContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contactId,
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
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *contactReadRepository) GetContactInOrganizationByEmail(ctx context.Context, tenant, organizationId, email string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))
	span.LogFields(log.String("email", email))

	cypher := `match (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization{id:$organizationId})<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e:Email{rawEmail:$email})
		return c`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"email":          email,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *contactReadRepository) GetContactById(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contactId)
	span.LogFields(log.String("contactId", contactId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contactId,
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
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *contactReadRepository) GetContactCountByOrganizations(ctx context.Context, tenant string, ids []string) (map[string]int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactCountByOrganizations")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WHERE o.id IN $ids
				WITH o
				OPTIONAL MATCH (o)--(:JobRole)--(c:Contact)
				RETURN o.id, count(c) as count`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		output[v.Values[0].(string)] = v.Values[1].(int64)
	}
	return output, err
}

func (r *contactReadRepository) GetContactsToFindEmail(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]TenantAndContactDetails, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsToFindEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("minutesFromLastContactUpdate", minutesFromLastContactUpdate))
	span.LogFields(log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(j:JobRole)--(o:Organization),
				(t)--(ts:TenantSettings)
				WHERE
					ts.enrichContacts = true AND
					NOT (c)-[:HAS]->(:Email) AND
					(c.firstName IS NOT NULL AND c.lastName IS NOT NULL AND c.firstName <> '' AND c.lastName <> '') AND
					(c.techFindEmailRequestedAt IS NULL) AND
					c.updatedAt < datetime() - duration({minutes: $minutesFromLastContactUpdate})
				WITH t, c, collect(o.id) as organizationIds
				RETURN t.name as tenant, c.id as contactId, head(organizationIds) as organizationId
				LIMIT $limit`
	params := map[string]any{
		"minutesFromLastContactUpdate": minutesFromLastContactUpdate,
		"limit":                        limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndContactDetails, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactDetails{
				Tenant:         v.Values[0].(string),
				ContractId:     v.Values[1].(string),
				OrganizationId: v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contactReadRepository) GetContactsToEnrichByEmail(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, limit int) ([]TenantAndContactId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsToEnrichByEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("minutesFromLastContactUpdate", minutesFromLastContactUpdate))
	span.LogFields(log.Int("minutesFromLastEnrichAttempt", minutesFromLastEnrichAttempt))
	span.LogFields(log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(j:JobRole)--(o:Organization),
				(t)--(ts:TenantSettings)
				WHERE
					ts.enrichContacts = true AND
					(c)-[:HAS]->(:Email) AND
					o.relationship IN $allowedOrgRelationships AND
					NOT o.stage IN $restrictedOrgStages AND
					(c.updatedAt < datetime() - duration({minutes: $minutesFromLastContactUpdate})) AND
					(c.techEnrichRequestedAt IS NULL OR c.techEnrichRequestedAt < datetime() - duration({minutes: $minutesFromLastEnrichAttempt})) AND
					c.enrichedAt IS NULL
				WITH t.name as tenant, c.id as contactId
				ORDER BY CASE WHEN c.techEnrichRequestedAt IS NULL THEN 0 ELSE 1 END, c.techEnrichRequestedAt ASC
				LIMIT $limit
				RETURN DISTINCT tenant, contactId`
	params := map[string]any{
		"minutesFromLastContactUpdate": minutesFromLastContactUpdate,
		"minutesFromLastEnrichAttempt": minutesFromLastEnrichAttempt,
		"allowedOrgRelationships":      []string{enum.Customer.String(), enum.Prospect.String()},
		"restrictedOrgStages":          []string{enum.Lead.String()},
		"limit":                        limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndContactId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactId{
				Tenant:     v.Values[0].(string),
				ContractId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
