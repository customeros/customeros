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
    input: Omit<MeetingInput, 'appSource'>,
  ) => Promise<UpdateMeetingMutation['meeting_Update'] | null>;
}

export const useUpdateMeeting = ({ meetingId, appSource }: Props): Result => {
  const [updateMeetingMutation, { loading, error, data }] =
    useUpdateMeetingMutation();
  const client = useApolloClient();

  const handleUpdateMeeting: Result['onUpdateMeeting'] = async (meeting) => {
    try {
      const response = await updateMeetingMutation({
        variables: {
          meetingId: meetingId,
          meetingInput: { ...meeting, appSource: appSource || 'OPENLINE' },
        },
        // update: handleUpdateCacheAfterAddingMeeting,
      });

      console.log('🏷️ ----- response: ', response);
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
            start
            end
            createdAt
            agenda
            agendaContentType
            recording
          }
        `,
      });

      console.log('🏷️ ----- meeting attendeedby: ', meeting.attendedBy);

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
              start
              end
              createdAt
              agenda
              agendaContentType
              recording
            }
          `,
          data: {
            ...data,
            // attendedBy: [{ contactParticipant: { ...meeting.attendedBy } }],
          },
        });
      }
      console.log('🏷️ ----- data: !!!!!!', data);

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
