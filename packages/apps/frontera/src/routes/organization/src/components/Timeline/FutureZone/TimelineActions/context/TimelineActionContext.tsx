import { useState, useContext, createContext, PropsWithChildren } from 'react';

import { observer } from 'mobx-react-lite';

import { useEvent } from '@shared/hooks/useEvent';

export const noop = () => undefined;
export type EditorType = 'email' | 'log-entry' | 'reminder' | null;
interface TimelineActionContextContextMethods {
  closeEditor: () => void;
  openedEditor: EditorType;
  showEditor: (editorType: EditorType) => void;
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

export const TimelineActionContextContextProvider = observer(
  ({ children }: PropsWithChildren) => {
    useEvent<{ openEditor: EditorType }>('openEmailEditor', (payload) => {
      setOpenedEditor(payload.openEditor);
    });

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
  },
);
