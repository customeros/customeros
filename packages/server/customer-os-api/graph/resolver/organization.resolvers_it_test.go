package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Organizations_FilterByNameLike(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateOrganization(driver, tenantName, "A closed organization")
	neo4jt.CreateOrganization(driver, tenantName, "OPENLINE")
	neo4jt.CreateOrganization(driver, tenantName, "the openline")
	neo4jt.CreateOrganization(driver, tenantName, "some other open organization")
	neo4jt.CreateOrganization(driver, tenantName, "OpEnLiNe")

	require.Equal(t, 5, neo4jt.GetCountOfNodes(driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizations struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizations)
	require.Nil(t, err)
	require.NotNil(t, organizations)
	pagedOrganizations := organizations.Organizations
	require.Equal(t, 2, pagedOrganizations.TotalPages)
	require.Equal(t, int64(4), pagedOrganizations.TotalElements)
	require.Equal(t, "OPENLINE", pagedOrganizations.Content[0].Name)
	require.Equal(t, "OpEnLiNe", pagedOrganizations.Content[1].Name)
	require.Equal(t, "some other open organization", pagedOrganizations.Content[2].Name)
}

func TestQueryResolver_Organization(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationInput := entity.OrganizationEntity{
		Name:        "Organization name",
		Description: "Organization description",
		Domain:      "Organization domain",
		Website:     "Organization_website.com",
		Industry:    "tech",
		IsPublic:    true,
	}
	organizationId1 := neo4jt.CreateFullOrganization(driver, tenantName, organizationInput)
	neo4jt.CreateOrganization(driver, tenantName, "otherOrganization")

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_by_id"),
		client.Var("organizationId", organizationId1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organization struct {
		Organization model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)
	require.Equal(t, organizationId1, organization.Organization.ID)
	require.Equal(t, organizationInput.Name, organization.Organization.Name)
	require.Equal(t, organizationInput.Description, *organization.Organization.Description)
	require.Equal(t, organizationInput.Domain, *organization.Organization.Domain)
	require.Equal(t, organizationInput.Website, *organization.Organization.Website)
	require.Equal(t, organizationInput.IsPublic, *organization.Organization.IsPublic)
	require.Equal(t, organizationInput.Industry, *organization.Organization.Industry)
	require.NotNil(t, organization.Organization.CreatedAt)
}

func TestQueryResolver_Organizations_WithLocationsAndPlaces(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(driver, tenantName, "OPENLINE")
	neo4jt.CreateOrganization(driver, tenantName, "some other organization")
	locationId1 := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Name:      "WORK",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
	})
	locationId2 := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Name:      "UNKNOWN",
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
	})
	placeInput := entity.PlaceEntity{
		Source:    entity.DataSourceOpenline,
		AppSource: "test",
		Country:   "testCountry",
		State:     "testState",
		City:      "testCity",
		Address:   "testAddress",
		Address2:  "testAddress2",
		Zip:       "testZip",
		Phone:     "testPhone",
		Fax:       "testFax",
	}
	placeId := neo4jt.CreatePlaceForLocation(driver, placeInput, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId1, locationId1)
	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId1, locationId2)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Location"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Place"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ASSOCIATED_WITH"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organizations_with_locations_and_places"),
		client.Var("page", 1),
		client.Var("limit", 3),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationsStruct struct {
		Organizations model.OrganizationPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationsStruct)
	require.Nil(t, err)

	organizations := organizationsStruct.Organizations
	require.NotNil(t, organizations)
	require.Equal(t, int64(1), organizations.TotalElements)
	require.Equal(t, 2, len(organizations.Content[0].Locations))

	var locationWithPlace, locationWithoutPlace *model.Location
	if organizations.Content[0].Locations[0].ID == locationId1 {
		locationWithPlace = organizations.Content[0].Locations[0]
		locationWithoutPlace = organizations.Content[0].Locations[1]
	} else {
		locationWithPlace = organizations.Content[0].Locations[1]
		locationWithoutPlace = organizations.Content[0].Locations[0]
	}

	require.Equal(t, locationId1, locationWithPlace.ID)
	require.Equal(t, "WORK", locationWithPlace.Name)
	require.NotNil(t, locationWithPlace.CreatedAt)
	require.NotNil(t, locationWithPlace.UpdatedAt)
	require.Equal(t, "test", *locationWithPlace.AppSource)
	require.Equal(t, model.DataSourceOpenline, *locationWithPlace.Source)
	require.NotNil(t, locationWithPlace.Place)

	place := locationWithPlace.Place
	require.Equal(t, placeId, place.ID)
	require.Equal(t, model.DataSourceOpenline, *place.Source)
	require.Equal(t, placeInput.Country, *place.Country)
	require.Equal(t, placeInput.City, *place.City)
	require.Equal(t, placeInput.State, *place.State)
	require.Equal(t, placeInput.Address, *place.Address)
	require.Equal(t, placeInput.Address2, *place.Address2)
	require.Equal(t, placeInput.Fax, *place.Fax)
	require.Equal(t, placeInput.Phone, *place.Phone)
	require.Equal(t, placeInput.Zip, *place.Zip)

	require.Equal(t, locationId2, locationWithoutPlace.ID)
	require.Equal(t, "UNKNOWN", locationWithoutPlace.Name)
	require.NotNil(t, locationWithoutPlace.CreatedAt)
	require.NotNil(t, locationWithoutPlace.UpdatedAt)
	require.Equal(t, "test", *locationWithoutPlace.AppSource)
	require.Equal(t, model.DataSourceOpenline, *locationWithoutPlace.Source)
	require.Nil(t, locationWithoutPlace.Place)
}

