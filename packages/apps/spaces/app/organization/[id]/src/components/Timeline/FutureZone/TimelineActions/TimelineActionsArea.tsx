import React from 'react';

import { cn } from '@ui/utils/cn';
import { useTimelineActionContext } from '@organization/src/components/Timeline/FutureZone/TimelineActions/context/TimelineActionContext';

import { EmailTimelineAction } from './email/EmailTimelineAction';
import { LogEntryTimelineAction } from './logger/LogEntryTimelineAction';

export const TimelineActionsArea: React.FC = () => {
  const { openedEditor } = useTimelineActionContext();

  return (
    <div
      className={cn(
        openedEditor !== null ? 'pt-6 pb-2' : 'pt-0 pb-8',
        'mt-[-16px] bg-[#F9F9FB] border-dashed border-t-[1px] border-gray-200',
      )}
    >
      {openedEditor === 'email' && <EmailTimelineAction />}
      {openedEditor === 'log-entry' && <LogEntryTimelineAction />}
    </div>
  );
};
