package repository

import (
	"context"
	"errors"
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	Relationship_Subsidiary = "SUBSIDIARY_OF"
)

type OrganizationRepository interface {
	CountCustomers(ctx context.Context, tenant string) (int64, error)
	GetPaginatedOrganizations(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedOrganizationsForContact(ctx context.Context, tenant, contactId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	MergeOrganizationPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string, sourceOfTruth neo4jentity.DataSource) error
	MergeOrganizationRelationsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string) error
	UpdateMergedOrganizationLabelsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, mergedOrganizationId string) error
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForJobRoles(ctx context.Context, tenant string, jobRoleIds []string) ([]*utils.DbNodeAndId, error)
	GetLinkedSubOrganizations(ctx context.Context, tenant string, parentOrganizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	GetLinkedParentOrganizations(ctx context.Context, tenant string, organizationIds []string, relationName string) ([]*utils.DbNodeWithRelationAndId, error)
	RemoveOwner(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetSuggestedMergePrimaryOrganizations(ctx context.Context, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error)
	GetOrganizations(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
}

type organizationRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext, database string) OrganizationRepository {
	return &organizationRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationRepository) CountCustomers(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CountCustomers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) where org.hide = false AND org.relationship = $customerRelationship
			RETURN count(org)`,
			map[string]any{
				"tenant":               tenant,
				"customerRelationship": neo4jenum.Customer,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	return dbRecord.(*db.Record).Values[0].(int64), nil
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

func (r *organizationRepository) MergeOrganizationPropertiesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, primaryOrganizationId, mergedOrganizationId string, sourceOfTruth neo4jentity.DataSource) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationPropertiesInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(primary:Organization {id:$primaryOrganizationId}),
			(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(merged:Organization {id:$mergedOrganizationId})
			SET primary.referenceId = CASE WHEN primary.referenceId is null OR primary.referenceId = '' THEN merged.referenceId ELSE primary.referenceId END, 
				primary.website = CASE WHEN primary.website is null OR primary.website = '' THEN merged.website ELSE primary.website END, 
				primary.industry = CASE WHEN primary.industry is null OR primary.industry = '' THEN merged.industry ELSE primary.industry END, 
				primary.subIndustry = CASE WHEN primary.subIndustry is null OR primary.subIndustry = '' THEN merged.subIndustry ELSE primary.subIndustry END, 
				primary.industryGroup = CASE WHEN primary.industryGroup is null OR primary.industryGroup = '' THEN merged.industryGroup ELSE primary.industryGroup END, 
				primary.name = CASE WHEN primary.name is null OR primary.name = '' THEN merged.name ELSE primary.name END, 
				primary.description = CASE WHEN primary.description is null OR primary.description = '' THEN merged.description ELSE primary.description END, 
				primary.isPublic = CASE WHEN primary.isPublic is null THEN merged.isPublic ELSE primary.isPublic END, 
				primary.employees = CASE WHEN primary.employees is null or primary.employees = 0 THEN merged.employees ELSE primary.employees END, 
				primary.market = CASE WHEN primary.market is null OR primary.market = '' THEN merged.market ELSE primary.market END, 
				primary.valueProposition = CASE WHEN primary.valueProposition is null OR primary.valueProposition = '' THEN merged.valueProposition ELSE primary.valueProposition END, 
				primary.targetAudience = CASE WHEN primary.targetAudience is null OR primary.targetAudience = '' THEN merged.targetAudience ELSE primary.targetAudience END, 
				primary.lastFundingRound = CASE WHEN primary.lastFundingRound is null OR primary.lastFundingRound = '' THEN merged.lastFundingRound ELSE primary.lastFundingRound END, 
				primary.lastFundingAmount = CASE WHEN primary.lastFundingAmount is null OR primary.lastFundingAmount = '' THEN merged.lastFundingAmount ELSE primary.lastFundingAmount END, 
				primary.note = CASE WHEN primary.note is null OR primary.note = '' THEN merged.note ELSE primary.note END, 
				primary.headquarters = CASE WHEN primary.headquarters is null OR primary.headquarters = '' THEN merged.headquarters ELSE primary.headquarters END, 
				primary.logoUrl = CASE WHEN primary.logoUrl is null OR primary.logoUrl = '' THEN merged.logoUrl ELSE primary.logoUrl END, 
				primary.iconUrl = CASE WHEN primary.iconUrl is null OR primary.iconUrl = '' THEN merged.iconUrl ELSE primary.iconUrl END,
				primary.employeeGrowthRate = CASE WHEN primary.employeeGrowthRate is null OR primary.employeeGrowthRate = '' THEN merged.employeeGrowthRate ELSE primary.employeeGrowthRate END,
				primary.yearFounded = CASE WHEN primary.yearFounded is null OR primary.yearFounded = 0 THEN merged.yearFounded ELSE primary.yearFounded END,
				primary.relationship = CASE WHEN primary.relationship is null THEN merged.relationship ELSE primary.relationship END,
				primary.stage = CASE WHEN primary.stage is null THEN merged.stage ELSE primary.stage END,
				primary.leadSource = CASE WHEN primary.leadSource is null OR primary.leadSource = '' THEN merged.leadSource ELSE primary.leadSource END,
				primary.icpFit = CASE WHEN primary.icpFit is null OR primary.leadSource = false THEN merged.icpFit ELSE primary.icpFit END,
				primary.sourceOfTruth=$sourceOfTruth,
				primary.updatedAt = datetime(),
				merged.updatedAt = datetime()
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
		" MATCH (merged)-[rel:LOGGED]->(n:LogEntry) "+
		" MERGE (primary)-[newRel:LOGGED]->(n) "+
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
		" MATCH (merged)-[rel:HAS]->(s:Social) "+
		" MERGE (primary)-[newRel:HAS]->(s) "+
		" ON CREATE SET newRel.primary=false, "+
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
		"				newRel.externalSource=rel.externalSource, "+
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
		" MATCH (merged)<-[rel:SUGGESTED_MERGE]-(org:Organization) "+
		" WHERE org.id <> primary.id "+
		" MERGE (primary)<-[newRel:SUGGESTED_MERGE]-(org) "+
		" ON CREATE SET newRel.suggestedBy = rel.suggestedBy, "+
		"				newRel.suggestedByInfo = rel.suggestedByInfo, "+
		"				newRel.confidence = rel.confidence, "+
		"				newRel.suggestedAt = rel.suggestedAt, "+
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
		" MATCH (merged)-[rel:HAS_CONTRACT]->(n:Contract) "+
		" MERGE (primary)-[newRel:HAS_CONTRACT]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_CONTRACT]->(n:DeletedContract) "+
		" MERGE (primary)-[newRel:HAS_CONTRACT]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_BILLING_PROFILE]->(n:BillingProfile) "+
		" MERGE (primary)-[newRel:HAS_BILLING_PROFILE]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_OPPORTUNITY]->(n:Opportunity) "+
		" MERGE (primary)-[newRel:HAS_OPPORTUNITY]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)-[rel:HAS_OPPORTUNITY]->(n:ArchivedOpportunity) "+
		" MERGE (primary)-[newRel:HAS_OPPORTUNITY]->(n) "+
		" ON CREATE SET newRel.mergedFrom = $mergedOrganizationId, "+
		"				newRel.createdAt = $now "+
		"			SET	rel.merged=true", params); err != nil {
		return err
	}

	if _, err := tx.Run(ctx, matchQuery+
		" WITH primary, merged "+
		" MATCH (merged)<-[rel:REMINDER_BELONGS_TO_ORGANIZATION]-(n:Reminder) "+
		" MERGE (primary)<-[newRel:REMINDER_BELONGS_TO_ORGANIZATION]-(n) "+
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *organizationRepository) RemoveOwner(ctx context.Context, tenant, organizationID string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
			OPTIONAL MATCH (:User)-[r:OWNS]->(org)
			SET org.updatedAt=datetime(), org.sourceOfTruth=$source
			DELETE r
			RETURN org`

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationID,
				"source":         neo4jentity.DataSourceOpenline,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *organizationRepository) GetSuggestedMergePrimaryOrganizations(ctx context.Context, organizationIds []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetSuggestedMergePrimaryOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[rel:SUGGESTED_MERGE]->(primaryOrg:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				WHERE org.id IN $organizationIds
				RETURN primaryOrg, rel, org.id 
				ORDER BY primaryOrg.name`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":          common.GetTenantFromContext(ctx),
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

func (r *organizationRepository) GetMinMaxRenewalForecastArr(ctx context.Context) (float64, float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetMinMaxRenewalForecastArr")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `CALL { MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
              RETURN min(o.renewalForecastArr) as min, max(o.renewalForecastArr) as max }
		      RETURN toFloat(min), toFloat(max)`
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
	}
	span.LogFields(log.String("cypher", cypher))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			record, err := queryResult.Single(ctx)
			if err != nil {
				return nil, err
			}
			return record, nil
		}
	})
	if err != nil {
		return 0, 0, err
	}
	if result != nil {
		record := result.(*db.Record)
		if record.Values[0] == nil || record.Values[1] == nil {
			return 0, 0, nil
		}
		minValue, minOk := record.Values[0].(float64)
		maxValue, maxOk := record.Values[1].(float64)
		if !minOk || !maxOk {
			return 0, 0, errors.New("unexpected type for min or max value")
		}
		return minValue, maxValue, nil
	} else {
		return 0, 0, nil
	}
}

func (r *organizationRepository) GetOrganizations(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
			WHERE org.id IN $ids
			RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*dbtype.Node), err
}
