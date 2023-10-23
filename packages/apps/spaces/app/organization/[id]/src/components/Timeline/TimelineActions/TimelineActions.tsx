import React from 'react';
import { useParams } from 'next/navigation';
import dynamic from 'next/dynamic';

import { Box } from '@ui/layout/Box';
import { TimelineActionLogEntryContextContextProvider } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';
import { TimelineActionButtons } from './TimelineActionButtons';
// import { TimelineActionEmailContextContextProvider } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { TimelineActionsArea } from './TimelineActionsArea';

const TimelineActionEmailContextContextProvider = dynamic(
  () =>
    import(
      '@organization/src/components/Timeline/TimelineActions/context/TimelineActionEmailContext'
    ),
  {
    ssr: false,
  },
);

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
        <Box bg='gray.25'>
          <TimelineActionButtons invalidateQuery={invalidateQuery} />
          <TimelineActionsArea />
        </Box>
      </TimelineActionLogEntryContextContextProvider>
    </TimelineActionEmailContextContextProvider>
  );
};
