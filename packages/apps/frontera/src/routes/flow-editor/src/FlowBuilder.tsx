import { useParams } from 'react-router-dom';
import React, { useState, MouseEvent, useCallback } from 'react';

import { observer } from 'mobx-react-lite';
import { OnConnect } from '@xyflow/system';
import { FlowActionType } from '@store/Flows/types';
import { FlowStore } from '@store/Flows/Flow.store';
import {
  Edge,
  Node,
  addEdge,
  ReactFlow,
  Background,
  MarkerType,
  OnNodeDrag,
  NodeChange,
  getIncomers,
  getOutgoers,
  useNodesState,
  useEdgesState,
  OnNodesDelete,
  OnEdgesDelete,
  OnNodesChange,
  OnBeforeDelete,
  FitViewOptions,
  applyNodeChanges,
  getConnectedEdges,
  SelectionDragHandler,
} from '@xyflow/react';

import { useStore } from '@shared/hooks/useStore';

import { nodeTypes } from './nodes';
import { BasicEdge } from './edges';
import { getHelperLines } from './utils';
import { useUndoRedo, useKeyboardShortcuts } from './hooks';
import {
  HelperLines,
  FlowSettingsPanel,
  FlowBuilderToolbar,
} from './components';

import '@xyflow/react/dist/style.css';
const edgeTypes = {
  baseEdge: BasicEdge,
};

