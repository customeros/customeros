import React, { useState } from 'react';
import { Timeline } from '../../ui-kit/organisms';
import { useOrganizationTimeline } from '../../../hooks/useOrganizationTimeline';

export const OrganizationTimeline = ({ id }: { id: string }) => {
  const { data, loading, error, fetchMore } = useOrganizationTimeline({
    organizationId: id,
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
      loggedActivities={[liveConversations, ...(data || [])]}
    />
  );
};
