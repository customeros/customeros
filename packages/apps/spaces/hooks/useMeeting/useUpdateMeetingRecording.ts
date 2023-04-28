import {
  UpdateMeetingMutation,
  useUpdateMeetingMutation,
  MeetingInput,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from 'apollo-cache';
import { GetContactTagsDocument } from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
export interface Props {
  meetingId: string;
  appSource: string;
}

export interface Result {
  onUpdateMeeting: (
    input: Pick<MeetingInput, 'recording'>,
  ) => Promise<UpdateMeetingMutation['meeting_Update'] | null>;
}

export const useUpdateMeetingRecording = ({ meetingId, appSource }: Props): Result => {
  const [updateMeetingMutation, { loading, error, data }] =
    useUpdateMeetingMutation();
  const client = useApolloClient();

  const handleUpdateMeeting: Result['onUpdateMeeting'] = async (meetingRecording) => {
    try {
      const response = await updateMeetingMutation({
        variables: {
          meetingId: meetingId,
          meetingInput: { ...meetingRecording, appSource: appSource || 'OPENLINE' },
        },
        // update: handleUpdateCacheAfterAddingMeeting,
      });
      if (response?.data?.meeting_Update.recording) {
        // call transcript api
        toast.success(`Meeting recording updated successfully`, {
          toastId: `update-meeting-${meetingId}`,
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
            startedAt
            endedAt
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

      // client.cache.writeFragment({
      //   id: `Contact:${contactId}`,
      //   fragment: gql`
      //     fragment Tags on Contact {
      //       id
      //       tags
      //     }
      //   `,
      //   data: {
      //     // @ts-expect-error revisit
      //     ...data.contact,
      //     // @ts-expect-error revisit
      //     tags: [...data.tags, response.data?.contact_AddTagById.tags],
      //   },
      // });
      // Update the cache with the new object

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
    onUpdateMeeting: handleUpdateMeeting,
  };
};
