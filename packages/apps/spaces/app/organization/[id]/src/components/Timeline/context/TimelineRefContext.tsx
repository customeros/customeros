import {
  useContext,
  createContext,
  PropsWithChildren,
  useRef,
  RefObject,
} from 'react';
import { VirtuosoHandle } from 'react-virtuoso';

interface TimelineRefContextMethods {
  virtuosoRef: RefObject<VirtuosoHandle> | null;
}

const TimelineRefContext = createContext<TimelineRefContextMethods>({
  virtuosoRef: null,
});

export const useTimelineRefContext = () => {
  return useContext(TimelineRefContext);
};

export const TimelineRefContextProvider = ({ children }: PropsWithChildren) => {
  const virtuosoRef = useRef<VirtuosoHandle>(null);

  return (
    <TimelineRefContext.Provider
      value={{
        virtuosoRef,
      }}
    >
      {children}
    </TimelineRefContext.Provider>
  );
};
