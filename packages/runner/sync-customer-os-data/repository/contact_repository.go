package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type ContactRepository interface {
	MergeContact(tenant string, syncDate time.Time, contact entity.ContactData) (string, error)
	MergePrimaryEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergeAdditionalEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergePrimaryPhoneNumber(tenant, contactId, phoneNumber, externalSystem string, createdAt time.Time) error
	SetOwnerRelationship(tenant, contactId, userExternalOwnerId, externalSystemId string) error
	MergeTextCustomField(tenant, contactId string, field entity.TextCustomField, createdAt time.Time) error
	MergeContactDefaultPlace(tenant, contactId string, contact entity.ContactData) error
	MergeContactType(tenant, contactId, contactTypeName string) error
	GetOrCreateContactId(tenant, email, firstName, lastName, source string) (string, error)
}

type contactRepository struct {
	driver *neo4j.Driver
}

func NewContactRepository(driver *neo4j.Driver) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) MergeContact(tenant string, syncDate time.Time, contact entity.ContactData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	// Create new Contact if it does not exist
	// If Contact exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateContact for incoming source, with a new relationship 'ALTERNATE'
	// Link Contact with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (c:Contact)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET r.externalId=$externalId, r.syncDate=$syncDate, c.id=randomUUID(), c.createdAt=$createdAt, " +
		"				c.source=$source, c.sourceOfTruth=$sourceOfTruth, c.appSource=$appSource, " +
		"				c.firstName=$firstName, c.lastName=$lastName,  " +
		" 				c:%s " +
		" ON MATCH SET 	r.syncDate = CASE WHEN c.sourceOfTruth=$sourceOfTruth THEN $syncDate ELSE r.syncDate END, " +
		"				c.firstName = CASE WHEN c.sourceOfTruth=$sourceOfTruth THEN $firstName ELSE c.firstName END, " +
		"				c.lastName = CASE WHEN c.sourceOfTruth=$sourceOfTruth THEN $lastName ELSE c.lastName END " +
		" WITH c, t " +
		" MERGE (c)-[:CONTACT_BELONGS_TO_TENANT]->(t) " +
		" WITH c " +
		" FOREACH (x in CASE WHEN c.sourceOfTruth <> $sourceOfTruth THEN [c] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateContact {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.firstName=$firstName, alt.lastName=$lastName " +
		" ) " +
		" RETURN c.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(
			query, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": contact.ExternalSystem,
				"externalId":     contact.ExternalId,
				"firstName":      contact.FirstName,
				"lastName":       contact.LastName,
				"syncDate":       syncDate,
				"createdAt":      contact.CreatedAt,
				"source":         contact.ExternalSystem,
				"sourceOfTruth":  contact.ExternalSystem,
				"appSource":      contact.ExternalSystem,
				"now":            time.Now().UTC(),
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

func (r *contactRepository) MergePrimaryEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" OPTIONAL MATCH (c)-[r:EMAILED_AT]->(e:Email) " +
		" SET r.primary=false " +
		" WITH c " +
		" MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
		" ON CREATE SET r.primary=true, e.id=randomUUID(), e.createdAt=$createdAt, e.updatedAt=$createdAt, e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, e:%s " +
		" ON MATCH SET r.primary=true, e.updatedAt=$now "

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
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

func (r *contactRepository) MergeAdditionalEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[r:EMAILED_AT]->(e:Email {email:$email}) " +
		" ON CREATE SET r.primary=false, e.id=randomUUID(), e.createdAt=$createdAt, e.updatedAt=$createdAt, e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, e:%s " +
		" ON MATCH SET r.primary=false, e.updatedAt=$now "

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
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

func (r *contactRepository) MergePrimaryPhoneNumber(tenant, contactId, e164, externalSystem string, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" OPTIONAL MATCH (c)-[r:CALLED_AT]->(p:PhoneNumber) " +
		" SET r.primary=false " +
		" WITH c " +
		" MERGE (c)-[r:CALLED_AT]->(p:PhoneNumber {e164: $e164}) " +
		" ON CREATE SET r.primary=true, p.id=randomUUID(), p.createdAt=$createdAt, p.source=$source, p.sourceOfTruth=$sourceOfTruth, p.appSource=$appSource, p:%s " +
		" ON MATCH SET r.primary=true"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "PhoneNumber_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"e164":          e164,
				"createdAt":     createdAt,
				"source":        externalSystem,
				"sourceOfTruth": externalSystem,
				"appSource":     externalSystem,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) SetOwnerRelationship(tenant, contactId, userExternalOwnerId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
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
		_, err = queryResult.Single()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *contactRepository) MergeTextCustomField(tenant, contactId string, field entity.TextCustomField, createdAt time.Time) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:TextField:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) " +
		" ON CREATE SET f.textValue=$value, f.id=randomUUID(), f.createdAt=$createdAt, " +
		"				f.source=$source, f.sourceOfTruth=$sourceOfTruth, f.appSource=$appSource, f:%s " +
		" ON MATCH SET 	f.textValue = CASE WHEN f.sourceOfTruth=$sourceOfTruth THEN $value ELSE f.textValue END " +
		" WITH f " +
		" FOREACH (x in CASE WHEN f.sourceOfTruth <> $sourceOfTruth THEN [f] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateCustomField:AlternateTextField {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.textValue=$value " +
		" ) "

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "CustomField_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"name":          field.Name,
				"value":         field.Value,
				"datatype":      "TEXT",
				"createdAt":     createdAt,
				"source":        field.ExternalSystem,
				"sourceOfTruth": field.ExternalSystem,
				"appSource":     field.ExternalSystem,
				"now":           time.Now().UTC(),
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactDefaultPlace(tenant, contactId string, contact entity.ContactData) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	// Create new Place and Location if it does not exist with given source property
	// If Place exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternatePlace for incoming source, with a new relationship 'ALTERNATE'
	// !!! Current assumption - there is single Location and place with source of externalSystem
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:ASSOCIATED_WITH]->(loc:Location {source:$source})-[:LOCATED_AT]->(p:Place {source:$source}) " +
		" ON CREATE SET p.id=randomUUID(), p.createdAt=$createdAt, p.updatedAt=$createdAt, p.sourceOfTruth=$sourceOfTruth, p.appSource=$appSource, " +
		"				p.country=$country, p.state=$state, p.city=$city, p.address=$address, p.zip=$zip, p.fax=$fax, p:%s," +
		"				loc.id=randomUUID(), loc.appSource=$appSource, loc.sourceOfTruth=$sourceOfTruth, loc.name=$locationName, " +
		"				loc.createdAt=$createdAt, loc.updatedAt=$createdAt, loc:%s " +
		" ON MATCH SET 	" +
		"             p.country = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $country ELSE p.country END, " +
		"             p.state = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $state ELSE p.state END, " +
		"             p.city = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $city ELSE p.city END, " +
		"             p.address = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $address ELSE p.address END, " +
		"             p.zip = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $zip ELSE p.zip END, " +
		"             p.fax = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $fax ELSE p.fax END, " +
		"             p.updatedAt = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $now ELSE p.updatedAt END, " +
		"             loc.updatedAt = CASE WHEN p.sourceOfTruth=$sourceOfTruth THEN $now ELSE loc.updatedAt END " +
		" WITH p " +
		" FOREACH (x in CASE WHEN p.sourceOfTruth <> $sourceOfTruth THEN [p] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternatePlace {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.state=$state, alt.city=$city, alt.address=$address, alt.zip=$zip, alt.fax=$fax " +
		") "

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Place_"+tenant, "Location_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"contactId":     contactId,
				"country":       contact.Country,
				"state":         contact.State,
				"city":          contact.City,
				"address":       contact.Address,
				"zip":           contact.Zip,
				"fax":           contact.Fax,
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

