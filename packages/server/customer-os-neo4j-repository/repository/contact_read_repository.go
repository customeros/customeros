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

type ContactsEnrichWorkEmail struct {
	Tenant             string
	ContactId          string
	ContactFirstName   string
	ContactLastName    string
	LinkedInUrl        string
	OrganizationId     string
	OrganizationName   string
	OrganizationDomain string
}

type TenantAndContactId struct {
	Tenant    string
	ContactId string
}

type ContactIdWithRequestId struct {
	Tenant    string
	ContactId string
	RequestId string
}

type ContactsEnrichedNotLinkedToOrganization struct {
	Tenant      string
	ContactId   string
	LinkedInUrl string
}

type ContactReadRepository interface {
	GetContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetContactsEnrichedNotLinkedToOrganization(ctx context.Context) ([]ContactsEnrichedNotLinkedToOrganization, error)
	GetContactsWithSocialUrl(ctx context.Context, tenant, socialUrl string) ([]*dbtype.Node, error)
	GetContactsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error)
	GetContactInOrganizationByEmail(ctx context.Context, tenant, organizationId, email string) (*neo4j.Node, error)
	GetContactCountByOrganizations(ctx context.Context, tenant string, ids []string) (map[string]int64, error)
	GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]ContactsEnrichWorkEmail, error)
	GetContactsToEnrichWithEmailFromBetterContact(ctx context.Context, limit int) ([]ContactIdWithRequestId, error)
	GetContactsToEnrichByEmail(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, minutesFromLastFailure, limit int) ([]TenantAndContactId, error)
	GetLinkedOrgDomains(ctx context.Context, tenant, contactId string) ([]string, error)
	GetContactsWithGroupEmail(ctx context.Context, limit int) ([]TenantAndContactId, error)
}

type contactReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContactReadRepository(driver *neo4j.DriverWithContext, database string) ContactReadRepository {
	return &contactReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contactReadRepository) GetContactsEnrichedNotLinkedToOrganization(ctx context.Context) ([]ContactsEnrichedNotLinkedToOrganization, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsEnrichedNotLinkedToOrganization")
	defer span.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(s:Social) 
			WHERE c.enrichedAt is not null and s.url =~ '.*linkedin.com.*' and not (c)--(:JobRole)
			RETURN DISTINCT t.name, c.id, s.url`,
			map[string]interface{}{}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make([]ContactsEnrichedNotLinkedToOrganization, 0)
	for _, v := range result.([]*neo4j.Record) {
		output = append(output,
			ContactsEnrichedNotLinkedToOrganization{
				Tenant:      v.Values[0].(string),
				ContactId:   v.Values[1].(string),
				LinkedInUrl: v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contactReadRepository) GetContactsWithSocialUrl(ctx context.Context, tenant, socialUrl string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsWithSocialUrl")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(s:Social) 
			WHERE s.url=$socialUrl
			RETURN DISTINCT c`,
			map[string]interface{}{
				"socialUrl": socialUrl,
				"tenant":    tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), err
}

func (r *contactReadRepository) GetContactsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsWithEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT c`,
			map[string]interface{}{
				"email":  email,
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), err
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
				OPTIONAL MATCH (o)--(:JobRole)--(c:Contact) WHERE c.hide IS NULL OR c.hide = false
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

func (r *contactReadRepository) GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]ContactsEnrichWorkEmail, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsToFindWorkEmailWithBetterContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("minutesFromLastContactUpdate", minutesFromLastContactUpdate))
	span.LogFields(log.Int("limit", limit))

	cypher := ` MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(j:JobRole)--(o:Organization)--(d:Domain), (t)--(ts:TenantSettings)
				WHERE
					ts.enrichContacts = true AND
					NOT (c)-[:HAS]->(:Email) AND
					(c.firstName IS NOT NULL AND c.lastName IS NOT NULL AND c.firstName <> '' AND c.lastName <> '') AND
					(c.techFindWorkEmailWithBetterContactRequestedAt IS NULL) AND
					c.updatedAt < datetime() - duration({minutes: 2})
				WITH t, c, o, d
				OPTIONAL MATCH (c)-[:HAS]->(s:Social)
				where s is null or s.url =~ '.*linkedin.com.*'
				RETURN t.name, c.id, c.firstName, c.lastName, CASE WHEN s is null THEN '' else s.url END, o.id, o.name, d.domain
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
	output := make([]ContactsEnrichWorkEmail, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			ContactsEnrichWorkEmail{
				Tenant:             v.Values[0].(string),
				ContactId:          v.Values[1].(string),
				ContactFirstName:   v.Values[2].(string),
				ContactLastName:    v.Values[3].(string),
				LinkedInUrl:        v.Values[4].(string),
				OrganizationId:     v.Values[5].(string),
				OrganizationName:   v.Values[6].(string),
				OrganizationDomain: v.Values[7].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contactReadRepository) GetContactsToEnrichWithEmailFromBetterContact(ctx context.Context, limit int) ([]ContactIdWithRequestId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsToEnrichWithEmailFromBetterContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("limit", limit))

	cypher := ` MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
				WHERE
					c.techFindWorkEmailWithBetterContactRequestId IS NOT NULL AND
					c.techFindWorkEmailWithBetterContactRequestedAt IS NOT NULL AND 
					c.techFindWorkEmailWithBetterContactCompletedAt is null AND
					c.techFindWorkEmailWithBetterContactRequestedAt < datetime() - duration({minutes: 2})
				RETURN t.name, c.id, c.techFindWorkEmailWithBetterContactRequestId ORDER BY c.techFindWorkEmailWithBetterContactRequestedAt asc
				LIMIT $limit`
	params := map[string]any{
		"limit": limit,
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
	output := make([]ContactIdWithRequestId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			ContactIdWithRequestId{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
				RequestId: v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contactReadRepository) GetContactsToEnrichByEmail(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, minutesFromLastFailure, limit int) ([]TenantAndContactId, error) {
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
					c.enrichedAt IS NULL AND
					(c.enrichedFailedAtScrapInPersonSearch IS NULL OR c.enrichedFailedAtScrapInPersonSearch < datetime() - duration({minutes: $minutesFromLastFailure})) AND
					o.relationship IN $allowedOrgRelationships AND
					NOT o.stage IN $restrictedOrgStages AND
					(c.updatedAt < datetime() - duration({minutes: $minutesFromLastContactUpdate})) AND
					(c.techEnrichRequestedAt IS NULL OR c.techEnrichRequestedAt < datetime() - duration({minutes: $minutesFromLastEnrichAttempt}))
				WITH t.name as tenant, c.id as contactId
				ORDER BY CASE WHEN c.techEnrichRequestedAt IS NULL THEN 0 ELSE 1 END, c.techEnrichRequestedAt ASC
				LIMIT $limit
				RETURN DISTINCT tenant, contactId`
	params := map[string]any{
		"minutesFromLastContactUpdate": minutesFromLastContactUpdate,
		"minutesFromLastEnrichAttempt": minutesFromLastEnrichAttempt,
		"minutesFromLastFailure":       minutesFromLastFailure,
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
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contactReadRepository) GetContactsWithGroupEmail(ctx context.Context, limit int) ([]TenantAndContactId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetContactsWithGroupEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
				WHERE
					(c.hide IS NULL OR c.hide = false) AND
					e.isRoleAccount = true
				RETURN DISTINCT t.name, c.id LIMIT $limit`
	params := map[string]any{
		"limit": limit,
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
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
