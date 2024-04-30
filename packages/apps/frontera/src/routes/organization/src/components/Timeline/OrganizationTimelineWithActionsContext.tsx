'use client';
import { TimelineActionContextContextProvider } from '@organization/components/Timeline/FutureZone/TimelineActions';

import { OrganizationTimeline } from './OrganizationTimeline';

export const OrganizationTimelineWithActionsContext = () => {
  return (
    <>
      <TimelineActionContextContextProvider>
        <OrganizationTimeline />
      </TimelineActionContextContextProvider>
    </>
  );
};
