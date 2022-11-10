package resolver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/integration_tests"
)

func cleanupAllData(driver *neo4j.Driver) {
	integration_tests.ExecuteWriteQuery(driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func createTenant(driver *neo4j.Driver, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant": tenant,
	})
}

func createUser(driver *neo4j.Driver, tenant string, user entity.UserEntity) {
	query := `
		MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {
				  id: randomUUID(),
				  firstName: $firstName,
				  lastName: $lastName,
				  email: $email,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:USER_BELONGS_TO_TENANT]->(t)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":    tenant,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
}

func createDefaultContact(driver *neo4j.Driver, tenant string) string {
	return createContact(driver, tenant, entity.ContactEntity{Title: "MR", FirstName: "first", LastName: "last"})
}

func createContact(driver *neo4j.Driver, tenant string, contact entity.ContactEntity) string {
	var contactId, _ = uuid.NewRandom()
	query := `
			MATCH (t:Tenant {name:$tenant})
			MERGE (c:Contact {
				  id: $contactId,
				  title: $title,
				  firstName: $firstName,
				  lastName: $lastName,
				  label: $label,
				  notes: $notes,
				  contactType: $contactType,
				  createdAt :datetime({timezone: 'UTC'})
				})-[:CONTACT_BELONGS_TO_TENANT]->(t)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"tenant":      tenant,
		"contactId":   contactId.String(),
		"title":       contact.Title,
		"firstName":   contact.FirstName,
		"lastName":    contact.LastName,
		"contactType": contact.ContactType,
		"notes":       contact.Notes,
		"label":       contact.Label,
	})
	return contactId.String()
}

func createDefaultFieldSet(driver *neo4j.Driver, contactId string) string {
	return createFieldSet(driver, contactId, entity.FieldSetEntity{Name: "name", Type: "type"})
}

func createFieldSet(driver *neo4j.Driver, contactId string, fieldSet entity.FieldSetEntity) string {
	var fieldSetId, _ = uuid.NewRandom()
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (s:FieldSet {
				  id: $fieldSetId,
				  type: $type,
				  name: $name
				})<-[:HAS_COMPLEX_PROPERTY {added:datetime({timezone: 'UTC'})}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":  contactId,
		"fieldSetId": fieldSetId.String(),
		"type":       fieldSet.Type,
		"name":       fieldSet.Name,
	})
	return fieldSetId.String()
}

func createDefaultTextFieldInSet(driver *neo4j.Driver, fieldSetId string) string {
	return createTextFieldInSet(driver, fieldSetId, entity.TextCustomFieldEntity{Name: "name", Value: "value"})
}

func createTextFieldInSet(driver *neo4j.Driver, fieldSetId string, textField entity.TextCustomFieldEntity) string {
	var fieldId, _ = uuid.NewRandom()
	query := `
			MATCH (s:FieldSet {id:$fieldSetId})
			MERGE (:TextCustomField {
				  id: $fieldId,
				  value: $value,
				  name: $name
				})<-[:HAS_TEXT_PROPERTY]-(s)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"fieldSetId": fieldSetId,
		"fieldId":    fieldId.String(),
		"name":       textField.Name,
		"value":      textField.Value,
	})
	return fieldId.String()
}

func addEmailToContact(driver *neo4j.Driver, contactId string, email string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:Email {
				  id: randomUUID(),
				  email: $email,
				  label: $label
				})<-[:EMAILED_AT {primary:$primary}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"email":     email,
		"label":     label,
	})
}

func addPhoneNumberToContact(driver *neo4j.Driver, contactId string, e164 string, primary bool, label string) {
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (:PhoneNumber {
				  id: randomUUID(),
				  e164: $e164,
				  label: $label
				})<-[:CALLED_AT {primary:$primary}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId": contactId,
		"primary":   primary,
		"e164":      e164,
		"label":     label,
	})
}

func getCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := integration_tests.ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
