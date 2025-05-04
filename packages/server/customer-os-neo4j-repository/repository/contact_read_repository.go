package neo4j_repository

import (
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

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

type TenantAndContactIdAndParams struct {
	Tenant    string
	ContactId string
	FieldStr1 string
}

type TenantAndContact struct {
	Tenant  string
	Contact *dbtype.Node
}

type ContactReadRepository interface {
	CountByTenant(ctx context.Context, tenant string) (int64, error)
	GetContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetContacts(ctx context.Context, tenant string, contactIds []string) ([]*dbtype.Node, error)
	GetContactsWithSocialUrl(ctx context.Context, tenant, socialUrl string) ([]*dbtype.Node, error)
	GetContactsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error)
	GetContactsByEmailAddresses(ctx context.Context, tenant string, emailAddresses []string) ([]*utils.DbNodeAndId, error)
	GetContactInOrganizationByEmail(ctx context.Context, tenant, organizationId, email string) (*neo4j.Node, error)
	GetActiveContactsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetContactCountByOrganizations(ctx context.Context, tenant string, ids []string) (map[string]int64, error)
	GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]ContactsEnrichWorkEmail, error)
	GetContactsToCheck(ctx context.Context, minutesSinceLastUpdate, hoursSinceLastCheck, limit int) ([]TenantAndContact, error)
	GetContactsByLinkedIn(ctx context.Context, tenant, url, alias, externalId string) ([]*dbtype.Node, error)
	GetDistinctContactRegions(ctx context.Context, tenant string) ([]string, error)
	GetDistinctContactCities(ctx context.Context, tenant string) ([]string, error)

	// cross tenant queries
	GetContactsToSetPrimaryJobRole(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsToEnrichWithEmailFromBetterContact(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsToEnrich(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsWithGroupOrSystemGeneratedEmail(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsWithEmailForNameUpdate(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsEnrichedNotLinkedToOrganization(ctx context.Context, delayFromPreviousAttemptDays, limit int) ([]TenantAndContactIdAndParams, error)
	GetContactsWithProfilePhotoUrlCrossTenant(ctx context.Context, profilePhotoUrl string) ([]TenantAndContactIdAndParams, error)
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

func (r *contactReadRepository) GetContactsEnrichedNotLinkedToOrganization(ctx context.Context, delayFromPreviousAttemptDays, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsEnrichedNotLinkedToOrganization")
	defer spans.Finish()

	spans.LogKV("delayFromPreviousAttemptDays", delayFromPreviousAttemptDays)
	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(s:Social)
			WHERE
				(c.hide IS NULL OR c.hide = false) AND
				(c.techLinkWithOrgRequestedAt IS NULL OR c.techLinkWithOrgRequestedAt < datetime() - duration({days: $delayDays})) AND
				c.enrichedAt IS NOT NULL AND
				s.url =~ '.*linkedin.com.*' AND 
				NOT (c)--(:JobRole)
			RETURN DISTINCT t.name, c.id, s.url limit $limit`
	params := map[string]interface{}{
		"limit":     limit,
		"delayDays": delayFromPreviousAttemptDays,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range result.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
				FieldStr1: v.Values[2].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsWithSocialUrl(ctx context.Context, tenant, socialUrl string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsWithSocialUrl")
	defer spans.Finish()

	spans.LogKV("socialUrl", socialUrl)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(s:Social) 
			WHERE s.url=$socialUrl
			RETURN DISTINCT c`
	params := map[string]any{
		"socialUrl": socialUrl,
		"tenant":    tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsWithEmail")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT c ORDER BY c.createdAt ASC`
	params := map[string]any{
		"tenant": tenant,
		"email":  email,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *contactReadRepository) GetContactsByEmailAddresses(ctx context.Context, tenant string, emailAddresses []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsByEmailAddresses")
	defer spans.Finish()

	spans.LogObjectAsJson("emailAddresses", emailAddresses)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE toLower(e.email)IN $emailAddresses OR toLower(e.rawEmail) IN $emailAddresses 
			RETURN c, COALESCE(CASE WHEN e.email IS NOT NULL AND e.email <> '' THEN e.email ELSE e.rawEmail END, e.email) as emailId`
	params := map[string]any{
		"tenant":         tenant,
		"emailAddresses": emailAddresses,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *contactReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *contactReadRepository) GetContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContact")
	defer spans.Finish()

	spans.TagEntity(contactId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contactId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *contactReadRepository) GetContacts(ctx context.Context, tenant string, contactIds []string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContacts")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact) WHERE c.id IN $ids RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"ids":    contactIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.LogKV("result.count", 0)
		return nil, err
	}
	nodes := result.([]*dbtype.Node)
	spans.LogKV("result.count", len(nodes))
	return nodes, err
}

func (r *contactReadRepository) GetContactInOrganizationByEmail(ctx context.Context, tenant, organizationId, email string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactById")
	defer spans.Finish()

	spans.LogKV("organizationId", organizationId)
	spans.LogKV("email", email)

	cypher := `match (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization{id:$organizationId})<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e:Email{rawEmail:$email})
		return c`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"email":          email,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactById")
	defer spans.Finish()

	spans.TagEntity(contactId)
	spans.LogKV("contactId", contactId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contactId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *contactReadRepository) GetActiveContactsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetActiveContactsForOrganizations")
	defer spans.Finish()

	spans.LogObjectAsJson("organizationIds", organizationIds)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)
								WHERE o.id IN $organizationIds
								MATCH (o)--(j:JobRole_%s)--(c:Contact_%s) 
								WHERE c.hide IS NULL OR c.hide = false
								RETURN c, o.id`, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *contactReadRepository) GetContactCountByOrganizations(ctx context.Context, tenant string, ids []string) (map[string]int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactCountByOrganizations")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WHERE o.id IN $ids
				WITH o
				OPTIONAL MATCH (o)--(:JobRole)--(c:Contact) WHERE c.hide IS NULL OR c.hide = false
				RETURN o.id, count(c) as count`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		output[v.Values[0].(string)] = v.Values[1].(int64)
	}
	return output, err
}

func (r *contactReadRepository) GetContactsToFindWorkEmailWithBetterContact(ctx context.Context, minutesFromLastContactUpdate, limit int) ([]ContactsEnrichWorkEmail, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsToFindWorkEmailWithBetterContact")
	defer spans.Finish()

	spans.LogKV("minutesFromLastContactUpdate", minutesFromLastContactUpdate)
	spans.LogKV("limit", limit)

	cypher := ` MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(j:JobRole)--(o:Organization)--(d:Domain), (t)--(ts:TenantSettings)
				WHERE
					ts.enrichContacts = true AND
					not o.stage = 'LEAD' and not o.stage = 'UNQUALIFIED' AND
					NOT (c)-[:HAS]->(:Email) AND
					(c.firstName IS NOT NULL AND c.lastName IS NOT NULL AND c.firstName <> '' AND c.lastName <> '') AND
					(c.techFindWorkEmailWithBetterContactRequestedAt IS NULL) AND
					c.updatedAt < datetime() - duration({minutes: $minutesFromLastContactUpdate})
				WITH t, c, o, d
				OPTIONAL MATCH (c)-[:HAS]->(s:Social)
				WHERE s IS NULL OR s.url =~ '.*linkedin.com.*'
				RETURN t.name, c.id, c.firstName, c.lastName, CASE WHEN s is null THEN '' else s.url END, o.id, o.name, d.domain
				ORDER BY c.createdAt ASC
				LIMIT $limit`
	params := map[string]any{
		"minutesFromLastContactUpdate": minutesFromLastContactUpdate,
		"limit":                        limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsToEnrichWithEmailFromBetterContact(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsToEnrichWithEmailFromBetterContact")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	minutesDelayFromUpdate := 2
	cypher := ` MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
				WHERE
					c.techFindWorkEmailWithBetterContactRequestId IS NOT NULL AND
					c.techFindWorkEmailWithBetterContactRequestId <> '' AND
					c.techFindWorkEmailWithBetterContactRequestedAt IS NOT NULL AND 
					c.techFindWorkEmailWithBetterContactCompletedAt is null AND
					c.techFindWorkEmailWithBetterContactRequestedAt < datetime() - duration({minutes: $minutesDelay})
				RETURN t.name, c.id, c.techFindWorkEmailWithBetterContactRequestId 
				ORDER BY
					CASE WHEN c.techUpdateWithWorkEmailRequestedAt IS NULL THEN 0 ELSE 1 END ASC, 
					c.techUpdateWithWorkEmailRequestedAt ASC
				LIMIT $limit`
	params := map[string]any{
		"limit":        limit,
		"minutesDelay": minutesDelayFromUpdate,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
				FieldStr1: v.Values[2].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsToEnrich(ctx context.Context, minutesFromLastContactUpdate, minutesFromLastEnrichAttempt, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsToEnrich")
	defer spans.Finish()

	spans.LogKV("minutesFromLastContactUpdate", minutesFromLastContactUpdate)
	spans.LogKV("minutesFromLastEnrichAttempt", minutesFromLastEnrichAttempt)
	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact),
				(t)--(ts:TenantSettings)
				WHERE 
					ts.enrichContacts = true AND
					c.enrichedAt IS NULL AND
					(c.techEnrichAttempts IS NULL OR c.techEnrichAttempts < $maxAttempts) AND
					(c.updatedAt < datetime() - duration({minutes: $minutesFromLastContactUpdate})) AND
					(c.techEnrichRequestedAt IS NULL OR c.techEnrichRequestedAt < datetime() - duration({minutes: $minutesFromLastEnrichAttempt}))
				WITH t, c
				OPTIONAL MATCH (c)-[:HAS]->(e:Email)
				WHERE
    				e.isRoleAccount = false AND
    				e.isSystemGenerated = false
				OPTIONAL MATCH (c)--(:JobRole)--(o:Organization)--(:Domain)
				OPTIONAL MATCH (c)--(s:Social)
				WHERE s.url CONTAINS 'linkedin.com'
				WITH t.name AS tenant, c.id AS contactId, c, e, o, s
				WHERE e IS NOT NULL OR o IS NOT NULL OR s IS NOT NULL
				ORDER BY CASE WHEN c.techEnrichRequestedAt IS NULL THEN 0 ELSE 1 END, c.techEnrichRequestedAt ASC
				RETURN DISTINCT tenant, contactId
				LIMIT $limit`
	params := map[string]any{
		"minutesFromLastContactUpdate": minutesFromLastContactUpdate,
		"minutesFromLastEnrichAttempt": minutesFromLastEnrichAttempt,
		"limit":                        limit,
		"maxAttempts":                  1,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsWithGroupOrSystemGeneratedEmail(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsWithGroupOrSystemGeneratedEmail")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
				WHERE
					(c.hide IS NULL OR c.hide = false) AND
					(e.isRoleAccount = true OR e.isSystemGenerated = true) 
				RETURN DISTINCT t.name, c.id LIMIT $limit`
	params := map[string]any{
		"limit": limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsWithEmailForNameUpdate(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsWithEmailForNameUpdate")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(e:Email)
				WHERE
					(c.hide IS NULL OR c.hide = false) AND
					(c.firstName IS NULL OR c.firstName = '') AND
					(c.lastName IS NULL OR c.lastName = '') AND
					e.email IS NOT NULL AND 
					e.email <> '' AND 
					c.updatedAt < datetime() - duration({minutes: 3})
				RETURN DISTINCT t.name, c.id, e.email LIMIT $limit`
	params := map[string]any{
		"limit": limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
				FieldStr1: v.Values[2].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsToCheck(ctx context.Context, minutesFromLastUpdate, hoursFromLastCheck, limit int) ([]TenantAndContact, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsToCheck")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("minutesFromLastUpdate", minutesFromLastUpdate)
	spans.LogKV("hoursFromLastCheck", hoursFromLastCheck)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
				WHERE
					(c.hide IS NULL OR c.hide = false) AND
					(c.techCheckedAt IS NULL OR c.checkedAt < datetime() - duration({hours: $hoursFromLastCheck})) AND
					c.updatedAt < datetime() - duration({minutes: $minutesFromLastUpdate})
					ORDER BY CASE WHEN c.techCheckedAt IS NULL THEN 0 ELSE 1 END, c.techCheckedAt ASC
				RETURN t.name, c LIMIT $limit`
	params := map[string]any{
		"limit":                 limit,
		"minutesFromLastUpdate": minutesFromLastUpdate,
		"hoursFromLastCheck":    hoursFromLastCheck,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContact, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContact{
				Tenant:  v.Values[0].(string),
				Contact: utils.ToPtr(v.Values[1].(dbtype.Node)),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetContactsByLinkedIn(ctx context.Context, tenant, url, alias, externalId string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetContactsByLinkedIn")
	defer spans.Finish()

	spans.LogKV("url", url)
	spans.LogKV("alias", alias)
	spans.LogKV("externalId", externalId)

	if !strings.Contains(url, "linkedin.com") {
		spans.LogKV("result.count", 0)
		return nil, nil
	}

	minimizedUrl := url
	// remove trailing / if any
	if minimizedUrl[len(minimizedUrl)-1] == '/' {
		minimizedUrl = minimizedUrl[:len(minimizedUrl)-1]
	}
	// remove all chars before linkedin.com
	if i := strings.Index(minimizedUrl, "linkedin.com/in"); i != -1 {
		minimizedUrl = minimizedUrl[i:]
	}
	minimizedUrlWithSlash := minimizedUrl + "/"

	cypher := `MATCH (:Tenant {name:$tenant})--(c:Contact)-[:HAS]->(s:Social)
				WHERE s.url ENDS WITH $url OR s.url ENDS WITH $urlWithSlash `
	if alias != "" {
		cypher += ` OR s.alias = $alias `
	}
	if externalId != "" {
		cypher += ` OR s.externalId = $externalId `
	}
	cypher += ` RETURN c`
	params := map[string]any{
		"tenant":       tenant,
		"url":          minimizedUrl,
		"urlWithSlash": minimizedUrlWithSlash,
		"alias":        alias,
		"externalId":   externalId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.LogKV("result.count", 0)
		return nil, err
	}
	nodes := result.([]*dbtype.Node)
	spans.LogKV("result.count", len(nodes))
	return nodes, err
}

func (r *contactReadRepository) CountByTenant(ctx context.Context, tenant string) (int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.CountByTenant")
	defer spans.Finish()

	cypher := `MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) WHERE c.hide = false or c.hide IS NULL
			RETURN count(c)`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	organizationsCount := dbRecord.(*db.Record).Values[0].(int64)
	spans.LogKV("result", organizationsCount)
	return organizationsCount, nil
}

func (r *contactReadRepository) GetContactsToSetPrimaryJobRole(ctx context.Context, limit int) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsToSetPrimaryJobRole")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(:Organization)
				WHERE c.hide IS NULL OR c.hide = false
				WITH t, c, sum(CASE WHEN j.primary = true THEN 1 ELSE 0 END) AS primaryCount
				WHERE primaryCount <> 1
				RETURN t.name, c.id ORDER BY c.updatedAt DESC LIMIT $limit`
	params := map[string]any{
		"limit": limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *contactReadRepository) GetDistinctContactRegions(ctx context.Context, tenant string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetDistinctContactRegions")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {active:true, name: $tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)--(l:Location) 
			WHERE c.hide = false AND l.region IS NOT NULL AND l.region <> '' RETURN DISTINCT l.region`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.LogKV("result.count", 0)
		return nil, err
	}

	return result.([]string), err
}

func (r *contactReadRepository) GetDistinctContactCities(ctx context.Context, tenant string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetDistinctContactCities")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (t:Tenant {active:true, name: $tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact_%s)--(l:Location) 
				WHERE c.hide = false and l.locality IS NOT NULL AND l.locality <> '' RETURN distinct l.locality`, tenant)
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.LogKV("result.count", 0)
		return nil, err
	}

	return result.([]string), err
}

func (r *contactReadRepository) GetContactsWithProfilePhotoUrlCrossTenant(ctx context.Context, profilePhotoUrl string) ([]TenantAndContactIdAndParams, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetContactsWithProfilePhotoUrlCrossTenant")
	defer spans.Finish()

	spans.LogKV("profilePhotoUrl", profilePhotoUrl)

	cypher := `MATCH (t:Tenant {active:true})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)
				WHERE c.profilePhotoUrl = $profilePhotoUrl
				RETURN t.name, c.id`
	params := map[string]any{
		"profilePhotoUrl": profilePhotoUrl,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	output := make([]TenantAndContactIdAndParams, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContactIdAndParams{
				Tenant:    v.Values[0].(string),
				ContactId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}
