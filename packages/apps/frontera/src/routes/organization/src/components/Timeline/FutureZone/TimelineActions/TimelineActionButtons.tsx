import { useParams } from 'react-router-dom';
import { useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { Mail01 } from '@ui/media/icons/Mail01';
import { useStore } from '@shared/hooks/useStore';
import { AlarmClockPlus } from '@ui/media/icons/AlarmClockPlus';
import { MessageChatSquare } from '@ui/media/icons/MessageChatSquare';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';
import { useTimelineActionEmailContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionEmailContext';
import { useTimelineActionLogEntryContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionLogEntryContext';
import {
  EditorType,
  useTimelineActionContext,
} from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

interface TimelineActionButtonsProps {
  invalidateQuery: () => void;
  activeEditor: 'log-entry' | 'email' | null;
  onClick: (activeEditor: 'log-entry' | 'email' | null) => void;
}

export const TimelineActionButtons = observer(
  ({ onClick, activeEditor, invalidateQuery }: TimelineActionButtonsProps) => {
    const store = useStore();
    const { id } = useParams();
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);

    const {
      checkCanExitSafely,
      showLogEntryConfirmationDialog,
      handleExitEditorAndCleanData: handleExitLogEntryEditorAndCleanData,
      onCreateLogEntry,
    } = useTimelineActionLogEntryContext();
    const {
      checkCanExitSafely: checkCanExitEmailSafely,
      showConfirmationDialog: showEmailConfirmationDialog,
      onCreateEmail,
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

    const handleToggleEditor = (
      targetEditor: 'email' | 'log-entry' | 'reminder',
    ) => {
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

    const handleConfirmLogEntry = () => {
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

    const handleConfirmEmail = () => {
      const handleSuccess = () => {
        handleExitEmailEditorAndCleanData();
        showEditor(openOnConfirm);
      };

      onCreateEmail(handleSuccess);
    };

    const toggleEmailEditor = () => {
      handleToggleEditor('email');
      onClick('email');
    };

    const handleEmail = () => {
      if (store.ui.dirtyEditor !== null) {
        store.ui.confirmAction(store.ui.dirtyEditor, toggleEmailEditor);
      } else {
        toggleEmailEditor();
      }
    };

    const handleLogEntry = () => {
      if (store.ui.dirtyEditor === 'log-entry') {
        store.ui.confirmAction('log-entry');
      } else {
        showEditor(null);
        activeEditor !== 'log-entry' ? onClick('log-entry') : onClick(null);
      }
    };

    return (
      <>
        <div className='relative border border-gray-200 p-2 gap-2 rounded-full bg-white top-0 left-6 z-1 transform translate-y-[5px] inline-flex'>
          <Button
            size='xs'
            variant='outline'
            onClick={handleEmail}
            className='rounded-3xl'
            data-test='timeline-email-button'
            leftIcon={<Mail01 color='inherit' />}
            colorScheme={openedEditor === 'email' ? 'primary' : 'gray'}
          >
            Email
          </Button>
          <Button
            size='xs'
            variant='outline'
            className='rounded-3xl'
            onClick={handleLogEntry}
            data-test='timeline-log-button'
            leftIcon={<MessageChatSquare color='inherit' />}
            colorScheme={activeEditor === 'log-entry' ? 'primary' : 'gray'}
          >
            Log
          </Button>
          <Button
            size='xs'
            variant='outline'
            className='rounded-3xl'
            data-test='timeline-reminder-button'
            leftIcon={<AlarmClockPlus color='inherit' />}
            colorScheme={openedEditor === 'reminder' ? 'primary' : 'gray'}
            onClick={() => {
              if (!id) return;
              store.reminders.create(id);
            }}
          >
            Reminder
          </Button>
        </div>

        <ConfirmDeleteDialog
          isLoading={false}
          colorScheme='primary'
          onClose={handleDiscard}
          confirmButtonLabel='Send'
          label={`Send this email?`}
          cancelButtonLabel='Discard'
          onConfirm={handleConfirmEmail}
          isOpen={showEmailConfirmationDialog}
          description={`You have typed an unsent email. Do you want to send it, or discard it?`}
        />

        <ConfirmDeleteDialog
          isLoading={false}
          colorScheme='primary'
          onClose={handleDiscard}
          label='Log this log entry?'
          confirmButtonLabel='Log it'
          cancelButtonLabel='Discard'
          onConfirm={handleConfirmLogEntry}
          isOpen={showLogEntryConfirmationDialog}
          description='You have typed an unlogged entry. Do you want to log it to the timeline, or discard it?'
        />
      </>
    );
  },
);
