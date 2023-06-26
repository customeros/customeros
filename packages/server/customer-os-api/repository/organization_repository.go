package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

const (
	Relationship_Subsidiary = "SUBSIDIARY_OF"
)

type OrganizationRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	GetOrganizationById(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetPaginatedOrganizations(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedOrganizationsForContact(ctx context.Context, tenant, contactId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) error
	LinkWithDomainsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, domains []string) error
	UnlinkFromDomainsNotInListInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, domains []string) error
	MergeOrganizationPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string, sourceOfTruth entity.DataSource) error
	MergeOrganizationRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string) error
	UpdateMergedOrganizationLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedOrganizationId string) error
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForJobRoles(ctx context.Context, tenant string, jobRoleIds []string) ([]*utils.DbNodeAndId, error)
	GetLinkedSubOrganizations(ctx context.Context, tenant string, parentOrganizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	GetLinkedParentOrganizations(ctx context.Context, tenant string, organizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	LinkSubOrganization(ctx context.Context, tenant, organizationId, subOrganizationId, subOrganizationType, relationName string) error
	UnlinkSubOrganization(ctx context.Context, tenant, organizationId, subOrganizationId, relationName string) error
	ReplaceOwner(ctx context.Context, tenant, organizationID, userID string) (*dbtype.Node, error)
	RemoveOwner(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	AddRelationship(ctx context.Context, tenant, organizationId, relationship string) (*dbtype.Node, error)
	RemoveRelationship(ctx context.Context, tenant, organizationId, relationship string) (*dbtype.Node, error)
	SetRelationshipWithStage(ctx context.Context, tenant, organizationId, relationship, stage string) (*dbtype.Node, error)
	RemoveRelationshipStage(ctx context.Context, tenant, organizationId, relationship string) (*dbtype.Node, error)
	ReplaceHealthIndicator(ctx context.Context, organizationId, healthIndicatorId string) (*dbtype.Node, error)
	RemoveHealthIndicator(ctx context.Context, organizationId string) (*dbtype.Node, error)

	GetAllOrganizationPhoneNumberRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)
	GetAllOrganizationEmailRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:randomUUID()})" +
		" ON CREATE SET org.name=$name, " +
		"				org.description=$description, " +
		"				org.readonly=false, " +
		"				org.website=$website, " +
		"				org.industry=$industry, " +
		"				org.isPublic=$isPublic, " +
		"				org.employees=$employees, " +
		"				org.market=$market, " +
		" 				org.source=$source, " +
		"				org.sourceOfTruth=$sourceOfTruth, " +
		"				org.appSource=$appSource, " +
		"				org.createdAt=$now, " +
		"				org.updatedAt=$now, " +
		" 				org:%s" +
		" RETURN org"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Organization_"+tenant),
		map[string]any{
			"tenant":        tenant,
			"name":          organization.Name,
			"description":   organization.Description,
			"readonly":      false,
			"website":       organization.Website,
			"industry":      organization.Industry,
			"isPublic":      organization.IsPublic,
			"employees":     organization.Employees,
			"market":        organization.Market,
			"source":        organization.Source,
			"sourceOfTruth": organization.SourceOfTruth,
			"appSource":     organization.AppSource,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *organizationRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query :=
		" MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})" +
			" SET 	org.name=$name, " +
			"		org.description=$description, " +
			"		org.website=$website, " +
			" 		org.industry=$industry, " +
			"		org.isPublic=$isPublic, " +
			"		org.employees=$employees, " +
			"		org.market=$market, " +
			"		org.sourceOfTruth=$sourceOfTruth," +
			"		org.updatedAt=datetime({timezone: 'UTC'}) " +
			" RETURN org"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organization.ID,
			"name":           organization.Name,
			"description":    organization.Description,
			"website":        organization.Website,
			"industry":       organization.Industry,
			"isPublic":       organization.IsPublic,
			"employees":      organization.Employees,
			"market":         organization.Market,
			"sourceOfTruth":  organization.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *organizationRepository) Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (org)-[:ASSOCIATED_WITH]->(l:Location)
			OPTIONAL MATCH (l)-[:LOCATED_AT]->(p:Place)	
            DETACH DELETE p, l, org`,
			map[string]interface{}{
				"organizationId": organizationId,
				"tenant":         tenant,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) GetOrganizationById(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN org`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), nil
}

