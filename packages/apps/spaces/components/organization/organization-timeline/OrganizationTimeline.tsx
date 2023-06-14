import React, { useState } from 'react';
import { Timeline, TimelineStatus } from '@spaces/organisms/timeline';
import { useOrganizationTimeline } from '@spaces/hooks/useOrganizationTimeline';
import { TimelineItemByType } from '@spaces/organisms/timeline/TimelineItemByType';

export const OrganizationTimeline = ({ id }: { id: string }) => {
  const { data, loading, error, fetchMore } = useOrganizationTimeline({
    organizationId: id,
  });

  const [prevDate, setPrevDate] = useState(null);

  if (error) {
    return <TimelineStatus status='timeline-error' />;
  }

  return (
    <Timeline
      mode='ORGANIZATION'
      loading={loading}
      onLoadMore={(containerRef) => {
        const newFromDate = data[0]?.createdAt || data[0]?.startedAt;
        if (!data[0] || prevDate === newFromDate) {
          return;
        }
        // todo remove me when switching to virtualized list
        containerRef.current.scrollTop = 400;
        setPrevDate(newFromDate);
        fetchMore({
          variables: {
            contactId: id,
            size: 10,
            from: newFromDate,
          },
        });
      }}
      noActivity={!data.length && !loading}
      loggedActivities={[...(data || [])]}
      id={id}
    />

  );
};
