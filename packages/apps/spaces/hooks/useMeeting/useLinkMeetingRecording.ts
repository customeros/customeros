import {
  GetContactTimelineDocument,
  NOW_DATE,
  useMeetingLinkRecordingMutation,
  useMeetingUnlinkRecordingMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import {
  GetContactTimelineQuery,
  MeetingLinkRecordingMutation,
  MeetingUnlinkRecordingMutation,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import FormData from 'form-data';
import axios from 'axios';

const { convert } = require('html-to-text');

export interface Props {
  meetingId: string;
  contactId?: string;
}

export interface Result {
  onLinkMeetingRecording: (
    attachmentId: string,
  ) => Promise<MeetingLinkRecordingMutation['meeting_LinkRecording'] | null>;

  onUnLinkMeetingRecording: (
    attachmentId: string,
  ) => Promise<
    MeetingUnlinkRecordingMutation['meeting_UnlinkRecording'] | null
  >;
}

export const useLinkMeetingRecording = ({
  meetingId,
  contactId,
}: Props): Result => {
  const [linkMeetingRecordingMutation, { loading, error, data }] =
    useMeetingLinkRecordingMutation();

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

  const handleLinkMeetingRecording: Result['onLinkMeetingRecording'] = async (
    attachmentId,
  ) => {
    try {
      const response = await linkMeetingRecordingMutation({
        variables: {
          meetingId,
          attachmentId,
        },

        //update: handleUpdateCacheAfterAddingMeeting,
      });
      console.log('Got response from update meeting mutation');
      console.log(response);
      if (response?.data?.meeting_LinkRecording.recording) {
        // call transcript api
        //move ot after transcript is done

        const request = new FormData();

        request.append('group_id', meetingId);
        request.append(
          'start',
          response?.data?.meeting_LinkRecording.meetingStartedAt.slice(0, -1) +
            '+00:00',
        );
        const users = [];
        const contacts = [];
        for (const participant of response?.data?.meeting_LinkRecording
          .attendedBy) {
          if (participant?.__typename === 'UserParticipant') {
            users.push(participant.userParticipant.id);
          } else if (participant?.__typename === 'ContactParticipant') {
            contacts.push(participant.contactParticipant.id);
          }
        }
        request.append('users', JSON.stringify(users));
        request.append('contacts', JSON.stringify(contacts));
        request.append(
          'topic',
          convert(response?.data?.meeting_LinkRecording.agenda),
        );
        request.append('type', 'meeting');
        request.append(
          'file_id',
          response?.data?.meeting_LinkRecording.recording,
        );

        // axios
        //   .post(`/transcription-api/transcribe`, request, {
        //     headers: {
        //       accept: `application/json`,
        //     },
        //   })
        //   .then((res) => {
        //     if (res.status === 200) {
        //       toast.success(`Meeting recording updated successfully`, {
        //         toastId: `update-meeting-${meetingId}`,
        //       });
        //     }
        //   });
      }

      toast.success(`Added meeting recording`);
      return response.data?.meeting_LinkRecording ?? null;
    } catch (err) {
      console.error(err);
      toast.error(
        `Something went wrong while attaching recording to the meeting`,
      );
      return null;
    }
  };

  return {
    onLinkMeetingRecording: handleLinkMeetingRecording,
    onUnLinkMeetingRecording: async (attachmentId) => {
      toast.error(`Removing recording not supported yet`);
      return null;
    },
  };
};
