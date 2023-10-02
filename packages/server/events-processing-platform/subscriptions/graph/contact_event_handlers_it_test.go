package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	contactAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	contactEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	contactModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	emailAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	emailEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphContactEventHandler_OnContactCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &GraphContactEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	contactAggregate := contactAggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())
	curTime := time.Now().UTC()

	contactDto := &contactModels.ContactDto{
		ID:          myContactId.String(),
		Tenant:      tenantName,
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
		Source:      commonModels.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"},
		CreatedAt:   &curTime,
		UpdatedAt:   &curTime,
	}
	event, err := contactEvents.NewContactCreateEvent(contactAggregate, contactDto, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbContactNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(contactProps, "appSource"))

}

func TestGraphContactEventHandler_OnContactCreateWithEmail(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &GraphContactEventHandler{
		Repositories: testDatabase.Repositories,
	}
	emailEventHandler := &GraphEmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	myMailId, _ := uuid.NewUUID()
	curTime := time.Now().UTC()

	contactAggregate := contactAggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())

	contactDto := &contactModels.ContactDto{
		ID:          myContactId.String(),
		Tenant:      tenantName,
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
		Source:      commonModels.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"},
		CreatedAt:   &curTime,
		UpdatedAt:   &curTime,
	}
	event, err := contactEvents.NewContactCreateEvent(contactAggregate, contactDto, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(props, "id"))

	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"

	event, err = emailEvents.NewEmailCreateEvent(emailAggregate, tenantName, email, commonModels.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	event, err = contactEvents.NewContactLinkEmailEvent(contactAggregate, tenantName, myMailId.String(), "work", true, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnEmailLinkToContact(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	emailNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, emailNode)
	emailProps := utils.GetPropsFromNode(*emailNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(emailProps, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(emailProps, "appSource"))

	dbContactNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(contactProps, "appSource"))

}

func TestGraphContactEventHandler_OnContactCreateWithEmailOutOfOrder(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &GraphContactEventHandler{
		Repositories: testDatabase.Repositories,
	}
	emailEventHandler := &GraphEmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	myMailId, _ := uuid.NewUUID()
	curTime := time.Now().UTC()

	contactAggregate := contactAggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())

	contactDto := &contactModels.ContactDto{
		ID:          myContactId.String(),
		Tenant:      tenantName,
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
		Source:      commonModels.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"},
		CreatedAt:   &curTime,
		UpdatedAt:   &curTime,
	}
	event, err := contactEvents.NewContactCreateEvent(contactAggregate, contactDto, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbNode)
	props := utils.GetPropsFromNode(*dbNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(props, "id"))

	event, err = contactEvents.NewContactLinkEmailEvent(contactAggregate, tenantName, myMailId.String(), "work", true, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnEmailLinkToContact(context.Background(), event)
	require.Nil(t, err)

	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"

	event, err = emailEvents.NewEmailCreateEvent(emailAggregate, tenantName, email, commonModels.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	emailNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, emailNode)
	emailProps := utils.GetPropsFromNode(*emailNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(emailProps, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(emailProps, "appSource"))

	dbContactNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(contactProps, "appSource"))

}

func TestGraphContactEventHandler_OnContactCreateWithVeryEmailOutOfOrder(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jt.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &GraphContactEventHandler{
		Repositories: testDatabase.Repositories,
	}
	emailEventHandler := &GraphEmailEventHandler{
		Repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	myMailId, _ := uuid.NewUUID()
	curTime := time.Now().UTC()

	contactAggregate := contactAggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())

	event, err := contactEvents.NewContactLinkEmailEvent(contactAggregate, tenantName, myMailId.String(), "work", true, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnEmailLinkToContact(context.Background(), event)
	require.Nil(t, err)

	contactDto := &contactModels.ContactDto{
		ID:          myContactId.String(),
		Tenant:      tenantName,
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
		Source:      commonModels.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "unit-test"},
		CreatedAt:   &curTime,
		UpdatedAt:   &curTime,
	}
	event, err = contactEvents.NewContactCreateEvent(contactAggregate, contactDto, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	emailAggregate := emailAggregate.NewEmailAggregateWithTenantAndID(tenantName, myMailId.String())
	email := "test@test.com"

	event, err = emailEvents.NewEmailCreateEvent(emailAggregate, tenantName, email, commonModels.Source{
		Source:        "N/A",
		SourceOfTruth: "N/A",
		AppSource:     "unit-test",
	}, curTime, curTime)
	require.Nil(t, err)
	err = emailEventHandler.OnEmailCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, testDatabase.Driver, "Email_"+tenantName), "Incorrect number of Email_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "EMAIL_ADDRESS_BELONGS_TO_TENANT"), "Incorrect number of EMAIL_ADDRESS_BELONGS_TO_TENANT relationships in Neo4j")
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")

	emailNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, myMailId.String())
	require.Nil(t, err)
	require.NotNil(t, emailNode)
	emailProps := utils.GetPropsFromNode(*emailNode)

	require.Equal(t, myMailId.String(), utils.GetStringPropOrEmpty(emailProps, "id"))
	require.Equal(t, email, utils.GetStringPropOrEmpty(emailProps, "rawEmail"))
	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(emailProps, "appSource"))

	dbContactNode, err := neo4jt.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "unit-test", utils.GetStringPropOrEmpty(contactProps, "appSource"))

}