export const FlowBuilder = observer(
  ({
    onHasNewChanges,
    showSidePanel,
    onToggleSidePanel,
  }: {
    showSidePanel: boolean;
    onHasNewChanges: () => void;
    onToggleSidePanel: (newState: boolean) => void;
  }) => {
    const store = useStore();
    const id = useParams().id as string;

    useKeyboardShortcuts(id, store);

    const flow = store.flows.value.get(id) as FlowStore;

    const [nodes, setNodes] = useNodesState(flow?.parsedNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(flow?.parsedEdges);
    const { takeSnapshot } = useUndoRedo();
    const [helperLineHorizontal, setHelperLineHorizontal] = useState<
      number | undefined
    >(undefined);
    const [helperLineVertical, setHelperLineVertical] = useState<
      number | undefined
    >(undefined);

    // const { screenToFlowPosition } = useReactFlow();
    const { ui } = useStore();

    const onConnect: OnConnect = useCallback(
      (params) => {
        takeSnapshot();

        setEdges((eds) =>
          addEdge(
            {
              ...params,
              type: 'baseEdge',
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

    const customApplyNodeChanges = useCallback(
      (changes: NodeChange[], nodes: Node[]): Node[] => {
        // reset the helper lines (clear existing lines, if any)
        setHelperLineHorizontal(undefined);
        setHelperLineVertical(undefined);

        // this will be true if it's a single node being dragged
        // inside we calculate the helper lines and snap position for the position where the node is being moved to
        if (
          changes.length === 1 &&
          changes[0].type === 'position' &&
          changes[0].dragging &&
          changes[0].position
        ) {
          const helperLines = getHelperLines(changes[0], nodes);

          // if we have a helper line, we snap the node to the helper line position
          // this is being done by manipulating the node position inside the change object
          changes[0].position.x =
            helperLines.snapPosition.x ?? changes[0].position.x;
          changes[0].position.y =
            helperLines.snapPosition.y ?? changes[0].position.y;

          // if helper lines are returned, we set them so that they can be displayed
          setHelperLineHorizontal(helperLines.horizontal);
          setHelperLineVertical(helperLines.vertical);
        }

        return applyNodeChanges(changes, nodes);
      },
      [],
    );

    const onNodesChange: OnNodesChange = useCallback(
      (changes) => {
        setNodes((nodes) => customApplyNodeChanges(changes, nodes));
      },
      [setNodes, customApplyNodeChanges],
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

    const onNodeDragStart: OnNodeDrag = useCallback(() => {
      takeSnapshot();
    }, [takeSnapshot]);

    const onSelectionDragStart: SelectionDragHandler = useCallback(() => {
      takeSnapshot();
    }, [takeSnapshot]);

    const onNodesDelete: OnNodesDelete = useCallback(
      (deleted) => {
        setEdges(
          deleted.reduce((acc, node) => {
            const incomers = getIncomers(node, nodes, edges);
            const outgoers = getOutgoers(node, nodes, edges);
            const connectedEdges = getConnectedEdges([node], edges);

            const remainingEdges = acc.filter(
              (edge) => !connectedEdges.includes(edge),
            );

            const createdEdges = incomers.flatMap(({ id: source }) =>
              outgoers.map(({ id: target }) => ({
                id: `${source}->${target}`,
                source,
                target,
                type: 'baseEdge',
                markerEnd: {
                  type: MarkerType.Arrow,
                  width: 20,
                  height: 20,
                },
              })),
            );

            return [...remainingEdges, ...createdEdges];
          }, edges),
        );
      },
      [nodes, edges],
    );

    const onEdgesDelete: OnEdgesDelete = useCallback(() => {
      takeSnapshot();
    }, [takeSnapshot]);

    const onNodesChangeHandler = useCallback(
      (changes: NodeChange[]) => {
        // this is hack to prevent removing initial edges automatically for some unknown yet reason

        const shouldProhibitChanges =
          changes.every((change) => change.type === 'remove') &&
          nodes.length === changes.length;

        if (shouldProhibitChanges) return;
        onNodesChange(changes);

        if (
          changes.some(
            (e) =>
              e.type === 'add' || e.type === 'remove' || e.type === 'replace',
          )
        ) {
          onHasNewChanges();
        }

        // Check if we need to open the side panel
        const shouldOpenSidePanel = changes.some((change) => {
          if (
            change.type === 'add' &&
            change.item.type === 'action' &&
            change.item.data.action === 'EMAIL_NEW'
          ) {
            // Check if this is the first email action
            const existingEmailNodes = nodes.filter(
              (node) =>
                node.type === 'action' && node.data.action === 'EMAIL_NEW',
            );

            return existingEmailNodes.length === 0;
          }

          return false;
        });

        if (shouldOpenSidePanel) {
          onToggleSidePanel(true);
        }
      },
      [nodes, onNodesChange],
    );

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
          onNodesDelete={onNodesDelete}
          onEdgesDelete={onEdgesDelete}
          onBeforeDelete={onBeforeDelete}
          onNodeDragStart={onNodeDragStart}
          onEdgeMouseLeave={onEdgeMouseLeave}
          onEdgeMouseEnter={onEdgeMouseEnter}
          // onConnectEnd={onConnectEnd}
          onNodesChange={onNodesChangeHandler}
          zoomOnPinch={!ui.flowCommandMenu.isOpen}
          zoomOnScroll={!ui.flowCommandMenu.isOpen}
          onSelectionDragStart={onSelectionDragStart}
          defaultViewport={{ zoom: 0.4, x: 50, y: 0 }}
          preventScrolling={!ui.flowCommandMenu.isOpen}
          proOptions={{
            hideAttribution: true,
          }}
          fitViewOptions={{
            padding: 0.95,
            includeHiddenNodes: false,
            minZoom: 0.1,
            maxZoom: 5,
          }}
          onClick={() => {
            if (ui.flowCommandMenu.isOpen) {
              ui.flowCommandMenu.setOpen(false);
            }

            if (showSidePanel) {
              onToggleSidePanel(false);
            }
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
          <HelperLines
            vertical={helperLineVertical}
            horizontal={helperLineHorizontal}
          />
          <Background />
          <FlowBuilderToolbar />
        </ReactFlow>
        {showSidePanel && <FlowSettingsPanel id={id} nodes={nodes} />}
      </>
    );
  },
);
