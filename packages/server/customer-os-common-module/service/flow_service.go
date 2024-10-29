package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto/events"
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
	FlowGetByParticipant(ctx context.Context, entityId string, entityType model.EntityType) (*neo4jentity.FlowEntity, error)
	FlowsGetListWithContact(ctx context.Context, contactIds []string) (*neo4jentity.FlowEntities, error)
	FlowsGetListWithSender(ctx context.Context, senderIds []string) (*neo4jentity.FlowEntities, error)
	FlowMerge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error)
	FlowChangeStatus(ctx context.Context, id string, status neo4jentity.FlowStatus) (*neo4jentity.FlowEntity, error)

	FlowActionGetStart(ctx context.Context, flowId string) (*neo4jentity.FlowActionEntity, error)
	FlowActionGetNext(ctx context.Context, actionId string) ([]*neo4jentity.FlowActionEntity, error)
	FlowActionGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowActionEntities, error)
	FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error)

	FlowParticipantGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowParticipantEntities, error)
	FlowParticipantById(ctx context.Context, flowParticipantId string) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantByEntity(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantAdd(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantDelete(ctx context.Context, flowParticipantId string) error

	FlowSenderGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowSenderEntities, error)
	FlowSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowSenderEntity, error)
	FlowSenderMerge(ctx context.Context, flowId string, input *neo4jentity.FlowSenderEntity) (*neo4jentity.FlowSenderEntity, error)
	FlowSenderDelete(ctx context.Context, flowSenderId string) error
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

func (s *flowService) FlowGetByParticipant(ctx context.Context, entityId string, entityType model.EntityType) (*neo4jentity.FlowEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowGetByParticipant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("entityId", entityId), log.String("entityType", entityType.String()))

	node, err := s.services.Neo4jRepositories.FlowActionReadRepository.GetFlowByEntity(ctx, entityId, entityType)
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

func (s *flowService) FlowsGetListWithSender(ctx context.Context, senderIds []string) (*neo4jentity.FlowEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowsGetListWithSender")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("senderIds", senderIds))

	data, err := s.services.Neo4jRepositories.FlowReadRepository.GetListWithSender(ctx, senderIds)
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

func (s *flowService) FlowMerge(ctx context.Context, tx *neo4j.ManagedTransaction, input *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error) {
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

	waitNodes := map[string]bool{}
	for _, v := range nodesMap {
		if v["data"] != nil {
			if v["data"].(map[string]interface{})["action"] == "WAIT" {
				waitNodes[v["id"].(string)] = true
			}
		}
	}

	edgesMap = s.removeWaitNodes(edgesMap, waitNodes)

	graph := &GraphTraversalIterative{
		nodes:   make(map[string]neo4jentity.FlowActionEntity),
		edges:   make(map[string][]string),
		visited: make(map[string]bool),
	}

	if input.Id != "" {
		flowEntity, err := s.FlowGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if flowEntity.Status != neo4jentity.FlowStatusInactive {
			return nil, nil
		}
	}

	var existing *neo4jentity.FlowEntity
	if input.Id != "" {
		existing, err = s.FlowGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if existing == nil {
			return nil, errors.New("flow not found")
		}
	}

	flowEntity, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		toStore := &neo4jentity.FlowEntity{}

		if input.Id == "" {
			toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelFlow)
			if err != nil {
				return nil, err
			}

			toStore.Status = neo4jentity.FlowStatusInactive
		} else {
			toStore, err = s.FlowGetById(ctx, input.Id)
			if err != nil {
				return nil, err
			}

			if toStore == nil {
				return nil, errors.New("flow not found")
			}

			if toStore.Status == neo4jentity.FlowStatusScheduling {
				return nil, errors.New("flow is in scheduling status")
			}
		}

		toStore.Name = input.Name
		toStore.Nodes = input.Nodes
		toStore.Edges = input.Edges

		_, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, toStore)
		if err != nil {
			return nil, err
		}

		//TODO this is not supporting live updates after scheduling
		//clear existing nodes in DB and reinsert them
		err := s.services.Neo4jRepositories.FlowActionWriteRepository.DeleteForFlow(ctx, &tx, toStore.Id)
		if err != nil {
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
							return nil, err
						}

						stored := mapper.MapDbNodeToFlowActionEntity(storedNode)

						if stored.Data.Action == neo4jentity.FlowActionTypeFlowStart {
							err = err
							if err != nil {
								return nil, err
							}
						}

						v["internalId"] = stored.Id

						jsoned, err := json.Marshal(&v)
						if err != nil {
							return nil, err
						}

						stored.Json = string(jsoned)

						_, err = s.services.Neo4jRepositories.FlowActionWriteRepository.Merge(ctx, &tx, stored)
						if err != nil {
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
			return nil, err
		}

		nodes, err := json.Marshal(&nodesMap)
		if err != nil {
			return nil, err
		}

		toStore.Nodes = string(nodes)
		_, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, toStore)
		if err != nil {
			return nil, err
		}

		return toStore, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	e := flowEntity.(*neo4jentity.FlowEntity)

	err = s.services.RabbitMQService.PublishEvent(ctx, e.Id, model.FLOW, events.FlowComputeParticipantsRequirements{})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (s *flowService) removeWaitNodes(edges []map[string]interface{}, waitNodes map[string]bool) []map[string]interface{} {
	newEdges := []map[string]interface{}{}
	targetMapping := make(map[string]string)

	// Step 1: Map out where each node points to
	for _, edge := range edges {
		source := edge["source"].(string)
		target := edge["target"].(string)
		targetMapping[source] = target
	}

	// Step 2: Rebuild edges, bypassing WAIT nodes
	for _, edge := range edges {
		source := edge["source"].(string)
		target := edge["target"].(string)

		// If the source is a WAIT node, skip processing this edge
		if waitNodes[source] {
			continue
		}

		// Recursively find the final target if there are multiple WAIT nodes in a row
		finalTarget := target
		for waitNodes[finalTarget] {
			finalTarget = targetMapping[finalTarget] // Keep following the target chain until it's not a WAIT node
		}

		// Create a new edge linking source directly to the final non-WAIT target
		newEdge := s.copyEdge(edge)
		newEdge["target"] = finalTarget
		newEdges = append(newEdges, newEdge)
	}

	return newEdges
}

