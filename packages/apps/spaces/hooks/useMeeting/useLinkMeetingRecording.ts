import { useMeetingLinkRecordingMutation } from './types';
import { toast } from 'react-toastify';
import {
  MeetingLinkRecordingMutation,
  MeetingUnlinkRecordingMutation,
} from '../../graphQL/__generated__/generated';
import FormData from 'form-data';
import axios from 'axios';
import { convert } from 'html-to-text';

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

  const handleLinkMeetingRecording: Result['onLinkMeetingRecording'] = async (
    attachmentId,
  ) => {
    try {
      const response = await linkMeetingRecordingMutation({
        variables: {
          meetingId,
          attachmentId,
        },
      });
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
        // eslint-disable-next-line no-unsafe-optional-chaining
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
          convert(response?.data?.meeting_LinkRecording?.agenda || ''),
        );
        request.append('type', 'meeting');
        request.append(
          'file_id',
          response?.data?.meeting_LinkRecording.recording.id,
        );

        axios
          .post(`/transcription-api/transcribe`, request, {
            headers: {
              accept: `application/json`,
            },
          })
          .then((res) => {
            if (res.status === 200) {
              toast.success(
                `Meeting recording transcription started successfully`,
                {
                  toastId: `update-meeting-${meetingId}`,
                },
              );
            }
          });
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
      toast.error(`Removing recording not supported yet`),
        {
          toastId: `remove-recording-error`,
        };
      return null;
    },
  };
};
