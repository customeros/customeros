package service

import (
	"testing"
)

func TestFlowExecutionService_FlowExecution_1(t *testing.T) {
	ctx := initContext()
	defer tearDownTestCase(ctx)(t)

	//var err error
	//
	//neo4jtest.CreateTenant(ctx, driver, tenantName)
	//require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelTenant))

	//Prepare data

	//A flow with a single email node
	//2 mailboxes
	//3 contacts
	//flow, err := CommonServices.FlowService.FlowMerge(ctx, nil, &neo4jentity.FlowEntity{
	//	Name:  "flow1",
	//	Nodes: ONE_EMAIL_FLOW,
	//	Edges: ONE_EMAIL_FLOW_EDGES,
	//})
	//require.NoError(t, err)
	//
	//require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlow))
	//require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowAction))
	//require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, model.NEXT.String()))
	//
	////mailboxes
	//mailbox1 := "mailbox1"
	//mailbox2 := "mailbox2"
	//
	//err = CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, tenantName, &postgresEntity.TenantSettingsMailbox{
	//	MailboxUsername:         mailbox1,
	//	MinMinutesBetweenEmails: 5,
	//	MaxMinutesBetweenEmails: 5,
	//})
	//require.NoError(t, err)
	//
	//err = CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, tenantName, &postgresEntity.TenantSettingsMailbox{
	//	MailboxUsername:         mailbox2,
	//	MinMinutesBetweenEmails: 5,
	//	MaxMinutesBetweenEmails: 5,
	//})
	//require.NoError(t, err)
	//
	//_, err = CommonServices.FlowService.FlowSenderMerge(ctx, flow.Id, &neo4jentity.FlowSenderEntity{
	//	UserId: &mailbox1,
	//})
	//require.NoError(t, err)
	//
	//_, err = CommonServices.FlowService.FlowSenderMerge(ctx, flow.Id, &neo4jentity.FlowSenderEntity{
	//	UserId: &mailbox2,
	//})
	//require.NoError(t, err)
	//
	//require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowSender))
	//
	////contacts
	//contactId1 := "1"
	//contactI2 := "2"
	//
	//err = CommonServices.Neo4jRepositories.ContactWriteRepository.CreateContact(ctx, tenantName, contactId1, repository.ContactFields{})
	//require.NoError(t, err)
	//
	//err = CommonServices.Neo4jRepositories.ContactWriteRepository.CreateContact(ctx, tenantName, contactI2, repository.ContactFields{})
	//require.NoError(t, err)
	//
	//require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelContact))
	//
	//_, err = CommonServices.FlowService.FlowParticipantAdd(ctx, flow.Id, contactId1, model.CONTACT)
	//require.NoError(t, err)
	//
	//_, err = CommonServices.FlowService.FlowParticipantAdd(ctx, flow.Id, contactI2, model.CONTACT)
	//require.NoError(t, err)
	//
	//require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowParticipant))

	//activate flow
	//_, err = CommonServices.FlowService.FlowChangeStatus(ctx, flow.Id, neo4jentity.FlowStatusActive)
	//require.NoError(t, err)

	//require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, model.NodeLabelFlowActionExecution))
	//
	////asserts
	//c1Executions, err := CommonServices.FlowExecutionService.getFlowActionExecutions(ctx, flow.Id, contactId1, model.CONTACT)
	//require.NoError(t, err)
	//
	//c2Executions, err := CommonServices.FlowExecutionService.getFlowActionExecutions(ctx, flow.Id, contactI2, model.CONTACT)
	//require.NoError(t, err)
	//
	//require.Equal(t, 1, len(c1Executions))
	//require.Equal(t, 1, len(c2Executions))
}
