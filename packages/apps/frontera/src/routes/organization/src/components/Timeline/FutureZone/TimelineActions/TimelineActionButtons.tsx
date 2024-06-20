import { useParams } from 'react-router-dom';
import { FC, useRef, useState, useEffect } from 'react';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { Send03 } from '@ui/media/icons/Send03';
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

export const TimelineActionButtons: FC<{ invalidateQuery: () => void }> =
  observer(({ invalidateQuery }) => {
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

    return (
      <>
        <div className='relative border border-gray-200 p-2 gap-2 rounded-full bg-white top-0 left-6 z-1 transform translate-y-[5px] inline-flex'>
          <Button
            variant='outline'
            onClick={() => handleToggleEditor('email')}
            size='xs'
            className='rounded-3xl'
            colorScheme={openedEditor === 'email' ? 'primary' : 'gray'}
            leftIcon={<Mail01 color='inherit' />}
          >
            Email
          </Button>
          <Button
            className='rounded-3xl'
            variant='outline'
            onClick={() => handleToggleEditor('log-entry')}
            size='xs'
            colorScheme={openedEditor === 'log-entry' ? 'primary' : 'gray'}
            leftIcon={<MessageChatSquare color='inherit' />}
          >
            Log
          </Button>
          <Button
            className='rounded-3xl'
            variant='outline'
            onClick={() => {
              if (!id) return;
              store.reminders.create(id);
            }}
            size='xs'
            colorScheme={openedEditor === 'reminder' ? 'primary' : 'gray'}
            leftIcon={<AlarmClockPlus color='inherit' />}
          >
            Reminder
          </Button>
        </div>

        <ConfirmDeleteDialog
          colorScheme='primary'
          label={`Send this email?`}
          description={`You have typed an unsent email. Do you want to send it, or discard it?`}
          confirmButtonLabel='Send'
          cancelButtonLabel='Discard'
          isOpen={showEmailConfirmationDialog}
          onClose={handleDiscard}
          onConfirm={handleConfirmEmail}
          isLoading={false}
          icon={<Send03 className='text-primary-700' />}
        />

        <ConfirmDeleteDialog
          colorScheme='primary'
          label='Log this log entry?'
          description='You have typed an unlogged entry. Do you want to log it to the timeline, or discard it?'
          confirmButtonLabel='Log it'
          cancelButtonLabel='Discard'
          isOpen={showLogEntryConfirmationDialog}
          onClose={handleDiscard}
          onConfirm={handleConfirmLogEntry}
          isLoading={false}
          icon={<MessageChatSquare className='text-primary-700' />}
        />
      </>
    );
  });
