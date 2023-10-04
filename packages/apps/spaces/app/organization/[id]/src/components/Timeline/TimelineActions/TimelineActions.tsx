import React from 'react';
import { useParams } from 'next/navigation';
import { VirtuosoHandle } from 'react-virtuoso';

import { Box } from '@ui/layout/Box';
import { TimelineActionLogEntryContextContextProvider } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';
import { TimelineActionButtons } from './TimelineActionButtons';
import { TimelineActionEmailContextContextProvider } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { TimelineActionsArea } from './TimelineActionsArea';

interface TimelineActionsProps {
  onScrollBottom: () => void;
  invalidateQuery: () => void;
  virtuosoRef?: React.RefObject<VirtuosoHandle>;
}

export const TimelineActions: React.FC<TimelineActionsProps> = ({
  virtuosoRef,
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
        virtuosoRef={virtuosoRef}
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
