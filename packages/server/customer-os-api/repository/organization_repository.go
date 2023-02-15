package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"golang.org/x/net/context"
)

type OrganizationRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	GetOrganizationForJobRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error)
	GetOrganizationById(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) (*dbtype.Node, error)
	GetPaginatedOrganizations(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetPaginatedOrganizationsForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) error
	LinkWithOrganizationTypeInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId, organizationTypeId string) error
	UnlinkFromOrganizationTypesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string) error
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
	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:randomUUID()})" +
		" ON CREATE SET org.name=$name, " +
		"				org.description=$description, " +
		"				org.readonly=false, " +
		" 				org.domain=$domain, " +
		"				org.website=$website, " +
		"				org.industry=$industry, " +
		"				org.isPublic=$isPublic, " +
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
			"domain":        organization.Domain,
			"website":       organization.Website,
			"industry":      organization.Industry,
			"isPublic":      organization.IsPublic,
			"source":        organization.Source,
			"sourceOfTruth": organization.SourceOfTruth,
			"appSource":     organization.AppSource,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *organizationRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error) {
	query :=
		" MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})" +
			" SET 	org.name=$name, " +
			"		org.description=$description, " +
			"		org.domain=$domain, " +
			"		org.website=$website, " +
			" 		org.industry=$industry, " +
			"		org.isPublic=$isPublic, " +
			"		org.sourceOfTruth=$sourceOfTruth," +
			"		org.updatedAt=datetime({timezone: 'UTC'}) " +
			" RETURN org"

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organization.ID,
			"name":           organization.Name,
			"description":    organization.Description,
			"domain":         organization.Domain,
			"website":        organization.Website,
			"industry":       organization.Industry,
			"isPublic":       organization.IsPublic,
			"sourceOfTruth":  organization.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *organizationRepository) Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) error {
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

func (r *organizationRepository) GetOrganizationForJobRole(ctx context.Context, session neo4j.SessionWithContext, tenant, roleId string) (*dbtype.Node, error) {
	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:JobRole {id:$roleId})-[:ROLE_IN]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN org`,
			map[string]any{
				"tenant": tenant,
				"roleId": roleId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	if len(dbRecords.([]*neo4j.Record)) == 0 {
		return nil, nil
	}
	return utils.NodePtr(dbRecords.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
}

func (r *organizationRepository) GetOrganizationById(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) (*dbtype.Node, error) {
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

func (r *organizationRepository) GetPaginatedOrganizations(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

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

func (r *organizationRepository) GetPaginatedOrganizationsForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("org")
		countParams := map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}
		utils.MergeMapToMap(filterParams, countParams)

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})--(org:Organization) "+
				" %s "+
				" RETURN count(org) as count", filterCypherStr),
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
			" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})--(org:Organization) "+
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

func (r *organizationRepository) LinkWithOrganizationTypeInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId, organizationTypeId string) error {
	queryResult, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$organizationTypeId})
			MERGE (org)-[r:IS_OF_TYPE]->(ot)
			RETURN r`,
		map[string]any{
			"tenant":             tenant,
			"organizationId":     organizationId,
			"organizationTypeId": organizationTypeId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *organizationRepository) UnlinkFromOrganizationTypesInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string) error {
	if _, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				(org)-[r:IS_OF_TYPE]->(:OrganizationType)
			DELETE r`,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
		}); err != nil {
		return err
	}
	return nil
}
