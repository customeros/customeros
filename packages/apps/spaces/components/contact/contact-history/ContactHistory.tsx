import React, { useEffect, useState } from 'react';
import { Timeline, TimelineStatus } from '@spaces/organisms/timeline';
import { useContactTimeline } from '@spaces/hooks/useContactTimeline';

export const ContactHistory = ({ id }: { id: string }) => {
  const { data, error, loading, fetchMore } = useContactTimeline({
    contactId: id,
  });
  const [prevDate, setPrevDate] = useState(null);

  if (error) {
    return <TimelineStatus status='timeline-error' />;
  }

  return (
    <Timeline
      mode='CONTACT'
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
      noActivity={!data?.length && !loading}
      id={id}
      loggedActivities={[...(data || [])]}
    />
  );
};

export default ContactHistory;
