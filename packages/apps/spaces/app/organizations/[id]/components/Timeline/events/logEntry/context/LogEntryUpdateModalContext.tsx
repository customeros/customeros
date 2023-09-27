import {
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
  useState,
} from 'react';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useForm } from 'react-inverted-form';
import { useQueryClient } from '@tanstack/react-query';
import { useUpdateLogEntryMutation } from '@organization/graphql/updateLogEntry.generated';
import {
  LogEntryUpdateFormDto,
  LogEntryUpdateFormDtoI,
} from './LogEntryUpdateFormDto';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { useTimelineEventPreviewContext } from '@organization/components/Timeline/preview/context/TimelineEventPreviewContext';
import { useSession } from 'next-auth/react';

export const noop = () => undefined;

interface LogEntryUpdateModalContextMethods {
  formId: string;
}

const LogEntryUpdateModalContext =
  createContext<LogEntryUpdateModalContextMethods>({
    formId: '',
  });

export const useLogEntryUpdateContext = () => {
  return useContext(LogEntryUpdateModalContext);
};

export const LogEntryUpdateModalContextProvider = ({
  children,
}: PropsWithChildren) => {
  const { modalContent, isModalOpen } = useTimelineEventPreviewContext();
  const [openedLogEntryId, setOpenedLogEntryId] = useState<null | string>(null);
  const event = modalContent as LogEntryWithAliases;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const formId = 'log-entry-update';
  const logEntryStartedAtValues = new LogEntryUpdateFormDto(event);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { data: session } = useSession();
  const isAuthor =
    event?.logEntryCreatedBy?.emails?.findIndex(
      (e) => session?.user?.email === e.email,
    ) !== -1;
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const {
    state: formState,
    setDefaultValues,
    reset,
  } = useForm<LogEntryUpdateFormDtoI>({
    formId,
    defaultValues: logEntryStartedAtValues,

    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_BLUR' && isAuthor) {
        updateLogEntryMutation.mutate({
          id: event.id,
          input: {
            contentType: 'text/html',
            ...LogEntryUpdateFormDto.toPayload({
              ...state.values,
              [action.payload.name]: action.payload.value,
            }),
          },
        });
      }
      return next;
    },
  });

  const updateLogEntryMutation = useUpdateLogEntryMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => queryClient.invalidateQueries(['GetTimeline.infinite']),
        500,
      );
    },
  });

  useEffect(() => {
    if (!isModalOpen && openedLogEntryId && isAuthor) {
      updateLogEntryMutation.mutate({
        id: openedLogEntryId,
        input: {
          ...LogEntryUpdateFormDto.toPayload({
            ...formState.values,
          }),
        },
      });
      setOpenedLogEntryId(null);
      reset();
    }
  }, [isModalOpen, openedLogEntryId]);

  useEffect(() => {
    if (event?.id && event.__typename === 'LogEntry') {
      setOpenedLogEntryId(event?.id);
      const newDefaults = new LogEntryUpdateFormDto(event);
      // console.log('newDefaults', newDefaults)
      setDefaultValues(newDefaults);
    }
  }, [event?.id]);

  return (
    <LogEntryUpdateModalContext.Provider
      value={{
        formId,
      }}
    >
      {children}
    </LogEntryUpdateModalContext.Provider>
  );
};
