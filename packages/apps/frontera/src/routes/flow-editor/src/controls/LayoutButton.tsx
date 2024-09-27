import React, { useCallback, useLayoutEffect } from 'react';

import { Edge, Node } from '@xyflow/react';
import { ElkNode, LayoutOptions } from 'elkjs/lib/elk-api';
import { useReactFlow, ControlButton } from '@xyflow/react';
import ELK, { ElkExtendedEdge } from 'elkjs/lib/elk.bundled.js';

import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02.tsx';

import '@xyflow/react/dist/style.css';

const elk = new ELK();

// Elk has a *huge* amount of options to configure. To see everything you can
// tweak check out:
//
// - https://www.eclipse.org/elk/reference/algorithms.html
// - https://www.eclipse.org/elk/reference/options.html
const elkOptions = {
  'elk.algorithm': 'layered',
  'elk.layered.spacing.nodeNodeBetweenLayers': '100',
  'elk.spacing.nodeNode': '80',
};

const getNodeWidth = (node: Node) => {
  if (node.type === 'trigger') {
    return {
      width: 300,
      height: 56,
    };
  }

  if (node.type === 'control') {
    return {
      width: 131,
      height: 56,
    };
  }

  if (node.type === 'wait') {
    return {
      width: 150,
      height: 56,
    };
  }

  if (node.type === 'action') {
    if (node.data.stepType === 'Wait') {
      return {
        width: 150,
        height: 56,
      };
    }

    return {
      width: 300,
      height: 56,
    };
  }

  return {
    width: 200,
    height: 75,
  };
};

export const getLayoutedElements = (
  nodes: Node[],
  edges: Edge[],
  options: LayoutOptions = {},
) => {
  const isHorizontal = options?.['elk.direction'] === 'RIGHT';
  const graph: ElkNode = {
    id: 'root',
    layoutOptions: options,
    children: nodes.map((node) => ({
      ...node,
      // Adjust the target and source handle positions based on the layout
      // direction.
      targetPosition: isHorizontal ? 'left' : 'top',
      sourcePosition: isHorizontal ? 'right' : 'bottom',

      // Hardcode a width and height for elk to use when layouting.
      width: node.width ?? getNodeWidth(node).width,
      height: node.height ?? getNodeWidth(node).height,
      properties: {
        'org.eclipse.elk.portConstraints': 'FIXED_ORDER',
      },
    })),
    edges: edges as unknown as ElkExtendedEdge[],
  };

  return elk
    .layout(graph)
    .then((layoutedGraph) => ({
      nodes: layoutedGraph?.children?.map((node) => ({
        ...node,
        // React Flow expects a position property on the node instead of `x`
        // and `y` fields.
        position: { x: node.x, y: node.y },
      })),

      edges: layoutedGraph.edges,
    }))
    .catch(console.error);
};

export const LayoutButton = () => {
  const { fitView, setNodes, setEdges, getNodes, getEdges } = useReactFlow();
  const nodes = getNodes();
  const edges = getEdges();
  const onLayout = useCallback(
    ({ direction }: { direction: 'DOWN' }) => {
      const opts = { 'elk.direction': direction, ...elkOptions };
      const ns = nodes;
      const es = edges;

      getLayoutedElements(ns, es, opts).then(
        // @ts-expect-error not for poc
        ({ nodes: layoutedNodes, edges: layoutedEdges }) => {
          setNodes(layoutedNodes);
          setEdges(layoutedEdges);

          window.requestAnimationFrame(() => fitView());
        },
      );
    },
    [nodes, edges],
  );

  // Calculate the initial layout on mount.
  useLayoutEffect(() => {
    onLayout({ direction: 'DOWN' });
  }, []);

  return (
    <ControlButton onClick={() => onLayout({ direction: 'DOWN' })}>
      <AlignHorizontalCentre02 />
    </ControlButton>
  );
};
