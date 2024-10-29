import React, { useState, useEffect, useCallback } from 'react';

import { Controls, useStore, useReactFlow, ControlButton } from '@xyflow/react';

import { ZoomIn } from '@ui/media/icons/ZoomIn.tsx';
import { ZoomOut } from '@ui/media/icons/ZoomOut.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Maximize02 } from '@ui/media/icons/Maximize02.tsx';
import { AlignHorizontalCentre02 } from '@ui/media/icons/AlignHorizontalCentre02.tsx';

import { useLayout } from '../hooks/useLayout.ts';

export const FlowBuilderToolbar = () => {
  const { zoomIn, zoomOut, fitView } = useReactFlow();
  const { organizeLayout } = useLayout();

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
            data-test='flow-zoom-in'
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
            data-test='flow-zoom-out'
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
            data-test='flow-fit-to-view'
            className='h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50'
          >
            <Maximize02 className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
      <Tooltip label={'Tidy up blocks'}>
        <div>
          <ControlButton
            onClick={organizeLayout}
            data-test='flow-tidy-up'
            className='rounded-r-lg h-[36px] w-[36px] hover:bg-gray-50 focus:bg-gray-50'
          >
            <AlignHorizontalCentre02 className='size-4 text-gray-500' />
          </ControlButton>
        </div>
      </Tooltip>
    </Controls>
  );
};
