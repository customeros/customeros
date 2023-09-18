import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
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

export const noop = () => undefined;

interface TimelineActionEmailContextContextMethods {
  showEmailEditor: () => void;
  closeEmailEditor: (prop?: { openLogEntry: () => void }) => void;
  handleExitEditorAndCleanData: () => void;
  closeConfirmationDialog: () => void;
  onCreateEmail: any;
  isEmailEditorOpen: boolean;
  remirrorProps: any;
  isSending: boolean;
  showConfirmationDialog: boolean;
  formId: string;
  state: any;
}

const TimelineActionEmailContextContext =
  createContext<TimelineActionEmailContextContextMethods>({
    showEmailEditor: noop,
    closeEmailEditor: noop,
    onCreateEmail: noop,
    handleExitEditorAndCleanData: noop,
    closeConfirmationDialog: noop,
    isEmailEditorOpen: false,
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

  const [isEmailEditorOpen, setShowEmailEditor] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [isSending, setIsSending] = useState(false);
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const formId = 'compose-email-timeline-footer';

  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: [],
    cc: [],
    bcc: [],
    subject: '',
    content: '',
  });
  const { state, reset } = useForm<ComposeEmailDtoI>({
    formId,
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);
  const handleEmailSendSuccess = () => {
    invalidateQuery();
    setIsSending(false);
    reset();
    setShowEmailEditor(false);
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
      session?.user?.email,
    );
  };

  const handleExitEditorAndCleanData = () => {
    reset();
    onClose();
    setShowEmailEditor(false);
  };

  const handleCloseEmailEditor = (prop?: { openLogEntry: () => void }) => {
    const isFormPristine = Object.values(state.fields)?.every(
      (e) => e.meta.pristine,
    );
    const isFormEmpty = Object.values(state.values)?.every((e) => !e.length);

    const showEmailEditorConfirmationDialog = !isFormPristine && !isFormEmpty;
    if (showEmailEditorConfirmationDialog) {
      onOpen();
    } else {
      handleExitEditorAndCleanData();
      prop?.openLogEntry();
    }
  };

  return (
    <TimelineActionEmailContextContext.Provider
      value={{
        showEmailEditor: () => setShowEmailEditor(true),
        closeEmailEditor: handleCloseEmailEditor,
        handleExitEditorAndCleanData,
        closeConfirmationDialog: onClose,
        onCreateEmail,
        isEmailEditorOpen,
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
