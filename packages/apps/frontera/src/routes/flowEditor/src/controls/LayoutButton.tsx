import React, { useCallback, useLayoutEffect } from 'react';

import ELK from 'elkjs/lib/elk.bundled.js';
import { useReactFlow, ControlButton } from '@xyflow/react';

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

const getLayoutedElements = (nodes, edges, options = {}) => {
  const isHorizontal = options?.['elk.direction'] === 'RIGHT';
  const graph = {
    id: 'root',
    layoutOptions: options,
    children: nodes.map((node) => ({
      ...node,
      // Adjust the target and source handle positions based on the layout
      // direction.
      targetPosition: isHorizontal ? 'left' : 'top',
      sourcePosition: isHorizontal ? 'right' : 'bottom',

      // Hardcode a width and height for elk to use when layouting.
      width: 200,
      height: 75,
    })),
    edges: edges,
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
      <AlignHorizontalCentre02 className='size-4' />
    </ControlButton>
  );
};
