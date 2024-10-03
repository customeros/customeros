import { useParams } from 'react-router-dom';
import React, { MouseEvent, useCallback } from 'react';

import { useKey, useKeys } from 'rooks';
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
} from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';

import { nodeTypes } from './nodes';
import { BasicEdge } from './edges';
import { FlowBuilderToolbar } from './components';

import '@xyflow/react/dist/style.css';
const edgeTypes = {
  baseEdge: BasicEdge,
};

export const FlowBuilder = observer(
  ({ onHasNewChanges }: { onHasNewChanges: () => void }) => {
    const store = useStore();
    const id = useParams().id as string;

    const flow = store.flows.value.get(id) as FlowStore;

    const [nodes, setNodes, onNodesChange] = useNodesState(flow?.parsedNodes);
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
              markerEnd: {
                type: MarkerType.Arrow,
                width: 20,
                height: 20,
              },
            },
            eds,
          ),
        );
      },
      [setEdges],
    );

    const onEdgeMouseLeave = (_event: MouseEvent, edge: Edge) => {
      const edgeId = edge.id;

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

    useKeys(['Shift', 'S'], (e) => {
      e.stopPropagation();
      e.preventDefault();

      store.ui.commandMenu.setContext({
        ids: [id || ''],
        entity: 'Flow',
      });
      store.ui.commandMenu.setType('ChangeFlowStatus');
      store.ui.commandMenu.setOpen(true);
    });

    useKeys(['Shift', 'R'], (e) => {
      e.stopPropagation();
      e.preventDefault();
      store.ui.commandMenu.setContext({
        ids: [id || ''],
        entity: 'Flow',
        property: 'name',
      });
      store.ui.commandMenu.setType('RenameFlow');
      store.ui.commandMenu.setOpen(true);
    });

    useKey(
      'Escape',
      () => {
        ui.flowCommandMenu.setOpen(false);
      },
      {
        when: ui.flowCommandMenu.isOpen,
      },
    );

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

    return (
      <>
        <ReactFlow
          snapToGrid
          maxZoom={5}
          nodes={nodes}
          edges={edges}
          minZoom={0.1}
          fitView={true}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
          onBeforeDelete={onBeforeDelete}
          onEdgeMouseLeave={onEdgeMouseLeave}
          onEdgeMouseEnter={onEdgeMouseEnter}
          zoomOnPinch={!ui.flowCommandMenu.isOpen}
          zoomOnScroll={!ui.flowCommandMenu.isOpen}
          defaultViewport={{ zoom: 0.4, x: 50, y: 0 }}
          preventScrolling={!ui.flowCommandMenu.isOpen}
          proOptions={{
            hideAttribution: true,
          }}
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
              duration: 0,
              nodes: nodes,
            };

            instance.fitView(fitViewOptions);
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
          // onConnectEnd={onConnectEnd}
          onNodesChange={(changes) => {
            // this is hack to prevent removing initial edges automatically for some unknown yet reason

            const shouldProhibitChanges =
              changes.every((change) => change.type === 'remove') &&
              nodes.length === changes.length;

            if (shouldProhibitChanges) return;
            onNodesChange(changes);

            if (
              changes.some(
                (e) =>
                  e.type === 'add' ||
                  e.type === 'remove' ||
                  e.type === 'replace',
              )
            ) {
              onHasNewChanges();
            }
          }}
          onNodeDoubleClick={(_event, node) => {
            if (node.type === 'wait' || node.type === 'action') {
              setNodes((nds) =>
                nds.map((n) =>
                  n.id === node.id
                    ? { ...n, data: { ...n.data, isEditing: true } }
                    : n,
                ),
              );

              return;
            }

            if (node.type === 'trigger') {
              ui.flowCommandMenu.setOpen(true);
              ui.flowCommandMenu.setType('TriggersHub');
              ui.flowCommandMenu.setContext({
                id: node.id,
                entity: 'Trigger',
              });

              return;
            }
          }}
        >
          <Background />
          <FlowBuilderToolbar />
        </ReactFlow>
      </>
    );
  },
);
