import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { TimelineEmailUsecase } from '@domain/usecases/email-composer/send-timeline-email.usecase';

import { useStore } from '@shared/hooks/useStore';
import { useChannel } from '@shared/hooks/useChannel';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import { TimelineActionLogEntryContextContextProvider } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';

import { TimelineActionsArea } from './TimelineActionsArea';
import { TimelineActionButtons } from './TimelineActionButtons';

interface TimelineActionsProps {
  invalidateQuery: () => void;
}

export const TimelineActions = observer(
  ({ invalidateQuery }: TimelineActionsProps) => {
    const store = useStore();

    const id = useParams()?.id as string;
    const [activeEditor, setActiveEditor] = useState<
      null | 'log-entry' | 'email'
    >(null);
    const { currentUserId } = useChannel(
      `finder:${store.session.value.tenant}`,
    );
    const orgId = useParams().id as string;

    const { virtuosoRef } = useTimelineRefContext();

    const updateTimelineCache = useUpdateCacheWithNewEvent(virtuosoRef);

    const handleEmailSendSuccess = (response: unknown) => {
      updateTimelineCache(response, ['GetTimeline.infinite']);
      invalidateQuery();
    };
    const emailUseCase = useMemo(
      () =>
        new TimelineEmailUsecase(
          orgId,
          [],
          currentUserId ?? '',
          handleEmailSendSuccess,
        ),
      [orgId],
    );

    return (
      <TimelineActionLogEntryContextContextProvider
        id={id}
        invalidateQuery={invalidateQuery}
      >
        <div className='bg-gray-25'>
          <TimelineActionButtons
            onClick={setActiveEditor}
            emailUseCase={emailUseCase}
            activeEditor={activeEditor}
            invalidateQuery={invalidateQuery}
          />
          <TimelineActionsArea
            emailUseCase={emailUseCase}
            activeEditor={activeEditor}
            hide={() => setActiveEditor(null)}
          />
        </div>
      </TimelineActionLogEntryContextContextProvider>
    );
  },
);
