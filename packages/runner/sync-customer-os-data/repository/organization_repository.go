package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"golang.org/x/net/context"
	"time"
)

type OrganizationRepository interface {
	MergeOrganization(ctx context.Context, tenant string, syncDate time.Time, organization entity.OrganizationData) (string, error)
	MergeOrganizationType(ctx context.Context, tenant, organizationId, organizationTypeName string) error
	MergeOrganizationDefaultPlace(ctx context.Context, tenant, organizationId string, organization entity.OrganizationData) error
	MergeOrganizationDomain(ctx context.Context, tenant string, organizationId string, domain string, externalSystem string) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) MergeOrganization(ctx context.Context, tenant string, syncDate time.Time, organization entity.OrganizationData) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Organization if it does not exist
	// If Organization exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateOrganization for incoming source, with a new relationship 'ALTERNATE'
	// Link Organization with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (org:Organization)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, r.syncDate=$syncDate, " +
		"				org.id=randomUUID(), org.createdAt=$createdAt, org.updatedAt=$createdAt, " +
		"               org.name=$name, org.description=$description, org.domain=$domain, " +
		"               org.website=$website, org.industry=$industry, org.isPublic=$isPublic, " +
		"				org.source=$source, org.sourceOfTruth=$sourceOfTruth, org.appSource=$appSource, " +
		"				org:%s " +
		" ON MATCH SET 	r.syncDate = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $syncDate ELSE r.syncDate END, " +
		"				org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $name ELSE org.name END, " +
		"				org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $description ELSE org.description END, " +
		"				org.domain = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $domain ELSE org.domain END, " +
		"				org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $website ELSE org.website END, " +
		"				org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $industry ELSE org.industry END, " +
		"				org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $isPublic ELSE org.isPublic END, " +
		"				org.updatedAt = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $now ELSE org.updatedAt END " +
		" WITH org, t " +
		" MERGE (org)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) " +
		" WITH org " +
		" FOREACH (x in CASE WHEN org.sourceOfTruth <> $sourceOfTruth THEN [org] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateOrganization {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" 		alt.name=$name, alt.description=$description, org.domain=$domain, org.website=$website, org.industry=$industry, org.isPublic=$isPublic " +
		") " +
		" RETURN org.id"

	dbRecord, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Organization_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": organization.ExternalSystem,
				"externalId":     organization.ExternalId,
				"syncDate":       syncDate,
				"name":           organization.Name,
				"description":    organization.Description,
				"createdAt":      organization.CreatedAt,
				"website":        organization.Website,
				"industry":       organization.Industry,
				"isPublic":       organization.IsPublic,
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
				"now":            time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single(ctx)
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

func (r *organizationRepository) MergeOrganizationType(ctx context.Context, tenant, organizationId, organizationTypeName string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (ot:OrganizationType {name:$organizationTypeName})-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET ot.id=randomUUID(), ot.createdAt=$now, ot.updatedAt=$now " +
		" WITH org, ot " +
		" MERGE (org)-[r:IS_OF_TYPE]->(ot) " +
		" return r"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":               tenant,
				"organizationId":       organizationId,
				"organizationTypeName": organizationTypeName,
				"now":                  time.Now().UTC(),
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

func (r *organizationRepository) MergeOrganizationDefaultPlace(ctx context.Context, tenant, organizationId string, organization entity.OrganizationData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location if it does not exist with given source and name
	// If Place exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternatePlace for incoming source, with a new relationship 'ALTERNATE'
	// !!! Current assumption - there is single Location with source of externalSystem and name per organization
	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (org)-[:ASSOCIATED_WITH]->(loc:Location {source:$source, name:$locationName}) " +
		" ON CREATE SET " +
		"	loc.country=$country, " +
		"	loc.region=$region, " +
		"	loc.locality=$locality, " +
		"	loc.address=$address, " +
		"	loc.address2=$address2, " +
		"	loc.zip=$zip, " +
		"	loc.phone=$phone, " +
		"	loc.id=randomUUID(), " +
		"	loc.appSource=$appSource, " +
		"	loc.sourceOfTruth=$sourceOfTruth, " +
		"	loc.createdAt=$createdAt, " +
		"	loc.updatedAt=$createdAt, " +
		"	loc:%s " +
		" ON MATCH SET 	" +
		"   loc.country = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $country ELSE loc.country END, " +
		"   loc.region = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $region ELSE loc.region END, " +
		"   loc.locality = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $locality ELSE loc.locality END, " +
		"	loc.address = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $address ELSE loc.address END, " +
		"	loc.address2 = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $address2 ELSE loc.address2 END, " +
		"	loc.zip = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $zip ELSE loc.zip END, " +
		"   loc.phone = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $phone ELSE loc.phone END, " +
		"   loc.updatedAt = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH loc, t " +
		" MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t) " +
		" WITH loc " +
		" FOREACH (x in CASE WHEN loc.sourceOfTruth <> $sourceOfTruth THEN [loc] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateLocation {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.region=$region, alt.locality=$locality, alt.address=$address, alt.address2=$address2, alt.zip=$zip, alt.phone=$phone " +
		") "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Location_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"country":        organization.Country,
				"region":         organization.Region,
				"locality":       organization.Locality,
				"address":        organization.Address,
				"address2":       organization.Address2,
				"zip":            organization.Zip,
				"phone":          organization.Phone,
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
				"locationName":   organization.DefaultLocationName,
				"createdAt":      organization.CreatedAt,
				"now":            time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) MergeOrganizationDomain(ctx context.Context, tenant string, organizationId string, domain string, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MERGE (d:Domain {domain:$domain}) 
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=$now,
								d.sourceOfTruth=$sourceOfTruth,
								d.appSource=$appSource,
								d.source=$source
				WITH d
				MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MERGE (org)-[rel:HAS_DOMAIN]->(d)
				RETURN rel`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"domain":         domain,
				"source":         externalSystem,
				"sourceOfTruth":  externalSystem,
				"appSource":      externalSystem,
				"now":            time.Now().UTC(),
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
