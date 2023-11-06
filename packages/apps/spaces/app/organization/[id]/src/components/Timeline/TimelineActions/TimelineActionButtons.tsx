import React, { FC, useRef, useState, useEffect } from 'react';

import { Box } from '@ui/layout/Box';
import { Button } from '@ui/form/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import { useTimelineActionEmailContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionEmailContext';
import { useTimelineActionLogEntryContext } from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionLogEntryContext';
import {
  EditorType,
  useTimelineActionContext,
} from '@organization/src/components/Timeline/TimelineActions/context/TimelineActionContext';

export const TimelineActionButtons: FC<{ invalidateQuery: () => void }> = ({
  invalidateQuery,
}) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const {
    checkCanExitSafely,
    showLogEntryConfirmationDialog,
    closeConfirmationDialog: closeLogEntryConfirmationDialog,
    handleExitEditorAndCleanData: handleExitLogEntryEditorAndCleanData,
    onCreateLogEntry,
  } = useTimelineActionLogEntryContext();
  const {
    checkCanExitSafely: checkCanExitEmailSafely,
    showConfirmationDialog: showEmailConfirmationDialog,
    closeConfirmationDialog: closeEmailConfirmationDialog,
    handleExitEditorAndCleanData: handleExitEmailEditorAndCleanData,
  } = useTimelineActionEmailContext();
  const { openedEditor, showEditor } = useTimelineActionContext();
  const [openOnConfirm, setOpenOnConfirm] = useState<null | EditorType>(null);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const handleToggleEditor = (targetEditor: 'email' | 'log-entry') => {
    if (openedEditor === null) {
      showEditor(targetEditor);

      return;
    }

    if (openedEditor === targetEditor) {
      const canClose =
        targetEditor === 'email'
          ? checkCanExitEmailSafely()
          : checkCanExitSafely();
      if (canClose) showEditor(null);

      return;
    }

    setOpenOnConfirm(targetEditor);
    const canClose =
      targetEditor === 'log-entry'
        ? checkCanExitEmailSafely()
        : checkCanExitSafely();

    if (canClose) {
      setOpenOnConfirm(null);
      showEditor(targetEditor);
    }
  };

  const handleDiscard = () => {
    if (showEmailConfirmationDialog) {
      handleExitEmailEditorAndCleanData();
    } else {
      handleExitLogEntryEditorAndCleanData();
    }

    showEditor(openOnConfirm);
  };
  const handleConfirm = () => {
    onCreateLogEntry({
      onSuccess: () => {
        handleExitLogEntryEditorAndCleanData();
        timeoutRef.current = setTimeout(() => {
          invalidateQuery();
        }, 500);
      },
      onSettled: () => {
        showEditor(openOnConfirm);
      },
    });
  };

  const handleCloseConfirmationModal = () => {
    setOpenOnConfirm(null);

    return showEmailConfirmationDialog
      ? closeEmailConfirmationDialog()
      : closeLogEntryConfirmationDialog();
  };

  return (
    <ButtonGroup
      position='sticky'
      border='1px solid'
      borderColor='gray.200'
      p='2'
      borderRadius='full'
      bg='white'
      top='0'
      left='6'
      zIndex='1'
      transform='translateY(5px)'
    >
      <Button
        variant='outline'
        onClick={() => handleToggleEditor('email')}
        borderRadius='3xl'
        size='xs'
        colorScheme={openedEditor === 'email' ? 'primary' : 'gray'}
        leftIcon={<Mail01 color='inherit' />}
      >
        Email
      </Button>
      <Button
        variant='outline'
        onClick={() => handleToggleEditor('log-entry')}
        borderRadius='3xl'
        size='xs'
        colorScheme={openedEditor === 'log-entry' ? 'primary' : 'gray'}
        leftIcon={<MessageChatSquare color='inherit' />}
      >
        Log
      </Button>

      <ConfirmDeleteDialog
        label={`Discard this email?`}
        description={`Saving draft log entries is not possible at the moment. Would you like to continue to discard this email?`}
        confirmButtonLabel={`Discard email`}
        isOpen={showEmailConfirmationDialog}
        onClose={handleCloseConfirmationModal}
        onConfirm={handleDiscard}
        isLoading={false}
      />

      <ConfirmDeleteDialog
        colorScheme='primary'
        label='Log this log entry?'
        description='You have typed an unlogged entry. Do you want to log it to the timeline, or discard it?'
        confirmButtonLabel='Log it'
        cancelButtonLabel='Discard'
        isOpen={showLogEntryConfirmationDialog}
        onClose={handleDiscard}
        onConfirm={handleConfirm}
        isLoading={false}
        icon={
          <Box>
            <MessageChatSquare color='primary.700' boxSize='inherit' />
          </Box>
        }
      />
    </ButtonGroup>
  );
};
