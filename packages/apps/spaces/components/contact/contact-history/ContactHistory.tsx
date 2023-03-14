import React, { useState } from 'react';
import { Timeline } from '../../ui-kit/organisms';
import { useContactTimeline } from '../../../hooks/useContactTimeline';

export const ContactHistory = ({ id }: { id: string }) => {
  const { data, loading, error, fetchMore } = useContactTimeline({
    contactId: id,
  });
  const [prevDate, setPrevDate] = useState(null);
  const liveConversations = {
    __typename: 'LiveConversation',
    source: 'LiveStream',
    createdAt: Date.now(),
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
      onLoadMore={() => {
        const newFromDate = data[0]?.createdAt || data[0]?.startedAt;
        if (!data[0] || prevDate === newFromDate) {
          return;
        }
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
      contactId={id}
      loggedActivities={[liveConversations, ...(data || [])]}
    />
  );
};

export default ContactHistory;
