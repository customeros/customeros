import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { Logger } from './components/Logger';
import { useTimelineActionContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionContext';

interface LogEntryTimelineActionProps {
  onScrollBottom: () => void;
}

export const LogEntryTimelineAction: React.FC<LogEntryTimelineActionProps> = ({
  onScrollBottom,
}) => {
  const { openedEditor } = useTimelineActionContext();
  const isLogEntryEditorOpen = openedEditor === 'log-entry';
  const virtuoso = useRef(null);

  useEffect(() => {
    if (isLogEntryEditorOpen) {
      onScrollBottom();
    }
  }, [isLogEntryEditorOpen]);

  return (
    <>
      {isLogEntryEditorOpen && (
        <SlideFade in={true}>
          <Box
            ref={virtuoso}
            borderRadius={'md'}
            boxShadow={'lg'}
            m={6}
            mt={2}
            p={6}
            pt={4}
            bg={'white'}
            border='1px solid'
            borderColor='gray.100'
          >
            <Logger />
          </Box>
        </SlideFade>
      )}
    </>
  );
};
