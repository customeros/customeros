import {
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
} from 'react';
import { useForm } from 'react-inverted-form';
import { VirtuosoHandle } from 'react-virtuoso';
import { useRemirror } from '@remirror/react';
import { useSession } from 'next-auth/react';
import {
  UseMutationOptions,
  useQueryClient,
  InfiniteData,
} from '@tanstack/react-query';

import { useDisclosure } from '@ui/utils';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  LogEntryFormDto,
  LogEntryFormDtoI,
} from '@organization/components/Timeline/TimelineActions/logger/LogEntryFormDto';
import {
  CreateLogEntryMutation,
  CreateLogEntryMutationVariables,
  useCreateLogEntryMutation,
} from '@organization/graphql/createLogEntry.generated';
import {
  GetTimelineQuery,
  useInfiniteGetTimelineQuery,
} from '@organization/graphql/getTimeline.generated';
import { DataSource } from '@graphql/types';

import { logEntryEditorExtensions } from './extensions';
import { useTimelineMeta } from '../../shared/state';

export const noop = () => undefined;

interface TimelineActionLogEntryContextContextMethods {
  checkCanExitSafely: () => boolean;
  closeConfirmationDialog: () => void;
  handleExitEditorAndCleanData: () => void;
  onCreateLogEntry: (
    options?: UseMutationOptions<
      CreateLogEntryMutation,
      unknown,
      CreateLogEntryMutationVariables,
      unknown
    >,
  ) => void;
  remirrorProps: any;
  isSaving: boolean;
  showLogEntryConfirmationDialog: boolean;
}

const TimelineActionLogEntryContextContext =
  createContext<TimelineActionLogEntryContextContextMethods>({
    checkCanExitSafely: () => false,
    onCreateLogEntry: noop,
    closeConfirmationDialog: noop,
    handleExitEditorAndCleanData: noop,
    remirrorProps: null,
    isSaving: false,
    showLogEntryConfirmationDialog: false,
  });

export const useTimelineActionLogEntryContext = () => {
  return useContext(TimelineActionLogEntryContextContext);
};

export const TimelineActionLogEntryContextContextProvider = ({
  children,
  invalidateQuery,
  id = '',
  virtuosoRef,
}: PropsWithChildren<{
  invalidateQuery: () => void;
  id: string;
  virtuosoRef?: React.RefObject<VirtuosoHandle>;
}>) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [timelineMeta] = useTimelineMeta();
  const session = useSession();

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const queryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );

  const logEntryValues = new LogEntryFormDto();
  const { state, reset, setDefaultValues } = useForm<LogEntryFormDtoI>({
    formId: 'organization-create-log-entry',
    defaultValues: logEntryValues,

    stateReducer: (_, _a, next) => {
      return next;
    },
  });
  const remirrorProps = useRemirror({
    extensions: logEntryEditorExtensions,
  });
  const handleResetEditor = () => {
    reset();
    setDefaultValues(logEntryValues);

    const context = remirrorProps.getContext();
    if (context) {
      context.commands.resetContent();
    }
  };

  const createLogEntryMutation = useCreateLogEntryMutation(client, {
    onMutate: async (payload) => {
      await queryClient.cancelQueries({ queryKey });

      const previousLogEntries =
        queryClient.getQueryData<InfiniteData<GetTimelineQuery>>(queryKey);

      const timelineEntries =
        previousLogEntries?.pages?.[0]?.organization?.timelineEvents;

      const newLogEntry = makeEmptyLogEntryWithAliases(
        session.data?.user?.name,
        payload.logEntry as any,
      );

      queryClient.setQueryData<InfiniteData<GetTimelineQuery>>(
        queryKey,
        (currentCache): InfiniteData<GetTimelineQuery> => {
          const nextCache = {
            ...currentCache,
            pages: currentCache?.pages?.map((p, idx) => {
              if (idx !== 0) return p;
              return {
                ...p,
                organization: {
                  ...p?.organization,
                  timelineEvents: [
                    newLogEntry,
                    ...(p?.organization?.timelineEvents ?? []),
                  ],
                  timelineEventsTotalCount:
                    p?.organization?.timelineEventsTotalCount + 1,
                },
              };
            }),
          } as InfiniteData<GetTimelineQuery>;

          return nextCache;
        },
      );

      virtuosoRef?.current?.scrollToIndex({
        index: (timelineEntries?.length ?? 0) + 1,
      });
      handleResetEditor();
      return { previousLogEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData(queryKey, context?.previousLogEntries);
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => invalidateQuery(), 2000);
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const onCreateLogEntry = (
    options?: UseMutationOptions<
      CreateLogEntryMutation,
      unknown,
      CreateLogEntryMutationVariables,
      unknown
    >,
  ) => {
    const logEntryPayload = LogEntryFormDto.toPayload({
      ...logEntryValues,
      tags: state.values.tags,
      content: state.values.content,
      contentType: state.values.contentType,
    });
    createLogEntryMutation.mutate(
      {
        organizationId: id,
        logEntry: logEntryPayload,
      },
      {
        ...(options ?? {}),
      },
    );
  };

  const handleExitEditorAndCleanData = () => {
    handleResetEditor();
    onClose();
  };

  const handleCheckCanExitSafely = () => {
    const { content } = state.values;
    const isContentEmpty = !content.length || content === `<p style=""></p>`;
    const showLogEntryEditorConfirmationDialog = !isContentEmpty;
    if (showLogEntryEditorConfirmationDialog) {
      onOpen();
      return false;
    } else {
      handleResetEditor();
      onClose();
      return true;
    }
  };

  return (
    <TimelineActionLogEntryContextContext.Provider
      value={{
        checkCanExitSafely: handleCheckCanExitSafely,
        handleExitEditorAndCleanData,
        closeConfirmationDialog: onClose,
        onCreateLogEntry,
        remirrorProps,
        isSaving: createLogEntryMutation.isLoading,
        showLogEntryConfirmationDialog: isOpen,
      }}
    >
      {children}
    </TimelineActionLogEntryContextContext.Provider>
  );
};

function makeEmptyLogEntryWithAliases(
  userName: string | null = '',
  data: LogEntryFormDto,
): LogEntryWithAliases {
  const { tags, ...rest } = data;
  return {
    __typename: 'LogEntry',
    id: Math.random().toString(),
    createdAt: '',
    updatedAt: '',
    sourceOfTruth: DataSource.Na,
    externalLinks: [],
    source: DataSource.Na,
    logEntryStartedAt: new Date().toISOString(),
    logEntryCreatedBy: {
      firstName: userName,
      lastName: '',
    } as any,
    tags: tags.map((t) => ({
      name: t.label,
      id: t.value,
      appSource: '',
      __typename: 'Tag',
      createdAt: '',
      source: DataSource.Na,
      updatedAt: '',
    })),
    ...rest,
  };
}
