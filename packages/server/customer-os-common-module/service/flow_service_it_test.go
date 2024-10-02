package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	ONE_EMAIL_FLOW = `
[
   {
      "id":"start",
      "internalId":"internal-start",
      "type":"trigger",
      "data":{
         "action":"FLOW_START",
         "entity":"CONTACT",
		 "triggerType":"RecordAddedManually"
      }
   },
   {
      "id":"email",
      "internalId":"internal-email",
      "type":"action",
      "data":{
         "action":"EMAIL_NEW",
         "waitBefore":0,
         "subject":"Nice subject here",
         "bodyTemplate":"Here would be the body of the email"
      }
   },
   {
      "id":"end",
      "internalId":"internal-end",
      "type":"trigger",
      "data":{
         "action":"FLOW_END"
      }
   }
]`

	ONE_EMAIL_FLOW_EDGES = `
[
   {
      "source":"start",
      "target":"email"
   },
   {
      "source":"email",
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
		Nodes: ONE_EMAIL_FLOW,
		Edges: ONE_EMAIL_FLOW_EDGES,
	})
	require.NoError(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}
