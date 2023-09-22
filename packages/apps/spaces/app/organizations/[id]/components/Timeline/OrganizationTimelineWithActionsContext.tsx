'use client';
import { TimelineActionContextContextProvider } from './TimelineActions/TimelineActionsContext/TimelineActionContext';
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
