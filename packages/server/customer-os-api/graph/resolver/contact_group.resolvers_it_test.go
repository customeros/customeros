package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_ContactGroupRemoveContact(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, tenantName)
	groupId := neo4jt.CreateContactGroup(driver, tenantName, "Group1")
	neo4jt.AddContactToGroup(driver, contactId1, groupId)
	neo4jt.AddContactToGroup(driver, contactId2, groupId)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "BELONGS_TO_GROUP"))

	rawResponse, err := c.RawPost(getQuery("remove_contact_from_group"),
		client.Var("contactId", contactId1),
		client.Var("groupId", groupId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		ContactGroupRemoveContact model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.ContactGroupRemoveContact.Result)

	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "BELONGS_TO_GROUP"))

	assertNeo4jLabels(t, driver, []string{"Tenant", "Contact", "Contact_" + tenantName, "ContactGroup"})
}

func TestMutationResolver_ContactGroupCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	rawResponse, err := c.RawPost(getQuery("create_contact_group"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactGroup struct {
		ContactGroupCreate model.ContactGroup
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactGroup)
	require.Nil(t, err)
	require.NotNil(t, contactGroup)
	require.NotNil(t, contactGroup.ContactGroupCreate.ID)
	require.Equal(t, "the contact group", contactGroup.ContactGroupCreate.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "ContactGroup_"+tenantName))

	assertNeo4jLabels(t, driver, []string{"Tenant", "ContactGroup", "ContactGroup_" + tenantName})
}

func TestQueryResolver_ContactGroups(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	group4 := neo4jt.CreateContactGroup(driver, tenantName, "GROUP4")
	group3 := neo4jt.CreateContactGroup(driver, tenantName, "group3")
	group2 := neo4jt.CreateContactGroup(driver, tenantName, "GROUP2")
	group1 := neo4jt.CreateContactGroup(driver, tenantName, "group1")

	rawResponse, err := c.RawPost(getQuery("get_contact_groups_default_sorting"))
	assertRawResponseSuccess(t, rawResponse, err)

	var groups struct {
		ContactGroups model.ContactGroupPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &groups)
	require.Nil(t, err)
	require.NotNil(t, groups.ContactGroups)
	require.Equal(t, 4, len(groups.ContactGroups.Content))
	require.Equal(t, group1, groups.ContactGroups.Content[0].ID)
	require.Equal(t, group2, groups.ContactGroups.Content[1].ID)
	require.Equal(t, group3, groups.ContactGroups.Content[2].ID)
	require.Equal(t, group4, groups.ContactGroups.Content[3].ID)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
}

func TestQueryResolver_ContactGroups_SortDescendingCaseSensitive(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	group1 := neo4jt.CreateContactGroup(driver, tenantName, "group1")
	group2 := neo4jt.CreateContactGroup(driver, tenantName, "GROUP2")
	group3 := neo4jt.CreateContactGroup(driver, tenantName, "group3")
	group4 := neo4jt.CreateContactGroup(driver, tenantName, "GROUP4")

	rawResponse, err := c.RawPost(getQuery("get_contact_groups_desc_sorting"))
	assertRawResponseSuccess(t, rawResponse, err)

	var groups struct {
		ContactGroups model.ContactGroupPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &groups)
	require.Nil(t, err)
	require.NotNil(t, groups.ContactGroups)
	require.Equal(t, 4, len(groups.ContactGroups.Content))
	require.Equal(t, group3, groups.ContactGroups.Content[0].ID)
	require.Equal(t, group1, groups.ContactGroups.Content[1].ID)
	require.Equal(t, group4, groups.ContactGroups.Content[2].ID)
	require.Equal(t, group2, groups.ContactGroups.Content[3].ID)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
}

