import React, { useState } from 'react';
import { useParams } from 'react-router-dom';

import { TimelineActionLogEntryContextContextProvider } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';

import { TimelineActionsArea } from './TimelineActionsArea';
import { TimelineActionButtons } from './TimelineActionButtons';
import TimelineActionEmailContextContextProvider from './context/TimelineActionEmailContext';

interface TimelineActionsProps {
  invalidateQuery: () => void;
}

export const TimelineActions: React.FC<TimelineActionsProps> = ({
  invalidateQuery,
}) => {
  const id = useParams()?.id as string;
  const [activeEditor, setActiveEditor] = useState<
    null | 'log-entry' | 'email'
  >(null);

  return (
    <TimelineActionEmailContextContextProvider
      id={id}
      invalidateQuery={invalidateQuery}
    >
      <TimelineActionLogEntryContextContextProvider
        id={id}
        invalidateQuery={invalidateQuery}
      >
        <div className='bg-gray-25'>
          <TimelineActionButtons
            onClick={setActiveEditor}
            activeEditor={activeEditor}
            invalidateQuery={invalidateQuery}
          />
          <TimelineActionsArea
            activeEditor={activeEditor}
            hide={() => setActiveEditor(null)}
          />
        </div>
      </TimelineActionLogEntryContextContextProvider>
    </TimelineActionEmailContextContextProvider>
  );
};
