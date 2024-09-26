import { useParams } from 'react-router-dom';
import React, { MouseEvent, useCallback } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { OnConnect } from '@xyflow/system';
import { FlowStore } from '@store/Flows/Flow.store.ts';
import {
  Edge,
  addEdge,
  ReactFlow,
  Background,
  MarkerType,
  useNodesState,
  useEdgesState,
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

    // Updates edge
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

  // const onConnectEnd: OnConnectEnd = useCallback(
  //   (event, connectionState) => {
  //     // when a connection is dropped on the pane it's not valid
  //     if (!connectionState.isValid) {
  //       // we need to remove the wrapper bounds, in order to get the correct position
  //       const id = Math.random();
  //       const { clientX, clientY } =
  //         'changedTouches' in event ? event.changedTouches[0] : event;
  //       const newNode = {
  //         id: `${id}-${nodes.length + 1}`,
  //
  //         position: screenToFlowPosition({
  //           x: clientX,
  //           y: clientY,
  //         }),
  //         data: { label: `Node ${id}` },
  //         origin: [0.5, 0.0],
  //         type: 'step',
  //       };
  //
  //       //@ts-expect-error not important at this moment
  //       setNodes((nds) => nds.concat(newNode));
  //       setEdges((eds) =>
  //         eds.concat({
  //           id: `e${id}-${newNode.id}`,
  //           target: newNode.id,
  //           type: 'baseEdge',
  //           source: connectionState.fromNode?.id ?? '',
  //         }),
  //       );
  //     }
  //   },
  //   [screenToFlowPosition],
  // );

  // Keyboard shortcuts
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

  return (
    <>
      <ReactFlow
        maxZoom={1}
        nodes={nodes}
        edges={edges}
        minZoom={0.1}
        fitView={false}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onEdgeMouseLeave={onEdgeMouseLeave}
        onEdgeMouseEnter={onEdgeMouseEnter}
        zoomOnPinch={!ui.flowCommandMenu.isOpen}
        zoomOnScroll={!ui.flowCommandMenu.isOpen}
        defaultViewport={{ zoom: 0.1, x: 50, y: 0 }}
        onClick={() => {
          if (ui.flowCommandMenu.isOpen) {
            ui.flowCommandMenu.setOpen(false);
          }
        }}
        fitViewOptions={{
          padding: 0.1,
          includeHiddenNodes: false,
          minZoom: 0.1,
          maxZoom: 1,
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
