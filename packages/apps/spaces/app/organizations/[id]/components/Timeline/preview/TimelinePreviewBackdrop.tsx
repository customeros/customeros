import { useState, useEffect, PropsWithChildren } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Card } from '@ui/presentation/Card';
import { ScaleFade } from '@ui/transitions/ScaleFade';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/context/TimelineEventPreviewContext';

interface TimelinePreviewBackdropProps extends PropsWithChildren {
  onCloseModal: () => void;
}

export const TimelinePreviewBackdrop = ({
  onCloseModal,
  children,
}: TimelinePreviewBackdropProps) => {
  const [isMounted, setIsMounted] = useState(false); // needed for delaying the backdrop filter
  const { isModalOpen, modalContent } = useTimelineEventPreviewContext();

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
      zIndex={1}
      cursor='pointer'
      backdropFilter='blur(3px)'
      justify='center'
      background={isMounted ? 'rgba(16, 24, 40, 0.45)' : 'rgba(16, 24, 40, 0)'}
      align='center'
      transition='all 0.1s linear'
      onClick={onCloseModal}
    >
      <ScaleFade
        in={isModalOpen}
        style={{
          position: 'absolute',
          marginInline: 'auto',
          top: '1rem',
          width: '544px',
          minWidth: '544px',
        }}
      >
        <Card
          size='lg'
          position='absolute'
          mx='auto'
          top='4'
          w='544px'
          minW='544px'
          cursor='default'
          onClick={(e) => e.stopPropagation()}
        >
          {children}
        </Card>
      </ScaleFade>
    </Flex>
  );
};
