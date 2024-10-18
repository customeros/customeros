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

	TWO_NEW_EMAILS_FLOW = `
[
  {
    "id": "tn-1",
    "data": {
      "action": "FLOW_START",
      "entity": "CONTACT",
      "triggerType": "RecordAddedManually"
    }
  },
  {
    "id": "WAIT-4",
    "type": "action",
    "data": {
      "action": "WAIT",
      "waitDuration": 3600,
      "nextStepId": "EMAIL_NEW-3"
    }
  },
  {
    "id": "EMAIL_NEW-3",
    "type": "action",
    "data": {
      "action": "EMAIL_NEW",
      "subject": "E1",
      "bodyTemplate": "<p dir=\"ltr\"><span style=\"white-space: pre-wrap;\">AA</span></p>",
      "waitBefore": 1,
      "waitStepId": "WAIT-4",
      "isEditing": false
    }
  },
  {
    "id": "WAIT-6",
    "type": "action",
    "data": {
      "action": "WAIT",
      "waitDuration": 3600,
      "nextStepId": "EMAIL_NEW-5",
      "isEditing": false
    }
  },
  {
    "id": "EMAIL_NEW-5",
    "type": "action",
    "data": {
      "action": "EMAIL_NEW",
      "subject": "E2",
      "bodyTemplate": "<p dir=\"ltr\"><span style=\"white-space: pre-wrap;\">BB</span></p>",
      "waitBefore": 1,
      "waitStepId": "WAIT-6",
      "isEditing": false
    }
  },
  {
	"id": "tn-2",
    "data": {
      "action": "FLOW_END"
    }
  }
]
`
	TWO_NEW_EMAILS_FLOW_EDGES = `
[
  {
    "id": "etn-1-WAIT-4",
    "source": "tn-1",
    "target": "WAIT-4"
  },
  {
    "id": "eWAIT-4-EMAIL_NEW-3",
    "source": "WAIT-4",
    "target": "EMAIL_NEW-3"
  },
  {
    "id": "eEMAIL_NEW-3-WAIT-6",
    "source": "EMAIL_NEW-3",
    "target": "WAIT-6"
  },
  {
    "id": "eWAIT-6-EMAIL_NEW-5",
    "source": "WAIT-6",
    "target": "EMAIL_NEW-5"
  },
  {
    "id": "eEMAIL_NEW-5-tn-2",
    "source": "EMAIL_NEW-5",
    "target": "tn-2"
  }
]
`
)

func TestFlowService_FlowMerge1(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	_, err := CommonServices.FlowService.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
		Name:  "flow1",
		Nodes: ONE_EMAIL_FLOW,
		Edges: ONE_EMAIL_FLOW_EDGES,
	})
	require.NoError(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}

func TestFlowService_FlowMerge2(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	_, err := CommonServices.FlowService.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
		Name:  "flow1",
		Nodes: TWO_NEW_EMAILS_FLOW,
		Edges: TWO_NEW_EMAILS_FLOW_EDGES,
	})
	require.NoError(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}
