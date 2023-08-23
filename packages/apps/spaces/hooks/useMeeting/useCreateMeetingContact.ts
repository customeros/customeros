import {
  GetContactTimelineDocument,
  NOW_DATE,
  Result,
  useCreateMeetingMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import { DataSource, GetContactTimelineQuery } from '@spaces/graphql';
import client from '../../apollo-client';

export interface Props {
  contactId?: string;
}
export const useCreateMeetingFromContact = ({ contactId }: Props): Result => {
  const [createMeetingMutation, { loading, error, data }] =
    useCreateMeetingMutation({
      onError: () => {
        toast.error(
          `Something went wrong while adding draft meeting to the timeline`,
        );
      },
      onCompleted: () => {
        setTimeout(() => ~{
        }, 0);
      },
    });

  const handleUpdateCacheAfterAddingMeeting = (
    cache: ApolloCache<any>,
    { data: { meeting_Create } }: any,
  ) => {
    const data: GetContactTimelineQuery | null = client.readQuery({
      query: GetContactTimelineDocument,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });

    const newMeeting = {
      createdAt: new Date(),
      meetingStartedAt: new Date(),
      meetingEndedAt: new Date(),
      agendaContentType: 'text/html',
      meetingCreatedBy: meeting_Create.createdBy,
      describedBy: [],
      includes: [],
      events: [],
      recording: null,
      source: DataSource.Openline,
      ...meeting_Create,
    };

    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [newMeeting],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [newMeeting],
      },
    };

    client.writeQuery({
      query: GetContactTimelineDocument,
      data: newData,
      variables: {
        contactId,
        from: NOW_DATE,
        size: 10,
      },
    });
  };

  const handleCreateMeetingFromContact: Result['onCreateMeeting'] = async (
    userId,
  ) => {
    return createMeetingMutation({
      variables: {
        meeting: {
          createdBy: [{ userId: userId }],
          attendedBy: [{ contactId: contactId }],
          appSource: 'OPENLINE',
          name: '',
          startedAt: new Date().toISOString(),
          endedAt: new Date().toISOString(),
          agenda: `<p>INTRODUCTION</p>
                     <p>DISCUSSION</p>
                     <p>NEXT STEPS</p>
                     `,
          agendaContentType: 'text/html',
          note: { html: '<p>Notes:</p>', appSource: 'OPENLINE' },
        },
      },
      update: handleUpdateCacheAfterAddingMeeting,
    });
  };

  return {
    onCreateMeeting: handleCreateMeetingFromContact,
  };
};
