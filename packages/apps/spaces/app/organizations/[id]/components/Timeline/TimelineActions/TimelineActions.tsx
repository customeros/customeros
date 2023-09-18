import React from 'react';
import { Box } from '@ui/layout/Box';
import { useParams } from 'next/navigation';
import { TimelineActionLogEntryContextContextProvider } from './TimelineActionsContext/TimelineActionLogEntryContext';
import { TimelineActionButtons } from './TimelineActionButtons';
import { LogEntryTimelineAction } from './logger/LogEntryTimelineAction';
import { TimelineActionEmailContextContextProvider } from './TimelineActionsContext/TimelineActionEmailContext';
import { EmailTimelineAction } from './email/EmailTimelineAction';

interface TimelineActionsProps {
  onScrollBottom: () => void;
  invalidateQuery: () => void;
}

export const TimelineActions: React.FC<TimelineActionsProps> = ({
  onScrollBottom,
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
        <Box bg='gray.25'>
          <TimelineActionButtons />
          <EmailTimelineAction onScrollBottom={onScrollBottom} />
          <LogEntryTimelineAction onScrollBottom={onScrollBottom} />
        </Box>
      </TimelineActionLogEntryContextContextProvider>
    </TimelineActionEmailContextContextProvider>
  );
};
