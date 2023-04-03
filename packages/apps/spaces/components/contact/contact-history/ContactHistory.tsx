import React, { useState } from 'react';
import { Timeline } from '../../ui-kit/organisms';
import { useContactTimeline } from '../../../hooks/useContactTimeline';
import { uuid4 } from '@sentry/utils';
import { TimelineStatus } from '../../ui-kit/organisms/timeline';

export const ContactHistory = ({ id }: { id: string }) => {
  const { data, contactName, loading, error, fetchMore } = useContactTimeline({
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
      loading={loading}
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
      contactName={contactName}
      loggedActivities={[...(data || []), liveInteractions]}
    />
  );
};

export default ContactHistory;
