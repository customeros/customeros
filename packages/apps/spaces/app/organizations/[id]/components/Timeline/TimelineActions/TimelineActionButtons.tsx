import React, { useState } from 'react';
import { Button } from '@ui/form/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import { Mail01 } from '@ui/media/icons/Mail01';
import { useTimelineActionLogEntryContext } from './TimelineActionsContext/TimelineActionLogEntryContext';
import { useTimelineActionEmailContext } from './TimelineActionsContext/TimelineActionEmailContext';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  EditorType,
  useTimelineActionContext,
} from './TimelineActionsContext/TimelineActionContext';

export const TimelineActionButtons = () => {
  const {
    checkCanExitSafely,
    showLogEntryConfirmationDialog,
    closeConfirmationDialog: closeLogEntryConfirmationDialog,
    handleExitEditorAndCleanData: handleExitLogEntryEditorAndCleanData,
  } = useTimelineActionLogEntryContext();
  const {
    checkCanExitSafely: checkCanExitEmailSafely,
    showConfirmationDialog: showEmailConfirmationDialog,
    closeConfirmationDialog: closeEmailConfirmationDialog,
    handleExitEditorAndCleanData: handleExitEmailEditorAndCleanData,
  } = useTimelineActionEmailContext();
  const { openedEditor, showEditor } = useTimelineActionContext();
  const [openOnConfirm, setOpenOnConfirm] = useState<null | EditorType>(null);

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
        label={`Discard this ${
          showEmailConfirmationDialog ? 'email' : 'log entry'
        }?`}
        description={`Saving draft log entries is not possible at the moment. Would you like to continue to discard this ${
          showEmailConfirmationDialog ? 'email' : 'entry'
        }?`}
        confirmButtonLabel={`Discard ${
          showEmailConfirmationDialog ? 'email' : 'entry'
        }`}
        isOpen={showLogEntryConfirmationDialog || showEmailConfirmationDialog}
        onClose={handleCloseConfirmationModal}
        onConfirm={handleDiscard}
        isLoading={false}
      />
    </ButtonGroup>
  );
};
