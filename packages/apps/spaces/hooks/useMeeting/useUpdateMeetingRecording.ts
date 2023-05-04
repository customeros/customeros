import {
  UpdateMeetingMutation,
  useUpdateMeetingMutation,
  MeetingInput,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import { GetContactTagsDocument } from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
const { convert } = require('html-to-text');

export interface Props {
  meetingId: string;
  appSource: string;
}
import axios from 'axios';
import FormData from 'form-data';


export interface Result {
  onUpdateMeetingRecording: (
    input: Pick<MeetingInput, 'recording'>,
  ) => Promise<UpdateMeetingMutation['meeting_Update'] | null>;
}

export const useUpdateMeetingRecording = ({
  meetingId,
  appSource,
}: Props): Result => {
  const [updateMeetingMutation, { loading, error, data }] =
    useUpdateMeetingMutation();
  const client = useApolloClient();

  const handleUpdateMeeting: Result['onUpdateMeetingRecording'] = async (
    meetingRecording,
  ) => {
    console.log('***********Inside useUpdateMeetingRecording***********');

    try {
      const response = await updateMeetingMutation({
        variables: {
          meetingId: meetingId,
          meetingInput: {
            ...meetingRecording,
            appSource: appSource || 'OPENLINE',
          },
        },
        // update: handleUpdateCacheAfterAddingMeeting,
      });
      console.log('Got response from update meeting mutation');
      console.log(response);
      if (response?.data?.meeting_Update.recording) {
        // call transcript api
        //move ot after transcript is done

        const request = new FormData();

        request.append('group_id', meetingId);
        request.append(
          'start',
          response?.data?.meeting_Update.meetingStartedAt.slice(0, -1) +
            '+00:00',
        );
        let users = [];
        let contacts = [];
        for (let participant of response?.data?.meeting_Update.attendedBy) {
          if (participant?.__typename === 'UserParticipant') {
            users.push(participant.userParticipant.id);
          } else if (participant?.__typename === 'ContactParticipant') {
            contacts.push(participant.contactParticipant.id);
          }
        }
        request.append('users', JSON.stringify(users));
        request.append('contacts', JSON.stringify(contacts));
        request.append('topic', convert(response?.data?.meeting_Update.agenda));
        request.append('type', 'meeting');
        request.append('file_id', response?.data?.meeting_Update.recording);

        axios
          .post(`/transcription-api/transcribe`, request, {
            headers: {
              accept: `application/json`,
            },
          })
          .then((res) => {
            if (res.status === 200) {
              toast.success(`Meeting recording updated successfully`, {
                toastId: `update-meeting-${meetingId}`,
              });
            }
          });
      }
      const data = client.cache.readFragment({
        id: `Meeting:${meetingId}`,
        fragment: gql`
          fragment MeetingUpdateFragment on Meeting {
            id
            attendedBy {
              ... on UserParticipant {
                userParticipant {
                  id
                }
              }
              ... on ContactParticipant {
                contactParticipant {
                  id
                }
              }
            }
            meetingCreatedBy: createdBy {
              ... on UserParticipant {
                userParticipant {
                  id
                }
              }
              ... on ContactParticipant {
                contactParticipant {
                  id
                }
              }
            }
            meetingStartedAt: startedAt
            meetingEndedAt: endedAt
            createdAt
            agenda
            agendaContentType
            recording
          }
        `,
      });

      if (data) {
        client.cache.writeFragment({
          id: `Meeting:${meetingId}`,
          fragment: gql`
            fragment MeetingUpdateFragment on Meeting {
              id
              attendedBy {
                ... on UserParticipant {
                  userParticipant {
                    id
                  }
                }
                ... on ContactParticipant {
                  contactParticipant {
                    id
                  }
                }
              }
              meetingCreatedBy: createdBy {
                ... on UserParticipant {
                  userParticipant {
                    id
                  }
                }
                ... on ContactParticipant {
                  contactParticipant {
                    id
                  }
                }
              }
              startedAt
              endedAt
              createdAt
              agenda
              agendaContentType
              recording
            }
          `,
          data: {
            ...data,
            // attendedBy: [{ contactParticipant: { ...meetingRecording.attendedBy } }],
          },
        });
      }

      return response.data?.meeting_Update ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while updating meeting `, {
        toastId: `update-meeting-${meetingId}`,
      });
      return null;
    }
  };

  return {
    onUpdateMeetingRecording: handleUpdateMeeting,
  };
};
