import { useState, useEffect, PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import { Card } from '@ui/presentation/Card';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/src/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface TimelinePreviewBackdropProps extends PropsWithChildren {
  onCloseModal?: () => void;
}

export const TimelinePreviewBackdrop = ({
  children,
  onCloseModal,
}: TimelinePreviewBackdropProps) => {
  const [isMounted, setIsMounted] = useState(false);
  const { isModalOpen, modalContent } = useTimelineEventPreviewStateContext();
  const { closeModal } = useTimelineEventPreviewMethodsContext();

  useEffect(() => {
    setIsMounted(isModalOpen);
  }, [isModalOpen]);

  if (!isModalOpen || !modalContent) {
    return null;
  }

  const avoidDefaultDomBehavior = (e: Event) => {
    e.preventDefault();
  };

  return (
    <Modal open={isModalOpen} modal={false} onOpenChange={closeModal}>
      <div
        className={cn(
          'absolute top-0 bottom-0 left-0 right-0 z-10 cursor-pointer flex justify-center align-middle transition-all duration-100 linear',
        )}
        id='timeline-preview-backdrop'
        style={{
          backgroundColor: isMounted
            ? 'rgba(16, 24, 40, 0.25)'
            : 'rgba(16, 24, 40, 0)',
          backdropFilter: isMounted ? 'blur(3px)' : 'blur(0)',
        }}
      >
        <ModalPortal container={document.getElementById('main-section')}>
          <ModalOverlay />
          <ModalContent
            className={cn(
              modalContent?.__typename === 'Invoice'
                ? 'w-[650px]'
                : 'w-[544px]',
              modalContent?.__typename === 'Invoice' ? 'h-[90vh]' : 'h-auto',
              'absolute top-4 min-w-[544px] bg-transparent',
            )}
            onPointerDownOutside={avoidDefaultDomBehavior}
            onInteractOutside={avoidDefaultDomBehavior}
            onOpenAutoFocus={(e) => e.preventDefault()}
          >
            <Card
              className={cn(
                modalContent?.__typename === 'Invoice'
                  ? 'w-[650px]'
                  : 'w-[544px]',
                modalContent?.__typename === 'Invoice' ? 'h-[90vh]' : 'h-auto',
                'absolute mx-auto top-4 min-w-[544px] cursor-default',
              )}
              id='timeline-preview-card'
              onMouseDown={(e) => {
                e.stopPropagation();
              }}
              onClick={(e) => e.stopPropagation()}
            >
              {children}
            </Card>
          </ModalContent>
        </ModalPortal>
      </div>
    </Modal>
  );
};
