import React, { useState, useEffect, useCallback } from 'react';

import { useKey } from 'rooks';
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

import SidePanel from './SidePanel';
import { BasicEdge } from './edges';
import { nodeTypes } from './Nodes.tsx';
import { Toolbar } from './controls/Toolbar.tsx';

import '@xyflow/react/dist/style.css';
const edgeTypes = {
  baseEdge: BasicEdge,
};

export const MarketingFlowBuilder = () => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [open, setOpen] = useState<string | false>(false);
  const { screenToFlowPosition } = useReactFlow();

  useEffect(() => {
    // Add the start node if it doesn't exist
    if (!nodes.some((node) => node?.type === 'startNode')) {
      setNodes([
        //@ts-expect-error not important at this moment
        {
          id: 'start-node',
          type: 'startNode',
          position: { x: 0, y: 0 },
          data: {
            triggerType: 'manual',
            interval: 1,
            unit: 'days',
          },
        },
      ]);
    }
  }, []);

  const onConnect: OnConnect = useCallback(
    (params) => {
      const edgeData = {
        triggerType: 'time',
        timeValue: 1,
        timeUnit: 'days',
      };

      //@ts-expect-error not important at this moment
      setEdges((eds) =>
        addEdge(
          {
            ...params,
            type: 'baseEdge',
            data: edgeData,
            style: { strokeWidth: 2 },
            markerEnd: { type: MarkerType.Arrow },
            label: `1 day`, // Default label
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
            //@ts-expect-error not important at this moment
            id: `e${id}-${newNode.id}`,
            target: newNode.id,
            type: 'baseEdge',
            source: connectionState?.fromNode?.id,
          }),
        );
      }
    },
    [screenToFlowPosition],
  );

  const editNode = useCallback((nodeId: string) => {
    setOpen(nodeId);
  }, []);

  const editEdge = useCallback(
    (edge: Edge) => {
      setOpen(edge.id);
    },
    [setEdges],
  );

  const addNode = useCallback(
    (type: string) => {
      const newNode = {
        id: `${type}-${nodes.length + 1}`,
        type: 'step',
        position: { x: Math.random() * 500, y: Math.random() * 500 },
        data: {
          color:
            type === 'emailNode'
              ? 'blue'
              : type === 'linkedInMessageNode'
              ? 'green'
              : 'yellow',
          subject: '',
        },
      };

      setNodes((nds) => nds.concat(newNode));
    },
    [nodes.length, setNodes],
  );

  // Keyboard shortcuts
  useKey(['S'], () => addNode('step'), {
    when: !open,
  });

  return (
    <>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onConnectEnd={onConnectEnd}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onEdgeDoubleClick={(_, edge) => editEdge(edge)}
        onNodeDoubleClick={(_, node) => editNode(node?.id)}
      >
        <Background />
        <Toolbar />
      </ReactFlow>
      <SidePanel
        open={open}
        nodeId={open}
        nodes={nodes}
        edges={edges}
        setOpen={setOpen}
        setNodes={setNodes}
        setEdges={setEdges}
      />
    </>
  );
};
