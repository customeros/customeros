import { useParams } from 'react-router-dom';
import React, { MouseEvent, useCallback } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { OnConnect } from '@xyflow/system';
import { FlowActionType } from '@store/Flows/types.ts';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import {
  Edge,
  addEdge,
  ReactFlow,
  Background,
  MarkerType,
  useNodesState,
  useEdgesState,
  OnBeforeDelete,
  FitViewOptions,
  NodeMouseHandler,
} from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';

import { nodeTypes } from './nodes';
import { BasicEdge } from './edges';
import { Toolbar } from './controls/Toolbar.tsx';

import '@xyflow/react/dist/style.css';
const edgeTypes = {
  baseEdge: BasicEdge,
};

export const FlowBuilder = observer(() => {
  const store = useStore();
  const id = useParams().id as string;

  const flow = store.flows.value.get(id) as FlowStore;

  const [nodes, _setNodes, onNodesChange] = useNodesState(flow?.parsedNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(flow?.parsedEdges);

  // const { screenToFlowPosition } = useReactFlow();
  const { ui } = useStore();

  const onConnect: OnConnect = useCallback(
    (params) => {
      const edgeData = {
        triggerType: 'time',
        timeValue: 1,
        timeUnit: 'days',
      };

      setEdges((eds) =>
        addEdge(
          {
            ...params,
            type: 'baseEdge',

            data: edgeData,
            markerEnd: { type: MarkerType.Arrow },
          },
          eds,
        ),
      );
    },
    [setEdges],
  );

  const onEdgeMouseLeave = (_event: MouseEvent, edge: Edge) => {
    const edgeId = edge.id;

    // Updates edge
    setEdges((prevElements) =>
      prevElements.map((element) =>
        element.id === edgeId
          ? {
              ...element,

              data: {
                ...element.data,
                isHovered: false,
              },
            }
          : element,
      ),
    );
  };

  const onEdgeMouseEnter = (_event: MouseEvent, edge: Edge) => {
    const edgeId = edge.id;

    setEdges((prevElements) =>
      prevElements.map((element) =>
        element.id === edgeId
          ? {
              ...element,

              data: {
                ...element.data,
                isHovered: true,
              },
            }
          : element,
      ),
    );
  };

  useKey(
    'Escape',
    () => {
      ui.flowCommandMenu.setOpen(false);
    },
    {
      when: ui.flowCommandMenu.isOpen,
    },
  );

  if (!store.flows.isBootstrapped) {
    return 'IS LOADING';
  }

  const onBeforeDelete: OnBeforeDelete = async (elements) => {
    const hasStartNode = elements.nodes.some(
      (e) => e.data?.action === FlowActionType.FLOW_START,
    );
    const hasEndNode = elements.nodes.some(
      (e) => e.data?.action === FlowActionType.FLOW_END,
    );

    const hasStartOrEndNode = hasStartNode || hasEndNode;

    return hasStartOrEndNode ? false : elements;
  };

  const onOpenTriggerHub = (id: string) => {
    ui.flowCommandMenu.setOpen(true);
    ui.flowCommandMenu.setType('TriggersHub');

    ui.flowCommandMenu.setContext({
      ...ui.flowCommandMenu.context,
      id,
    });
  };

  const onOpenTriggerHubDropdown: NodeMouseHandler = (event, node) => {
    event.preventDefault();
    event.stopPropagation();
    onOpenTriggerHub(node.id);
  };

  return (
    <>
      <ReactFlow
        snapToGrid
        maxZoom={3}
        nodes={nodes}
        edges={edges}
        minZoom={0.1}
        fitView={false}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        preventScrolling={false}
        zoomActivationKeyCode={'91'}
        onBeforeDelete={onBeforeDelete}
        onEdgeMouseLeave={onEdgeMouseLeave}
        onEdgeMouseEnter={onEdgeMouseEnter}
        onNodeClick={onOpenTriggerHubDropdown}
        zoomOnPinch={!ui.flowCommandMenu.isOpen}
        zoomOnScroll={!ui.flowCommandMenu.isOpen}
        defaultViewport={{ zoom: 0.4, x: 50, y: 0 }}
        onClick={() => {
          if (ui.flowCommandMenu.isOpen) {
            ui.flowCommandMenu.setOpen(false);
          }
        }}
        fitViewOptions={{
          padding: 0.1,
          includeHiddenNodes: false,
          minZoom: 0.1,
          maxZoom: 5,
        }}
        onInit={(instance) => {
          const fitViewOptions: FitViewOptions = {
            padding: 0.1,
            maxZoom: 1,
            minZoom: 1,
            duration: 150,
            nodes: nodes,
          };

          instance.fitView(fitViewOptions);
        }}
        // onConnectEnd={onConnectEnd}
        onNodesChange={(changes) => {
          // this is hack to prevent removing initial edges automatically for some unknown yet reason

          const shouldProhibitChanges =
            changes.every((change) => change.type === 'remove') &&
            nodes.length === changes.length;

          if (shouldProhibitChanges) return;
          onNodesChange(changes);
        }}
        onEdgesChange={(changes) => {
          // this is hack to prevent removing initial edges automatically for some unknown yet reason

          const shouldProhibitChanges =
            changes.every((change) => change.type === 'remove') &&
            edges.length === changes.length;

          if (shouldProhibitChanges) {
            return;
          }
          onEdgesChange(changes);
        }}
      >
        <Background />
        <Toolbar />
      </ReactFlow>
    </>
  );
});