func TestQueryResolver_Organization_WithNotes_ById(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "test org")
	userId := neo4jt.CreateDefaultUserWithId(driver, tenantName, testUserId)
	noteId1 := neo4jt.CreateNoteForOrganization(driver, tenantName, organizationId, "note1")
	noteId2 := neo4jt.CreateNoteForOrganization(driver, tenantName, organizationId, "note2")
	neo4jt.NoteCreatedByUser(driver, noteId1, userId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "User"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Note"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "NOTED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "CREATED"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_notes_by_id"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)

	notes := searchedOrganization.Organization.Notes.Content
	require.Equal(t, 2, len(notes))
	var noteWithUser, noteWithoutUser *model.Note
	if noteId1 == notes[0].ID {
		noteWithUser = notes[0]
		noteWithoutUser = notes[1]
	} else {
		noteWithUser = notes[1]
		noteWithoutUser = notes[0]
	}
	require.Equal(t, noteId1, noteWithUser.ID)
	require.Equal(t, "note1", noteWithUser.HTML)
	require.NotNil(t, noteWithUser.CreatedAt)
	require.NotNil(t, noteWithUser.CreatedBy)
	require.Equal(t, userId, noteWithUser.CreatedBy.ID)
	require.Equal(t, "first", noteWithUser.CreatedBy.FirstName)
	require.Equal(t, "last", noteWithUser.CreatedBy.LastName)

	require.Equal(t, noteId2, noteWithoutUser.ID)
	require.Equal(t, "note2", noteWithoutUser.HTML)
	require.NotNil(t, noteWithoutUser.CreatedAt)
	require.Nil(t, noteWithoutUser.CreatedBy)
}

func TestMutationResolver_OrganizationCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationTypeId := neo4jt.CreateOrganizationType(driver, tenantName, "COMPANY")

	// Ensure that the tenant and organization type nodes were created in the database.
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "OrganizationType"))
	require.Equal(t, 2, neo4jt.GetTotalCountOfNodes(driver))

	// Call the "create_organization" mutation.
	rawResponse, err := c.RawPost(getQuery("organization/create_organization"),
		client.Var("organizationTypeId", organizationTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	// Unmarshal the response data into the "organization" struct.
	var organization struct {
		Organization_Create model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)

	// Assign the organization to a shorter variable for easier reference.
	createdOrganization := organization.Organization_Create

	// Ensure that the organization was created correctly.
	require.NotNil(t, createdOrganization.ID)
	require.NotNil(t, createdOrganization.CreatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdOrganization.CreatedAt)
	require.NotNil(t, createdOrganization.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), createdOrganization.UpdatedAt)
	require.Equal(t, "organization name", createdOrganization.Name)
	require.Equal(t, "organization description", *createdOrganization.Description)
	require.Equal(t, "organization domain", *createdOrganization.Domain)
	require.Equal(t, "organization website", *createdOrganization.Website)
	require.Equal(t, "organization industry", *createdOrganization.Industry)
	require.Equal(t, true, *createdOrganization.IsPublic)
	require.Equal(t, organizationTypeId, createdOrganization.OrganizationType.ID)
	require.Equal(t, "COMPANY", createdOrganization.OrganizationType.Name)
	require.Equal(t, model.DataSourceOpenline, createdOrganization.Source)
	require.Equal(t, model.DataSourceOpenline, createdOrganization.SourceOfTruth)
	require.Equal(t, "test", createdOrganization.AppSource)

	// Check the number of nodes and relationships in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization_"+tenantName))

	// Check the labels on the nodes in the Neo4j database
	assertNeo4jLabels(t, driver, []string{"Tenant", "OrganizationType", "Organization", "Organization_" + tenantName})
}

