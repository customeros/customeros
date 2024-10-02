import React, { useState, useEffect, useCallback } from 'react';

import {
  Controls,
  useStore,
  useReactFlow,
  ControlButton,
  FitViewOptions,
} from '@xyflow/react';

import { ZoomIn } from '@ui/media/icons/ZoomIn.tsx';
import { ZoomOut } from '@ui/media/icons/ZoomOut.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Maximize02 } from '@ui/media/icons/Maximize02.tsx';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02.tsx';

import { useLayoutedElements } from '../hooks';

export const FlowBuilderToolbar = () => {
  const { zoomIn, zoomOut, fitView, setNodes, setEdges, getNodes, getEdges } =
    useReactFlow();
  const { getLayoutedElements } = useLayoutedElements();

  const zoom = useStore((store) => store.transform[2]);
  const [canZoomIn, setCanZoomIn] = useState(true);
  const [canZoomOut, setCanZoomOut] = useState(true);

  useEffect(() => {
    setCanZoomIn(zoom < 5);
    setCanZoomOut(zoom > 0.1);
  }, [zoom]);

  const handleZoomIn = useCallback(() => {
    zoomIn();
  }, [zoomIn]);

  const handleZoomOut = useCallback(() => {
    zoomOut();
  }, [zoomOut]);

  const handleFitView = () => fitView({ duration: 800 });

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
  }, [nodes, edges, getLayoutedElements, fitView, setNodes, setEdges]);

  return (
    <Controls
      showZoom={false}
      showFitView={false}
      position='bottom-left'
      showInteractive={false}
      orientation='horizontal'
      className='bg-white rounded-lg border border-grayModern-300'
    >
      <Tooltip label={canZoomIn && 'Zoom in'}>
        <div>
          <ControlButton
            disabled={!canZoomIn}
            onClick={handleZoomIn}
            className={`rounded-l-lg h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            <ZoomIn className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
      <Tooltip label={'Zoom out'}>
        <div>
          <ControlButton
            disabled={!canZoomOut}
            onClick={handleZoomOut}
            className={`h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            <ZoomOut className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
      <Tooltip label={'Fit to view'}>
        <div>
          <ControlButton
            onClick={handleFitView}
            className='h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50'
          >
            <Maximize02 className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
      <Tooltip label={'Tidy up blocks'}>
        <div>
          <ControlButton
            onClick={onLayout}
            className='rounded-r-lg h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50'
          >
            <AlignHorizontalCentre02 className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
    </Controls>
  );
};
