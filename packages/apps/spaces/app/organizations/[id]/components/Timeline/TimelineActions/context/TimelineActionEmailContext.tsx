import React, {
  useState,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useForm } from 'react-inverted-form';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useDisclosure } from '@ui/utils';
import { handleSendEmail } from '@organization/components/Timeline/events/email/compose-email/utils';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { useSearchParams } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { useTimelineActionContext } from './TimelineActionContext';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useTimelineMeta } from '@organization/components/Timeline/shared/state';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/hooks/updateCacheWithNewEvent';
import { LogEntry } from '@graphql/types';

export const noop = () => undefined;

interface TimelineActionEmailContextContextMethods {
  checkCanExitSafely: () => boolean;
  handleExitEditorAndCleanData: () => void;
  closeConfirmationDialog: () => void;
  onCreateEmail: () => void;
  remirrorProps: any;
  isSending: boolean;
  showConfirmationDialog: boolean;
  formId: string;
  state: any;
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
  invalidateQuery: () => void;
  id: string;
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
  const updateTimelineCache = useUpdateCacheWithNewEvent<LogEntry>();

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
    const params = new URLSearchParams(searchParams ?? '');

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
