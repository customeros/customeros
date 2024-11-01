package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
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
	GetOrganizationByDomain(ctx context.Context, tenant, domain string) (*dbtype.Node, error)
	GetOrganizationBySocialUrl(ctx context.Context, tenant, socialUrl string) (*dbtype.Node, error)
	GetForApiCache(ctx context.Context, tenant string, skip, limit int) ([]map[string]interface{}, error)
	GetPatchesForApiCache(ctx context.Context, tenant string, lastPatchTimestamp time.Time) ([]map[string]interface{}, error)
	GetAllForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error)
	GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error)
	GetOrganizationsForEnrichByDomain(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error)
	GetOrganizationsForAdjustIndustry(ctx context.Context, delayInMinutes, limit int, validIndustries []string) ([]TenantAndOrganizationId, error)
	GetOrganizationsForUpdateLastTouchpoint(ctx context.Context, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error)
	GetLatestOrganizationWithJobRoleForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodePairAndId, error)
	GetHiddenOrganizationIds(ctx context.Context, tenant string, hiddenAfter time.Time) ([]string, error)
	GetMergedOrganizationIds(ctx context.Context, tenant string, mergedAfter time.Time) ([]string, error)
	GetOrganizationsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.CountByTenant")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) where org.hide = false
			RETURN count(org)`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span.LogFields(log.Int64("result", organizationsCount))
	return organizationsCount, nil
}

func (r *organizationReadRepository) GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationIdsConnectedToInteractionEvent")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

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
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), err
}

func (r *organizationReadRepository) GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByOpportunityId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId))

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
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByContactId(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByContactId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contactId", contactId))

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}), 
				(t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				WHERE (c)-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(org) 
			RETURN org limit 1`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByContractId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_CONTRACT]->(c:Contract {id:$id})
			RETURN org limit 1`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByInvoiceId(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByInvoiceId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("invoiceId", invoiceId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(inv:Invoice {id:$invoiceId})<-[:HAS_INVOICE]-(c:Contract)<-[:HAS_CONTRACT]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			RETURN org`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByCustomerOsId(ctx context.Context, tenant, customerOsId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByCustomerOsId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("customerOsId", customerOsId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {customerOsId:$customerOsId})
			RETURN org`
	params := map[string]any{
		"tenant":       tenant,
		"customerOsId": customerOsId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByIdOrCustomerOsId(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByIdOrCustomerOsId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("id", id))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
			WHERE org.customerOsId = $id OR org.id = $id RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByReferenceId(ctx context.Context, tenant, referenceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByReferenceId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("referenceId", referenceId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {referenceId:$referenceId}) RETURN org`
	params := map[string]any{
		"tenant":      tenant,
		"referenceId": referenceId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByDomain(ctx context.Context, tenant, domain string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByDomain")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("domain", domain))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain{domain:$domain}) RETURN o limit 1`

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, map[string]any{
			"tenant": tenant,
			"domain": domain,
		}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))

	return result.(*dbtype.Node), err
}

func (r *organizationReadRepository) GetOrganizationBySocialUrl(ctx context.Context, tenant, socialUrl string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationBySocialUrl")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("socialUrl", socialUrl)

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
	span.LogKV("urlWithoutSlash", urlWithoutSlash, "urlWithSlash", urlWithSlash)

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
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))

	return result.(*dbtype.Node), err
}

