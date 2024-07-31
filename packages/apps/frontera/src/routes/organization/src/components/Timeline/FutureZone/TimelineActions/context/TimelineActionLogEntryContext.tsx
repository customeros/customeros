import { useForm } from 'react-inverted-form';
import {
  useRef,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { useRemirror } from '@remirror/react';
import { useQueryClient, UseMutationOptions } from '@tanstack/react-query';

import { useStore } from '@shared/hooks/useStore';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { LogEntry, DataSource, LogEntryInput } from '@graphql/types';
import { LogEntryWithAliases } from '@organization/components/Timeline/types';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useTimelineRefContext } from '@organization/components/Timeline/context/TimelineRefContext';
import { useUpdateCacheWithNewEvent } from '@organization/components/Timeline/PastZone/hooks/updateCacheWithNewEvent';
import {
  LogEntryFormDto,
  LogEntryFormDtoI,
} from '@organization/components/Timeline/FutureZone/TimelineActions/logger/LogEntryFormDto';
import {
  CreateLogEntryMutation,
  useCreateLogEntryMutation,
  CreateLogEntryMutationVariables,
} from '@organization/graphql/createLogEntry.generated';

import { useTimelineMeta } from '../../../state';
import { logEntryEditorExtensions } from './extensions';

export const noop = () => undefined;

interface TimelineActionLogEntryContextContextMethods {
  isSaving: boolean;
  // TODO: type this correctly
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  remirrorProps: any;
  checkCanExitSafely: () => boolean;
  closeConfirmationDialog: () => void;
  showLogEntryConfirmationDialog: boolean;
  handleExitEditorAndCleanData: () => void;
  onCreateLogEntry: (
    options?: UseMutationOptions<
      CreateLogEntryMutation,
      unknown,
      CreateLogEntryMutationVariables,
      unknown
    > & { payload?: LogEntryInput },
  ) => void;
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
  id: string;
  invalidateQuery: () => void;
}>) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const [timelineMeta] = useTimelineMeta();
  const store = useStore();

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
        store.session.value.profile.name,
        // TODO: type this correctly
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
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
    > & { payload?: LogEntryInput },
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
        logEntry: options?.payload ?? logEntryPayload,
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
        isSaving: createLogEntryMutation.isPending,
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
      // TODO: this object should contain defaults for all required properties of User
      // this will in exchange solve the type issues
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any,
    tags: tags.map((t) => ({
      name: t.label,
      id: t.value,
      appSource: '',
      __typename: 'Tag',
      createdAt: '',
      source: DataSource.Na,
      updatedAt: '',
      metadata: {
        id: t.value,
        source: DataSource.Openline,
        sourceOfTruth: DataSource.Openline,
        appSource: 'organization',
        created: new Date().toISOString(),
        lastUpdated: new Date().toISOString(),
      },
    })),
    ...rest,
  };
}