func (r *organizationRepository) GetPaginatedOrganizations(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetPaginatedOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("org")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) "+
				" %s "+
				" RETURN count(org) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) "+
				" %s "+
				" RETURN org "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sorting.SortingCypherFragment("org")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *organizationRepository) GetPaginatedOrganizationsForContact(ctx context.Context, tenant, contactId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetPaginatedOrganizationsForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("org")
		countParams := map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}
		utils.MergeMapToMap(filterParams, countParams)

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(org:Organization) "+
				" %s "+
				" RETURN count(distinct(org)) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"skip":      skip,
			"limit":     limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[:WORKS_AS]->(:JobRole)-[:ROLE_IN]->(org:Organization) "+
				" %s "+
				" RETURN distinct(org) "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sorting.SortingCypherFragment("org")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *organizationRepository) LinkWithDomainsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, domains []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.LinkWithDomainsInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			(d:Domain) WHERE d.domain IN $domains
			MERGE (org)-[:HAS_DOMAIN]->(d)`,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"domains":        domains,
		})
	return err
}

func (r *organizationRepository) UnlinkFromDomainsNotInListInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, domains []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UnlinkFromDomainsNotInListInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})-[rel:HAS_DOMAIN]->(d:Domain)
			WHERE NOT d.domain IN $domains
			DELETE rel`,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"domains":        domains,
		})
	return err
}

