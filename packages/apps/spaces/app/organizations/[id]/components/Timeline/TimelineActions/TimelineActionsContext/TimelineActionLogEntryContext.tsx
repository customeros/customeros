import {
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
} from 'react';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  LogEntryFormDto,
  LogEntryFormDtoI,
} from '@organization/components/Timeline/TimelineActions/logger/LogEntryFormDto';
import { useForm } from 'react-inverted-form';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useCreateLogEntryMutation } from '@organization/graphql/createLogEntry.generated';
import { useDisclosure } from '@ui/utils';
import { useTimelineActionContext } from './TimelineActionContext';

export const noop = () => undefined;

interface TimelineActionLogEntryContextContextMethods {
  checkCanExitSafely: () => boolean;
  closeConfirmationDialog: () => void;
  handleExitEditorAndCleanData: () => void;
  onCreateLogEntry: () => void;
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

  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { closeEditor } = useTimelineActionContext();

  const logEntryValues = new LogEntryFormDto();
  const { state, reset, setDefaultValues } = useForm<LogEntryFormDtoI>({
    formId: 'organization-create-log-entry',
    defaultValues: logEntryValues,

    stateReducer: (_, _a, next) => {
      return next;
    },
  });
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
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
    onSuccess: () => {
      timeoutRef.current = setTimeout(() => invalidateQuery(), 500);
      handleResetEditor();
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const onCreateLogEntry = () => {
    const logEntryPayload = LogEntryFormDto.toPayload({
      ...logEntryValues,
      tags: state.values.tags,
      content: state.values.content,
      contentType: state.values.contentType,
    });
    createLogEntryMutation.mutate({
      organizationId: id,

      logEntry: logEntryPayload,
    });
  };

  const handleExitEditorAndCleanData = () => {
    handleResetEditor();
    onClose();
    closeEditor();
  };

  const handleCheckCanExitSafely = () => {
    const { content, tags } = state.values;

    const isContentEmpty = !content.length || content === `<p style=""></p>`;

    const showLogEntryEditorConfirmationDialog =
      !!tags.length || !isContentEmpty;
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
