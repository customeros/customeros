import { TimelineRefContextProvider } from '@organization/components/Timeline/context/TimelineRefContext';
import { TimelineEventPreviewContextContextProvider } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

interface TimelineContextsProviderProps {
  id: string;
  children: React.ReactNode;
}

export const TimelineContextsProvider = ({
  children,
  id,
}: TimelineContextsProviderProps) => {
  return (
    <TimelineRefContextProvider>
      <TimelineEventPreviewContextContextProvider id={id}>
        {children}
      </TimelineEventPreviewContextContextProvider>
    </TimelineRefContextProvider>
  );
};
