import { useEffect } from 'react';

import { useKey } from 'rooks';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { ComposeEmailContainer } from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmailContainer';
import { useTimelineActionContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

export const EmailTimelineAction = ({
  emailUseCase,
}: {
  emailUseCase: TimelineEmailUsecase;
}) => {
  const { virtuosoRef } = useTimelineRefContext();
  const { closeEditor } = useTimelineActionContext();

  useEffect(() => {
    virtuosoRef?.current?.scrollBy({ top: 300 });
  }, [virtuosoRef]);

  const handleClose = () => {
    const canClose = emailUseCase.canExitSafely;

    if (canClose) {
      closeEditor();
    }
  };

  const handleDiscard = () => {
    emailUseCase.resetEditor();
    closeEditor();
  };

  useKey('Escape', () => {
    handleClose();
  });

  return (
    <div className='rounded-md shadow-lg m-6 mt-2 bg-white border border-gray-100 max-w-[800px]'>
      <div className='rounded-b-2xl py-4 px-6 overflow-visible pt-1 bg-white rounded-lg max-h-[100%]'>
        <ComposeEmailContainer
          modal={false}
          onClose={handleClose}
          onDiscard={handleDiscard}
          emailUseCase={emailUseCase}
        />
      </div>
    </div>
  );
};
