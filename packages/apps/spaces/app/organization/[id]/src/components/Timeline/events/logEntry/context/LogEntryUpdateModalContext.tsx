import { useForm } from 'react-inverted-form';
import {
  useRef,
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useSession } from 'next-auth/react';

import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { LogEntryWithAliases } from '@organization/src/components/Timeline/types';
import { useTimelineMeta } from '@organization/src/components/Timeline/shared/state';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { useUpdateLogEntryMutation } from '@organization/src/graphql/updateLogEntry.generated';
import { useUpdateCacheWithExistingEvent } from '@organization/src/components/Timeline/hooks/useCacheExistingEvent';
import { useTimelineEventPreviewStateContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';

import {
  LogEntryUpdateFormDto,
  LogEntryUpdateFormDtoI,
} from './LogEntryUpdateFormDto';

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
  const { modalContent, isModalOpen } = useTimelineEventPreviewStateContext();
  const [openedLogEntryId, setOpenedLogEntryId] = useState<null | string>(null);
  const event = modalContent as LogEntryWithAliases;
  const client = getGraphQLClient();
  const formId = 'log-entry-update';
  const logEntryStartedAtValues = new LogEntryUpdateFormDto(event);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const updateTimelineCache = useUpdateCacheWithExistingEvent();
  const [timelineMeta] = useTimelineMeta();
  const queryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );
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
    onSuccess: (data, variables, context) => {
      const mappedData = {
        logEntryStartedAt: variables?.input?.startedAt,
        content: variables?.input?.content,
        contentType: variables?.input?.contentType,
      };

      updateTimelineCache({ id: variables.id, ...mappedData }, queryKey);
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