func (r *organizationRepository) MergeOrganizationPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string, sourceOfTruth entity.DataSource) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationPropertiesInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(primary:Organization {id:$primaryOrganizationId}),
			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(merged:Organization {id:$mergedOrganizationId})
			SET primary.website = CASE WHEN primary.website is null OR primary.website = '' THEN merged.website ELSE primary.website END, 
				primary.industry = CASE WHEN primary.industry is null OR primary.industry = '' THEN merged.industry ELSE primary.industry END, 
				primary.name = CASE WHEN primary.name is null OR primary.name = '' THEN merged.name ELSE primary.name END, 
				primary.description = CASE WHEN primary.description is null OR primary.description = '' THEN merged.description ELSE primary.description END, 
				primary.isPublic = CASE WHEN primary.isPublic is null THEN merged.isPublic ELSE primary.isPublic END, 
				primary.employees = CASE WHEN primary.employees is null or primary.employees = 0 THEN merged.employees ELSE primary.employees END, 
				primary.market = CASE WHEN primary.market is null OR primary.market = '' THEN merged.market ELSE primary.market END, 
				primary.sourceOfTruth=$sourceOfTruth,
				primary.updatedAt = $now
			`,
		map[string]any{
			"tenant":                tenant,
			"primaryOrganizationId": primaryOrganizationId,
			"mergedOrganizationId":  mergedOrganizationId,
			"sourceOfTruth":         string(sourceOfTruth),
			"now":                   utils.Now(),
		})
	return err
}

func (r *organizationRepository) MergeOrganizationRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationRelationsInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	matchQuery := "MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(primary:Organization {id:$primaryOrganizationId}), " +
		"(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(merged:Organization {id:$mergedOrganizationId})"

	params := map[string]any{
		"tenant":                tenant,
		"primaryOrganizationId": primaryOrganizationId,
		"mergedOrganizationId":  mergedOrganizationId,
		"now":                   utils.Now(),
	}

	if _, err := tx.Run(ctx, matchQuery+" "+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_DOMAIN]->(d:Domain) "+
		" MERGE (primary)-[newRel:HAS_DOMAIN]->(d)"+
		" ON CREATE SET	newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged"+
		" MATCH (merged)<-[rel:ROLE_IN]-(jb:JobRole) "+
		" MERGE (primary)<-[newRel:ROLE_IN]-(jb) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:ASSOCIATED_WITH]->(loc:Location) "+
		" MERGE (primary)-[newRel:ASSOCIATED_WITH]->(loc) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:NOTED]->(n:Note) "+
		" MERGE (primary)-[newRel:NOTED]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:CREATED]->(n:Note) "+
		" MERGE (primary)-[newRel:CREATED]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS]->(e:Email) "+
		" MERGE (primary)-[newRel:HAS]->(e) "+
		" ON CREATE SET newRel.primary=false, "+
		"				newRel.label=rel.label, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS]->(p:PhoneNumber) "+
		" MERGE (primary)-[newRel:HAS]->(p) "+
		" ON CREATE SET newRel.primary=false, "+
		"				newRel.label=rel.label, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:TAGGED]->(t:Tag) "+
		" MERGE (primary)-[newRel:TAGGED]->(t) "+
		" ON CREATE SET newRel.taggedAt=rel.taggedAt, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:REPORTED_BY]-(i:Issue) "+
		" MERGE (primary)<-[newRel:REPORTED_BY]-(i) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:IS_LINKED_WITH]->(ext:ExternalSystem) "+
		" MERGE (primary)-[newRel:IS_LINKED_WITH {externalId:rel.externalId}]->(ext) "+
		" ON CREATE SET newRel.syncDate=rel.syncDate, "+
		"				newRel.externalUrl=rel.externalUrl, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:SUBSIDIARY_OF]->(org:Organization) "+
		" WHERE org.id <> primary.id "+
		" MERGE (primary)-[newRel:SUBSIDIARY_OF]->(org) "+
		" ON CREATE SET newRel.type = rel.type, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:SUBSIDIARY_OF]-(org:Organization) "+
		" WHERE org.id <> primary.id "+
		" MERGE (primary)<-[newRel:SUBSIDIARY_OF]-(org) "+
		" ON CREATE SET newRel.type = rel.type, "+
		"				newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:SENT_BY]-(i:InteractionEvent) "+
		" MERGE (primary)<-[newRel:SENT_BY]-(i) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now, "+
		"				newRel.type = rel.type "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:SENT_TO]-(i:InteractionEvent) "+
		" MERGE (primary)<-[newRel:SENT_TO]-(i) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now, "+
		"				newRel.type = rel.type "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" OPTIONAL MATCH (primary)<-[:OWNS]-(existing:User) "+
		" WITH primary, merged, existing "+
		" WHERE existing IS NULL "+
		" MATCH (merged)<-[rel:OWNS]-(u:User) "+
		" MERGE (primary)<-[newRel:OWNS]-(u) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" OPTIONAL MATCH (primary)-[:HAS_INDICATOR]->(existing:HealthIndicator) "+
		" WITH primary, merged, existing "+
		" WHERE existing IS NULL "+
		" MATCH (merged)-[rel:HAS_INDICATOR]->(h:HealthIndicator) "+
		" MERGE (primary)-[newRel:HAS_INDICATOR]->(h) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:IS]->(or:OrganizationRelationship) "+
		" MERGE (primary)-[newRel:IS]->(or) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_STAGE]->(ors:OrganizationRelationshipStage)<-[:HAS_STAGE]-(or:OrganizationRelationship) "+
		" WITH primary, merged, ors, or, rel "+
		" WHERE NOT (primary)-[:HAS_STAGE]->(:OrganizationRelationshipStage)<-[:HAS_STAGE]-(or) "+
		" MERGE (primary)-[newRel:HAS_STAGE]->(ors) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MERGE (merged)-[rel:IS_MERGED_INTO]->(primary) "+
		" ON CREATE SET rel.mergedAt=$now", params); err != nil {
		return err
	}

	return nil
}

func (r *organizationRepository) UpdateMergedOrganizationLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedOrganizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateMergedOrganizationLabelsInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}) " +
		" SET org:MergedOrganization:%s " +
		" REMOVE org:Organization:%s"

	_, err := tx.Run(ctx, fmt.Sprintf(query, "MergedOrganization_"+tenant, "Organization_"+tenant),
		map[string]any{
			"tenant":         tenant,
			"organizationId": mergedOrganizationId,
		})
	return err
}

func (r *organizationRepository) GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE e.id IN $emailIds
			RETURN o, e.id as emailId ORDER BY o.name`,
			map[string]any{
				"tenant":   tenant,
				"emailIds": emailIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationRepository) GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
			WHERE p.id IN $phoneNumberIds
			RETURN o, p.id as phoneNumberId ORDER BY o.name`,
			map[string]any{
				"tenant":         tenant,
				"phoneNumberIds": phoneNumberIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationRepository) GetAllForJobRoles(ctx context.Context, tenant string, jobRoleIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllForJobRoles")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("jobRoleIds", fmt.Sprintf("%v", jobRoleIds)))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:ROLE_IN]-(j:JobRole)
				WHERE j.id IN $jobRoleIds
				RETURN o, j.id`

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":     tenant,
				"jobRoleIds": jobRoleIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationRepository) GetLinkedSubOrganizations(ctx context.Context, tenant string, parentOrganizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetLinkedSubOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization)<-[rel:%s]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
								WHERE parent.id IN $parentOrganizationIds
								RETURN org, rel, parent.id ORDER BY org.name`, relationName)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":                tenant,
				"parentOrganizationIds": parentOrganizationIds,
			}); err != nil {
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

func (r *organizationRepository) GetLinkedParentOrganizations(ctx context.Context, tenant string, organizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetLinkedParentOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization)-[rel:%s]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			WHERE sub.id IN $organizationIds
			RETURN org, rel, sub.id ORDER BY org.name`, relationName)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIds,
			}); err != nil {
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

func (r *organizationRepository) LinkSubOrganization(ctx context.Context, tenant, organizationId, subOrganizationId, subOrganizationType, relationName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.LinkSubOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$organizationId})," +
		" (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(sub:Organization {id:$subOrganizationId}) " +
		" MERGE (parent)<-[rel:%s]-(sub) " +
		" ON CREATE SET rel.type=$type " +
		" ON MATCH SET rel.type=$type " +
		" return parent.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, relationName),
			map[string]interface{}{
				"tenant":            tenant,
				"organizationId":    organizationId,
				"subOrganizationId": subOrganizationId,
				"type":              subOrganizationType,
				"now":               utils.Now(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *organizationRepository) UnlinkSubOrganization(ctx context.Context, tenant, organizationId, subOrganizationId, relationName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UnlinkSubOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(parent:Organization {id:$organizationId})<-[rel:%s]-(sub:Organization {id:$subOrganizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)" +
		" DELETE rel "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, relationName),
			map[string]interface{}{
				"tenant":            tenant,
				"organizationId":    organizationId,
				"subOrganizationId": subOrganizationId,
				"now":               utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) GetAllOrganizationPhoneNumberRelationships(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllOrganizationPhoneNumberRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[rel:HAS]->(p:PhoneNumber)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, org.id, p.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *organizationRepository) GetAllOrganizationEmailRelationships(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllOrganizationEmailRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[rel:HAS]->(e:Email)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, org.id, e.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *organizationRepository) ReplaceOwner(ctx context.Context, tenant, organizationID, userID string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.ReplaceOwner")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[rel:OWNS]->(org)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=$now, org.sourceOfTruth=$source			
			RETURN org`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"userId":         userID,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) RemoveOwner(ctx context.Context, tenant, organizationID string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[r:OWNS]->(org)
			SET org.updatedAt=$now, org.sourceOfTruth=$source
			DELETE r
			RETURN org`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) AddRelationship(ctx context.Context, tenant, organizationID, relationship string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.AddRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
			(or:OrganizationRelationship {name:$relationship})
			MERGE (org)-[:IS]->(or)
			SET org.updatedAt=$now, org.sourceOfTruth=$source			
			RETURN org`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"relationship":   relationship,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) RemoveRelationship(ctx context.Context, tenant, organizationID, relationship string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (org)-[rel:IS]->(or:OrganizationRelationship {name:$relationship})
			OPTIONAL MATCH (org)-[rel_stage:HAS_STAGE]->(:OrganizationRelationshipStage)<-[:HAS_STAGE]-(or)
			SET org.updatedAt=$now, org.sourceOfTruth=$source			
			DELETE rel, rel_stage
			RETURN distinct(org)`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"relationship":   relationship,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) SetRelationshipWithStage(ctx context.Context, tenant, organizationID, relationship, stage string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetRelationshipWithStage")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			MATCH (or:OrganizationRelationship {name:$relationship})
			MERGE (org)-[:IS]->(or)
			WITH t, org, or
			CALL { 
				WITH org, or
				MATCH (or)-[:HAS_STAGE]->(:OrganizationRelationshipStage {name:$stage}),
			    (or)-[:HAS_STAGE]->(existing:OrganizationRelationshipStage)<-[rel:HAS_STAGE]-(org)
			    WHERE existing.name <> $stage
			    DELETE rel
			}
			WITH t, org, or
			CALL {
				WITH t, org, or
				MATCH (t)<-[:STAGE_BELONGS_TO_TENANT]-(ors:OrganizationRelationshipStage {name:$stage})<-[:HAS_STAGE]-(or)
				MERGE (org)-[:HAS_STAGE]->(ors)
			}
			WITH org
			SET org.updatedAt=$now, org.sourceOfTruth=$source
			RETURN org`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"relationship":   relationship,
				"stage":          stage,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) RemoveRelationshipStage(ctx context.Context, tenant, organizationID, relationship string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.RemoveRelationship")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (org)-[:IS]->(or:OrganizationRelationship {name:$relationship})
			OPTIONAL MATCH (org)-[rel_stage:HAS_STAGE]->(:OrganizationRelationshipStage)<-[:HAS_STAGE]-(or)
			SET org.updatedAt=$now, org.sourceOfTruth=$source			
			DELETE rel_stage
			RETURN distinct(org)`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"relationship":   relationship,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateLastTouchpointByOrganizationId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("touchpointId", touchpointId), log.Object("touchpointAt", touchpointAt))

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId`

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"touchpointAt":   touchpointAt,
				"touchpointId":   touchpointId,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) ReplaceHealthIndicator(ctx context.Context, organizationId, healthIndicatorId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.ReplaceHealthIndicator")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (org)-[rel:HAS_INDICATOR]->(:HealthIndicator)
			DELETE rel
			WITH org, t
			MATCH (t)<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator {id:$healthIndicatorId})
			MERGE (org)-[:HAS_INDICATOR]->(h)
			SET org.updatedAt=$now, org.sourceOfTruth=$source			
			RETURN org`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":            common.GetTenantFromContext(ctx),
				"organizationId":    organizationId,
				"healthIndicatorId": healthIndicatorId,
				"source":            entity.DataSourceOpenline,
				"now":               utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) RemoveHealthIndicator(ctx context.Context, organizationID string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.RemoveHealthIndicator")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:HealthIndicator)<-[r:HAS_INDICATOR]->(org)
			SET org.updatedAt=$now, org.sourceOfTruth=$source
			DELETE r
			RETURN org`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         common.GetTenantFromContext(ctx),
				"organizationId": organizationID,
				"source":         entity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
