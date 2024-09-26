package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

const tenantName = "openline"

const (
	f1 = `
[
   {
      "id":"start-1",
      "type":"trigger",
      "data":{
         "action":"FLOW_START",
         "entity":"CONTACT"
      }
   },
   {
      "id":"linkedin-1",
      "type":"action",
      "data":{
         "action":"LINKEDIN_CONNECTION_REQUEST",
         "messageTemplate":"Here would be the body of the email"
      }
   },
   {
      "id":"email-1",
      "type":"action",
      "data":{
         "action":"EMAIL_NEW",
         "waitBefore":120,
         "subject":"Nice subject here",
         "bodyTemplate":"Here would be the body of the email"
      }
   },
   {
      "id":"email-2",
      "type":"action",
      "data":{
         "action":"EMAIL_REPLY",
         "waitBefore":120,
         "bodyTemplate":"Here would be the body of the email"
      }
   },
   {
      "id":"end",
      "type":"trigger",
      "data":{
         "action":"FLOW_END"
      }
   }
]`

	e1 = `
[
   {
      "source":"start-1",
      "target":"email-1"
   },
   {
      "source":"start-1",
      "target":"linkedin-1"
   },
   {
      "source":"email-1",
      "target":"email-2"
   },
   {
      "source":"email-2",
      "target":"end"
   },
   {
      "source":"linkedin-1",
      "target":"end"
   }
]`
)

func TestFlowService_FlowMerge1(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	_, err := CommonServices.FlowService.FlowMerge(ctx, &neo4jentity.FlowEntity{
		Name:  "flow1",
		Nodes: f1,
		Edges: e1,
	})
	require.NoError(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	require.Equal(t, 5, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}
