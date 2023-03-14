import React, { useEffect } from 'react';
import { OrganizationTimelineSkeleton } from './skeletons';
import { Timeline } from '../../ui-kit/organisms';
import { useOrganizationTimelineData } from '../../../hooks/useOrganization/useOrganizationTimeline';

export const OrganizationTimeline = ({ id }: { id: string }) => {
  const {
    data,
    loading: orgLoading,
    error,
  } = useOrganizationTimelineData({
    id,
  });

  const [loading, setLoading] = React.useState<boolean>(true);
  const [notes, setNotes] = React.useState<any>([]);
  const [tickets, setTickets] = React.useState<any>([]);
  const [conversations, setConversations] = React.useState<any>([]);

  useEffect(() => {
    if (!orgLoading && data) {
      let ticketsData = [] as any;
      let notesData = [...data.notes.content] as any;
      let conversationsData = [] as any;

      data.contacts.content.forEach((contact: any) => {
        if (contact.notes && contact.notes.content) {
          notesData = [...notesData, ...contact.notes.content];
        }
        if (contact.tickets) {
          ticketsData = [...ticketsData, ...contact.tickets];
        }
        if (contact.conversations && contact.conversations.content) {
          conversationsData = [
            ...conversationsData,
            ...contact.conversations.content,
          ];
        }
      });

      setTickets(ticketsData);
      setNotes(notesData);
      setConversations(conversationsData);
      setLoading(false);
    }
  }, [orgLoading, data?.notes.content.length]); // fixme after adding new timeline

  const noHistoryItemsAvailable =
    !loading &&
    notes.length == 0 &&
    conversations.length == 0 &&
    tickets.length == 0;

  const getSortedItems = (
    data1: Array<any> | undefined,
    data2: Array<any> | undefined,
    data3: Array<any> | undefined,
    data4: Array<any> | undefined,
  ) => {
    const data = [
      ...(data1 || []),
      ...(data2 || []),
      ...(data3 || []),
      ...(data4 || []),
    ];
    return data.sort((a, b) => {
      return Date.parse(a?.createdAt) - Date.parse(b?.createdAt);
    });
  };

  if (loading) {
    return <OrganizationTimelineSkeleton />;
  }
  if (error) {
    return null;
  }

  return (
    <Timeline
      loading={loading}
      noActivity={noHistoryItemsAvailable}
      contactId={id}
      loggedActivities={getSortedItems(notes, conversations, tickets, [])}
    />
  );
};
