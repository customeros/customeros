package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_ContactTypes(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other")
	contactTypeId1 := neo4jt.CreateContactType(driver, tenantName, "first")
	contactTypeId2 := neo4jt.CreateContactType(driver, tenantName, "second")
	neo4jt.CreateContactType(driver, "other", "contact type for other tenant")

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "ContactType"))

	rawResponse, err := c.RawPost(getQuery("get_contact_types"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactTypes []model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.Equal(t, 2, len(contactType.ContactTypes))
	require.Equal(t, contactTypeId1, contactType.ContactTypes[0].ID)
	require.Equal(t, "first", contactType.ContactTypes[0].Name)
	require.Equal(t, contactTypeId2, contactType.ContactTypes[1].ID)
	require.Equal(t, "second", contactType.ContactTypes[1].Name)
}

func TestMutationResolver_ContactTypeCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "otherTenantName")

	rawResponse, err := c.RawPost(getQuery("create_contact_type"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactType_Create model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.NotNil(t, contactType.ContactType_Create.ID)
	require.Equal(t, "the contact type", contactType.ContactType_Create.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType_"+tenantName))

	assertNeo4jLabels(t, driver, []string{"Tenant", "ContactType", "ContactType_" + tenantName})
}

func TestMutationResolver_ContactTypeUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactTypeId := neo4jt.CreateContactType(driver, tenantName, "original type")

	rawResponse, err := c.RawPost(getQuery("update_contact_type"),
		client.Var("contactTypeId", contactTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactType struct {
		ContactType_Update model.ContactType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactType)
	require.Nil(t, err)
	require.NotNil(t, contactType)
	require.Equal(t, contactTypeId, contactType.ContactType_Update.ID)
	require.Equal(t, "updated type", contactType.ContactType_Update.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))
}

func TestMutationResolver_ContactTypeDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactTypeId := neo4jt.CreateContactType(driver, tenantName, "the type")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactType"))

	rawResponse, err := c.RawPost(getQuery("delete_contact_type"),
		client.Var("contactTypeId", contactTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		ContactType_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.ContactType_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "ContactType"))
}
