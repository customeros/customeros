package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type OrganizationRepository interface {
	GetMatchedOrganizationId(ctx context.Context, tenant string, organization entity.OrganizationData) (string, error)
	MergeOrganization(ctx context.Context, tenant string, syncDate time.Time, organization entity.OrganizationData) error
	MergeOrganizationRelationship(ctx context.Context, tenant, organizationId, relationship, externalSystem string) error
	MergeOrganizationLocation(ctx context.Context, tenant, organizationId string, organization entity.OrganizationData) error
	MergeOrganizationDomain(ctx context.Context, tenant, organizationId, domain, externalSystem string) error
	MergePhoneNumber(ctx context.Context, tenant, organizationId, phoneNumber, externalSystem string, createdAt time.Time) error
	MergeEmail(ctx context.Context, tenant, organizationId, email, externalSystem string, createdAt time.Time) error
	LinkToParentOrganizationAsSubsidiary(ctx context.Context, tenant, organizationId, externalSystem string, parentOrganizationDtls *entity.ParentOrganization) error
	SetOwner(ctx context.Context, tenant, contactId, userExternalOwnerId, externalSystem string) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) GetMatchedOrganizationId(ctx context.Context, tenant string, organization entity.OrganizationData) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o1:Organization)-[:IS_LINKED_WITH {externalId:$organizationExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o2:Organization)-[:HAS_DOMAIN]->(d:Domain)
					WHERE d.domain in $domains
				with coalesce(o1, o2) as organization
				where organization is not null
				return organization.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":                 tenant,
				"externalSystem":         organization.ExternalSystem,
				"organizationExternalId": organization.ExternalId,
				"domains":                organization.Domains,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	orgIDs := dbRecords.([]*db.Record)
	if len(orgIDs) == 1 {
		return orgIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *organizationRepository) MergeOrganization(ctx context.Context, tenant string, syncDate time.Time, organization entity.OrganizationData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Organization if it does not exist
	// If Organization exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateOrganization for incoming source, with a new relationship 'ALTERNATE'
	// Link Organization with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (org:Organization {id:$orgId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET org.createdAt=$createdAt, " +
		"				org.updatedAt=$updatedAt, " +
		"               org.tenantOrganization=false, " +
		"               org.name=$name, " +
		"				org.description=$description, " +
		"               org.website=$website, " +
		"				org.industry=$industry, " +
		"				org.isPublic=$isPublic, " +
		"				org.employees=$employees, " +
		"				org.source=$source, " +
		"				org.sourceOfTruth=$sourceOfTruth, " +
		"				org.appSource=$appSource, " +
		"				org:%s " +
		" ON MATCH SET 	org.name = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.name is null OR org.name = '' THEN $name ELSE org.name END, " +
		"				org.description = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.description is null OR org.description = '' THEN $description ELSE org.description END, " +
		"				org.website = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.website is null OR org.website = '' THEN $website ELSE org.website END, " +
		"				org.industry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.industry is null OR org.industry = '' THEN $industry ELSE org.industry END, " +
		"				org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $isPublic ELSE org.isPublic END, " +
		"				org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $employees ELSE org.employees END, " +
		"				org.updatedAt = $now " +
		" WITH org, ext " +
		" MERGE (org)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH org " +
		" FOREACH (x in CASE WHEN org.sourceOfTruth <> $sourceOfTruth THEN [org] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateOrganization {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" 		alt.name=$name, alt.description=$description, alt.website=$website, alt.industry=$industry, alt.isPublic=$isPublic " +
		") " +
		" RETURN org.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Organization_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"orgId":          organization.Id,
				"externalSystem": organization.ExternalSystem,
				"externalId":     organization.ExternalId,
				"externalUrl":    organization.ExternalUrl,
				"syncDate":       syncDate,
				"name":           organization.Name,
				"description":    organization.Description,
				"createdAt":      organization.CreatedAt,
				"updatedAt":      organization.UpdatedAt,
				"website":        organization.Website,
				"industry":       organization.Industry,
				"isPublic":       organization.IsPublic,
				"employees":      organization.Employees,
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
				"now":            utils.Now(),
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

func (r *organizationRepository) MergeOrganizationRelationship(ctx context.Context, tenant, organizationId, relationship, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// TODO set where condition after re-run sync for existing tenants
	//WHERE org.sourceOfTruth=$sourceOfTruth
	query := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
				
				WITH org
		 		MATCH (or:OrganizationRelationship {name:$relationship}) 
		 		MERGE (org)-[:IS]->(or) 
				SET org.updatedAt=$now`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"relationship":   relationship,
				"now":            utils.Now(),
				"sourceOfTruth":  externalSystem,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *organizationRepository) MergeOrganizationLocation(ctx context.Context, tenant, organizationId string, organization entity.OrganizationData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location if it does not exist with given source property
	// If Place exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternatePlace for incoming source, with a new relationship 'ALTERNATE'
	// !!! Current assumption - there is single Location with source of externalSystem per organization
	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (org)-[:ASSOCIATED_WITH]->(loc:Location {source:$source}) " +
		" ON CREATE SET " +
		"	loc.name=$locationName, " +
		"	loc.country=$country, " +
		"	loc.region=$region, " +
		"	loc.locality=$locality, " +
		"	loc.address=$address, " +
		"	loc.address2=$address2, " +
		"	loc.zip=$zip, " +
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
		"   loc.updatedAt = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH loc, t " +
		" MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t) " +
		" WITH loc " +
		" FOREACH (x in CASE WHEN loc.sourceOfTruth <> $sourceOfTruth THEN [loc] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateLocation {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.region=$region, alt.locality=$locality, alt.address=$address, alt.address2=$address2, alt.zip=$zip " +
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
				"source":         organization.ExternalSystem,
				"sourceOfTruth":  organization.ExternalSystem,
				"appSource":      organization.ExternalSystem,
				"locationName":   organization.LocationName,
				"createdAt":      organization.CreatedAt,
				"now":            utils.Now(),
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
				"appSource":      externalSystem,
				"now":            utils.Now(),
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

func (r *organizationRepository) MergePhoneNumber(ctx context.Context, tenant, organizationId, phoneNumber, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (p:PhoneNumber {rawPhoneNumber: $phoneNumber})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				p.id=randomUUID(), " +
		"				p.createdAt=$now, " +
		"				p.updatedAt=$now, " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.appSource=$appSource, " +
		"				p:%s " +
		" WITH DISTINCT o, p " +
		" MERGE (o)-[rel:HAS]->(p) " +
		" ON CREATE SET rel.primary=false, p.updatedAt=$now, o.updatedAt=$now "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "PhoneNumber_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"phoneNumber":    phoneNumber,
				"createdAt":      createdAt,
				"source":         externalSystem,
				"sourceOfTruth":  externalSystem,
				"appSource":      externalSystem,
				"now":            utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) MergeEmail(ctx context.Context, tenant, organizationId, email, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				e.id=randomUUID(), " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.appSource=$appSource, " +
		"				e:%s " +
		" WITH DISTINCT o, e " +
		" MERGE (o)-[rel:HAS]->(e) " +
		" ON CREATE SET rel.primary=false, e.updatedAt=$now, o.updatedAt=$now "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"email":          email,
				"createdAt":      createdAt,
				"source":         externalSystem,
				"sourceOfTruth":  externalSystem,
				"appSource":      externalSystem,
				"now":            utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) LinkToParentOrganizationAsSubsidiary(ctx context.Context, tenant, organizationId, externalSystem string, parentOrganizationDtls *entity.ParentOrganization) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$parentExternalId}]-(parent:Organization),
				(t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MERGE (org)-[rel:SUBSIDIARY_OF]->(parent)
				ON CREATE SET rel.type=$type`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":           tenant,
				"externalSystem":   externalSystem,
				"parentExternalId": parentOrganizationDtls.ExternalId,
				"organizationId":   organizationId,
				"type":             parentOrganizationDtls.Type,
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) SetOwner(ctx context.Context, tenant, organizationId, userExternalOwnerId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// TODO set where condition after re-run sync for existing tenants
	//WHERE org.sourceOfTruth=$sourceOfTruth
	query := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})

			WITH org, t
			OPTIONAL MATCH (:User)-[r:OWNS]->(org)
			DELETE r
			WITH org, t
			MATCH (u:User)-[:IS_LINKED_WITH {externalOwnerId:$userExternalOwnerId}]->(e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=$now`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":              tenant,
				"organizationId":      organizationId,
				"sourceOfTruth":       externalSystem,
				"externalSystemId":    externalSystem,
				"userExternalOwnerId": userExternalOwnerId,
				"now":                 utils.Now(),
			})
		return nil, err
	})
	return err
}
