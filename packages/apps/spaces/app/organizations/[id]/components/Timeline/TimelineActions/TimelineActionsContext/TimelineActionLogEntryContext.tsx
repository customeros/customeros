import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
} from 'react';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  LogEntryDto,
  LogEntryDtoI,
} from '@organization/components/Timeline/TimelineActions/logger/LogEntry.dto';
import { useForm } from 'react-inverted-form';
import { useRemirror } from '@remirror/react';
import { basicEditorExtensions } from '@ui/form/RichTextEditor/extensions';
import { useCreateLogEntryMutation } from '@organization/graphql/createLogEntry.generated';
import { useDisclosure } from '@ui/utils';

export const noop = () => undefined;

interface TimelineActionLogEntryContextContextMethods {
  showLogEntryEditor: () => void;
  closeLogEntryEditor: (prop?: { openEmailEditor: () => void }) => void;
  closeConfirmationDialog: () => void;
  handleExitEditorAndCleanData: () => void;
  onCreateLogEntry: any;
  isLogEntryEditorOpen: boolean;
  remirrorProps: any;
  isSaving: boolean;
  showConfirmationDialog: boolean;
}

const TimelineActionLogEntryContextContext =
  createContext<TimelineActionLogEntryContextContextMethods>({
    showLogEntryEditor: noop,
    closeLogEntryEditor: noop,
    onCreateLogEntry: noop,
    closeConfirmationDialog: noop,
    handleExitEditorAndCleanData: noop,
    isLogEntryEditorOpen: false,
    remirrorProps: null,
    isSaving: false,
    showConfirmationDialog: false,
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

  const [isLogEntryEditorOpen, setEditorOpen] = useState(false);
  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const logEntryValues: LogEntryDtoI = new LogEntryDto();
  const { state, reset } = useForm<LogEntryDtoI>({
    formId: 'organization-create-log-entry',
    defaultValues: logEntryValues,

    stateReducer: (_, _a, next) => {
      return next;
    },
  });
  const remirrorProps = useRemirror({
    extensions: basicEditorExtensions,
  });
  const createLogEntryMutation = useCreateLogEntryMutation(client, {
    onSuccess: () => {
      reset();
      invalidateQuery();
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
    const logEntryPayload = LogEntryDto.toPayload({
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

  const handleExitEditorAndCleanData = (prop?: any) => {
    reset();
    onClose();
    setEditorOpen(false);
    console.log('🏷️ ----- : HERE');
    setTimeout(() => {
      console.log('🏷️ ----- : insite', prop);
      prop?.openEmailEditor();
    }, 100);
  };

  const handleCloseEmailEditor = (prop?: { openEmailEditor: () => void }) => {
    const isFormPristine = Object.values(state.fields)?.every(
      (e) => e.meta.pristine,
    );
    const isFormEmpty = Object.values(state.values)?.every((e) => !e.length);

    const showEmailEditorConfirmationDialog = !isFormPristine && !isFormEmpty;
    if (showEmailEditorConfirmationDialog) {
      onOpen();
    } else {
      handleExitEditorAndCleanData(prop);
    }
  };

  return (
    <TimelineActionLogEntryContextContext.Provider
      value={{
        showLogEntryEditor: () => setEditorOpen(true),
        closeLogEntryEditor: handleCloseEmailEditor,
        handleExitEditorAndCleanData,
        closeConfirmationDialog: onClose,
        onCreateLogEntry,
        isLogEntryEditorOpen,
        remirrorProps,
        isSaving: createLogEntryMutation.isLoading,
        showConfirmationDialog: isOpen,
      }}
    >
      {children}
    </TimelineActionLogEntryContextContext.Provider>
  );
};
