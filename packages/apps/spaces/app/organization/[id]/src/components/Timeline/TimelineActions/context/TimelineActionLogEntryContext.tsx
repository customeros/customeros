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
import { LogEntryWithAliases } from '@organization/src/components/Timeline/types';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  LogEntryFormDto,
  LogEntryFormDtoI,
} from '@organization/src/components/Timeline/TimelineActions/logger/LogEntryFormDto';
import {
  CreateLogEntryMutation,
  CreateLogEntryMutationVariables,
  useCreateLogEntryMutation,
} from '@organization/src/graphql/createLogEntry.generated';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { DataSource, LogEntry } from '@graphql/types';

import { logEntryEditorExtensions } from './extensions';
import { useTimelineMeta } from '../../shared/state';
import { useUpdateCacheWithNewEvent } from '@organization/src/components/Timeline/hooks/updateCacheWithNewEvent';
import { useTimelineRefContext } from '@organization/src/components/Timeline/context/TimelineRefContext';

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
}: PropsWithChildren<{
  invalidateQuery: () => void;
  id: string;
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
  const { virtuosoRef } = useTimelineRefContext();
  const updateTimelineCache = useUpdateCacheWithNewEvent<LogEntry>(virtuosoRef);

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
      const newLogEntry = makeEmptyLogEntryWithAliases(
        session.data?.user?.name,
        payload.logEntry as any,
      );

      const prevTimelineEvents = await updateTimelineCache(
        newLogEntry,
        queryKey,
      );
      handleResetEditor();
      return { prevTimelineEvents };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData(queryKey, context?.prevTimelineEvents);
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
