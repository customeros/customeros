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
	MergePrimaryEmail(tenant, contactId, email string) error
	MergeAdditionalEmail(tenant, contactId, email string) error
	MergePrimaryPhoneNumber(tenant, contactId, phoneNumber string) error
	SetOwnerRelationship(tenant, contactId, userExternalId, externalSystemId string) error
	MergeTextCustomField(tenant, contactId string, field entity.TextCustomField) error
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

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(
			"MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) "+
				" MERGE (c:Contact)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) "+
				" ON CREATE SET r.externalId=$externalId, c.id=randomUUID(), "+
				"				c.firstName=$firstName, c.lastName=$lastName, r.syncDate=$syncDate, c.readonly=$readonly, "+
				" 				c.createdAt=$createdAt, c:%s "+
				" ON MATCH SET 	c.firstName=$firstName, c.lastName=$lastName, r.syncDate=$syncDate, c.readonly=$readonly "+
				" WITH c, t "+
				" MERGE (c)-[:CONTACT_BELONGS_TO_TENANT]->(t) "+
				" RETURN c.id", "Contact_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": contact.ExternalSystem,
				"externalId":     contact.ExternalId,
				"firstName":      contact.FirstName,
				"lastName":       contact.LastName,
				"syncDate":       syncDate,
				"readonly":       contact.Readonly,
				"createdAt":      contact.CreatedAt,
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

func (r *contactRepository) MergePrimaryEmail(tenant, contactId, email string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" OPTIONAL MATCH (c)-[r:EMAILED_AT]->(e:Email) " +
		" SET r.primary=false " +
		" WITH c " +
		" MERGE (c)-[r:EMAILED_AT]->(e:Email {email: $email}) " +
		" ON CREATE SET r.primary=true, e.id=randomUUID(), e:%s " +
		" ON MATCH SET r.primary=true"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"email":     email,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeAdditionalEmail(tenant, contactId, email string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[r:EMAILED_AT]->(e:Email {email:$email}) " +
		" ON CREATE SET r.primary=false, e.id=randomUUID(), e:%s " +
		" ON MATCH SET r.primary=false"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"email":     email,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergePrimaryPhoneNumber(tenant, contactId, e164 string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[r:CALLED_AT]->(p:PhoneNumber)
			SET r.primary=false
			WITH c
			MERGE (c)-[r:CALLED_AT]->(p:PhoneNumber {e164: $e164})
            ON CREATE SET r.primary=true, p.id=randomUUID()
            ON MATCH SET r.primary=true`,
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"e164":      e164,
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

func (r *contactRepository) MergeTextCustomField(tenant, contactId string, field entity.TextCustomField) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
			MERGE (f:TextField:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) 
			ON CREATE SET f.textValue=$value, f.id=randomUUID(), f.source=$source
			ON MATCH SET f.textValue=$value, f.source=$source
			`,
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"name":      field.Name,
				"value":     field.Value,
				"source":    field.Source,
				"datatype":  "TEXT",
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactAddress(tenant, contactId string, contact entity.ContactData) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:LOCATED_AT]->(a:Address {source:$source}) " +
		" ON CREATE SET a.id=randomUUID(), a.source=$source, " +
		"	a.country=$country, a.state=$state, a.city=$city, a.address=$address, a.zip=$zip, a.fax=$fax, a:%s " +
		" ON MATCH SET 	" +
		"   a.country=$country, a.state=$state, a.city=$city, a.address=$address, a.zip=$zip, a.fax=$fax"

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Address_"+tenant),
			map[string]interface{}{
				"tenant":    tenant,
				"contactId": contactId,
				"source":    contact.ExternalSystem,
				"country":   contact.Country,
				"state":     contact.State,
				"city":      contact.City,
				"address":   contact.Address,
				"zip":       contact.Zip,
				"fax":       contact.Fax,
			})
		return nil, err
	})
	return err
}

func (r *contactRepository) MergeContactType(tenant, contactId, contactTypeName string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (ct:ContactType {name:$contactTypeName})-[:CONTACT_TYPE_BELONGS_TO_TENANT]->(t)
			MERGE (c)-[r:IS_OF_TYPE]->(ct)
			return r`,
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
