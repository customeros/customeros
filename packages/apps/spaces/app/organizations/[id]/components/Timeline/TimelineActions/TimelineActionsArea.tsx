import React from 'react';
import { Box } from '@ui/layout/Box';
import { LogEntryTimelineAction } from './logger/LogEntryTimelineAction';
import { useTimelineActionContext } from './TimelineActionsContext/TimelineActionContext';
import { EmailTimelineAction } from './email/EmailTimelineAction';

interface TimelineActionsAreaProps {
  onScrollBottom: () => void;
}

export const TimelineActionsArea: React.FC<TimelineActionsAreaProps> = ({
  onScrollBottom,
}) => {
  const { openedEditor } = useTimelineActionContext();

  return (
    <Box
      bg={'#F9F9FB'}
      borderTopColor='gray.200'
      pt={openedEditor !== null ? 6 : 0}
      pb={openedEditor !== null ? 2 : 8}
      mt={-4}
    >
      <EmailTimelineAction onScrollBottom={onScrollBottom} />
      <LogEntryTimelineAction onScrollBottom={onScrollBottom} />
    </Box>
  );
};
