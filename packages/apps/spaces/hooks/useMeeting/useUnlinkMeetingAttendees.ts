import {
  GetContactTimelineDocument,
  NOW_DATE,
  useUnlinkMeetingAttendeeMutation,
  UnlinkMeetingAttendeeMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import { GetContactTimelineQuery } from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  meetingId: string;
  contactId?: string;
}

export interface Result {
  onUnlinkMeetingAttendee: (
    fileId: string,
  ) => Promise<
    UnlinkMeetingAttendeeMutation['meeting_UnlinkAttendedBy'] | null
  >;
}

export const useUnlinkMeetingAttendee = ({
  meetingId,
  contactId,
}: Props): Result => {
  const [unlinkMeetingAttendeeMutation, { loading, error, data }] =
    useUnlinkMeetingAttendeeMutation();
  const loggedInUserData = useRecoilValue(userData);

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

    if (data === null) {
      client.writeQuery({
        query: GetContactTimelineDocument,
        data: {
          contact: {
            contactId,
            timelineEvents: [meeting_Create],
          },
          variables: { contactId, from: NOW_DATE, size: 10 },
        },
      });
      return;
    }

    const newData = {
      contact: {
        ...data.contact,
        timelineEvents: [
          ...(data.contact?.timelineEvents || []),
          meeting_Create,
        ],
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

  const handleUnlinkMeetingAttendee: Result['onUnlinkMeetingAttendee'] = async (
    participantId,
  ) => {
    try {
      const response = await unlinkMeetingAttendeeMutation({
        variables: {
          meetingId,
          meetingParticipant: participantId,
        },

        ////@ts-expect-error fixme
        // update: handleUpdateCacheAfterAddingMeeting,
      });
      console.log('🏷️ ----- response: unlink attendee ', response);
      // toast.success(`Added draft meeting to the timeline`);
      return response.data?.meeting_UnlinkAttendedBy ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while removing attendee from meeting`);
      return null;
    }
  };

  return {
    onUnlinkMeetingAttendee: handleUnlinkMeetingAttendee,
  };
};
