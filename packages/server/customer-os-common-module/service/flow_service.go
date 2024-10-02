package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FlowService interface {
	FlowGetList(ctx context.Context) (*neo4jentity.FlowEntities, error)
	FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowGetByActionId(ctx context.Context, flowActionId string) (*neo4jentity.FlowEntity, error)
	FlowGetByContactId(ctx context.Context, flowContactId string) (*neo4jentity.FlowEntity, error)
	FlowsGetListWithContact(ctx context.Context, contactIds []string) (*neo4jentity.FlowEntities, error)
	FlowMerge(ctx context.Context, entity *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error)
	FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error)

	FlowActionGetStart(ctx context.Context, flowId string) (*neo4jentity.FlowActionEntity, error)
	FlowActionGetNext(ctx context.Context, actionId string) ([]*neo4jentity.FlowActionEntity, error)
	FlowActionGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowActionEntities, error)
	FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error)

	FlowContactGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowContactEntities, error)
	FlowContactGetById(ctx context.Context, id string) (*neo4jentity.FlowContactEntity, error)
	FlowContactGetByContactId(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error)
	FlowContactAdd(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error)
	FlowContactDelete(ctx context.Context, flowContactId string) error

	FlowActionSenderGetList(ctx context.Context, actionIds []string) (*neo4jentity.FlowActionSenderEntities, error)
	FlowActionSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowActionSenderEntity, error)
	FlowActionSenderMerge(ctx context.Context, flowActionId string, input *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error)
	FlowActionSenderDelete(ctx context.Context, flowActionSenderId string) error
}

type flowService struct {
	services *Services
}

func NewFlowService(services *Services) FlowService {
	return &flowService{
		services: services,
	}
}

func (s *flowService) FlowGetList(ctx context.Context) (*neo4jentity.FlowEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowReadRepository.GetList(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowEntity(v)
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowGetByActionId(ctx context.Context, flowActionId string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetByActionId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionId", flowActionId))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetFlowByActionId(ctx, flowActionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowGetByContactId(ctx context.Context, flowContactId string) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowContactId", flowContactId))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetFlowByContactId(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowsGetListWithContact(ctx context.Context, contactIds []string) (*neo4jentity.FlowEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowsGetListWithContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("contactIds", contactIds))

	data, err := s.services.Neo4jRepositories.FlowReadRepository.GetListWithContact(ctx, contactIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowEntities, 0)
	for _, v := range data {
		e := mapper.MapDbNodeToFlowEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowMerge(ctx context.Context, input *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	var err error

	//unmarshal the Nodes and Edges
	var nodesMap []map[string]interface{}
	err = json.Unmarshal([]byte(input.Nodes), &nodesMap)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	var edgesMap []map[string]interface{}
	err = json.Unmarshal([]byte(input.Edges), &edgesMap)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	graph := &GraphTraversalIterative{
		nodes:   make(map[string]neo4jentity.FlowActionEntity),
		edges:   make(map[string][]string),
		visited: make(map[string]bool),
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	flowEntity, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		toStore := &neo4jentity.FlowEntity{}

		if input.Id == "" {
			toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlow)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			toStore.Status = neo4jentity.FlowStatusInactive
		} else {
			toStore, err = s.FlowGetById(ctx, input.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			if toStore == nil {
				tracing.TraceErr(span, errors.New("flow not found"))
				return nil, errors.New("flow not found")
			}
		}

		toStore.Name = input.Name
		toStore.Nodes = input.Nodes
		toStore.Edges = input.Edges

		_, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		//populate the nodes
		if nodesMap != nil && len(nodesMap) > 0 {
			for _, v := range nodesMap {

				e := neo4jentity.FlowActionEntity{}

				if v["internalId"] != nil {
					e.Id = v["internalId"].(string)
				}

				if v["id"] != nil {
					e.ExternalId = v["id"].(string)
				}

				if v["type"] != nil {
					e.Type = v["type"].(string)
				}

				if v["data"] != nil {

					for k, v2 := range v["data"].(map[string]interface{}) {
						if v2 != nil {
							if k == "action" {
								e.Data.Action = neo4jentity.GetFlowActionType(v2.(string))
							} else if k == "entity" {
								t := v2.(string)
								e.Data.Entity = &t
							} else if k == "triggerType" {
								t := v2.(string)
								e.Data.TriggerType = &t
							} else if k == "waitBefore" {
								e.Data.WaitBefore = int64(v2.(float64))
							} else if k == "subject" {
								t := v2.(string)
								e.Data.Subject = &t
							} else if k == "bodyTemplate" {
								t := v2.(string)
								e.Data.BodyTemplate = &t
							} else if k == "messageTemplate" {
								t := v2.(string)
								e.Data.MessageTemplate = &t
							}
						}
					}

					//exclude nodes not supported
					if e.Data.Action == neo4jentity.FlowActionTypeFlowStart ||
						e.Data.Action == neo4jentity.FlowActionTypeFlowEnd ||
						e.Data.Action == neo4jentity.FlowActionTypeWait ||
						e.Data.Action == neo4jentity.FlowActionTypeEmailNew ||
						e.Data.Action == neo4jentity.FlowActionTypeEmailReply ||
						e.Data.Action == neo4jentity.FlowActionTypeLinkedinConnectionRequest ||
						e.Data.Action == neo4jentity.FlowActionTypeLinkedinMessage {

						if e.Id == "" {
							e.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlowAction)
							if err != nil {
								tracing.TraceErr(span, err)
								return nil, err
							}
						}

						storedNode, err := s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, &tx, &e)
						if err != nil {
							tracing.TraceErr(span, err)
							return nil, err
						}

						err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
							FromEntityId:   toStore.Id,
							FromEntityType: model.FLOW,
							Relationship:   model.HAS,
							ToEntityId:     e.Id,
							ToEntityType:   model.FLOW_ACTION,
						})
						if err != nil {
							tracing.TraceErr(span, err)
							return nil, err
						}

						stored := mapper.MapDbNodeToFlowActionEntity(storedNode)

						if stored.Data.Action == neo4jentity.FlowActionTypeFlowStart {
							err = err
							if err != nil {
								tracing.TraceErr(span, err)
								return nil, err
							}
						}

						v["internalId"] = stored.Id

						jsoned, err := json.Marshal(&v)
						if err != nil {
							tracing.TraceErr(span, err)
							return nil, err
						}

						stored.Json = string(jsoned)

						_, err = s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, &tx, stored)
						if err != nil {
							tracing.TraceErr(span, err)
							return nil, err
						}

						graph.nodes[stored.ExternalId] = *stored
					}
				}
			}
		}

		// Populate the edges (adjacency list)
		for _, v := range edgesMap {
			source := v["source"].(string)
			target := v["target"].(string)

			if source != "" && target != "" {
				graph.edges[source] = append(graph.edges[source], target)
			}
		}

		//get the start nodes and traverse the graph
		err = s.TraverseInputGraph(ctx, &tx, graph)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		nodes, err := json.Marshal(&nodesMap)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		toStore.Nodes = string(nodes)
		_, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return toStore, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return flowEntity.(*neo4jentity.FlowEntity), nil
}

