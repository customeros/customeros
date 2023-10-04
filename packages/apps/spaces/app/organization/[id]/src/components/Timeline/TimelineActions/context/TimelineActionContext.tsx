import React, {
  useState,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

export const noop = () => undefined;
export type EditorType = 'email' | 'log-entry' | null;
interface TimelineActionContextContextMethods {
  showEditor: (editorType: EditorType) => void;
  closeEditor: () => void;
  openedEditor: EditorType;
}

const TimelineActionContextContext =
  createContext<TimelineActionContextContextMethods>({
    showEditor: noop,
    closeEditor: noop,
    openedEditor: null,
  });

export const useTimelineActionContext = () => {
  return useContext(TimelineActionContextContext);
};

export const TimelineActionContextContextProvider = ({
  children,
}: PropsWithChildren) => {
  const [openedEditor, setOpenedEditor] = useState<EditorType>(null);

  return (
    <TimelineActionContextContext.Provider
      value={{
        showEditor: setOpenedEditor,
        closeEditor: () => setOpenedEditor(null),
        openedEditor,
      }}
    >
      {children}
    </TimelineActionContextContext.Provider>
  );
};
