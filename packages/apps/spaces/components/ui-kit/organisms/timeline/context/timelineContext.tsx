import type { PropsWithChildren, RefObject } from 'react';
import { createContext, useCallback, useRef } from 'react';

export const noop = () => undefined;

interface TimelineContextMethods {
  timelineContainerRef: RefObject<HTMLDivElement> | null;
  onScrollToBottom: () => void;
}

export const TimelineContext = createContext<TimelineContextMethods>({
  timelineContainerRef: null,
  onScrollToBottom: noop,
});

export const TimelineContextProvider = ({ children }: PropsWithChildren) => {
  const timelineContainerRef = useRef<HTMLDivElement>(null);
  const handleScrollToBottom = useCallback(() => {
    timelineContainerRef?.current?.scrollTo({
      top: timelineContainerRef?.current?.scrollHeight,
    });
  }, [timelineContainerRef]);

  return (
    <TimelineContext.Provider
      value={{
        timelineContainerRef,
        onScrollToBottom: handleScrollToBottom,
      }}
    >
      {children}
    </TimelineContext.Provider>
  );
};
