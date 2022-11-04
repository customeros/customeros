package resolver

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/integration_tests"
)

func createTenant(driver *neo4j.Driver, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	integration_tests.ExecuteWriteQuery(driver, query, map[string]interface{}{
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
	integration_tests.ExecuteWriteQuery(driver, query, map[string]interface{}{
		"tenant":    tenant,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
}
