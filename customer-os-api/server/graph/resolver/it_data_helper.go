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

func getCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := integration_tests.ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
