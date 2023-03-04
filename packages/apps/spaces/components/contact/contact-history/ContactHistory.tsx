import React, {useEffect, useState} from 'react';
import { useContactNotes } from '../../../hooks/useContactNote/useContactNotes';
import { Timeline } from '../../ui-kit/organisms';
import { useContactConversations } from '../../../hooks/useContactConversations';
import { useContactTickets } from '../../../hooks/useContact/useContactTickets';
import {gql} from "@apollo/client";
import {GraphQLClient} from "graphql-request";

export const ContactHistory = ({ id }: { id: string }) => {
  const {
    data: notes,
    loading: notesLoading,
    error: notesError,
  } = useContactNotes({ id });
  const {
    data: tickets,
    loading: ticketsLoading,
    error: ticketsError,
  } = useContactTickets({ id });
  // TODO add pagination support
  const {
    data: conversations,
    loading: conversationsLoading,
    error: conversationsError,
  } = useContactConversations({ id });

  const query = gql`query GetActionsForContact($id: ID!, $from: Time!, $to: Time! ) {
    contact(id: $id) {
      id
      firstName
      lastName
      createdAt
      actions(from:$from, to:$to) {
        ... on PageViewAction {
          id
          application
          startedAt
          endedAt
          engagedTime
          pageUrl
          pageTitle
          orderInSession
          sessionId
        }
      }
    }
  }`

  const [actions, setActions] = useState<any[] | undefined>(undefined);
  const [actionsLoading, setActionsLoading] = useState<boolean>(true);

  useEffect(() => {
    const from = new Date(1970, 0, 1).toISOString();
    const to = new Date().toISOString()
    const client = new GraphQLClient(`/customer-os-api/query`);
    client.request(query, {id, from, to}).then((response) => {
      console.log(response)
      if (response && response.contact) {
        response.contact.actions.forEach((action: any) => {
          action.__typename = 'PageViewAction';
        });
        setActions(response.contact.actions);
        setActionsLoading(false);
      } else {
        setActions([]);
        setActionsLoading(false);
      }
    });
  }, [id]);

  const liveConversations = {__typename: "LiveConversation", source: "LiveStream", createdAt: Date.now()}

  console.log('in history');

  const noHistoryItemsAvailable =
      !notesLoading &&
      notes?.notes?.content.length == 0 &&
      !conversationsLoading &&
      conversations?.conversations?.content.length == 0 &&
      !ticketsLoading &&
      tickets?.tickets?.length == 0 &&
      !actionsLoading &&
      actions?.length == 0;

  const getSortedItems = (
      data1: Array<any> | undefined,
      data2: Array<any> | undefined,
      data3: Array<any> | undefined,
      data4: Array<any> | undefined,
  ) => {
    const data = [...(data1 || []), ...(data2 || []), ...(data3 || []), ...(data4 || [])];
    return data.sort((a, b) => {
      return Date.parse(a?.createdAt) - Date.parse(b?.createdAt);
    });
  };

  return (
    <Timeline
      loading={notesLoading || conversationsLoading || ticketsLoading}
      noActivity={noHistoryItemsAvailable}
      contactId={id}
      loggedActivities={[liveConversations, ...getSortedItems(
          notes?.notes.content,
          conversations?.conversations.content,
          tickets?.tickets,
          actions
      )]}
    />
  );
};

export default ContactHistory;
