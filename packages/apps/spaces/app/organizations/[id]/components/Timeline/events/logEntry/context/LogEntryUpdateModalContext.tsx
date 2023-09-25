import {
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
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

  const event = modalContent as LogEntryWithAliases;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const formId = 'log-entry-update';
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const updateLogEntryMutation = useUpdateLogEntryMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => queryClient.invalidateQueries(['GetTimeline.infinite']),
        500,
      );
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const logEntryStartedAtValues = new LogEntryUpdateFormDto(event);

  const { state: formState } = useForm<LogEntryUpdateFormDtoI>({
    formId,
    defaultValues: logEntryStartedAtValues,

    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        updateLogEntryMutation.mutate({
          id: event.id,
          input: {
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

  useEffect(() => {
    if (isModalOpen && event.__typename === 'LogEntry') {
      updateLogEntryMutation.mutate({
        id: event.id,
        input: {
          ...LogEntryUpdateFormDto.toPayload({
            ...formState.values,
          }),
        },
      });
    }
  }, [isModalOpen]);

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