type GraphTraversalIterative struct {
	nodes   map[string]neo4jentity.FlowActionEntity
	edges   map[string][]string // adjacency list of edges
	visited map[string]bool
}

func (s *flowService) FindStartNodes(graph *GraphTraversalIterative) []string {
	var startNodes []string
	for _, node := range graph.nodes {
		if node.Data.Action == neo4jentity.FlowActionTypeFlowStart {
			startNodes = append(startNodes, node.ExternalId)
		}
	}
	return startNodes
}

// TraverseBFS traverses the graph iteratively using BFS (Breadth-First Search)
func (s *flowService) TraverseInputGraph(ctx context.Context, tx *neo4j.ManagedTransaction, graph *GraphTraversalIterative) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.TraverseBFS")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	queue := s.FindStartNodes(graph)

	for len(queue) > 0 {
		// Dequeue the first node
		currentNode := queue[0]
		queue = queue[1:]

		// Skip if already visited
		if graph.visited[currentNode] {
			continue
		}
		graph.visited[currentNode] = true

		// Get the next nodes and add them to the queue for further exploration
		nextNodes := graph.edges[currentNode]
		queue = append(queue, nextNodes...)

		// Process the current node and its edges (relationship batch)
		err := s.ProcessNode(ctx, tx, graph, currentNode, nextNodes)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *flowService) ProcessNode(ctx context.Context, tx *neo4j.ManagedTransaction, graph *GraphTraversalIterative, nodeId string, batch []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.ProcessNode")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	// Process relationships in batch
	for _, nextNodeId := range batch {
		fmt.Printf("Creating relationship: %s -> %s\n", nodeId, nextNodeId)

		var currentNodeInternalId, nextNodeInternalId string

		for _, v := range graph.nodes {
			if v.ExternalId == nodeId {
				currentNodeInternalId = v.Id
			}
			if v.ExternalId == nextNodeId {
				nextNodeInternalId = v.Id
			}
		}

		if currentNodeInternalId == "" || nextNodeInternalId == "" {
			tracing.TraceErr(span, errors.New("internal ids not found"))
			return errors.New("internal ids not found")
		}

		err := s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, tx, tenant, repository.LinkDetails{
			FromEntityId:   currentNodeInternalId,
			FromEntityType: model.FLOW_ACTION,
			Relationship:   model.NEXT,
			ToEntityId:     nextNodeInternalId,
			ToEntityType:   model.FLOW_ACTION,
		})

		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *flowService) FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowChangeStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.FlowReadRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if node == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	flow := mapper.MapDbNodeToFlowEntity(node)

	if flow.Status == status {
		return flow, nil
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		//if status == neo4jentity.FlowStatusActive {
		//	flowContactList, err := s.FlowContactGetList(ctx, []string{flow.Id})
		//	if err != nil {
		//		tracing.TraceErr(span, err)
		//		return nil, err
		//	}
		//
		//	for _, v := range *flowContactList {
		//		err := s.services.FlowExecutionService.ScheduleFlow(ctx, &tx, flow.Id, v.ContactId, model.CONTACT)
		//		if err != nil {
		//			tracing.TraceErr(span, err)
		//			return nil, err
		//		}
		//	}
		//}

		flow.Status = status

		node, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, flow)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowEntity(node), nil
}

