import { cn } from '@ui/utils/cn';
import { useTimelineActionContext } from '@organization/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

import { EmailTimelineAction } from './email/EmailTimelineAction';
import { LogEntryTimelineAction } from './logger/LogEntryTimelineAction';

interface TimelineActionsAreaProps {
  hide: () => void;
  activeEditor: 'log-entry' | null;
}

export const TimelineActionsArea = ({
  hide,
  activeEditor,
}: TimelineActionsAreaProps) => {
  const { openedEditor } = useTimelineActionContext();

  return (
    <div
      className={cn(
        openedEditor !== null || activeEditor !== null
          ? 'pt-6 pb-2'
          : 'pt-0 pb-8',
        'mt-[-16px] bg-[#F9F9FB] border-dashed border-t-[1px] border-gray-200',
      )}
    >
      {openedEditor === 'email' && <EmailTimelineAction />}
      {activeEditor === 'log-entry' && <LogEntryTimelineAction hide={hide} />}
    </div>
  );
};
