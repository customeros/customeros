package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type OrganizationRepository interface {
	Create(tx neo4j.Transaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	Update(tx neo4j.Transaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error)
	FindOrganizationForRole(session neo4j.Session, tenant, roleId string) (*dbtype.Node, error)
	GetOrganizationById(session neo4j.Session, tenant, organizationId string) (*dbtype.Node, error)
	GetPaginatedOrganizations(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	Delete(session neo4j.Session, tenant, organizationId string) error
	LinkWithOrganizationTypeInTx(tx neo4j.Transaction, tenant, organizationId, organizationTypeId string) error
	UnlinkFromOrganizationTypesInTx(tx neo4j.Transaction, tenant, organizationId string) error
}

type organizationRepository struct {
	driver *neo4j.Driver
}

func NewOrganizationRepository(driver *neo4j.Driver) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) Create(tx neo4j.Transaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error) {
	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:randomUUID()})" +
		" ON CREATE SET org.name=$name, org.description=$description, org.readonly=false," +
		" org.domain=$domain, org.website=$website, org.industry=$industry, org.isPublic=$isPublic, " +
		" org.source=$source, org.sourceOfTruth=$sourceOfTruth," +
		" org:%s" +
		" RETURN org"

	queryResult, err := tx.Run(fmt.Sprintf(query, "Organization_"+tenant),
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
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *organizationRepository) Update(tx neo4j.Transaction, tenant string, organization entity.OrganizationEntity) (*dbtype.Node, error) {
	query :=
		" MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})" +
			" SET org.name=$name, org.description=$description, org.domain=$domain, org.website=$website, " +
			" org.industry=$industry, org.isPublic=$isPublic, org.sourceOfTruth=$sourceOfTruth " +
			" RETURN org"

	queryResult, err := tx.Run(query,
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
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *organizationRepository) Delete(session neo4j.Session, tenant, organizationId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (org)-[:LOCATED_AT]->(a:Address)
            DETACH DELETE a, org`,
			map[string]interface{}{
				"organizationId": organizationId,
				"tenant":         tenant,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) FindOrganizationForRole(session neo4j.Session, tenant, roleId string) (*dbtype.Node, error) {
	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Role {id:$roleId})-[:WORKS]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN org`,
			map[string]any{
				"tenant": tenant,
				"roleId": roleId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
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

func (r *organizationRepository) GetOrganizationById(session neo4j.Session, tenant, organizationId string) (*dbtype.Node, error) {
	dbRecord, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			RETURN org`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Single()
		}
	})
	if err != nil {
		return nil, err
	}
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), nil
}

func (r *organizationRepository) GetPaginatedOrganizations(session neo4j.Session, tenant string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("org")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)

		queryResult, err := tx.Run(fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) "+
				" %s "+
				" RETURN count(org) as count", filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) "+
				" %s "+
				" RETURN org "+
				" %s "+
				" SKIP $skip LIMIT $limit", filterCypherStr, sorting.SortingCypherFragment("org")),
			params)
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *organizationRepository) LinkWithOrganizationTypeInTx(tx neo4j.Transaction, tenant, organizationId, organizationTypeId string) error {
	queryResult, err := tx.Run(`
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
	_, err = queryResult.Single()
	return err
}

func (r *organizationRepository) UnlinkFromOrganizationTypesInTx(tx neo4j.Transaction, tenant, organizationId string) error {
	if _, err := tx.Run(`
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