func (r *organizationReadRepository) GetForApiCache(ctx context.Context, tenant string, skip, limit int) ([]map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetForApiCache")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("skip", skip), log.Object("limit", limit))

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

	span.LogFields(log.String("query", cypher))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetPatchesForApiCache")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

	span.LogFields(log.String("query", cypher))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForInvoices")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(:Contract)<-[:HAS_CONTRACT]-(o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				WHERE i.id IN $invoiceIds
				RETURN o, i.id`
	params := map[string]any{
		"tenant":     tenant,
		"invoiceIds": invoiceIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForSlackChannels")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("slackChannelIds", slackChannelIds))

	cypher := `MATCH (t:Tenant {name:$tenant})-[:ORGANIZATION_BELONGS_TO_TENANT]->(o:Organization)
				WHERE o.slackChannelId IN $slackChannelIds
				RETURN o, i.id`
	params := map[string]any{
		"tenant":          tenant,
		"slackChannelIds": slackChannelIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForOpportunities")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_OPPORTUNITY]->(op:Opportunity)
				WHERE op.id IN $opportunityIds
				RETURN org, op.id`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit))

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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsWithWebsiteAndWithoutDomains")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit), log.Int("delayInMinutes", delayInMinutes))

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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForEnrichByDomain(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForEnrichByDomain")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit), log.Int("delayInMinutes", delayInMinutes))

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
	output := make([]TenantAndOrganizationIdExtended, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationIdExtended{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
				Param1:         v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForAdjustIndustry(ctx context.Context, delayInMinutes, limit int, validIndustries []string) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForAdjustIndustry")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit), log.Int("delayInMinutes", delayInMinutes))

	cypher := `MATCH (t:Tenant {active:true})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
				WHERE org.hide = false AND
						org.industry <> '' AND
						org.industry IS NOT NULL AND
						NOT org.industry IN $validIndustries AND
						org.updatedAt < datetime() - duration({minutes: $minutesFromUpdate}) AND
						(org.techIndustryCheckedAt IS NULL OR org.techIndustryCheckedAt < datetime() - duration({minutes: $delayInMinutes}))
				RETURN t.name, org.id
				ORDER BY CASE WHEN org.techIndustryCheckedAt IS NULL THEN 0 ELSE 1 END, org.techIndustryCheckedAt ASC
				LIMIT $limit`

	params := map[string]any{
		"limit":             limit,
		"delayInMinutes":    delayInMinutes,
		"validIndustries":   validIndustries,
		"minutesFromUpdate": 2,
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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForUpdateLastTouchpoint(ctx context.Context, limit, delayFromPreviousCheckMin int) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForUpdateLastTouchpoint")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit))

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
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetLatestOrganizationWithJobRoleForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodePairAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetLatestOrganizationWithJobRoleForContacts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("contactIds", contactIds))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(o:Organization)
				WHERE c.id IN $contactIds
				WITH c, o, j
				ORDER BY 
  					CASE 
    					WHEN j.endedAt IS NULL THEN 1 
    					ELSE 0 
  					END DESC, 
  					j.endedAt DESC, 
  					j.startedAt DESC
				WITH c, COLLECT(o)[0] AS latestOrganization, COLLECT(j)[0] AS latestJobRole
				RETURN latestOrganization, latestJobRole, c.id`
	params := map[string]any{
		"tenant":     tenant,
		"contactIds": contactIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodePairAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodePairAndId))))
	return result.([]*utils.DbNodePairAndId), err
}

func (r *organizationReadRepository) GetHiddenOrganizationIds(ctx context.Context, tenant string, hiddenAfter time.Time) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetHiddenOrganizationIds")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("hiddenAfter", hiddenAfter.String()))

	cypher := `MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
				WHERE org.hide = true AND org.hiddenAt > $hiddenAfter
				RETURN org.id ORDER BY org.hiddenAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"hiddenAfter": hiddenAfter,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), err
}

func (r *organizationReadRepository) GetMergedOrganizationIds(ctx context.Context, tenant string, mergedAfter time.Time) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetMergedOrganizationIds")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("mergedAfter", mergedAfter.String()))

	cypher := `MATCH (org:MergedOrganization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
				WHERE org.updatedAt >= $mergedAfter
				RETURN org.id ORDER BY org.updatedAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"mergedAfter": mergedAfter,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), err
}

func (r *organizationReadRepository) GetOrganizationsWithEmail(ctx context.Context, tenant, email string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactReadRepository.GetOrganizationsWithEmail")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err
}
