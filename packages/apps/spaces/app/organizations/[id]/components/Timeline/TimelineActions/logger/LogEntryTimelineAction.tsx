import React, { useEffect, useRef } from 'react';
import { SlideFade } from '@ui/transitions/SlideFade';
import { Box } from '@ui/layout/Box';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { Logger } from '@organization/components/Timeline/TimelineActions/logger/components/Logger';
import { useTimelineActionLogEntryContext } from '../TimelineActionsContext/TimelineActionLogEntryContext';

interface LogEntryTimelineActionProps {
  onScrollBottom: () => void;
}

export const LogEntryTimelineAction: React.FC<LogEntryTimelineActionProps> = ({
  onScrollBottom,
}) => {
  const {
    isLogEntryEditorOpen,
    showConfirmationDialog,
    closeConfirmationDialog,
    handleExitEditorAndCleanData,
  } = useTimelineActionLogEntryContext();
  const virtuoso = useRef(null);

  useEffect(() => {
    if (isLogEntryEditorOpen) {
      onScrollBottom();
    }
  }, [isLogEntryEditorOpen]);

  return (
    <>
      <Box
        bg={'#F9F9FB'}
        borderTop='1px dashed'
        borderTopColor='gray.200'
        pt={isLogEntryEditorOpen ? 6 : 0}
        pb={isLogEntryEditorOpen ? 2 : 8}
        mt={-4}
      >
        {isLogEntryEditorOpen && (
          <SlideFade in={true}>
            <Box
              ref={virtuoso}
              borderRadius={'md'}
              boxShadow={'lg'}
              m={6}
              mt={2}
              p={6}
              bg={'white'}
              border='1px solid var(--gray-100, #F2F4F7)'
            >
              <Logger />
            </Box>
          </SlideFade>
        )}
      </Box>
      <ConfirmDeleteDialog
        label='Discard this log entry?'
        description='Saving draft log entries is not possible at the moment. Would you like to continue to discard this entry?'
        confirmButtonLabel='Discard entry'
        isOpen={showConfirmationDialog}
        onClose={closeConfirmationDialog}
        onConfirm={() => handleExitEditorAndCleanData()}
        isLoading={false}
      />
    </>
  );
};
