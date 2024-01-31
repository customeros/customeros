import { useRef, useState, useEffect, PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Card } from '@ui/presentation/Card';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

interface TimelinePreviewBackdropProps extends PropsWithChildren {
  onCloseModal?: () => void;
}

export const TimelinePreviewBackdrop = ({
  children,
  onCloseModal,
}: TimelinePreviewBackdropProps) => {
  const mouseTarget = useRef<string | null>(null);
  const [isMounted, setIsMounted] = useState(false); // needed for delaying the backdrop filter
  const { isModalOpen, modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  useEffect(() => {
    setIsMounted(isModalOpen);
  }, [isModalOpen]);

  if (!isModalOpen || !modalContent) {
    return null;
  }

  return (
    <Flex
      position='absolute'
      top='0'
      bottom='0'
      left='0'
      right='0'
      zIndex={10}
      cursor='pointer'
      backdropFilter='blur(3px)'
      justify='center'
      id='timeline-preview-backdrop'
      background={isMounted ? 'rgba(16, 24, 40, 0.45)' : 'rgba(16, 24, 40, 0)'}
      align='center'
      transition='all 0.1s linear'
      onMouseDown={(e) => {
        e.stopPropagation();
        mouseTarget.current = e.currentTarget.id;
      }}
      onMouseUp={() => {
        if (mouseTarget?.current === 'timeline-preview-backdrop') {
          closeModal();
          onCloseModal?.();
        } else {
          mouseTarget.current = null;
        }
      }}
    >
      <ScaleFade
        in={isModalOpen}
        style={{
          position: 'absolute',
          marginInline: 'auto',
          top: '1rem',
          width: modalContent?.__typename === 'Invoice' ? '650px' : '544px',
          height: modalContent?.__typename === 'Invoice' ? '90vh' : 'auto',
          minWidth: '544px',
        }}
      >
        <Card
          size='lg'
          position='absolute'
          mx='auto'
          top='4'
          w={modalContent?.__typename === 'Invoice' ? '650px' : '544px'}
          h={modalContent?.__typename === 'Invoice' ? '90vh' : 'auto'}
          minW='544px'
          cursor='default'
          id='timeline-preview-card'
          onMouseDown={(e) => {
            e.stopPropagation();
          }}
          onClick={(e) => e.stopPropagation()}
        >
          {children}
        </Card>
      </ScaleFade>
    </Flex>
  );
};