func (r *contactRepository) MergeContactType(tenant, contactId, contactTypeName string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (ct:ContactType {name:$contactTypeName})-[:CONTACT_TYPE_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET ct.id=randomUUID() " +
		" WITH c, ct " +
		" MERGE (c)-[r:IS_OF_TYPE]->(ct) " +
		" return r"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(query,
			map[string]interface{}{
				"tenant":          tenant,
				"contactId":       contactId,
				"contactTypeName": contactTypeName,
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

func (r *contactRepository) GetOrCreateContactId(tenant, email, firstName, lastName, source string) (string, error) {
	session := (*r.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	record, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(fmt.Sprintf(
			" MATCH (t:Tenant {name:$tenant}) "+
				" MERGE (e:Email {email: $email})<-[r:EMAILED_AT]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET r.primary=true, e.id=randomUUID(), e.createdAt=$createdAt, e.updatedAt=$createdAt, "+
				"				e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, "+
				"				c.id=randomUUID(), c.firstName=$firstName, c.lastName=$lastName, "+
				"				c.createdAt=$createdAt, c.updatedAt=$createdAt, "+
				"				c.source=$source, c.sourceOfTruth=$sourceOfTruth, c.appSource=$appSource, "+
				"               c:%s, e:%s "+
				" RETURN c.id", "Contact_"+tenant, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"email":         email,
				"firstName":     firstName,
				"lastName":      lastName,
				"source":        source,
				"sourceOfTruth": source,
				"appSource":     source,
				"createdAt":     time.Now().UTC(),
			})
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})

	return record.(*db.Record).Values[0].(string), err
}
