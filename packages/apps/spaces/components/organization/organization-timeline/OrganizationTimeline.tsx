import React, { useEffect } from 'react';
import { useRouter } from 'next/router';
import classNames from 'classnames';
import { Phone, Envelope } from '../../ui-kit';
import { OrganizationTimelineSkeleton } from './skeletons';
import { useOrganizationContacts } from '../../../hooks/useOrganization';
import styles from './organization-timeline.module.scss';
import { ContactTags } from '../../contact/contact-tags';
import { getContactDisplayName } from '../../../utils';
import { useContactNotes } from '../../../hooks/useContactNote/useContactNotes';
import { useContactTickets } from '../../../hooks/useContact/useContactTickets';
import { useContactConversations } from '../../../hooks/useContactConversations';
import {
  Contact,
  useGetContactNotesQuery,
} from '../../../graphQL/__generated__/generated';
import { Timeline } from '../../ui-kit/organisms';
import { useOrganizationTimelineData } from '../../../hooks/useOrganization/useOrganizationTimeline';
import { OrganizationContactsSkeleton } from '../organization-contacts/skeletons';

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

  console.log(data);

  useEffect(() => {
    if (!orgLoading && data) {
      let ticketsData = [] as any;
      let notesData = [] as any;
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
  }, [orgLoading]);

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