// Helper function to copy an edge
func (s *flowService) copyEdge(edge map[string]interface{}) map[string]interface{} {
	newEdge := make(map[string]interface{})
	for key, value := range edge {
		newEdge[key] = value
	}
	return newEdge
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

		//exclude the link between start and end nodes
		if graph.nodes[currentNode].Data.Action == neo4jentity.FlowActionTypeFlowStart && len(nextNodes) == 1 && graph.nodes[nextNodes[0]].Data.Action == neo4jentity.FlowActionTypeFlowEnd {
			continue
		}

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

	startInitialSchedule := false

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		if flow.Status == neo4jentity.FlowStatusInactive && status == neo4jentity.FlowStatusActive {
			flow.Status = neo4jentity.FlowStatusScheduling

			startInitialSchedule = true
		} else {
			flow.Status = status
		}

		node, err = s.services.Neo4jRepositories.FlowWriteRepository.Merge(ctx, &tx, flow)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return nil, nil
	})

	if startInitialSchedule {
		err := s.services.RabbitMQService.PublishEvent(ctx, flow.Id, model.FLOW, events.FlowInitialSchedule{})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

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

func (s *flowService) FlowParticipantGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowParticipantEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowParticipantGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("flowIds", flowIds))

	nodes, err := s.services.Neo4jRepositories.FlowParticipantReadRepository.GetList(ctx, flowIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowParticipantEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowParticipantEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowParticipantById(ctx context.Context, flowParticipantId string) (*neo4jentity.FlowParticipantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowParticipantById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowParticipantId", flowParticipantId))

	node, err := s.services.Neo4jRepositories.FlowParticipantReadRepository.GetById(ctx, flowParticipantId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowParticipantEntity(node), nil
}

func (s *flowService) FlowParticipantByEntity(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowParticipantByEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.String("entityId", entityId), log.String("entityType", entityType.String()))

	identified, err := s.services.Neo4jRepositories.FlowParticipantReadRepository.Identify(ctx, flowId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowParticipantEntity(identified), nil
}

func (s *flowService) FlowParticipantAdd(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowParticipantAdd")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.String("entityId", entityId), log.String("entityType", entityType.String()))

	flow, err := s.FlowGetById(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	identified, err := s.services.Neo4jRepositories.FlowParticipantReadRepository.Identify(ctx, flowId, entityId, entityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if identified == nil {

		//validation section
		if entityType == model.CONTACT {
			contactNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, common.GetTenantFromContext(ctx), entityId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get contact")
			}
			if contactNode == nil {
				return nil, errors.New("contact not found")
			}
		}

		e, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
			toStore := neo4jentity.FlowParticipantEntity{
				Status:     neo4jentity.FlowParticipantStatusOnHold,
				EntityId:   entityId,
				EntityType: entityType,
			}

			toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowParticipant)
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate id")
			}

			identified, err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Merge(ctx, &tx, &toStore)
			if err != nil {
				return nil, errors.Wrap(err, "failed to merge flow participant")
			}
			entity := mapper.MapDbNodeToFlowParticipantEntity(identified)

			err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, common.GetTenantFromContext(ctx), repository.LinkDetails{
				FromEntityId:   flowId,
				FromEntityType: model.FLOW,
				Relationship:   model.HAS,
				ToEntityId:     entity.Id,
				ToEntityType:   model.FLOW_PARTICIPANT,
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to link flow to flow participant")
			}

			err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, common.GetTenantFromContext(ctx), repository.LinkDetails{
				FromEntityId:   entity.Id,
				FromEntityType: model.FLOW_PARTICIPANT,
				Relationship:   model.HAS,
				ToEntityId:     entityId,
				ToEntityType:   entityType,
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to link flow participant to entity")
			}

			requirements, err := s.services.FlowExecutionService.GetFlowRequirements(ctx, flowId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get flow requirements")
			}

			err = s.services.FlowExecutionService.UpdateParticipantFlowRequirements(ctx, &tx, entity, requirements)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update participant flow requirements")
			}

			if flow.Status == neo4jentity.FlowStatusActive && entity.Status == neo4jentity.FlowParticipantStatusReady {
				err := s.services.FlowExecutionService.ScheduleFlow(ctx, &tx, flowId, entity)
				if err != nil {
					return nil, errors.Wrap(err, "failed to schedule flow")
				}
			}

			return entity, nil
		})

		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return e.(*neo4jentity.FlowParticipantEntity), nil
	} else {
		return mapper.MapDbNodeToFlowParticipantEntity(identified), nil
	}
}

