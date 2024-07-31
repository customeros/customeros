import { useState, useEffect, PropsWithChildren } from 'react';

import { cn } from '@ui/utils/cn';
import {
  Modal,
  ModalPortal,
  ModalContent,
  ModalScrollBody,
} from '@ui/overlay/Modal/Modal';
import {
  useTimelineEventPreviewStateContext,
  useTimelineEventPreviewMethodsContext,
} from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface TimelinePreviewBackdropProps extends PropsWithChildren {
  onCloseModal?: () => void;
}

export const TimelinePreviewBackdrop = ({
  children,
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
    <Modal modal={false} open={isModalOpen} onOpenChange={closeModal}>
      <div
        onClick={closeModal}
        id='timeline-preview-backdrop'
        className={cn(
          'absolute top-0 bottom-0 left-0 right-0 z-40 cursor-pointer flex justify-center align-middle transition-all duration-100 linear',
        )}
        style={{
          backgroundColor: isMounted
            ? 'rgba(16, 24, 40, 0.25)'
            : 'rgba(16, 24, 40, 0)',
          backdropFilter: isMounted ? 'blur(3px)' : 'blur(0)',
        }}
      >
        <ModalPortal container={document.getElementById('main-section')}>
          <ModalContent
            placement='top'
            onInteractOutside={avoidDefaultDomBehavior}
            onOpenAutoFocus={(e) => e.preventDefault()}
            onPointerDownOutside={avoidDefaultDomBehavior}
            className={cn(
              modalContent?.__typename === 'Invoice' ? 'h-[90vh]' : 'h-auto',
              'absolute top-4 min-w-[544px] z-50 rounded-2xl max-w-fit',
            )}
          >
            <ModalScrollBody
              id='timeline-preview-card'
              onClick={(e) => e.stopPropagation()}
              className='mx-auto top-4 cursor-default bg-transparent p-0'
              onMouseDown={(e) => {
                e.stopPropagation();
              }}
            >
              {children}
            </ModalScrollBody>
          </ModalContent>
        </ModalPortal>
      </div>
    </Modal>
  );
};
