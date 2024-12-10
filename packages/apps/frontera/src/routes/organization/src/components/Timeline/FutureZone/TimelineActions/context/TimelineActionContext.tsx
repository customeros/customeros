import {
  useState,
  useEffect,
  useContext,
  createContext,
  PropsWithChildren,
} from 'react';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

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
    const store = useStore();
    const [openedEditor, setOpenedEditor] = useState<EditorType>(null);

    useEffect(() => {
      if (store.ui.openEmailEditor) {
        setOpenedEditor('email');
      } else {
        store.ui.setEmailAdress('');
      }
    }, [store.ui.openEmailEditor]);

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
