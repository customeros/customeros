import React, { useState } from 'react';
import { Timeline, TimelineStatus } from '@spaces/organisms/timeline';
import { useContactTimeline } from '@spaces/hooks/useContactTimeline';
import { uuid4 } from '@sentry/utils';

export const ContactHistory = ({ id }: { id: string }) => {
  const { data, error, fetchMore } = useContactTimeline({
    contactId: id,
  });
  const [prevDate, setPrevDate] = useState(null);
  const liveInteractions = {
    __typename: 'LiveEventTimelineItem',
    source: 'LiveStream',
    createdAt: Date.now(),
    id: uuid4(),
  };
  if (error) {
    return <TimelineStatus status='timeline-error' />;
  }
  return (
    <Timeline
      mode='CONTACT'
      loading={false}
      onLoadMore={(containerRef) => {
        const newFromDate = data[0]?.createdAt || data[0]?.startedAt;
        if (!data[0] || prevDate === newFromDate) {
          return;
        }
        // todo remove me when switching to virtualized list
        containerRef.current.scrollTop = 100;
        setPrevDate(newFromDate);
        fetchMore({
          variables: {
            contactId: id,
            size: 10,
            from: newFromDate,
          },
        });
      }}
      noActivity={!data.length}
      id={id}
      loggedActivities={[...(data || []), liveInteractions]}
    />
  );
};

export default ContactHistory;
