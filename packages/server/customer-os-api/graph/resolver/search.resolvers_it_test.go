package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_SearchBasic(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateFullTextBasicSearchIndexes(driver, tenantName)

	keyword := "abc"

	notExpectedContactId := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		FirstName: "x",
		LastName:  "y",
	})
	expectedSecondContactId := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		FirstName: "Matching by last name",
		LastName:  "abcdefgh",
	})
	expectedFirstContactId := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		FirstName: "abdd",
		LastName:  "Matching by first name",
	})

	neo4jt.CreateOrganization(driver, tenantName, "THATISNOTMATCHING")
	perfectMatchOgranizationId := neo4jt.CreateOrganization(driver, tenantName, "abc")

	expectedPartialMatchedEmailId := neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, notExpectedContactId, "abd@openline.ai", false, "WORK")
	neo4jt.AddEmailTo(driver, entity.CONTACT, tenantName, expectedSecondContactId, "xxx@yyy.zzz", false, "WORK")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Contact_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization_"+tenantName))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Email"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Email_"+tenantName))

	rawResponse, err := c.RawPost(getQuery("search/search_basic"),
		client.Var("keyword", keyword))
	assertRawResponseSuccess(t, rawResponse, err)

	logrus.Print(perfectMatchOgranizationId)
	logrus.Print(expectedFirstContactId)
	logrus.Print(expectedPartialMatchedEmailId)

	searchBasicResult := rawResponse.Data.(map[string]interface{})["search_Basic"]
	require.NotNil(t, searchBasicResult)
	require.Equal(t, 4, len(searchBasicResult.([]interface{})))
	require.NotNil(t, searchBasicResult.([]interface{})[0].(map[string]interface{})["score"])
	require.NotNil(t, searchBasicResult.([]interface{})[1].(map[string]interface{})["score"])
	require.NotNil(t, searchBasicResult.([]interface{})[2].(map[string]interface{})["score"])
	require.NotNil(t, searchBasicResult.([]interface{})[3].(map[string]interface{})["score"])
	require.NotNil(t, searchBasicResult.([]interface{})[0].(map[string]interface{})["result"])
	require.NotNil(t, searchBasicResult.([]interface{})[1].(map[string]interface{})["result"])
	require.NotNil(t, searchBasicResult.([]interface{})[2].(map[string]interface{})["result"])
	require.NotNil(t, searchBasicResult.([]interface{})[3].(map[string]interface{})["result"])

	organization := searchBasicResult.([]interface{})[0].(map[string]interface{})["result"]
	require.Equal(t, "Organization", organization.(map[string]interface{})["__typename"])
	require.Equal(t, perfectMatchOgranizationId, organization.(map[string]interface{})["id"])

	email := searchBasicResult.([]interface{})[1].(map[string]interface{})["result"]
	require.Equal(t, "Email", email.(map[string]interface{})["__typename"])
	require.Equal(t, expectedPartialMatchedEmailId, email.(map[string]interface{})["id"])

	firstContact := searchBasicResult.([]interface{})[2].(map[string]interface{})["result"]
	require.Equal(t, "Contact", firstContact.(map[string]interface{})["__typename"])
	require.Equal(t, expectedFirstContactId, firstContact.(map[string]interface{})["id"])

	secondContact := searchBasicResult.([]interface{})[3].(map[string]interface{})["result"]
	require.Equal(t, "Contact", secondContact.(map[string]interface{})["__typename"])
	require.Equal(t, expectedSecondContactId, secondContact.(map[string]interface{})["id"])
}
