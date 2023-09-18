import React from 'react';
import { Button } from '@ui/form/Button';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import { Mail01 } from '@ui/media/icons/Mail01';
import { useTimelineActionLogEntryContext } from './TimelineActionsContext/TimelineActionLogEntryContext';
import { useTimelineActionEmailContext } from './TimelineActionsContext/TimelineActionEmailContext';

export const TimelineActionButtons = () => {
  const { showLogEntryEditor, isLogEntryEditorOpen, closeLogEntryEditor } =
    useTimelineActionLogEntryContext();
  const { isEmailEditorOpen, closeEmailEditor, showEmailEditor } =
    useTimelineActionEmailContext();

  const handleToggleEmailEditor = () => {
    if (!isEmailEditorOpen) {
      if (isLogEntryEditorOpen) {
        closeLogEntryEditor({
          openEmailEditor: showEmailEditor,
        });
        return;
      }
      showEmailEditor();
    } else {
      closeEmailEditor();
    }
  };

  const handleToggleLogger = () => {
    if (!isLogEntryEditorOpen) {
      if (isEmailEditorOpen) {
        closeEmailEditor({
          openLogEntry: showLogEntryEditor,
        });
        return;
      }

      showLogEntryEditor();
    } else {
      closeLogEntryEditor();
    }
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
        onClick={() => handleToggleEmailEditor()}
        borderRadius='3xl'
        size='xs'
        colorScheme={isEmailEditorOpen ? 'primary' : 'gray'}
        leftIcon={<Mail01 color='inherit' />}
      >
        Email
      </Button>
      <Button
        variant='outline'
        onClick={handleToggleLogger}
        borderRadius='3xl'
        size='xs'
        colorScheme={isLogEntryEditorOpen ? 'primary' : 'gray'}
        leftIcon={<MessageChatSquare color='inherit' />}
      >
        Log
      </Button>
    </ButtonGroup>
  );
};
