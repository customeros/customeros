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

type ContactRepository interface {
	GetMatchedContactId(ctx context.Context, tenant string, contact entity.ContactData) (string, error)
	MergeContact(ctx context.Context, tenant string, syncDate time.Time, contact entity.ContactData) error
	MergePrimaryEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergeAdditionalEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergePrimaryPhoneNumber(ctx context.Context, tenant, contactId, phoneNumber, externalSystem string, createdAt time.Time) error
	SetOwnerRelationship(ctx context.Context, tenant, contactId, userExternalOwnerId, externalSystemId string) error
	MergeTextCustomField(ctx context.Context, tenant, contactId string, field entity.TextCustomField) error
	MergeContactDefaultPlace(ctx context.Context, tenant, contactId string, contact entity.ContactData) error
	MergeTagForContact(ctx context.Context, tenant, contactId, tagName, sourceApp string) error
	LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationExternalId, source string) error
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetMatchedContactId(ctx context.Context, tenant string, contact entity.ContactData) (string, error) {
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
				"phoneNumber":       contact.PhoneNumber,
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
		" 				c:%s " +
		" ON MATCH SET 	c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR c.firstName is null OR c.firstName = '' THEN $firstName ELSE c.firstName END, " +
		"				c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth  OR c.lastName is null OR c.lastName = '' THEN $lastName ELSE c.lastName END, " +
		"				c.updatedAt = $now " +
		" WITH c, ext " +
		" MERGE (c)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH c " +
		" FOREACH (x in CASE WHEN c.sourceOfTruth <> $sourceOfTruth THEN [c] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateContact {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.firstName=$firstName, alt.lastName=$lastName " +
		" ) " +
		" RETURN c.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			query, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"contactId":      contact.Id,
				"externalSystem": contact.ExternalSystem,
				"externalId":     contact.ExternalId,
				"externalUrl":    contact.ExternalUrl,
				"firstName":      contact.FirstName,
				"lastName":       contact.LastName,
				"syncDate":       syncDate,
				"createdAt":      contact.CreatedAt,
				"updatedAt":      contact.UpdatedAt,
				"source":         contact.ExternalSystem,
				"sourceOfTruth":  contact.ExternalSystem,
				"appSource":      contact.ExternalSystem,
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

func (r *contactRepository) MergePrimaryEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error {
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
				"appSource":     externalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeAdditionalEmail(ctx context.Context, tenant, contactId, email, externalSystem string, createdAt time.Time) error {
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
				"appSource":     externalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergePrimaryPhoneNumber(ctx context.Context, tenant, contactId, phoneNumber, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" OPTIONAL MATCH (c)-[rel:HAS]->(p:PhoneNumber) " +
		" SET rel.primary=false " +
		" WITH DISTINCT c, t " +
		" MERGE (p:PhoneNumber {rawPhoneNumber: $phoneNumber})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				p.id=randomUUID(), " +
		"				p.createdAt=$now, " +
		"				p.updatedAt=$now, " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.appSource=$appSource, " +
		"				p:%s " +
		" WITH DISTINCT c, p " +
		" MERGE (c)-[rel:HAS]->(p) " +
		" ON CREATE SET rel.primary=true, p.updatedAt=$now, c.updatedAt=$now " +
		" ON MATCH SET rel.primary=true, c.updatedAt=$now "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "PhoneNumber_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"phoneNumber":   phoneNumber,
				"createdAt":     createdAt,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     externalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) SetOwnerRelationship(ctx context.Context, tenant, contactId, userExternalOwnerId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (u:User)-[:IS_LINKED_WITH {externalOwnerId:$userExternalOwnerId}]->(e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			MERGE (u)-[r:OWNS]->(c)
			return r`,
			map[string]interface{}{
				"tenant":              tenant,
				"contactId":           contactId,
				"externalSystemId":    externalSystemId,
				"userExternalOwnerId": userExternalOwnerId,
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
				"createdAt":     field.CreatedAt,
				"source":        field.ExternalSystem,
				"sourceOfTruth": field.ExternalSystem,
				"appSource":     field.ExternalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactDefaultPlace(ctx context.Context, tenant, contactId string, contact entity.ContactData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Location and Location if it does not exist with given source property and namd
	// If Location exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateLocation for incoming source, with a new relationship 'ALTERNATE'
	// !!! Current assumption - there is single Location with source of externalSystem and name per contact
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:ASSOCIATED_WITH]->(loc:Location {source:$source, name:$locationName}) " +
		" ON CREATE SET " +
		"	loc.country=$country, " +
		"	loc.region=$region, " +
		"	loc.locality=$locality, " +
		"	loc.address=$address, " +
		"	loc.zip=$zip, " +
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
		"             loc.address = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $address ELSE loc.address END, " +
		"             loc.zip = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $zip ELSE loc.zip END, " +
		"             loc.updatedAt = CASE WHEN loc.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH loc, t " +
		" MERGE (loc)-[:LOCATION_BELONGS_TO_TENANT]->(t) " +
		" WITH loc " +
		" FOREACH (x in CASE WHEN loc.sourceOfTruth <> $sourceOfTruth THEN [loc] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternatePlace {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.region=$region, alt.locality=$locality, alt.address=$address, alt.zip=$zip " +
		") "

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Location_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"country":       contact.Country,
				"region":        contact.Region,
				"locality":      contact.Locality,
				"address":       contact.Address,
				"zip":           contact.Zip,
				"createdAt":     contact.CreatedAt,
				"source":        contact.ExternalSystem,
				"sourceOfTruth": contact.ExternalSystem,
				"appSource":     contact.ExternalSystem,
				"locationName":  contact.DefaultLocationName,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeTagForContact(ctx context.Context, tenant, contactId, tagName, source string) error {
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
				"now":       time.Now().UTC(),
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

func (r *contactRepository) LinkContactWithOrganization(ctx context.Context, tenant, contactId, organizationExternalId, source string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MATCH (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization) " +
		" MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) " +
		" ON CREATE SET j.id=randomUUID(), " +
		"				j.source=$source, " +
		"				j.sourceOfTruth=$sourceOfTruth, " +
		"				j.appSource=$appSource, " +
		"				j.createdAt=$now, " +
		"				j.updatedAt=$now, " +
		"				j:%s "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "JobRole_"+tenant),
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"externalSystemId":       source,
				"organizationExternalId": organizationExternalId,
				"now":                    time.Now().UTC(),
				"source":                 source,
				"sourceOfTruth":          source,
				"appSource":              source,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// TODO implement removing outdated linked companies
