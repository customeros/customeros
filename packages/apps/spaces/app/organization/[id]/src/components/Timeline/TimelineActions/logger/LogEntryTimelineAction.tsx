import React, { useEffect, useMemo, useRef } from 'react';
import { Box } from '@ui/layout/Box';
import { Logger } from './components/Logger';
import { useTimelineActionContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionContext';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';

export const LogEntryTimelineAction: React.FC = () => {
  const { virtuosoRef } = useTimelineRefContext();

  const { openedEditor } = useTimelineActionContext();
  const isLogEntryEditorOpen = useMemo(
    () => openedEditor === 'log-entry',
    [openedEditor],
  );
  const logEntryWrapperRef = useRef(null);

  useEffect(() => {
    if (isLogEntryEditorOpen) {
      virtuosoRef?.current?.scrollBy({ top: 300 });
    }
  }, [isLogEntryEditorOpen, virtuosoRef]);

  return (
    <>
      {isLogEntryEditorOpen && (
        <Box
          ref={logEntryWrapperRef}
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
      )}
    </>
  );
};
