import { useForm } from 'react-inverted-form';
import {
  useRef,
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useStore } from '@shared/hooks/useStore';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useUpdateLogEntryMutation } from '@organization/graphql/updateLogEntry.generated';
import { useUpdateCacheWithExistingEvent } from '@organization/components/Timeline/PastZone/hooks/useCacheExistingEvent';
import { useTimelineEventPreviewStateContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
  const store = useStore();

  const isAuthor =
    event?.logEntryCreatedBy?.emails?.findIndex(
      (e) => store.session.value?.profile.email === e.email,
    ) !== -1;

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const { setDefaultValues, reset } = useForm<LogEntryUpdateFormDtoI>({
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
    onSuccess: (_data, variables, _context) => {
      const mappedData = {
        logEntryStartedAt: variables?.input?.startedAt,
        // content: variables?.input?.content,
        // contentType: variables?.input?.contentType,
      };

      updateTimelineCache(
        { id: variables.id, ...mappedData } as LogEntryWithAliases,
        queryKey,
      );
    },
  });

  useEffect(() => {
    if (!isModalOpen && openedLogEntryId && isAuthor) {
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
