package neo4j_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
)

const (
	Relationship_Subsidiary = "SUBSIDIARY_OF"
)

type TenantAndOrganizationId struct {
	Tenant         string
	OrganizationId string
}

type TenantAndOrganizationIdExtended struct {
	Tenant         string
	OrganizationId string
	Param1         string
}

type TenantAndOrganization struct {
	Tenant       string
	Organization *dbtype.Node
}

type OrganizationReadRepository interface {
	CountByTenant(ctx context.Context, tenant string) (int64, error)
	GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error)
	GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	GetOrganizationByContactId(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOrganizationByInvoiceId(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error)
	GetOrganizationByCustomerOsId(ctx context.Context, tenant, customerOsId string) (*dbtype.Node, error)
	GetOrganizationByReferenceId(ctx context.Context, tenant, referenceId string) (*dbtype.Node, error)
	GetOrganizationByIdOrCustomerOsId(ctx context.Context, tenant, id string) (*dbtype.Node, error)
	GetOrganizationByDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, domain string) (*dbtype.Node, error)
	GetOrganizationBySocialUrl(ctx context.Context, tenant, socialUrl string) (*dbtype.Node, error)
	GetOrganizationsByLinkedIn(ctx context.Context, tenant, url, alias, externalId string) ([]*dbtype.Node, error)
	GetForApiCache(ctx context.Context, tenant string, skip, limit int) ([]map[string]interface{}, error)
	GetPatchesForApiCache(ctx context.Context, tenant string, lastPatchTimestamp time.Time) ([]map[string]interface{}, error)
	GetAllForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error)
	GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error)
	GetOrganizationsForEnrichByDomain(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error)
	GetOrganizationsForUpdateLastTouchpoint(ctx context.Context, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error)
	GetOrganizationsForIcpCheck(ctx context.Context, tenants []string, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error)
	GetPrimaryOrganizationsWithJobRoleForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodePairAndId, error)
	GetHiddenOrganizationIds(ctx context.Context, tenant string, hiddenAfter time.Time) ([]string, error)
	GetMergedOrganizationIds(ctx context.Context, tenant string, mergedAfter time.Time) ([]string, error)
	GetOrganizationsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error)
	GetActiveOrganizationIdsByDomain(ctx context.Context, tenant string, domains []string) (map[string]string, error)
	GetLinkedSubOrganizations(ctx context.Context, tenant string, parentOrganizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	GetLinkedParentOrganizations(ctx context.Context, tenant string, organizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	GetOrganizationsByDomainAcrossAllTenants(ctx context.Context, domain string) ([]TenantAndOrganizationId, error)
	GetOrganizationsByStage(ctx context.Context, tenant string, stage commonenum.OrganizationStage) ([]*dbtype.Node, error)
}

type organizationReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationReadRepository(driver *neo4j.DriverWithContext, database string) OrganizationReadRepository {
	return &organizationReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *organizationReadRepository) CountByTenant(ctx context.Context, tenant string) (int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.CountByTenant")
	defer spans.Finish()

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) where org.hide = false
			RETURN count(org)`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("query", cypher)
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

func (r *organizationReadRepository) GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganization")
	defer spans.Finish()

	spans.TagEntity(organizationId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationIdsConnectedToInteractionEvent")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}),
				(t:Tenant {name:$tenant})
				CALL {
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:REPORTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				UNION 
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:SUBMITTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				}
				RETURN distinct orgId`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]string)))
	return result.([]string), err
}

