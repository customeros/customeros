package resolver

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_ContactGroups(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	group4 := createContactGroup(driver, tenantName, "GROUP4")
	group3 := createContactGroup(driver, tenantName, "group3")
	group2 := createContactGroup(driver, tenantName, "GROUP2")
	group1 := createContactGroup(driver, tenantName, "group1")

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

	require.Equal(t, 4, getCountOfNodes(driver, "ContactGroup"))
}

func TestQueryResolver_ContactGroups_SortDescendingCaseSensitive(t *testing.T) {
	defer setupTestCase()(t)
	createTenant(driver, tenantName)

	group1 := createContactGroup(driver, tenantName, "group1")
	group2 := createContactGroup(driver, tenantName, "GROUP2")
	group3 := createContactGroup(driver, tenantName, "group3")
	group4 := createContactGroup(driver, tenantName, "GROUP4")

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

	require.Equal(t, 4, getCountOfNodes(driver, "ContactGroup"))
}