func TestQueryResolver_Contacts_ForContactGroup(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other tenant")
	group1 := neo4jt.CreateContactGroup(driver, tenantName, "Group1")
	group2 := neo4jt.CreateContactGroup(driver, tenantName, "Group2")
	group3 := neo4jt.CreateContactGroup(driver, tenantName, "Group3")

	contact1InGroup1 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "first",
		LastName:  "contact",
	})
	contact2InGroup1 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "second",
		LastName:  "contact",
	})
	contact3InGroup2 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "third",
		LastName:  "contact",
	})
	contact4InGroups1And2 := neo4jt.CreateContact(driver, tenantName, entity.ContactEntity{
		Title:     "MR",
		FirstName: "forth",
		LastName:  "contact",
	})
	neo4jt.AddContactToGroup(driver, contact1InGroup1, group1)
	neo4jt.AddContactToGroup(driver, contact2InGroup1, group1)
	neo4jt.AddContactToGroup(driver, contact3InGroup2, group2)
	neo4jt.AddContactToGroup(driver, contact4InGroups1And2, group1)
	neo4jt.AddContactToGroup(driver, contact4InGroups1And2, group2)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
	require.Equal(t, 5, neo4jt.GetCountOfRelationships(driver, "BELONGS_TO_GROUP"))

	rawResponse, err := c.RawPost(getQuery("get_contact_groups_with_contacts"))
	assertRawResponseSuccess(t, rawResponse, err)

	var groups struct {
		ContactGroups model.ContactGroupPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &groups)
	require.Nil(t, err)
	require.NotNil(t, groups.ContactGroups)
	require.Equal(t, int64(3), groups.ContactGroups.TotalElements)
	require.Equal(t, 1, groups.ContactGroups.TotalPages)
	require.Equal(t, 3, len(groups.ContactGroups.Content))

	firstGroup := groups.ContactGroups.Content[0]
	require.Equal(t, group1, firstGroup.ID)
	require.Equal(t, int64(3), firstGroup.Contacts.TotalElements)
	require.Equal(t, 1, firstGroup.Contacts.TotalPages)
	require.Equal(t, 3, len(firstGroup.Contacts.Content))
	require.Equal(t, contact1InGroup1, firstGroup.Contacts.Content[0].ID)
	require.Equal(t, contact4InGroups1And2, firstGroup.Contacts.Content[1].ID)
	require.Equal(t, contact2InGroup1, firstGroup.Contacts.Content[2].ID)

	secondGroup := groups.ContactGroups.Content[1]
	require.Equal(t, group2, secondGroup.ID)
	require.Equal(t, int64(2), secondGroup.Contacts.TotalElements)
	require.Equal(t, 1, secondGroup.Contacts.TotalPages)
	require.Equal(t, 2, len(secondGroup.Contacts.Content))
	require.Equal(t, contact4InGroups1And2, secondGroup.Contacts.Content[0].ID)
	require.Equal(t, contact3InGroup2, secondGroup.Contacts.Content[1].ID)

	thirdGroup := groups.ContactGroups.Content[2]
	require.Equal(t, group3, thirdGroup.ID)
	require.Equal(t, int64(0), thirdGroup.Contacts.TotalElements)
	require.Equal(t, 0, thirdGroup.Contacts.TotalPages)
	require.Equal(t, 0, len(thirdGroup.Contacts.Content))
}

func TestQueryResolver_ContactGroups_MultipleFiltersByName(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	groupAcceptFilterCaseInsensitive := neo4jt.CreateContactGroup(driver, tenantName, "aA")
	groupAcceptFilterCaseSensitive := neo4jt.CreateContactGroup(driver, tenantName, "_ABC_")
	groupRejectFilterCaseSensitive := neo4jt.CreateContactGroup(driver, tenantName, "_ABc_")
	groupRejectFilterNegation := neo4jt.CreateContactGroup(driver, tenantName, "ABC_test")

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "ContactGroup"))

	rawResponse, err := c.RawPost(getQuery("get_contact_groups_filter_by_name"))
	assertRawResponseSuccess(t, rawResponse, err)

	var groups struct {
		ContactGroups model.ContactGroupPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &groups)
	require.Nil(t, err)
	require.NotNil(t, groups.ContactGroups)
	require.Equal(t, 2, len(groups.ContactGroups.Content))
	require.Equal(t, groupAcceptFilterCaseSensitive, groups.ContactGroups.Content[0].ID)
	require.Equal(t, groupAcceptFilterCaseInsensitive, groups.ContactGroups.Content[1].ID)
	require.Equal(t, 1, groups.ContactGroups.TotalPages)
	require.Equal(t, int64(2), groups.ContactGroups.TotalElements)
	// suppress unused warnings
	require.NotNil(t, groupRejectFilterCaseSensitive)
	require.NotNil(t, groupRejectFilterNegation)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "ContactGroup"))
}