func TestMutationResolver_OrganizationUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "some organization")
	organizationTypeIdOrig := neo4jt.CreateOrganizationType(driver, tenantName, "ORIG")
	organizationTypeIdUpdate := neo4jt.CreateOrganizationType(driver, tenantName, "UPDATED")
	neo4jt.SetOrganizationTypeForOrganization(driver, organizationTypeIdOrig, organizationTypeIdOrig)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "OrganizationType"))

	rawResponse, err := c.RawPost(getQuery("organization/update_organization"),
		client.Var("organizationId", organizationId),
		client.Var("organizationTypeId", organizationTypeIdUpdate))
	assertRawResponseSuccess(t, rawResponse, err)

	var organization struct {
		Organization_Update model.Organization
	}
	err = decode.Decode(rawResponse.Data.(map[string]any), &organization)
	require.Nil(t, err)
	require.NotNil(t, organization)
	updatedOrganization := organization.Organization_Update
	require.Equal(t, organizationId, updatedOrganization.ID)
	require.NotNil(t, updatedOrganization.UpdatedAt)
	require.NotEqual(t, utils.GetEpochStart(), updatedOrganization.UpdatedAt)
	require.Equal(t, "updated name", updatedOrganization.Name)
	require.Equal(t, "updated description", *updatedOrganization.Description)
	require.Equal(t, "updated domain", *updatedOrganization.Domain)
	require.Equal(t, "updated website", *updatedOrganization.Website)
	require.Equal(t, "updated industry", *updatedOrganization.Industry)
	require.Equal(t, true, *updatedOrganization.IsPublic)
	require.Equal(t, organizationTypeIdUpdate, updatedOrganization.OrganizationType.ID)
	require.Equal(t, "UPDATED", updatedOrganization.OrganizationType.Name)
	require.Equal(t, model.DataSourceOpenline, updatedOrganization.SourceOfTruth)

	// Check still single organization node exists after update, no new node created
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
}

func TestMutationResolver_OrganizationDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)

	organizationId := neo4jt.CreateOrganization(driver, tenantName, "LLC LLC")
	locationId := neo4jt.CreateLocation(driver, tenantName, entity.LocationEntity{
		Source: "manual",
	})
	neo4jt.CreatePlaceForLocation(driver, entity.PlaceEntity{
		Country: "aCountry",
	}, locationId)
	neo4jt.OrganizationAssociatedWithLocation(driver, organizationId, locationId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Location"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Place"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "ASSOCIATED_WITH"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))

	rawResponse, err := c.RawPost(getQuery("organization/delete_organization"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		Organization_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.Organization_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Location"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Place"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "LOCATED_AT"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "ASSOCIATED_WITH"))

	assertNeo4jLabels(t, driver, []string{"Tenant"})
}

func TestQueryResolver_Organization_WithRoles_ById(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "some organization")
	role1 := neo4jt.ContactWorksForOrganization(driver, contactId1, organizationId, "CTO", false)
	role2 := neo4jt.ContactWorksForOrganization(driver, contactId2, organizationId, "CEO", true)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "JobRole"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "ROLE_IN"))
	require.Equal(t, 2, neo4jt.GetCountOfRelationships(driver, "WORKS_AS"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_job_roles_by_id"),
		client.Var("organizationId", organizationId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)

	roles := searchedOrganization.Organization.JobRoles
	require.Equal(t, 2, len(roles))
	var cto, ceo *model.JobRole
	ceo = roles[0]
	cto = roles[1]
	require.Equal(t, role1, cto.ID)
	require.Equal(t, "CTO", *cto.JobTitle)
	require.Equal(t, false, cto.Primary)
	require.Equal(t, contactId1, cto.Contact.ID)

	require.Equal(t, role2, ceo.ID)
	require.Equal(t, "CEO", *ceo.JobTitle)
	require.Equal(t, true, ceo.Primary)
	require.Equal(t, contactId2, ceo.Contact.ID)
}

func TestQueryResolver_Organization_WithContacts_ById(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	organizationId := neo4jt.CreateOrganization(driver, tenantName, "organization1")
	organizationId2 := neo4jt.CreateOrganization(driver, tenantName, "organization2")
	contactId1 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId3 := neo4jt.CreateDefaultContact(driver, tenantName)
	contactId4 := neo4jt.CreateDefaultContact(driver, tenantName)
	neo4jt.LinkContactWithOrganization(driver, contactId1, organizationId)
	neo4jt.LinkContactWithOrganization(driver, contactId2, organizationId)
	neo4jt.LinkContactWithOrganization(driver, contactId3, organizationId)
	neo4jt.LinkContactWithOrganization(driver, contactId4, organizationId2)

	require.Equal(t, 4, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(driver, "Organization"))
	require.Equal(t, 4, neo4jt.GetCountOfRelationships(driver, "CONTACT_OF"))

	rawResponse, err := c.RawPost(getQuery("organization/get_organization_with_contacts_by_id"),
		client.Var("organizationId", organizationId),
		client.Var("limit", 1),
		client.Var("page", 1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedOrganization struct {
		Organization model.Organization
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedOrganization)
	require.Nil(t, err)
	require.Equal(t, organizationId, searchedOrganization.Organization.ID)
	require.Equal(t, 3, searchedOrganization.Organization.Contacts.TotalPages)
	require.Equal(t, int64(3), searchedOrganization.Organization.Contacts.TotalElements)

	contacts := searchedOrganization.Organization.Contacts.Content
	require.Equal(t, 1, len(contacts))
	require.Equal(t, contactId1, contacts[0].ID)
}
