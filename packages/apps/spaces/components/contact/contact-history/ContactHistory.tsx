import React, { useState } from 'react';
import { Timeline } from '../../ui-kit/organisms';
import { useContactTimeline } from '../../../hooks/useContactTimeline';
import { uuid4 } from '@sentry/utils';

export const ContactHistory = ({ id }: { id: string }) => {
  const { data, contactName, loading, error, fetchMore } = useContactTimeline({
    contactId: id,
  });
  const [prevDate, setPrevDate] = useState(null);
  const liveConversations = {
    __typename: 'LiveConversation',
    source: 'LiveStream',
    createdAt: Date.now(),
    id: uuid4(),
  };

  if (error) {
    return (
      <div>
        <h1>Oops! Timeline error</h1>
      </div>
    );
  }

  return (
    <Timeline
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
      noActivity={!data}
      id={id}
      contactName={contactName}
      loggedActivities={[liveConversations, ...(data || [])]}
    />
  );
};

export default ContactHistory;
