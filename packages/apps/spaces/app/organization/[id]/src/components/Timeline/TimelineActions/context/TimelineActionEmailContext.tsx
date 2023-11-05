'use client';

import { useForm } from 'react-inverted-form';
import { useSearchParams } from 'next/navigation';
import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useSession } from 'next-auth/react';
import { useRemirror } from '@remirror/react';

import { useDisclosure } from '@ui/utils';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { handleSendEmail } from '@organization/src/components/Timeline/events/email/compose-email/utils';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';
import { useUpdateCacheWithNewEvent } from '@organization/src/components/Timeline/hooks/updateCacheWithNewEvent';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail.dto';

import { useTimelineActionContext } from './TimelineActionContext';

export const noop = () => undefined;

interface TimelineActionEmailContextContextMethods {
  state: any;
  formId: string;
  remirrorProps: any;
  isSending: boolean;
  onCreateEmail: () => void;
  showConfirmationDialog: boolean;
  checkCanExitSafely: () => boolean;
  closeConfirmationDialog: () => void;
  handleExitEditorAndCleanData: () => void;
}

const TimelineActionEmailContextContext =
  createContext<TimelineActionEmailContextContextMethods>({
    checkCanExitSafely: () => false,
    onCreateEmail: noop,
    handleExitEditorAndCleanData: noop,
    closeConfirmationDialog: noop,
    remirrorProps: null,
    isSending: false,
    showConfirmationDialog: false,
    formId: '',
    state: null,
  });

export const useTimelineActionEmailContext = () => {
  return useContext(TimelineActionEmailContextContext);
};

export const TimelineActionEmailContextContextProvider = ({
  children,
  invalidateQuery,
  id = '',
}: PropsWithChildren<{
  id: string;
  invalidateQuery: () => void;
}>) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const searchParams = useSearchParams();
  const { data: session } = useSession();

  const [isSending, setIsSending] = useState(false);
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const { closeEditor } = useTimelineActionContext();
  const [timelineMeta] = useTimelineMeta();

  const queryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
  const { virtuosoRef } = useTimelineRefContext();
  const updateTimelineCache = useUpdateCacheWithNewEvent(virtuosoRef);
  const formId = 'compose-email-timeline-footer';

  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: [],
    cc: [],
    bcc: [],
    subject: '',
    content: '',
  });
  const { state, reset, setDefaultValues } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
  });
  const handleResetEditor = () => {
    setDefaultValues(defaultValues);
    const context = remirrorProps.getContext();
    if (context) {
      context.commands.resetContent();
    }
    reset();
  };

  const handleEmailSendSuccess = async (response: any) => {
    await updateTimelineCache(response, queryKey);

    // no timeout needed is this case as the event id is created when this is called
    invalidateQuery();
    setIsSending(false);
    handleResetEditor();
  };
  const handleEmailSendError = () => {
    setIsSending(false);
  };
  const onCreateEmail = () => {
    const to = [...state.values.to].map(({ value }) => value);
    const cc = [...state.values.cc].map(({ value }) => value);
    const bcc = [...state.values.bcc].map(({ value }) => value);
    const params = new URLSearchParams(searchParams?.toString() ?? '');

    setIsSending(true);
    const id = params.get('events');

    return handleSendEmail(
      state.values.content,
      to,
      cc,
      bcc,
      id,
      state.values.subject,
      handleEmailSendSuccess,
      handleEmailSendError,
      session?.user,
    );
  };

  const handleExitEditorAndCleanData = () => {
    handleResetEditor();

    onClose();
    closeEditor();
  };

  const handleCheckCanExitSafely = () => {
    const { content, ...values } = state.values;

    const isFormEmpty = !content.length || content === `<p style=""></p>`;
    const areFieldsEmpty = Object.values(values).every((e) => !e.length);
    const showEmailEditorConfirmationDialog = !isFormEmpty || !areFieldsEmpty;
    if (showEmailEditorConfirmationDialog) {
      onOpen();

      return false;
    } else {
      handleResetEditor();
      onClose();

      return true;
    }
  };

  useEffect(() => {
    const handleCloseOnEsc = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handleCheckCanExitSafely();
      }
    };
    if (isOpen) {
      document.addEventListener('keydown', handleCloseOnEsc);
    }
    if (!isOpen) {
      document.removeEventListener('keydown', handleCloseOnEsc);
    }

    return () => {
      document.removeEventListener('keydown', handleCloseOnEsc);
    };
  }, [isOpen]);

  return (
    <TimelineActionEmailContextContext.Provider
      value={{
        checkCanExitSafely: handleCheckCanExitSafely,
        handleExitEditorAndCleanData,
        closeConfirmationDialog: onClose,
        onCreateEmail,
        remirrorProps,
        isSending,
        showConfirmationDialog: isOpen,
        formId,
        state,
      }}
    >
      {children}
    </TimelineActionEmailContextContext.Provider>
  );
};

export default TimelineActionEmailContextContextProvider;
