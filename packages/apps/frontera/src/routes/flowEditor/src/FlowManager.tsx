import React, { useCallback } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { OnConnect } from '@xyflow/system';
import {
  Edge,
  addEdge,
  ReactFlow,
  Background,
  MarkerType,
  useReactFlow,
  OnConnectEnd,
  useNodesState,
  useEdgesState,
} from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';

import { BasicEdge } from './edges';
import { nodeTypes } from './Nodes.tsx';
import { Toolbar } from './controls/Toolbar.tsx';

import '@xyflow/react/dist/style.css';
const edgeTypes = {
  baseEdge: BasicEdge,
};

const initialNodes = [
  {
    id: 'tn-1',
    type: 'trigger',

    position: { x: 250, y: 100 },
    data: {
      triggerEntity: undefined,
      triggerType: undefined,
    },
  },
  {
    id: 'tn-2',
    type: 'trigger',
    position: { x: 320, y: 300 },
    data: {
      triggerEntity: undefined,
      triggerType: 'EndFlow',
    },
  },
];

const initialEdges: Edge[] = [
  {
    id: 'e1-2',
    source: 'tn-1',
    target: 'tn-2',
    selected: false,
    selectable: true,
    focusable: true,
    interactionWidth: 60,
    markerEnd: {
      type: MarkerType.ArrowClosed,
      width: 60,
      height: 60,
      color: '#FF0072',
    },
    type: 'baseEdge', // You can change this to other types like 'default', 'straight', etc.
  },
];

export const MarketingFlowBuilder = observer(() => {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const { screenToFlowPosition } = useReactFlow();
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

  const onConnectEnd: OnConnectEnd = useCallback(
    (event, connectionState) => {
      // when a connection is dropped on the pane it's not valid
      if (!connectionState.isValid) {
        // we need to remove the wrapper bounds, in order to get the correct position
        const id = Math.random();
        const { clientX, clientY } =
          'changedTouches' in event ? event.changedTouches[0] : event;
        const newNode = {
          id: `${id}-${nodes.length + 1}`,

          position: screenToFlowPosition({
            x: clientX,
            y: clientY,
          }),
          data: { label: `Node ${id}` },
          origin: [0.5, 0.0],
          type: 'step',
        };

        //@ts-expect-error not important at this moment
        setNodes((nds) => nds.concat(newNode));
        setEdges((eds) =>
          eds.concat({
            id: `e${id}-${newNode.id}`,
            target: newNode.id,
            type: 'baseEdge',
            source: connectionState.fromNode?.id ?? '',
          }),
        );
      }
    },
    [screenToFlowPosition],
  );

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

  return (
    <>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onConnectEnd={onConnectEnd}
        onNodesChange={(changes) => {
          const shouldProhibitChanges =
            changes.every((change) => change.type === 'remove') &&
            nodes.length === 2 &&
            nodes.every((e) => e.type === 'trigger');

          if (shouldProhibitChanges) return;

          onNodesChange(changes);
        }}
        onEdgesChange={(changes) => {
          // this is hack to prevent removing initial edges automatically for some unknown yet reason

          const shouldProhibitChanges =
            changes.every((change) => change.type === 'remove') &&
            edges.length === 1;

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
