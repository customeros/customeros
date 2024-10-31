import { useForm } from 'react-inverted-form';
import { useSearchParams } from 'react-router-dom';
import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { observer } from 'mobx-react-lite';
import { render } from '@react-email/render';

import { useStore } from '@shared/hooks/useStore';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { EmailTemplate } from '@shared/components/EmailTemplate/EmailTemplate.tsx';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmail.dto';

import { useTimelineActionContext } from './TimelineActionContext';

export const noop = () => undefined;

// TODO: type those any props accordingly
interface TimelineActionEmailContextContextMethods {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  state: any;
  formId: string;
  isSending: boolean;
  showConfirmationDialog: boolean;
  checkCanExitSafely: () => boolean;
  closeConfirmationDialog: () => void;
  handleExitEditorAndCleanData: () => void;
  onCreateEmail: (handleSuccess?: () => void) => void;
}

const TimelineActionEmailContextContext =
  createContext<TimelineActionEmailContextContextMethods>({
    checkCanExitSafely: () => false,
    onCreateEmail: noop,
    handleExitEditorAndCleanData: noop,
    closeConfirmationDialog: noop,
    isSending: false,
    showConfirmationDialog: false,
    formId: '',
    state: null,
  });

export const useTimelineActionEmailContext = () => {
  return useContext(TimelineActionEmailContextContext);
};

export const TimelineActionEmailContextContextProvider = observer(
  ({
    children,
    invalidateQuery,
  }: PropsWithChildren<{
    id: string;
    invalidateQuery: () => void;
  }>) => {
    const { open: isOpen, onOpen, onClose } = useDisclosure();
    const [searchParams] = useSearchParams();
    const store = useStore();

    const [isSending, setIsSending] = useState(false);
    const { closeEditor } = useTimelineActionContext();
    const [timelineMeta] = useTimelineMeta();

    const queryKey = useInfiniteGetTimelineQuery.getKey(
      timelineMeta.getTimelineVariables,
    );
    const { virtuosoRef } = useTimelineRefContext();
    const updateTimelineCache = useUpdateCacheWithNewEvent(virtuosoRef);
    const formId = 'compose-email-timeline-footer';

    const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
      from: '',
      fromProvider: '',
      to: [],
      cc: [],
      bcc: [],
      subject: '',
      content: '',
    });
    const { state, reset, setDefaultValues } = useForm<ComposeEmailDtoI>({
      formId,
      defaultValues,

      stateReducer: (_, __, next) => {
        return next;
      },
    });

    const handleResetEditor = () => {
      setDefaultValues(defaultValues);
      reset();
    };

    const handleEmailSendSuccess = (response: unknown) => {
      updateTimelineCache(response, queryKey);

      // no timeout needed is this case as the event id is created when this is called
      invalidateQuery();
      setIsSending(false);
      handleResetEditor();
    };

    const handleEmailSendError = () => {
      setIsSending(false);
    };

    const prepareEmailContent = async (bodyHtml: string) => {
      try {
        const emailHtml = await render(<EmailTemplate bodyHtml={bodyHtml} />, {
          pretty: true,
        });

        return {
          html: emailHtml,
        };
      } catch (error) {
        store.ui.toastError(
          'Unable to process email content',
          'email-content-parsing-error',
        );
      }
    };

    const onCreateEmail = async (handleSuccess = () => {}) => {
      const from = state.values.from?.value ?? '';
      const fromProvider = state.values.from?.provider ?? '';
      const to = [...state.values.to].map(({ value }) => value);
      const cc = [...state.values.cc].map(({ value }) => value);
      const bcc = [...state.values.bcc].map(({ value }) => value);
      const params = new URLSearchParams(searchParams?.toString() ?? '');

      setIsSending(true);

      const id = params.get('events') ?? undefined;

      const handleSendSuccess = (response: unknown) => {
        handleEmailSendSuccess(response);
        handleSuccess?.();
      };

      try {
        const emailContent = await prepareEmailContent(state.values.content);

        if (emailContent?.html) {
          store.mail.send(
            {
              from,
              fromProvider,
              to,
              cc,
              bcc,
              replyTo: id,
              content: emailContent.html,
              subject: state.values.subject,
            },
            {
              onSuccess: (r) => handleSendSuccess(r),
              onError: handleEmailSendError,
            },
          );
        }
      } catch (error) {
        console.error('Error saving email:', error);
        store.ui.toastError('Error saving email', 'email-save-error');
      }
    };

    const handleExitEditorAndCleanData = () => {
      handleResetEditor();

      onClose();
      closeEditor();
    };

    const handleCheckCanExitSafely = () => {
      const { content, ...values } = state.values;

      const isFormEmpty = !content.length || content === `<p style=""></p>`;
      const areFieldsEmpty =
        !values.from ||
        !values.fromProvider ||
        !values.to ||
        values.to.length === 0;
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
          isSending,
          showConfirmationDialog: isOpen,
          formId,
          state,
        }}
      >
        {children}
      </TimelineActionEmailContextContext.Provider>
    );
  },
);

export default TimelineActionEmailContextContextProvider;