func (s *flowService) FlowActionGetStart(ctx context.Context, flowId string) (*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetStart")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetStartAction(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionEntity(node), nil
}

func (s *flowService) FlowActionGetNext(ctx context.Context, actionId string) ([]*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetNext")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("actionId", actionId))

	nodes, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetNext(ctx, actionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make([]*neo4jentity.FlowActionEntity, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionEntity(v)
		entities = append(entities, e)
	}

	return entities, nil
}

func (s *flowService) FlowActionGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowActionEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("flowIds", flowIds))

	nodes, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetList(ctx, flowIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowActionEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionEntity(node), nil
}

func (s *flowService) FlowContactGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowContactEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("flowIds", flowIds))

	nodes, err := s.services.Neo4jRepositories.FlowContactReadRepository.GetList(ctx, flowIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowContactEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowContactEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowContactGetById(ctx context.Context, id string) (*neo4jentity.FlowContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowContactReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowContactEntity(node), nil
}

func (s *flowService) FlowContactGetByContactId(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactGetByContactId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.String("contactId", contactId))

	identified, err := s.services.Neo4jRepositories.FlowContactReadRepository.Identify(ctx, flowId, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowContactEntity(identified), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowSequenceContact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Contact)
// (fs:FlowSequenceContact {id: $id})-[:HAS]->(:Email)
func (s *flowService) FlowContactAdd(ctx context.Context, flowId, contactId string) (*neo4jentity.FlowContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactLink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	flow, err := s.FlowGetById(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	identified, err := s.services.Neo4jRepositories.FlowContactReadRepository.Identify(ctx, flowId, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if identified == nil {

		contactNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, common.GetTenantFromContext(ctx), contactId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if contactNode == nil {
			tracing.TraceErr(span, errors.New("contact not found"))
			return nil, errors.New("contact not found")
		}

		toStore := neo4jentity.FlowContactEntity{
			ContactId: contactId,
			Status:    neo4jentity.FlowContactStatusPending,
		}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		toStore.ContactId = contactId

		identified, err = s.services.Neo4jRepositories.FlowContactWriteRepository.Merge(ctx, &toStore)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		entity := mapper.MapDbNodeToFlowContactEntity(identified)

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   flowId,
			FromEntityType: model.FLOW,
			Relationship:   model.HAS,
			ToEntityId:     entity.Id,
			ToEntityType:   model.FLOW_CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   entity.Id,
			FromEntityType: model.FLOW_CONTACT,
			Relationship:   model.HAS,
			ToEntityId:     contactId,
			ToEntityType:   model.CONTACT,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return mapper.MapDbNodeToFlowContactEntity(identified), nil
}

func (s *flowService) FlowContactDelete(ctx context.Context, flowContactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowContactUnlink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowContactId", flowContactId))

	tenant := common.GetTenantFromContext(ctx)

	entity, err := s.FlowContactGetById(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow sequence contact not found"))
		return errors.New("flow sequence contact not found")
	}

	flow, err := s.FlowGetByContactId(ctx, flowContactId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return errors.New("flow not found")
	}

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   flow.Id,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     entity.Id,
		ToEntityType:   model.FLOW_CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   entity.Id,
		FromEntityType: model.FLOW_CONTACT,
		Relationship:   model.HAS,
		ToEntityId:     entity.ContactId,
		ToEntityType:   model.CONTACT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.FlowContactWriteRepository.Delete(ctx, entity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *flowService) FlowActionSenderGetList(ctx context.Context, actionIds []string) (*neo4jentity.FlowActionSenderEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	nodes, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetList(ctx, actionIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowActionSenderEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowActionSenderEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowActionSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowActionSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowActionSenderEntity(node), nil
}

// (fs:FlowSequence {id: $id})-[:HAS]->(:FlowActionSender)
func (s *flowService) FlowActionSenderMerge(ctx context.Context, flowActionId string, input *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	action, err := s.FlowActionGetById(ctx, flowActionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if action == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return nil, errors.New("flow action not found")
	}

	isNew := input.Id == ""
	var toStore *neo4jentity.FlowActionSenderEntity

	if input.Id == "" {
		toStore = &neo4jentity.FlowActionSenderEntity{}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowActionSender)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else {
		toStore, err = s.FlowActionSenderGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if toStore == nil {
			tracing.TraceErr(span, errors.New("flow action sender not found"))
			return nil, errors.New("flow action sender not found")
		}
	}

	toStore.Mailbox = input.Mailbox
	toStore.UserId = input.UserId

	node, err := s.services.Neo4jRepositories.FlowActionSenderWriteRepository.Merge(ctx, toStore)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if isNew {
		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   flowActionId,
			FromEntityType: model.FLOW_ACTION,
			Relationship:   model.HAS,
			ToEntityId:     toStore.Id,
			ToEntityType:   model.FLOW_ACTION_SENDER,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	//TODO LINK WITH USER and UNLINK WITH PREVIOUS IF CHANGED

	return mapper.MapDbNodeToFlowActionSenderEntity(node), nil
}

func (s *flowService) FlowActionSenderDelete(ctx context.Context, flowActionSenderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowActionSenderDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowActionSenderId", flowActionSenderId))

	actionSender, err := s.FlowActionSenderGetById(ctx, flowActionSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if actionSender == nil {
		tracing.TraceErr(span, errors.New("flow action sender not found"))
		return errors.New("flow action sender not found")
	}

	actionNode, err := s.services.Neo4jRepositories.FlowActionSenderReadRepository.GetFlowActionBySenderId(ctx, flowActionSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if actionNode == nil {
		tracing.TraceErr(span, errors.New("flow action not found"))
		return errors.New("flow action not found")
	}

	action := mapper.MapDbNodeToFlowActionEntity(actionNode)

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   action.Id,
		FromEntityType: model.FLOW_ACTION,
		Relationship:   model.HAS,
		ToEntityId:     actionSender.Id,
		ToEntityType:   model.FLOW_ACTION_SENDER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//TODO unlink from user
	err = s.services.Neo4jRepositories.FlowContactWriteRepository.Delete(ctx, actionSender.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

//func (s *flowService) GetFlowSequenceSenders(ctx context.Context, tenant, sequenceId string, page, limit int) (*utils.Pagination, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenders")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	pageResult := utils.Pagination{
//		Page:  page,
//		Limit: limit,
//	}
//
//	count, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Count(ctx, tenant, sequenceId)
//	if err != nil {
//		return nil, err
//	}
//
//	entities, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Get(ctx, tenant, sequenceId, page, limit)
//	if err != nil {
//		return nil, err
//	}
//
//	pageResult.SetTotalRows(count)
//	pageResult.SetRows(entities)
//
//	return &pageResult, nil
//}
//
//func (s *flowService) GetFlowSequenceSenderById(ctx context.Context, tenant, id string) (*neo4jentity.FlowActionSenderEntity, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.GetFlowSequenceSenderById")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	span.LogFields(log.String("id", id))
//
//	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.GetById(ctx, tenant, id)
//	if err != nil {
//		return nil, err
//	}
//
//	return entity, nil
//}
//
//func (s *flowService) StoreFlowSequenceSender(ctx context.Context, tenant string, entity *neo4jentity.FlowActionSenderEntity) (*neo4jentity.FlowActionSenderEntity, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.StoreFlowSequenceSender")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	entity, err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Store(ctx, tenant, entity)
//	if err != nil {
//		return nil, err
//	}
//
//	return entity, nil
//}
//
//func (s *flowService) DeleteFlowSequenceSender(ctx context.Context, tenant, id string) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.DeleteFlowSequenceSender")
//	defer span.Finish()
//	tracing.SetDefaultServiceSpanTags(ctx, span)
//
//	span.LogFields(log.String("id", id))
//
//	err := s.services.PostgresRepositories.FlowSequenceSenderRepository.Delete(ctx, tenant, id)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
