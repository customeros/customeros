import { useMemo, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase.ts';

import { useStore } from '@shared/hooks/useStore';
import { useChannel } from '@shared/hooks/useChannel';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { ComposeEmailContainer } from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmailContainer';
import { useTimelineActionContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';
import { useTimelineActionEmailContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionEmailContext';

export const EmailTimelineAction = () => {
  const { virtuosoRef } = useTimelineRefContext();
  const { closeEditor } = useTimelineActionContext();
  const store = useStore();
  const orgId = useParams().id as string;

  useEffect(() => {
    virtuosoRef?.current?.scrollBy({ top: 300 });
  }, [virtuosoRef]);

  const { currentUserId } = useChannel(`finder:${store.session.value.tenant}`);

  const emailUseCase = useMemo(
    () => new TimelineEmailUsecase(orgId, [], currentUserId),
    [orgId],
  );

  const handleClose = () => {
    const canClose = emailUseCase.canExitSafely;

    if (canClose) {
      closeEditor();
    }
  };

  return (
    <div className='rounded-md shadow-lg m-6 mt-2 bg-white border border-gray-100 max-w-[800px]'>
      <ComposeEmailContainer
        modal={false}
        onClose={handleClose}
        onDiscard={closeEditor}
        emailUseCase={emailUseCase}
      />
    </div>
  );
};
