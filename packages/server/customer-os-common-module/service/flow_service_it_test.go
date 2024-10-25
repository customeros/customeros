package service

import (
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
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
[{"$H":499,"data":{"action":"FLOW_END"},"height":56,"id":"tn-2","internalId":"32120956-2d2c-47ce-a2e7-bce004f3dcd2","measured":{"height":56,"width":131},"position":{"x":96.5,"y":1884},"properties":{"org.eclipse.elk.portConstraints":"FIXED_ORDER"},"sourcePosition":"bottom","targetPosition":"top","type":"control","width":131,"x":96.5,"y":1884},{"$H":343,"data":{"action":"WAIT","nextStepId":"EMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c","waitDuration":1440},"height":56,"id":"WAIT-502d8197-c9be-46f6-82c7-1407253db77c","measured":{"height":56,"width":156},"position":{"x":84,"y":168},"type":"wait","width":156,"x":84,"y":168},{"$H":345,"data":{"action":"EMAIL_NEW","bodyTemplate":"","subject":"","waitBefore":1440,"waitStepId":"WAIT-502d8197-c9be-46f6-82c7-1407253db77c"},"height":56,"id":"EMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c","internalId":"6cd3584a-8a7b-4e86-97b5-cee725b307fb","measured":{"height":56,"width":300},"position":{"x":12,"y":1416},"type":"action","width":300,"x":12,"y":1416},{"$H":497,"data":{"action":"FLOW_START","entity":"CONTACT","triggerType":"RecordAddedManually"},"height":56,"id":"tn-1","internalId":"25b66565-5a0b-4724-9467-7f2568b924e8","measured":{"height":56,"width":300},"position":{"x":12,"y":12},"properties":{"org.eclipse.elk.portConstraints":"FIXED_ORDER"},"sourcePosition":"bottom","targetPosition":"top","type":"trigger","width":300,"x":12,"y":12},{"$H":377,"data":{"action":"WAIT","nextStepId":"EMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74","waitDuration":1440},"height":56,"id":"WAIT-7b225580-8376-423f-8dd4-ae4c4431adb4","measured":{"height":56,"width":156},"position":{"x":84,"y":1572},"type":"wait","width":156,"x":84,"y":1572},{"$H":379,"data":{"action":"EMAIL_NEW","bodyTemplate":"","subject":"","waitBefore":1440,"waitStepId":"WAIT-7b225580-8376-423f-8dd4-ae4c4431adb4"},"height":56,"id":"EMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74","internalId":"a22cb46c-a01f-4de5-ab8b-d0897f4dbedd","measured":{"height":56,"width":300},"position":{"x":12,"y":1728},"type":"action","width":300,"x":12,"y":1728},{"$H":427,"data":{"action":"WAIT","nextStepId":"EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2","waitDuration":1440},"height":56,"id":"WAIT-0838708e-7a9a-4616-b2d7-872ece9fe673","measured":{"height":56,"width":156},"position":{"x":84,"y":324},"type":"wait","width":156,"x":84,"y":324},{"$H":429,"data":{"action":"EMAIL_NEW","bodyTemplate":"","subject":"","waitBefore":1440,"waitStepId":"WAIT-0838708e-7a9a-4616-b2d7-872ece9fe673"},"height":56,"id":"EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2","internalId":"9377e6c3-f38a-4d12-a37e-526e423fad34","measured":{"height":56,"width":300},"position":{"x":12,"y":480},"type":"action","width":300,"x":12,"y":480},{"$H":565,"data":{"action":"WAIT","nextStepId":"EMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d","waitDuration":1440},"height":56,"id":"WAIT-13935715-164d-429e-8bfb-817d0e12a488","measured":{"height":56,"width":156},"position":{"x":84,"y":636},"type":"wait","width":156,"x":84,"y":636},{"$H":567,"data":{"action":"EMAIL_NEW","bodyTemplate":"","subject":"","waitBefore":1440,"waitStepId":"WAIT-13935715-164d-429e-8bfb-817d0e12a488"},"height":56,"id":"EMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d","internalId":"6ea3503d-a938-48cd-ab40-1664480cda63","measured":{"height":56,"width":300},"position":{"x":12,"y":792},"type":"action","width":300,"x":12,"y":792},{"$H":761,"data":{"action":"WAIT","nextStepId":"EMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a","waitDuration":1440},"height":56,"id":"WAIT-8fbdad46-4557-4635-93a9-43da7710d232","measured":{"height":56,"width":156},"position":{"x":84,"y":948},"type":"wait","width":156,"x":84,"y":948},{"$H":763,"data":{"action":"EMAIL_NEW","bodyTemplate":"","subject":"","waitBefore":1440,"waitStepId":"WAIT-8fbdad46-4557-4635-93a9-43da7710d232"},"height":56,"id":"EMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a","internalId":"7785a516-854d-42a3-9ab4-793a824f78ce","measured":{"height":56,"width":300},"position":{"x":12,"y":1104},"type":"action","width":300,"x":12,"y":1104},{"$H":883,"data":{"action":"WAIT","waitDuration":1440},"height":56,"id":"WAIT-02bac3cf-de50-41dc-85ca-2ca1e03b76f9","measured":{"height":56,"width":156},"position":{"x":84,"y":1260},"type":"wait","width":156,"x":84,"y":1260}]
`
	TWO_NEW_EMAILS_FLOW_EDGES = `
[{"id":"etn-1-WAIT-502d8197-c9be-46f6-82c7-1407253db77c","source":"tn-1","target":"WAIT-502d8197-c9be-46f6-82c7-1407253db77c","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"eEMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c-WAIT-7b225580-8376-423f-8dd4-ae4c4431adb4","source":"EMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c","target":"WAIT-7b225580-8376-423f-8dd4-ae4c4431adb4","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20},"data":{"isHovered":false}},{"id":"eWAIT-7b225580-8376-423f-8dd4-ae4c4431adb4-EMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74","source":"WAIT-7b225580-8376-423f-8dd4-ae4c4431adb4","target":"EMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20},"data":{"isHovered":false}},{"id":"eEMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74-tn-2","source":"EMAIL_NEW-d03ae5f7-416e-433d-9990-c9d746424d74","target":"tn-2","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20},"data":{"isHovered":false}},{"id":"eWAIT-502d8197-c9be-46f6-82c7-1407253db77c-WAIT-0838708e-7a9a-4616-b2d7-872ece9fe673","source":"WAIT-502d8197-c9be-46f6-82c7-1407253db77c","target":"WAIT-0838708e-7a9a-4616-b2d7-872ece9fe673","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"eWAIT-0838708e-7a9a-4616-b2d7-872ece9fe673-EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2","source":"WAIT-0838708e-7a9a-4616-b2d7-872ece9fe673","target":"EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"eWAIT-13935715-164d-429e-8bfb-817d0e12a488-EMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d","source":"WAIT-13935715-164d-429e-8bfb-817d0e12a488","target":"EMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20},"data":{"isHovered":false}},{"id":"eEMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d-WAIT-8fbdad46-4557-4635-93a9-43da7710d232","source":"EMAIL_NEW-64c2984a-3ba1-44f4-8209-af07d8a6475d","target":"WAIT-8fbdad46-4557-4635-93a9-43da7710d232","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"eWAIT-8fbdad46-4557-4635-93a9-43da7710d232-EMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a","source":"WAIT-8fbdad46-4557-4635-93a9-43da7710d232","target":"EMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20},"data":{"isHovered":false}},{"id":"eEMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a-WAIT-02bac3cf-de50-41dc-85ca-2ca1e03b76f9","source":"EMAIL_NEW-2c3d61c3-2b9f-41ee-a5a9-e0ffeb09bb5a","target":"WAIT-02bac3cf-de50-41dc-85ca-2ca1e03b76f9","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"eWAIT-02bac3cf-de50-41dc-85ca-2ca1e03b76f9-EMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c","source":"WAIT-02bac3cf-de50-41dc-85ca-2ca1e03b76f9","target":"EMAIL_NEW-7b208629-7562-45cb-a78a-9b8b3073122c","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}},{"id":"EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2->WAIT-13935715-164d-429e-8bfb-817d0e12a488","source":"EMAIL_NEW-015d2eb3-84c4-444c-bb47-f988811bbda2","target":"WAIT-13935715-164d-429e-8bfb-817d0e12a488","type":"baseEdge","markerEnd":{"type":"arrow","width":20,"height":20}}]
`
)

func TestFlowService_FlowMerge1(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	//neo4jtest.CreateTenant(ctx, driver, tenantName)
	//
	//_, err := CommonServices.FlowService.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
	//	Name:  "flow1",
	//	Nodes: ONE_EMAIL_FLOW,
	//	Edges: ONE_EMAIL_FLOW_EDGES,
	//})
	//require.NoError(t, err)
	//
	//require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	//require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	//require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}

func TestFlowService_FlowMerge2(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	//_, err := CommonServices.FlowService.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
	//	Name:  "flow1",
	//	Nodes: TWO_NEW_EMAILS_FLOW,
	//	Edges: TWO_NEW_EMAILS_FLOW_EDGES,
	//})
	//require.NoError(t, err)
	//
	//require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))
	//require.Equal(t, 7, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	//require.Equal(t, 6, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
}