func (r *organizationReadRepository) GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByOpportunityId")
	defer spans.Finish()

	spans.LogKV("opportunityId", opportunityId)

	cypher := `MATCH (op:Opportunity {id:$id})
				MATCH (t:Tenant {name:$tenant})
				OPTIONAL MATCH (op)<-[:HAS_OPPORTUNITY]-(:Contract)<-[:HAS_CONTRACT]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				OPTIONAL MATCH (op)<-[:HAS_OPPORTUNITY]-(directOrg:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			WITH COALESCE(org, directOrg) as organization 
			WHERE organization IS NOT NULL RETURN organization`
	params := map[string]any{
		"tenant": tenant,
		"id":     opportunityId,
	}
	spans.LogKV("query", cypher)
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
		spans.TraceError(err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	} else {
		spans.LogKV("result.found", true)
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByContactId(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByContactId")
	defer spans.Finish()

	spans.LogKV("contactId", contactId)

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), 
				(t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				WHERE (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(org) 
			RETURN org limit 1`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	spans.LogKV("query", cypher)
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
		spans.TraceError(err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	} else {
		spans.LogKV("result.found", true)
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByContractId")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_CONTRACT]->(c:Contract {id:$id})
			RETURN org limit 1`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
	}
	spans.LogKV("query", cypher)
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
		spans.TraceError(err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	} else {
		spans.LogKV("result.found", true)
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByInvoiceId(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByInvoiceId")
	defer spans.Finish()

	spans.LogKV("invoiceId", invoiceId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(inv:Invoice {id:$invoiceId})<-[:HAS_INVOICE]-(c:Contract)<-[:HAS_CONTRACT]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			RETURN org`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
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
		spans.TraceError(err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	} else {
		spans.LogKV("result.found", true)
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByCustomerOsId(ctx context.Context, tenant, customerOsId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByCustomerOsId")
	defer spans.Finish()

	spans.LogKV("customerOsId", customerOsId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {customerOsId:$customerOsId})
			RETURN org`
	params := map[string]any{
		"tenant":       tenant,
		"customerOsId": customerOsId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByIdOrCustomerOsId(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId")
	defer spans.Finish()

	spans.LogKV("id", id)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
			WHERE org.customerOsId = $id OR org.id = $id RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByReferenceId(ctx context.Context, tenant, referenceId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByReferenceId")
	defer spans.Finish()

	spans.LogKV("referenceId", referenceId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {referenceId:$referenceId}) RETURN org`
	params := map[string]any{
		"tenant":      tenant,
		"referenceId": referenceId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByDomain(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, domain string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationByDomain")
	defer spans.Finish()

	spans.LogKV("domain", domain)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain{domain:$domain}) RETURN o limit 1`
	params := map[string]any{
		"tenant": tenant,
		"domain": domain,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if len(result.([]*dbtype.Node)) == 0 {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return result.([]*dbtype.Node)[0], err
}

func (r *organizationReadRepository) GetOrganizationBySocialUrl(ctx context.Context, tenant, socialUrl string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationBySocialUrl")
	defer spans.Finish()

	spans.LogKV("socialUrl", socialUrl)

	if socialUrl == "" || socialUrl == "/" {
		return nil, nil
	}

	var urlWithoutSlash, urlWithSlash string
	if len(socialUrl) > 0 && socialUrl[len(socialUrl)-1] == '/' {
		urlWithSlash = socialUrl
		urlWithoutSlash = socialUrl[:len(socialUrl)-1]
	} else {
		urlWithoutSlash = socialUrl
		urlWithSlash = socialUrl + "/"
	}
	spans.LogKV("urlWithoutSlash", urlWithoutSlash, "urlWithSlash", urlWithSlash)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(s:Social)
				WHERE s.url = $urlWithoutSlash OR s.url = $urlWithSlash
			   	RETURN o LIMIT 1`
	params := map[string]any{
		"tenant":          tenant,
		"urlWithoutSlash": urlWithoutSlash,
		"urlWithSlash":    urlWithSlash,
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)

	return result.(*dbtype.Node), err
}

func (r *organizationReadRepository) GetOrganizationsByLinkedIn(ctx context.Context, tenant, url, alias, externalId string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsByLinkedIn")
	defer spans.Finish()

	spans.LogKV("url", url, "alias", alias, "externalId", externalId)

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
	if i := strings.Index(minimizedUrl, "linkedin.com/company"); i != -1 {
		minimizedUrl = minimizedUrl[i:]
	}
	minimizedUrlWithSlash := minimizedUrl + "/"

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(s:Social)
				WHERE s.url ENDS WITH $url OR s.url ENDS WITH $urlWithSlash `
	if alias != "" {
		cypher += ` OR s.alias = $alias `
	}
	if externalId != "" {
		cypher += ` OR s.externalId = $externalId `
	}
	cypher += ` RETURN o ORDER by o.createdAt`
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

func (r *organizationReadRepository) GetForApiCache(ctx context.Context, tenant string, skip, limit int) ([]map[string]interface{}, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetForApiCache")
	defer spans.Finish()

	spans.LogKV("skip", skip, "limit", limit)

	cypher := ` MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WHERE o.hide = false
				
				OPTIONAL MATCH (o)-[:HAS_CONTRACT|HAS|TAGGED|SUBSIDIARY_OF]->(related)
				OPTIONAL MATCH (o)<-[:SUBSIDIARY_OF]-(sub:Organization)
				OPTIONAL MATCH (o)<-[:ROLE_IN]-(:JobRole)<-[:WORKS_AS]-(c:Contact)
				OPTIONAL MATCH (o)<-[:OWNS]-(u:User)
				
				WITH o,  
					collect(DISTINCT c.id) AS contactList,
     				collect(DISTINCT CASE WHEN related:Contract THEN related.id END) AS contractList,
     				collect(DISTINCT CASE WHEN related:Social THEN related.id END) AS socialList,
     				collect(DISTINCT CASE WHEN related:Tag THEN related.id END) AS tagList,
     				collect(DISTINCT CASE WHEN related:Organization THEN related.id END) AS parentList,
     				collect(DISTINCT sub.id) AS subsidiaryList,
     				u.id AS ownerId
				
				RETURN o, contactList, contractList, socialList, tagList, subsidiaryList, parentList, ownerId
				ORDER BY o.createdAt DESC
				SKIP $skip LIMIT $limit`
	params := map[string]any{
		"tenant": tenant,
		"skip":   skip,
		"limit":  limit,
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
		return nil, err
	}

	var results []map[string]interface{}
	if result != nil {
		for _, v := range result.([]*neo4j.Record) {
			organization := v.Values[0]
			contactList := v.Values[1]
			contractList := v.Values[2]
			socialList := v.Values[3]
			tagList := v.Values[4]
			subsidiaryList := v.Values[5]
			parentList := v.Values[6]
			ownerId := v.Values[7]

			record := map[string]interface{}{
				"organization":   organization,
				"contactList":    contactList,
				"contractList":   contractList,
				"socialList":     socialList,
				"tagList":        tagList,
				"subsidiaryList": subsidiaryList,
				"parentList":     parentList,
				"ownerId":        ownerId,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *organizationReadRepository) GetPatchesForApiCache(ctx context.Context, tenant string, lastPatchTimestamp time.Time) ([]map[string]interface{}, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetPatchesForApiCache")
	defer spans.Finish()

	cypher := ` MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)
				where o.updatedAt > $lastPatchTimestamp

				optional match (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)
				optional match (o)-[:HAS_CONTRACT]->(ctr:Contract)
				optional match (o)-[:HAS]->(s:Social)
				optional match (o)-[:TAGGED]->(t:Tag)
				optional match (o)<-[:SUBSIDIARY_OF]-(sub:Organization)
				optional match (o)-[:SUBSIDIARY_OF]->(par:Organization)
				
				optional match (o)<-[:OWNS]-(u:User)
				
				with o, 
				collect(c) as contactList, 
				collect(ctr) as contractList, 
				collect(s) as socialList, 
				collect(t) as tagList,
				collect(sub) as subsidiaryList,
				collect(par) as parentList,
				u.id as ownerId
				
				with o, 
				reduce(l = [], c in contactList | l + c.id) as contactList, 
				reduce(l = [], c in contractList | l + c.id) as contractList, 
				reduce(l = [], c in socialList | l + c.id) as socialList, 
				reduce(l = [], c in tagList | l + c.id) as tagList, 
				reduce(l = [], c in subsidiaryList | l + c.id) as subsidiaryList, 
				reduce(l = [], c in parentList | l + c.id) as parentList, 
				ownerId
				
				return o, contactList, contractList, socialList, tagList, subsidiaryList, parentList, ownerId ORDER BY o.createdAt DESC `
	params := map[string]any{
		"tenant":             tenant,
		"lastPatchTimestamp": lastPatchTimestamp,
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
		return nil, err
	}

	var results []map[string]interface{}
	if result != nil {
		for _, v := range result.([]*neo4j.Record) {
			organization := v.Values[0]
			contactList := v.Values[1]
			contractList := v.Values[2]
			socialList := v.Values[3]
			tagList := v.Values[4]
			subsidiaryList := v.Values[5]
			parentList := v.Values[6]
			ownerId := v.Values[7]

			record := map[string]interface{}{
				"organization":   organization,
				"contactList":    contactList,
				"contractList":   contractList,
				"socialList":     socialList,
				"tagList":        tagList,
				"subsidiaryList": subsidiaryList,
				"parentList":     parentList,
				"ownerId":        ownerId,
			}

			results = append(results, record)
		}
	}

	return results, nil
}

func (r *organizationReadRepository) GetAllForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetAllForInvoices")
	defer spans.Finish()

	spans.LogKV("invoiceIds", invoiceIds)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(:Contract)<-[:HAS_CONTRACT]-(o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				WHERE i.id IN $invoiceIds
				RETURN o, i.id`
	params := map[string]any{
		"tenant":     tenant,
		"invoiceIds": invoiceIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
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

func (r *organizationReadRepository) GetAllForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetAllForContracts")
	defer spans.Finish()

	spans.LogKV("contractIds", contractIds)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)
				WHERE c.id IN $contractIds
				RETURN o, c.id`
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
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

func (r *organizationReadRepository) GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetAllForSlackChannels")
	defer spans.Finish()

	spans.LogKV("slackChannelIds", slackChannelIds)

	cypher := `MATCH (t:Tenant {name:$tenant})-[:ORGANIZATION_BELONGS_TO_TENANT]->(o:Organization)
				WHERE o.slackChannelId IN $slackChannelIds
				RETURN o, o.slackChannelId`
	params := map[string]any{
		"tenant":          tenant,
		"slackChannelIds": slackChannelIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
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

func (r *organizationReadRepository) GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetAllForOpportunities")
	defer spans.Finish()

	spans.LogKV("opportunityIds", opportunityIds)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_OPPORTUNITY]->(op:Opportunity)
				WHERE op.id IN $opportunityIds
				RETURN org, op.id`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
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

func (r *organizationReadRepository) GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_CONTRACT]-(c:Contract)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity) 
				WITH t, org, collect(c) as contracts, collect(op) as ops 
					WHERE ALL(c IN contracts WHERE c.status = $liveStatus) 
				UNWIND ops AS op
				WITH t, org, min(op.renewedAt) AS minOpRenewalDate 
					WHERE date(org.derivedNextRenewalAt) < date(minOpRenewalDate) 
				RETURN t.name, org.id LIMIT $limit`
	params := map[string]any{
		"liveStatus": neo4jenum.ContractStatusLive.String(),
		"limit":      limit,
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsWithWebsiteAndWithoutDomains")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("delayInMinutes", delayInMinutes)

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) 
				WHERE NOT (org)-[:HAS_DOMAIN]->(:Domain) AND 
						org.website IS NOT NULL AND 
						org.website <> "" AND 
						(org.techDomainCheckedAt IS NULL OR org.techDomainCheckedAt < datetime() - duration({minutes: $delayInMinutes}))
				WITH t.name as tenant, org.id as orgId
				ORDER BY CASE WHEN org.techDomainCheckedAt IS NULL THEN 0 ELSE 1 END, org.techDomainCheckedAt ASC
				LIMIT $limit 
				RETURN tenant, orgId`
	params := map[string]any{
		"limit":          limit,
		"delayInMinutes": delayInMinutes,
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForEnrichByDomain(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsForEnrichByDomain")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("delayInMinutes", delayInMinutes)

	cypher := `MATCH (t:Tenant {active:true})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {primary:true})
				WHERE 	org.enrichedAt IS NULL AND
						org.hide = false AND
						(org.techEnrichAttempts IS NULL OR org.techEnrichAttempts < $maxAttempts OR org.techEnrichRequestedAt IS NULL) AND
						(org.techEnrichRequestedAt IS NULL OR org.techEnrichRequestedAt < datetime() - duration({minutes: $delayInMinutes}))
				WITH t.name as tenant, org.id as orgId, d.domain as domain
				ORDER BY CASE WHEN org.techEnrichRequestedAt IS NULL THEN 0 ELSE 1 END, org.techEnrichRequestedAt ASC
				LIMIT $limit
				RETURN tenant, orgId, domain`

	params := map[string]any{
		"limit":          limit,
		"delayInMinutes": delayInMinutes,
		"maxAttempts":    1,
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
	output := make([]TenantAndOrganizationIdExtended, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationIdExtended{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
				Param1:         v.Values[2].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForUpdateLastTouchpoint(ctx context.Context, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsForUpdateLastTouchpoint")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("delayFromPreviousCheckMin", delayFromPreviousCheckMin)

	cypher := `MATCH (t:Tenant {active:true})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
				WHERE org.hide = false AND
				(org.techLastTouchpointRequestedAt IS NULL OR org.techLastTouchpointRequestedAt < datetime() - duration({minutes: $delayFromPreviousCheckMin}))
				RETURN t.name, org.id
				ORDER BY CASE WHEN org.techLastTouchpointRequestedAt IS NULL THEN 0 ELSE 1 END, org.techLastTouchpointRequestedAt ASC
				LIMIT $limit`

	params := map[string]any{
		"limit":                     limit,
		"delayFromPreviousCheckMin": delayFromPreviousCheckMin,
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForIcpCheck(ctx context.Context, tenants []string, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsForIcpCheck")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("delayFromPreviousCheckMin", delayFromPreviousCheckMin)

	cypher := `MATCH (t:Tenant {active:true})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
				WHERE t.name IN $tenants AND 
				org.hide = false AND
				org.techIcpCheckedAt IS NULL AND
				(org.icpFit IS NULL OR org.icpFit = $icpNotSet) AND
				org.stage = $leadStage AND
				org.createdAt < datetime() - duration({minutes: $delayFromCreatedAt}) AND
				(org.techIcpCheckRequestedAt IS NULL OR org.techIcpCheckRequestedAt < datetime() - duration({minutes: $delayFromPreviousCheckMin}))
				RETURN t.name, org.id
				ORDER BY CASE WHEN org.techIcpCheckRequestedAt IS NULL THEN 0 ELSE 1 END, org.techIcpCheckRequestedAt ASC
				LIMIT $limit`

	params := map[string]any{
		"tenants":                   tenants,
		"limit":                     limit,
		"delayFromPreviousCheckMin": delayFromPreviousCheckMin,
		"delayFromCreatedAt":        5,
		"icpNotSet":                 commonenum.IcpNotSet.String(),
		"leadStage":                 commonenum.Lead.String(),
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetPrimaryOrganizationsWithJobRoleForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodePairAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts")
	defer spans.Finish()

	spans.LogKV("contactIds", contactIds)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:WORKS_AS]->(j:JobRole {primary:true})-[:ROLE_IN]->(o:Organization)
				WHERE c.id IN $contactIds AND o.hide = false
				RETURN o, j, c.id`
	params := map[string]any{
		"tenant":     tenant,
		"contactIds": contactIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodePairAndId(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodePairAndId)))
	return result.([]*utils.DbNodePairAndId), err
}

func (r *organizationReadRepository) GetHiddenOrganizationIds(ctx context.Context, tenant string, hiddenAfter time.Time) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetHiddenOrganizationIds")
	defer spans.Finish()

	spans.LogKV("hiddenAfter", hiddenAfter.String())

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
				WHERE org.hide = true AND org.hiddenAt > $hiddenAfter
				RETURN org.id ORDER BY org.hiddenAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"hiddenAfter": hiddenAfter,
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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]string)))
	return result.([]string), err
}

func (r *organizationReadRepository) GetMergedOrganizationIds(ctx context.Context, tenant string, mergedAfter time.Time) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetMergedOrganizationIds")
	defer spans.Finish()

	spans.LogKV("mergedAfter", mergedAfter.String())

	cypher := `MATCH (org:MergedOrganization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
				WHERE org.updatedAt >= $mergedAfter
				RETURN org.id ORDER BY org.updatedAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"mergedAfter": mergedAfter,
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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]string)))
	return result.([]string), err
}

func (r *organizationReadRepository) GetOrganizationsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ContactReadRepository.GetOrganizationsWithEmail")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT o`,
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
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *organizationReadRepository) GetActiveOrganizationIdsByDomain(ctx context.Context, tenant string, domains []string) (map[string]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetActiveOrganizationIdsByDomain")
	defer spans.Finish()

	spans.LogKV("domains", strings.Join(domains, ","))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain)
				WHERE d.domain IN $domains AND o.hide = false
				RETURN d.domain as domain, o.id as orgId`
	params := map[string]any{
		"tenant":  tenant,
		"domains": domains,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	output := make(map[string]string)
	for _, v := range records.([]*neo4j.Record) {
		output[v.Values[0].(string)] = v.Values[1].(string)
	}
	return output, err
}

func (r *organizationReadRepository) GetLinkedSubOrganizations(ctx context.Context, tenant string, parentOrganizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetLinkedSubOrganizations")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization)<-[rel:%s]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
								WHERE parent.id IN $parentOrganizationIds
								RETURN org, rel, parent.id ORDER BY org.name`, relationName)
	params := map[string]any{
		"tenant":                tenant,
		"parentOrganizationIds": parentOrganizationIds,
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

func (r *organizationReadRepository) GetLinkedParentOrganizations(ctx context.Context, tenant string, organizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetLinkedParentOrganizations")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization)-[rel:%s]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			WHERE sub.id IN $organizationIds
			RETURN org, rel, sub.id ORDER BY org.name`, relationName)
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
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

func (r *organizationReadRepository) GetOrganizationsByDomainAcrossAllTenants(ctx context.Context, domain string) ([]TenantAndOrganizationId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsByDomainAcrossAllTenants")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {active:true})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
				RETURN t.name, org.id`

	params := map[string]any{
		"domain": domain,
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsByStage(ctx context.Context, tenant string, stage commonenum.OrganizationStage) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "OrganizationReadRepository.GetOrganizationsByStage")
	defer spans.Finish()

	spans.LogKV("stage", stage.String())

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE org.hide = false AND org.stage = $stage
			RETURN org ORDER BY org.createdAt DESC`
	params := map[string]any{
		"tenant": tenant,
		"stage":  stage.String(),
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
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*dbtype.Node)))

	return result.([]*dbtype.Node), nil
}
