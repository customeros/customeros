import React from 'react';
import { Box } from '@ui/layout/Box';
import { useParams } from 'next/navigation';
import { TimelineActionLogEntryContextContextProvider } from '@organization/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';
import { TimelineActionButtons } from './TimelineActionButtons';
import { TimelineActionEmailContextContextProvider } from '@organization/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { TimelineActionsArea } from './TimelineActionsArea';

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
          <TimelineActionButtons invalidateQuery={invalidateQuery} />
          <TimelineActionsArea onScrollBottom={onScrollBottom} />
        </Box>
      </TimelineActionLogEntryContextContextProvider>
    </TimelineActionEmailContextContextProvider>
  );
};
