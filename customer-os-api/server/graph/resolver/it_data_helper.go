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

func createTenantUser(driver *neo4j.Driver, tenant string, user entity.TenantUserEntity) {
	query := `
		MATCH (t:Tenant {name:$tenant})
			MERGE (u:TenantUser {
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

func createDefaultFieldsSet(driver *neo4j.Driver, contactId string) string {
	return createFieldsSet(driver, contactId, entity.FieldsSetEntity{Name: "name", Type: "type"})
}

func createFieldsSet(driver *neo4j.Driver, contactId string, fieldsSet entity.FieldsSetEntity) string {
	var fieldsSetId, _ = uuid.NewRandom()
	query := `
			MATCH (c:Contact {id:$contactId})
			MERGE (s:FieldsSet {
				  id: $fieldsSetId,
				  type: $type,
				  name: $name
				})<-[:HAS_COMPLEX_PROPERTY {added:datetime({timezone: 'UTC'})}]-(c)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"contactId":   contactId,
		"fieldsSetId": fieldsSetId.String(),
		"type":        fieldsSet.Type,
		"name":        fieldsSet.Name,
	})
	return fieldsSetId.String()
}

func createDefaultTextFieldInSet(driver *neo4j.Driver, fieldsSetId string) string {
	return createTextFieldInSet(driver, fieldsSetId, entity.TextCustomFieldEntity{Name: "name", Value: "value"})
}

func createTextFieldInSet(driver *neo4j.Driver, fieldsSetId string, textField entity.TextCustomFieldEntity) string {
	var fieldId, _ = uuid.NewRandom()
	query := `
			MATCH (s:FieldsSet {id:$fieldsSetId})
			MERGE (:TextCustomField {
				  id: $fieldId,
				  value: $value,
				  name: $name
				})<-[:HAS_TEXT_PROPERTY]-(s)`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]any{
		"fieldsSetId": fieldsSetId,
		"fieldId":     fieldId.String(),
		"name":        textField.Name,
		"value":       textField.Value,
	})
	return fieldId.String()
}

func getCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := integration_tests.ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