func (s *flowService) FlowParticipantDelete(ctx context.Context, flowParticipantId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowParticipantDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowParticipantId", flowParticipantId))

	tenant := common.GetTenantFromContext(ctx)

	entity, err := s.FlowParticipantById(ctx, flowParticipantId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if entity == nil {
		tracing.TraceErr(span, errors.New("flow sequence contact not found"))
		return errors.New("flow sequence contact not found")
	}

	flow, err := s.FlowGetByParticipant(ctx, entity.EntityId, entity.EntityType)
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
		ToEntityType:   model.FLOW_PARTICIPANT,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, tenant, repository.LinkDetails{
		FromEntityId:   entity.Id,
		FromEntityType: model.FLOW_PARTICIPANT,
		Relationship:   model.HAS,
		ToEntityId:     entity.EntityId,
		ToEntityType:   entity.EntityType,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Delete(ctx, entity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//TODO
	//remove scheduled events???

	return nil
}

func (s *flowService) FlowSenderGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowSenderEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSenderGetList")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("flowIds", flowIds))

	nodes, err := s.services.Neo4jRepositories.FlowSenderReadRepository.GetList(ctx, flowIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	entities := make(neo4jentity.FlowSenderEntities, 0)
	for _, v := range nodes {
		e := mapper.MapDbNodeToFlowSenderEntity(v.Node)
		e.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *e)
	}

	return &entities, nil
}

func (s *flowService) FlowSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSenderGetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	node, err := s.services.Neo4jRepositories.FlowSenderReadRepository.GetById(ctx, id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSenderEntity(node), nil
}

func (s *flowService) FlowSenderMerge(ctx context.Context, flowId string, input *neo4jentity.FlowSenderEntity) (*neo4jentity.FlowSenderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSenderMerge")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.Object("input", input))

	flow, err := s.FlowGetById(ctx, flowId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if flow == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return nil, errors.New("flow not found")
	}

	isNew := input.Id == ""
	var toStore *neo4jentity.FlowSenderEntity

	if input.Id == "" {
		toStore = &neo4jentity.FlowSenderEntity{}
		toStore.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, common.GetTenantFromContext(ctx), model.NodeLabelFlowSender)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	} else {
		toStore, err = s.FlowSenderGetById(ctx, input.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if toStore == nil {
			tracing.TraceErr(span, errors.New("flow sender not found"))
			return nil, errors.New("flow sender not found")
		}
	}

	toStore.UserId = input.UserId

	node, err := s.services.Neo4jRepositories.FlowSenderWriteRepository.Merge(ctx, toStore)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if isNew {
		err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
			FromEntityId:   flowId,
			FromEntityType: model.FLOW,
			Relationship:   model.HAS,
			ToEntityId:     toStore.Id,
			ToEntityType:   model.FLOW_SENDER,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	//TODO LINK WITH USER and UNLINK WITH PREVIOUS IF CHANGED
	err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   toStore.Id,
		FromEntityType: model.FLOW_SENDER,
		Relationship:   model.HAS,
		ToEntityId:     *toStore.UserId,
		ToEntityType:   model.USER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToFlowSenderEntity(node), nil
}

func (s *flowService) FlowSenderDelete(ctx context.Context, flowSenderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowService.FlowSenderDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("flowSenderId", flowSenderId))

	flowSender, err := s.FlowSenderGetById(ctx, flowSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowSender == nil {
		tracing.TraceErr(span, errors.New("flow sender not found"))
		return errors.New("flow sender not found")
	}

	flowNode, err := s.services.Neo4jRepositories.FlowSenderReadRepository.GetFlowBySenderId(ctx, flowSenderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowNode == nil {
		tracing.TraceErr(span, errors.New("flow not found"))
		return errors.New("flow not found")
	}

	flow := mapper.MapDbNodeToFlowActionEntity(flowNode)

	//todo use TX

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   flow.Id,
		FromEntityType: model.FLOW,
		Relationship:   model.HAS,
		ToEntityId:     flowSender.Id,
		ToEntityType:   model.FLOW_SENDER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CommonWriteRepository.Unlink(ctx, nil, common.GetTenantFromContext(ctx), repository.LinkDetails{
		FromEntityId:   flowSender.Id,
		FromEntityType: model.FLOW_SENDER,
		Relationship:   model.HAS,
		ToEntityId:     *flowSender.UserId,
		ToEntityType:   model.USER,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.FlowParticipantWriteRepository.Delete(ctx, flowSender.Id)
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
