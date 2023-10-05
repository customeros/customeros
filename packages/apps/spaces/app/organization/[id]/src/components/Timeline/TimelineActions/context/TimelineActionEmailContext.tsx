import React, {
  useState,
  useContext,
  createContext,
  PropsWithChildren,
  useEffect,
} from 'react';

import { useForm } from 'react-inverted-form';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useDisclosure } from '@ui/utils';
import { handleSendEmail } from '@organization/src/components/Timeline/events/email/compose-email/utils';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { useSearchParams } from 'next/navigation';
import { useSession } from 'next-auth/react';
import { useTimelineActionContext } from './TimelineActionContext';

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

  const handleEmailSendSuccess = () => {
    invalidateQuery();
    setIsSending(false);
    handleResetEditor();
    onClose();
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
      session?.user?.email,
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
