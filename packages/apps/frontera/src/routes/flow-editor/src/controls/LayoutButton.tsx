import React, { useCallback } from 'react';

import { FitViewOptions } from '@xyflow/react';
import { useReactFlow, ControlButton } from '@xyflow/react';

import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02.tsx';

import { useLayoutedElements } from '../hooks';

import '@xyflow/react/dist/style.css';

export const LayoutButton = () => {
  const { fitView, setNodes, setEdges, getNodes, getEdges } = useReactFlow();
  const { getLayoutedElements } = useLayoutedElements();
  const nodes = getNodes();
  const edges = getEdges();
  const onLayout = useCallback(() => {
    const ns = nodes;
    const es = edges;

    getLayoutedElements(ns, es).then(
      // @ts-expect-error not for poc
      ({ nodes: layoutedNodes, edges: layoutedEdges }) => {
        window.requestAnimationFrame(() => {
          const fitViewOptions: FitViewOptions = {
            padding: 0.1,
            maxZoom: 1,
            minZoom: 1,
            duration: 10,
            nodes: layoutedNodes,
          };

          fitView(fitViewOptions);
        });
        setNodes(layoutedNodes);
        setEdges(layoutedEdges);
      },
    );
  }, [nodes, edges]);

  return (
    <ControlButton onClick={() => onLayout()}>
      <AlignHorizontalCentre02 />
    </ControlButton>
  );
};
