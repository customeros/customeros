import {
  GetContactTimelineDocument,
  NOW_DATE,
  Result,
  useCreateMeetingMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import {
  DataSource,
  GetContactTimelineQuery,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  contactId?: string;
}
export const useCreateMeetingFromContact = ({ contactId }: Props): Result => {
  const [createMeetingMutation, { loading, error, data }] =
    useCreateMeetingMutation();

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
      ...meeting_Create,
      createdAt: new Date(),
      agenda: '',
      agendaContentType: 'text/html',
      meetingCreatedBy: meeting_Create.createdBy,
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
    try {
      const response = await createMeetingMutation({
        variables: {
          meeting: {
            createdBy: [{ userID: userId, type: 'user' }],
            attendedBy: [{ contactID: contactId, type: 'contact' }],
            appSource: 'OPENLINE',
            name: '',
            start: new Date().toISOString(),
            end: new Date().toISOString(),
            agenda: '',
            agendaContentType: 'text/html',
          },
        },

        update: handleUpdateCacheAfterAddingMeeting,
      });

      if (response.data?.meeting_Create.id) {
        toast.success(`Added draft meeting to the timeline`, {
          toastId: `draft-meeting-added-${response.data?.meeting_Create.id}`,
        });
      }

      return response.data?.meeting_Create ?? null;
    } catch (err) {
      console.error(err);
      toast.error(
        `Something went wrong while adding draft meeting to the timeline`,
      );
      return null;
    }
  };

  return {
    onCreateMeeting: handleCreateMeetingFromContact,
  };
};
