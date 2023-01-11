package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type ContactRepository interface {
	MergeContact(tenant string, syncDate time.Time, contact entity.ContactData) (string, error)
	MergePrimaryEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergeAdditionalEmail(tenant, contactId, email, externalSystem string, createdAt time.Time) error
	MergePrimaryPhoneNumber(tenant, contactId, phoneNumber, externalSystem string, createdAt time.Time) error
	SetOwnerRelationship(tenant, contactId, userExternalId, externalSystemId string) error
	MergeTextCustomField(tenant, contactId string, field entity.TextCustomField, createdAt time.Time) error
	MergeContactAddress(tenant, contactId string, contact entity.ContactData) error
	MergeContactType(tenant, contactId, contactTypeName string) error
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
		" FOREACH (x in CASE WHEN c.sourceOfTruth <> $source THEN [c] ELSE [] END | " +
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
		" ON CREATE SET r.primary=true, e.id=randomUUID(), e.createdAt=$createdAt, e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, e:%s " +
		" ON MATCH SET r.primary=true"

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
		" ON CREATE SET r.primary=false, e.id=randomUUID(), e.createdAt=$createdAt, e.source=$source, e.sourceOfTruth=$sourceOfTruth, e.appSource=$appSource, e:%s " +
		" ON MATCH SET r.primary=false"

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

func (r *contactRepository) SetOwnerRelationship(tenant, contactId, userExternalId, externalSystemId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
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
		" ON CREATE SET f.textValue=$value, f.id=randomUUID(), f.createdAt=$createdAt, f.source=$source, f.sourceOfTruth=$sourceOfTruth, f.appSource=$appSource, f:%s " +
		" ON MATCH SET f.textValue=$value"

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
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactAddress(tenant, contactId string, contact entity.ContactData) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	// Create new Address if it does not exist with given source property
	// If Address exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateAddress for incoming source, with a new relationship 'ALTERNATE'
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:LOCATED_AT]->(a:Address {source:$source}) " +
		" ON CREATE SET a.id=randomUUID(), a.createdAt=$createdAt, a.source=$source, a.sourceOfTruth=$sourceOfTruth, a.appSource=$appSource, " +
		"				a.country=$country, a.state=$state, a.city=$city, a.address=$address, a.zip=$zip, a.fax=$fax, a:%s " +
		" ON MATCH SET 	" +
		"             a.country = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $country ELSE a.country END, " +
		"             a.state = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $state ELSE a.state END, " +
		"             a.city = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $city ELSE a.city END, " +
		"             a.address = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $address ELSE a.address END, " +
		"             a.zip = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $zip ELSE a.zip END, " +
		"             a.fax = CASE WHEN a.sourceOfTruth=$sourceOfTruth THEN $fax ELSE a.fax END " +
		" WITH a " +
		" FOREACH (x in CASE WHEN a.sourceOfTruth <> $source THEN [a] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateAddress {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, " +
		" alt.country=$country, alt.state=$state, alt.city=$city, alt.address=$address, alt.zip=$zip, alt.fax=$fax " +
		") "

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Address_"+tenant),
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
