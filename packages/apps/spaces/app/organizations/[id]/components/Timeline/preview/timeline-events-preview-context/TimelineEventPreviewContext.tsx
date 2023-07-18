import type { PropsWithChildren, RefObject } from 'react';
import { createContext, useCallback, useRef } from 'react';

export const noop = () => undefined;

interface TimelineEventPreviewContextContextMethods {
  TimelineEventPreviewContextContainerRef: RefObject<HTMLDivElement> | null;
  onScrollToBottom: () => void;
}

export const TimelineEventPreviewContextContext =
  createContext<TimelineEventPreviewContextContextMethods>({
    TimelineEventPreviewContextContainerRef: null,
    onScrollToBottom: noop,
  });

export const TimelineEventPreviewContextContextProvider = ({
  children,
}: PropsWithChildren) => {
  const TimelineEventPreviewContextContainerRef = useRef<HTMLDivElement>(null);
  const handleScrollToBottom = useCallback(() => {
    TimelineEventPreviewContextContainerRef?.current?.scrollTo({
      top: TimelineEventPreviewContextContainerRef?.current?.scrollHeight,
    });
  }, [TimelineEventPreviewContextContainerRef]);

  return (
    <TimelineEventPreviewContextContext.Provider
      value={{
        TimelineEventPreviewContextContainerRef,
        onScrollToBottom: handleScrollToBottom,
      }}
    >
      {children}
    </TimelineEventPreviewContextContext.Provider>
  );
};
