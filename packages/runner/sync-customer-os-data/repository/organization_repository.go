package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type OrganizationRepository interface {
	GetMatchedOrganizationId(ctx context.Context, tenant string, organization entity.OrganizationData) (string, error)
	MergeOrganization(ctx context.Context, tenant string, syncDate time.Time, organization entity.OrganizationData) error
	MergeOrganizationRelationshipAndStage(ctx context.Context, tenant, organizationId, relationship, stage, externalSystem string) error
	MergeOrganizationLocation(ctx context.Context, tenant, organizationId string, organization entity.OrganizationData) error
	MergeOrganizationDomain(ctx context.Context, tenant, organizationId, domain, externalSystem string) error
	MergePhoneNumber(ctx context.Context, tenant, organizationId, phoneNumber, externalSystem string, createdAt time.Time) error
	MergeEmail(ctx context.Context, tenant, organizationId, email, externalSystem string, createdAt time.Time) error
	LinkToParentOrganizationAsSubsidiary(ctx context.Context, tenant, organizationId, externalSystem string, parentOrganizationDtls *entity.ParentOrganization) error
	SetOwnerByOwnerExternalId(ctx context.Context, tenant, contactId, userExternalOwnerId, externalSystem string) error
	SetOwnerByUserExternalId(ctx context.Context, tenant, contactId, userExternalId, externalSystem string) error
	CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error)
	UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error
	GetOrganizationIdsForContact(ctx context.Context, tenant, contactId string) ([]string, error)
	GetOrganizationIdsForContactByExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string) ([]string, error)
	GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
	GetAllDomainLinksCrossTenantsNotSynced(ctx context.Context, size int) ([]*neo4j.Record, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetMatchedOrganizationId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
		"				org.subIndustry=$subIndustry, " +
		"				org.industryGroup=$industryGroup, " +
		"				org.targetAudience=$targetAudience, " +
		"				org.valueProposition=$valueProposition, " +
		"				org.lastFundingRound=$lastFundingRound, " +
		"				org.lastFundingAmount=$lastFundingAmount, " +
		"				org.market=$market, " +
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
		"				org.subIndustry = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.subIndustry is null OR org.subIndustry = '' THEN $subIndustry ELSE org.subIndustry END, " +
		"				org.industryGroup = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.industryGroup is null OR org.industryGroup = '' THEN $industryGroup ELSE org.industryGroup END, " +
		"				org.targetAudience = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.targetAudience is null OR org.targetAudience = '' THEN $targetAudience ELSE org.targetAudience END, " +
		"				org.valueProposition = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.valueProposition is null OR org.valueProposition = '' THEN $valueProposition ELSE org.valueProposition END, " +
		"				org.lastFundingRound = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.lastFundingRound is null OR org.lastFundingRound = '' THEN $lastFundingRound ELSE org.lastFundingRound END, " +
		"				org.lastFundingAmount = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.lastFundingAmount is null OR org.lastFundingAmount = '' THEN $lastFundingAmount ELSE org.lastFundingAmount END, " +
		"				org.market = CASE WHEN org.sourceOfTruth=$sourceOfTruth OR org.market is null OR org.market = '' THEN $market ELSE org.market END, " +
		"				org.isPublic = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $isPublic ELSE org.isPublic END, " +
		"				org.employees = CASE WHEN org.sourceOfTruth=$sourceOfTruth THEN $employees ELSE org.employees END, " +
		"				org.updatedAt = $now " +
		" WITH org, ext " +
		" MERGE (org)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl, r.externalSource=$externalSource " +
		" ON MATCH SET r.syncDate=$syncDate, r.externalSource=$externalSource " +
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
				"tenant":            tenant,
				"orgId":             organization.Id,
				"externalSystem":    organization.ExternalSystem,
				"externalId":        organization.ExternalId,
				"externalUrl":       organization.ExternalUrl,
				"syncDate":          syncDate,
				"name":              organization.Name,
				"description":       organization.Description,
				"createdAt":         utils.TimePtrFirstNonNilNillableAsAny(organization.CreatedAt),
				"updatedAt":         utils.TimePtrFirstNonNilNillableAsAny(organization.UpdatedAt),
				"website":           organization.Website,
				"industry":          organization.Industry,
				"subIndustry":       organization.SubIndustry,
				"industryGroup":     organization.IndustryGroup,
				"targetAudience":    organization.TargetAudience,
				"valueProposition":  organization.ValueProposition,
				"lastFundingRound":  organization.LastFundingRound,
				"lastFundingAmount": organization.LastFundingAmount,
				"market":            organization.Market,
				"isPublic":          organization.IsPublic,
				"employees":         organization.Employees,
				"source":            organization.ExternalSystem,
				"sourceOfTruth":     organization.ExternalSystem,
				"appSource":         constants.AppSourceSyncCustomerOsData,
				"externalSource":    organization.ExternalSourceTable,
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

func (r *organizationRepository) MergeOrganizationRelationshipAndStage(ctx context.Context, tenant, organizationId, relationship, stage, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationRelationshipAndStage")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
				WHERE org.sourceOfTruth=$sourceOfTruth			
				WITH org, t
		 		MATCH (or:OrganizationRelationship {name:$relationship})
		 		MERGE (org)-[:IS]->(or) 
				SET org.updatedAt=$now
				WITH org, or, t
				MATCH (t)<-[:STAGE_BELONGS_TO_TENANT]-(os:OrganizationRelationshipStage {name:$stage})<-[:HAS_STAGE]-(or)
				WHERE NOT (org)-[:HAS_STAGE]->(:OrganizationRelationshipStage)<-[:HAS_STAGE]-(or)
				MERGE (org)-[:HAS_STAGE]->(os)`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"organizationId": organizationId,
				"relationship":   relationship,
				"stage":          stage,
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationLocation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location if it does not exist with given source property
	// If Location exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateLocation for incoming source, with a new relationship 'ALTERNATE'
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
				"appSource":      constants.AppSourceSyncCustomerOsData,
				"locationName":   organization.LocationName,
				"createdAt":      utils.TimePtrFirstNonNilNillableAsAny(organization.CreatedAt),
				"now":            utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) MergeOrganizationDomain(ctx context.Context, tenant string, organizationId string, domain string, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeOrganizationDomain")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MERGE (d:Domain {domain:toLower($domain)}) 
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
				"appSource":      constants.AppSourceSyncCustomerOsData,
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergePhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
				"appSource":      constants.AppSourceSyncCustomerOsData,
				"now":            utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) MergeEmail(ctx context.Context, tenant, organizationId, email, externalSystem string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.MergeEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
				"appSource":      constants.AppSourceSyncCustomerOsData,
				"now":            utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) LinkToParentOrganizationAsSubsidiary(ctx context.Context, tenant, organizationId, externalSystem string, parentOrganizationDtls *entity.ParentOrganization) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.LinkToParentOrganizationAsSubsidiary")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *organizationRepository) SetOwnerByOwnerExternalId(ctx context.Context, tenant, organizationId, userExternalOwnerId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetOwnerByOwnerExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			WHERE org.sourceOfTruth=$sourceOfTruth
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

func (r *organizationRepository) SetOwnerByUserExternalId(ctx context.Context, tenant, organizationId, userExternalId, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.SetOwnerByOwnerExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			WHERE org.sourceOfTruth=$sourceOfTruth
			WITH org, t
			OPTIONAL MATCH (:User)-[r:OWNS]->(org)
			DELETE r
			WITH org, t
			MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[:OWNS]->(org)
			SET org.updatedAt=$now`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":           tenant,
				"organizationId":   organizationId,
				"sourceOfTruth":    externalSystem,
				"externalSystemId": externalSystem,
				"userExternalId":   userExternalId,
				"now":              utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *organizationRepository) CalculateAndGetLastTouchpoint(ctx context.Context, tenant string, organizationId string) (*time.Time, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CalculateAndGetLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	params := map[string]any{
		"tenant":                             tenant,
		"now":                                utils.Now(),
		"organizationId":                     organizationId,
		"nodeLabels":                         []string{"InteractionSession", "Issue", "Conversation", "InteractionEvent", "Meeting"},
		"excludeInteractionEventContentType": []string{"x-openline-transcript-element"},
		"contactRelationTypes":               []string{"HAS_ACTION", "PARTICIPATES", "SENT_TO", "SENT_BY", "PART_OF", "REPORTED_BY", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
		"organizationRelationTypes":          []string{"REPORTED_BY", "SENT_TO", "SENT_BY"},
		"emailAndPhoneRelationTypes":         []string{"SENT_TO", "SENT_BY", "PART_OF", "DESCRIBES", "ATTENDED_BY", "CREATED_BY"},
	}

	query := `MATCH (o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		CALL { ` +
		// get all timeline events for the organization contatcs
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact), 
		p = (c)-[*1..2]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $contactRelationTypes)
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events directly for the organization
		` WITH o MATCH (o), 
		p = (o)-[*1]-(a:TimelineEvent) 
		WHERE all(r IN relationships(p) WHERE type(r) in $organizationRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 
		AND (NOT "InteractionEvent" in labels(a) or "InteractionEvent" in labels(a) AND NOT a.contentType IN $excludeInteractionEventContentType)
		AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization contacts' emails and phone numbers
		` WITH o MATCH (o)<-[:ROLE_IN]-(j:JobRole)<-[:WORKS_AS]-(c:Contact)-[:HAS]->(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $emailAndPhoneRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
		RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		UNION ` +
		// get all timeline events for the organization emails and phone numbers
		` WITH o MATCH (o)-[:HAS]->(e), 
		p = (e)-[*1..2]-(a:TimelineEvent) 
		WHERE ('Email' in labels(e) OR 'PhoneNumber' in labels(e)) 
		AND all(r IN relationships(p) WHERE type(r) in $emailAndPhoneRelationTypes)
		AND size([label IN labels(a) WHERE label IN $nodeLabels | 1]) > 0 AND coalesce(a.startedAt, a.updatedAt, a.createdAt) <= $now
	 	RETURN a as timelineEvent ORDER BY coalesce(a.startedAt, a.updatedAt, a.createdAt) DESC LIMIT 1 
		} 
		RETURN coalesce(timelineEvent.startedAt, timelineEvent.updatedAt, timelineEvent.createdAt), timelineEvent.id ORDER BY coalesce(timelineEvent.startedAt, timelineEvent.updatedAt, timelineEvent.createdAt) DESC LIMIT 1`

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, "", err
	}

	if len(records.([]*neo4j.Record)) > 0 {
		return utils.TimePtr(records.([]*neo4j.Record)[0].Values[0].(time.Time)), records.([]*neo4j.Record)[0].Values[1].(string), nil
	}
	return nil, "", nil
}

func (r *organizationRepository) UpdateLastTouchpoint(ctx context.Context, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.UpdateLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId`

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

func (r *organizationRepository) GetOrganizationIdsForContact(ctx context.Context, tenant, contactId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdsForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)--(:JobRole)--(:Contact {id:$contactId})
		RETURN org.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	orgIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		orgIDs = append(orgIDs, v.Values[0].(string))
	}
	return orgIDs, nil
}

func (r *organizationRepository) GetOrganizationIdsForContactByExternalId(ctx context.Context, tenant, contactExternalId, externalSystem string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationIdsForContactByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)--(:JobRole)--(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		RETURN org.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"contactExternalId": contactExternalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	orgIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		orgIDs = append(orgIDs, v.Values[0].(string))
	}
	return orgIDs, nil
}

func (r *organizationRepository) GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant)
 			WHERE (org.syncedWithEventStore is null or org.syncedWithEventStore=false)
			RETURN org, t.name limit $size`,
			map[string]any{
				"size": size,
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

func (r *organizationRepository) GetAllDomainLinksCrossTenantsNotSynced(ctx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[rel:HAS_DOMAIN]->(d:Domain)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false) AND org.syncedWithEventStore=true AND d.domain <> "" 
			RETURN org.id, t.name, d.domain limit $size`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
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
	return result.([]*neo4j.Record), err
}
