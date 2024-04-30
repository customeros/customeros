import React from 'react';
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

  return (
    <TimelineActionEmailContextContextProvider
      id={id}
      invalidateQuery={invalidateQuery}
    >
      <TimelineActionLogEntryContextContextProvider
        id={id}
        invalidateQuery={invalidateQuery}
      >
        <div id='mariana' className='bg-gray-25'>
          <TimelineActionButtons invalidateQuery={invalidateQuery} />
          <TimelineActionsArea />
        </div>
      </TimelineActionLogEntryContextContextProvider>
    </TimelineActionEmailContextContextProvider>
  );
};
