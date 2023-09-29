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

type ContactRepository interface {
	GetMatchedContactId(ctx context.Context, tenant, primaryPhoneNumber string, contact entity.ContactData) (string, error)
	MergeContact(ctx context.Context, tenant string, syncDate time.Time, contact entity.ContactData) error
	MergePrimaryEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergeAdditionalEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergePhoneNumber(ctx context.Context, tenant, contactId, externalSystem string, createdAt time.Time, phoneNumber entity.PhoneNumber) error
	SetOwnerByOwnerExternalId(ctx context.Context, tenant, contactId, ownerExternalId, externalSystemId string) error
	SetOwnerByUserExternalId(ctx context.Context, tenant, contactId, userExternalId, externalSystemId string) error
	MergeTextCustomField(ctx context.Context, tenant, contactId string, field entity.TextCustomField) error
	MergeContactLocation(ctx context.Context, tenant, contactId string, contact entity.ContactData) error
	MergeTagForContact(ctx context.Context, tenant, contactId, tagName, sourceApp string) error
	LinkContactWithOrganizationByExternalId(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error
	LinkContactWithOrganizationByInternalId(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error
	LinkContactWithOrganizationByDomain(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error
	GetContactIdsForEmail(ctx context.Context, tenant, emailId string) ([]string, error)
	GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
	GetContactIdById(ctx context.Context, tenant, id string) (string, error)
	GetContactIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetJobRoleId(ctx context.Context, tenant, contactId, organizationId string) (string, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetMatchedContactId(ctx context.Context, tenant, primaryPhoneNumber string, contact entity.ContactData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetMatchedContactId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c1:Contact)-[:IS_LINKED_WITH {externalId:$contactExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c2:Contact),
						(c2)-[:HAS]->(e2:Email),
						(c2)-[:HAS]->(p2:PhoneNumber)
					WHERE e2.rawEmail in $emails AND size($emails) > 0 AND p2.rawPhoneNumber=$phoneNumber AND $phoneNumber <> ''
				OPTIONAL MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c3:Contact)-[:HAS]->(e3:Email)
					WHERE e3.rawEmail in $emails AND size($emails) > 0
				with coalesce(c1, c2, c3) as contacts
				where contacts is not null
				return contacts.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    contact.ExternalSystem,
				"contactExternalId": contact.ExternalId,
				"emails":            contact.EmailsForUnicity(),
				"phoneNumber":       primaryPhoneNumber,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	contactIDs := dbRecords.([]*db.Record)
	if len(contactIDs) == 1 {
		return contactIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *contactRepository) MergeContact(ctx context.Context, tenant string, syncDate time.Time, contact entity.ContactData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Contact if it does not exist
	// If Contact exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateContact for incoming source, with a new relationship 'ALTERNATE'
	// Link Contact with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET c.createdAt=$createdAt, " +
		"				c.updatedAt=$updatedAt, " +
		"				c.source=$source, " +
		"				c.sourceOfTruth=$sourceOfTruth, " +
		"				c.appSource=$appSource, " +
		"				c.firstName=$firstName, " +
		"				c.lastName=$lastName,  " +
		"				c.timezone=$timezone,  " +
		"				c.profilePhotoUrl=$profilePhotoUrl,  " +
		"				c.name=$name,  " +
		" 				c:%s " +
		" ON MATCH SET 	c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END, " +
		"				c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth  OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END, " +
		"				c.name = CASE WHEN c.sourceOfTruth=$sourceOfTruth  OR c.name is null OR c.name = '' THEN $name ELSE c.name END, " +
		"				c.timezone = CASE WHEN c.sourceOfTruth=$sourceOfTruth  OR c.timezone is null OR c.timezone = '' THEN $timezone ELSE c.timezone END, " +
		"				c.profilePhotoUrl = CASE WHEN c.sourceOfTruth=$sourceOfTruth  OR c.profilePhotoUrl is null OR c.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE c.profilePhotoUrl END, " +
		"				c.updatedAt = $now " +
		" WITH c, ext " +
		" MERGE (c)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH c " +
		" FOREACH (x in CASE WHEN c.sourceOfTruth <> $sourceOfTruth THEN [c] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateContact {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.firstName=$firstName, alt.lastName=$lastName, alt.name=$name, alt.timezone=$timezone, alt.profilePhotoUrl=$profilePhotoUrl " +
		" ) " +
		" RETURN c.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			query, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":          tenant,
				"contactId":       contact.Id,
				"externalSystem":  contact.ExternalSystem,
				"externalId":      contact.ExternalId,
				"externalUrl":     contact.ExternalUrl,
				"firstName":       contact.FirstName,
				"lastName":        contact.LastName,
				"name":            contact.Name,
				"timezone":        contact.Timezone,
				"profilePhotoUrl": contact.ProfilePhotoUrl,
				"syncDate":        syncDate,
				"createdAt":       utils.TimePtrFirstNonNilNillableAsAny(contact.CreatedAt),
				"updatedAt":       utils.TimePtrFirstNonNilNillableAsAny(contact.UpdatedAt),
				"source":          contact.ExternalSystem,
				"sourceOfTruth":   contact.ExternalSystem,
				"appSource":       constants.AppSourceSyncCustomerOsData,
				"now":             utils.Now(),
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

func (r *contactRepository) MergePrimaryEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergePrimaryEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" OPTIONAL MATCH (c)-[rel:HAS]->(:Email) " +
		" SET rel.primary=false " +
		" WITH DISTINCT c, t " +
		" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				e.id=randomUUID(), " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.appSource=$appSource, " +
		"				e:%s " +
		" WITH DISTINCT c, e " +
		" MERGE (c)-[rel:HAS]->(e) " +
		" ON CREATE SET rel.primary=true " +
		" ON MATCH SET rel.primary=true, e.updatedAt=$now "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"email":         email,
				"createdAt":     createdAt,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeAdditionalEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeAdditionalEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				e.id=randomUUID(), " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.appSource=$appSource, " +
		"				e:%s " +
		" WITH DISTINCT c, e " +
		" MERGE (c)-[rel:HAS]->(e) " +
		" ON CREATE SET rel.primary=false " +
		" ON MATCH SET rel.primary=false, e.updatedAt=$now "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"email":         email,
				"createdAt":     createdAt,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergePhoneNumber(ctx context.Context, tenant, contactId, externalSystem string, createdAt time.Time, phoneNumber entity.PhoneNumber) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergePhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		 MERGE (p:PhoneNumber {rawPhoneNumber: $phoneNumber})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) 
		 ON CREATE SET 
						p.id=randomUUID(), 
						p.createdAt=$now, 
						p.updatedAt=$now, 
						p.source=$source, 
						p.sourceOfTruth=$sourceOfTruth, 
						p.appSource=$appSource, 
						p:PhoneNumber_%s 
		 WITH DISTINCT c, p 
		 MERGE (c)-[rel:HAS]->(p) 
		 ON CREATE SET rel.primary=$primary, rel.label=$label 
		 ON MATCH SET 
			rel.label = CASE WHEN rel.label = '' OR rel.label IS NULL THEN $label ELSE rel.label END`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"phoneNumber":   phoneNumber.Number,
				"primary":       phoneNumber.Primary,
				"label":         phoneNumber.Label,
				"createdAt":     createdAt,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) SetOwnerByOwnerExternalId(ctx context.Context, tenant, contactId, ownerExternalId, externalSystemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.SetOwnerByOwnerExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (u:User)-[:IS_LINKED_WITH {externalIdSecond:$userExternalOwnerId}]->(e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[r:OWNS]->(c)
			return r`,
			map[string]interface{}{
				"tenant":              tenant,
				"contactId":           contactId,
				"externalSystemId":    externalSystemId,
				"userExternalOwnerId": ownerExternalId,
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

func (r *contactRepository) SetOwnerByUserExternalId(ctx context.Context, tenant, contactId, userExternalId, externalSystemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.SetOwnerByOwnerExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[r:OWNS]->(c)
			return r`,
			map[string]interface{}{
				"tenant":           tenant,
				"contactId":        contactId,
				"externalSystemId": externalSystemId,
				"userExternalId":   userExternalId,
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

func (r *contactRepository) MergeTextCustomField(ctx context.Context, tenant, contactId string, field entity.TextCustomField) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeTextCustomField")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:TextField:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) " +
		" ON CREATE SET f.textValue=$value, f.id=randomUUID(), f.createdAt=$createdAt, f.updatedAt=$createdAt, " +
		"				f.source=$source, f.sourceOfTruth=$sourceOfTruth, f.appSource=$appSource, f:%s " +
		" ON MATCH SET 	f.textValue = CASE WHEN f.sourceOfTruth=$sourceOfTruth THEN $value ELSE f.textValue END," +
		"				f.updatedAt = CASE WHEN f.sourceOfTruth=$sourceOfTruth THEN $now ELSE f.updatedAt END " +
		" WITH f " +
		" FOREACH (x in CASE WHEN f.sourceOfTruth <> $sourceOfTruth THEN [f] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateCustomField:AlternateTextField {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.textValue=$value " +
		" ) "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "CustomField_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"name":          field.Name,
				"value":         field.Value,
				"datatype":      "TEXT",
				"createdAt":     utils.TimePtrFirstNonNilNillableAsAny(field.CreatedAt),
				"source":        field.ExternalSystem,
				"sourceOfTruth": field.ExternalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactLocation(ctx context.Context, tenant, contactId string, contact entity.ContactData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeContactLocation")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location if it does not exist with given source property
	// If Location exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateLocation for incoming source, with a new relationship 'ALTERNATE'
	// !!! Current assumption - there is single Location with source of externalSystem per contact
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:ASSOCIATED_WITH]->(loc:Location {source:$source}) " +
		" ON CREATE SET " +
		"	loc.name=$locationName, " +
		"	loc.country=$country, " +
		"	loc.region=$region, " +
		"	loc.locality=$locality, " +
		"	loc.street=$street, " +
		"	loc.address=$address, " +
		"	loc.zip=$zip, " +
		"	loc.postalCode=$postalCode, " +
		"	loc.id=randomUUID(), " +
		"	loc.appSource=$appSource, " +
		"	loc.sourceOfTruth=$sourceOfTruth, " +
		"	loc.createdAt=$createdAt, " +
		"	loc.updatedAt=$createdAt, " +
		"	loc:%s " +
		" ON MATCH SET 	" +
		"             loc.country = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $country ELSE loc.country END, " +
		"             loc.region = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $region ELSE loc.region END, " +
		"             loc.locality = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $locality ELSE loc.locality END, " +
		"             loc.street = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $street ELSE loc.street END, " +
		"             loc.address = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $address ELSE loc.address END, " +
		"             loc.zip = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $zip ELSE loc.zip END, " +
		"             loc.postalCode = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $postalCode ELSE loc.postalCode END, " +
		"             loc.updatedAt = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH loc, t " +
		" MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t) " +
		" WITH loc " +
		" FOREACH (x in CASE WHEN loc.sourceOfTruth <> $sourceOfTruth THEN [loc] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateLocation {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.region=$region, alt.locality=$locality, alt.address=$address, " +
		" alt.zip=$zip, alt.postalCode=$postalCode, alt.street=$street " +
		") "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Location_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"country":       contact.Country,
				"region":        contact.Region,
				"locality":      contact.Locality,
				"street":        contact.Street,
				"address":       contact.Address,
				"zip":           contact.Zip,
				"postalCode":    contact.PostalCode,
				"createdAt":     utils.TimePtrFirstNonNilNillableAsAny(contact.CreatedAt),
				"source":        contact.ExternalSystem,
				"sourceOfTruth": contact.ExternalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"locationName":  contact.LocationName,
				"now":           utils.Now(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeTagForContact(ctx context.Context, tenant, contactId, tagName, source string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.MergeTagForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (tag:Tag {name:$tagName})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET tag.id=randomUUID(), " +
		"				tag.createdAt=$now, " +
		"				tag.updatedAt=$now, " +
		"				tag.source=$source," +
		"				tag.appSource=$source," +
		"				tag:%s  " +
		" WITH DISTINCT c, tag " +
		" MERGE (c)-[r:TAGGED]->(tag) " +
		"	ON CREATE SET r.taggedAt=$now " +
		" return r"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Tag_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"tagName":   tagName,
				"source":    source,
				"now":       utils.Now(),
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

func (r *contactRepository) LinkContactWithOrganizationByExternalId(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkContactWithOrganizationByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		 MATCH (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization) 
		 MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
		 ON CREATE SET j.id=randomUUID(), 
						j.source=$source, 
						j.jobTitle=$jobTitle, 
						j.sourceOfTruth=$sourceOfTruth, 
						j.appSource=$appSource, 
						j.createdAt=$now, 
						j.updatedAt=$now, 
						j:JobRole_%s 
		 ON MATCH SET j.jobTitle = CASE WHEN (j.sourceOfTruth=$sourceOfTruth AND $jobTitle <> '') OR j.jobTitle is null OR j.jobTitle = '' THEN $jobTitle ELSE j.jobTitle END`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"externalSystemId":       externalSystemId,
				"jobTitle":               org.JobTitle,
				"organizationExternalId": org.ExternalId,
				"now":                    utils.Now(),
				"source":                 externalSystemId,
				"sourceOfTruth":          externalSystemId,
				"appSource":              constants.AppSourceSyncCustomerOsData,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *contactRepository) LinkContactWithOrganizationByInternalId(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkContactWithOrganizationByInternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		 MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}) 
		 MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
		 ON CREATE SET j.id=randomUUID(), 
						j.source=$source, 
						j.jobTitle=$jobTitle, 
						j.sourceOfTruth=$sourceOfTruth, 
						j.appSource=$appSource, 
						j.createdAt=$now, 
						j.updatedAt=$now, 
						j:JobRole_%s 
		 ON MATCH SET j.jobTitle = CASE WHEN (j.sourceOfTruth=$sourceOfTruth AND $jobTitle <> '') OR j.jobTitle is null OR j.jobTitle = '' THEN $jobTitle ELSE j.jobTitle END`, tenant)
	span.LogFields(log.String("query", query))

	return utils.ExecuteQuery(ctx, *r.driver, query, map[string]interface{}{
		"tenant":         tenant,
		"contactId":      contactId,
		"organizationId": org.Id,
		"jobTitle":       org.JobTitle,
		"now":            utils.Now(),
		"source":         externalSystemId,
		"sourceOfTruth":  externalSystemId,
		"appSource":      constants.AppSourceSyncCustomerOsData,
	})
}

func (r *contactRepository) LinkContactWithOrganizationByDomain(ctx context.Context, tenant, contactId, externalSystemId string, org entity.ReferencedOrganization) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.LinkContactWithOrganizationByDomain")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		 MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain {domain:$domain})
		 MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
		 ON CREATE SET j.id=randomUUID(), 
						j.source=$source, 
						j.jobTitle=$jobTitle, 
						j.sourceOfTruth=$sourceOfTruth, 
						j.appSource=$appSource, 
						j.createdAt=$now, 
						j.updatedAt=$now, 
						j:JobRole_%s 
		ON MATCH SET j.jobTitle = CASE WHEN (j.sourceOfTruth=$sourceOfTruth AND $jobTitle <> '') OR j.jobTitle is null OR j.jobTitle = '' THEN $jobTitle ELSE j.jobTitle END`, tenant)
	span.LogFields(log.String("query", query))

	return utils.ExecuteQuery(ctx, *r.driver, query, map[string]interface{}{
		"tenant":        tenant,
		"contactId":     contactId,
		"domain":        org.Domain,
		"jobTitle":      org.JobTitle,
		"now":           utils.Now(),
		"source":        externalSystemId,
		"sourceOfTruth": externalSystemId,
		"appSource":     constants.AppSourceSyncCustomerOsData,
	})
}

func (r *contactRepository) GetContactIdsForEmail(ctx context.Context, tenant, emailId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdsForEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:HAS]->(:Email {id:$emailId})
		RETURN c.id`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":  tenant,
				"emailId": emailId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return []string{}, err
	}
	contactIDs := make([]string, 0)
	for _, v := range dbRecords.([]*db.Record) {
		contactIDs = append(contactIDs, v.Values[0].(string))
	}
	return contactIDs, nil
}

func (r *contactRepository) GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant)
 			WHERE (c.syncedWithEventStore is null or c.syncedWithEventStore=false)
			RETURN c, t.name limit $size`,
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

func (r *contactRepository) GetContactIdById(ctx context.Context, tenant, id string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				return c.id order by c.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":    tenant,
			"contactId": id,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *contactRepository) GetContactIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetContactIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return c.id order by c.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *contactRepository) GetJobRoleId(ctx context.Context, tenant, contactId, organizationId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactRepository.GetJobRoleId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})
				MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
				MATCH (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) 
				return j.id order by j.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":         tenant,
			"contactId":      contactId,
			"organizationId": organizationId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
