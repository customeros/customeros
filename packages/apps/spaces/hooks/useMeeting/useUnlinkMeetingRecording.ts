import {
  GetContactTimelineDocument,
  NOW_DATE,
  useMeetingUnlinkRecordingMutation,
  MeetingUnlinkRecordingMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import { GetContactTimelineQuery } from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  meetingId: string;
  contactId?: string;
}

export interface Result {
  onUnlinkMeetingRecording: (
    fileId: string,
  ) => Promise<
    MeetingUnlinkRecordingMutation['meeting_UnlinkRecording'] | null
  >;
}

export const useUnlinkMeetingRecording = ({
  meetingId,
  contactId,
}: Props): Result => {
  const [unlinkMeetingRecordingMutation, { loading, error, data }] =
    useMeetingUnlinkRecordingMutation();
  const loggedInUserData = useRecoilValue(userData);

  const handleUpdateCacheAfterUnlinkingRecording = (
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

  const handleUnlinkMeetingRecording: Result['onUnlinkMeetingRecording'] =
    async (attachmentId) => {
      try {
        const response = await unlinkMeetingRecordingMutation({
          variables: {
            meetingId,
            attachmentId,
          },

          update: handleUpdateCacheAfterUnlinkingRecording,
        });

        return response.data?.meeting_UnlinkRecording ?? null;
      } catch (err) {
        console.error(err);
        toast.error(`removing a recording from a meeting is not supported yet`);
        return null;
      }
    };

  return {
    onUnlinkMeetingRecording: handleUnlinkMeetingRecording,
  };
};
