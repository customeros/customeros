package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type OrganizationRepository interface {
	MergeOrganization(tenant string, syncDate time.Time, organization entity.OrganizationData) (string, error)
	MergeOrganizationAddress(tenant, organizationId string, organization entity.OrganizationData) error
	MergeOrganizationType(tenant, organizationId, organizationTypeName string) error
}

type organizationRepository struct {
	driver *neo4j.Driver
}

func NewOrganizationRepository(driver *neo4j.Driver) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) MergeOrganizationAddress(tenant, organizationId string, organization entity.OrganizationData) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (org)-[:LOCATED_AT]->(a:Address {source:$source}) " +
		" ON CREATE SET a.id=randomUUID(), a.appSource=$appSource, a.sourceOfTruth=$sourceOfTruth, " +
		"	a.country=$country, a.state=$state, a.city=$city, a.address=$address, " +
		"	a.address2=$address2, a.zip=$zip, a.phone=$phone, a:%s " +
		" ON MATCH SET 	a.country=$country, a.state=$state, a.city=$city, a.address=$address, " +
		"	a.address2=$address2, a.zip=$zip, a.phone=$phone"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Address_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"country":        organization.Country,
				"state":          organization.State,
				"city":           organization.City,
				"address":        organization.Address,
				"address2":       organization.Address2,
				"zip":            organization.Zip,
				"phone":          organization.Phone,
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) MergeOrganization(tenant string, syncDate time.Time, organization entity.OrganizationData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (org:Organization)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, org.id=randomUUID(), org.createdAt=$createdAt, " +
		"               org.name=$name, org.description=$description, r.syncDate=$syncDate, org.readonly=$readonly, " +
		"               org.domain=$domain, org.website=$website, org.industry=$industry, org.isPublic=$isPublic, " +
		"				org.source=$source, org.sourceOfTruth=$sourceOfTruth, org.appSource=$appSource, " +
		"				org:%s " +
		" ON MATCH SET org.name=$name, org.description=$description, r.syncDate=$syncDate, org.readonly=$readonly, " +
		"              org.domain=$domain, org.website=$website, org.industry=$industry, org.isPublic=$isPublic " +
		" WITH org, t " +
		" MERGE (org)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) " +
		" RETURN org.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Organization_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": organization.ExternalSystem,
				"externalId":     organization.ExternalId,
				"syncDate":       syncDate,
				"name":           organization.Name,
				"description":    organization.Description,
				"readonly":       organization.Readonly,
				"createdAt":      organization.CreatedAt,
				"domain":         organization.Domain,
				"website":        organization.Website,
				"industry":       organization.Industry,
				"isPublic":       organization.IsPublic,
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(string), nil
}

func (r *organizationRepository) MergeOrganizationType(tenant, organizationId, organizationTypeName string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (ot:OrganizationType {name:$organizationTypeName})-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET ot.id=randomUUID() " +
		" WITH org, ot " +
		" MERGE (org)-[r:IS_OF_TYPE]->(ot) " +
		" return r"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(query,
			map[string]interface{}{
				"tenant":               tenant,
				"organizationId":       organizationId,
				"organizationTypeName": organizationTypeName,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
