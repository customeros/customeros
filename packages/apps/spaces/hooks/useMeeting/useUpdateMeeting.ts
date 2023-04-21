import {
  UpdateMeetingMutation,
  useUpdateMeetingMutation,
  GetContactTimelineQuery,
  GetOrganizationTimelineDocument,
  MeetingInput,
  NOW_DATE,
} from './types';
import { ApolloCache } from 'apollo-cache';
import client from '../../apollo-client';
import { toast } from 'react-toastify';
export interface Props {
  meetingId?: string;
}

export interface Result {
  onUpdateMeeting: (
    input: MeetingInput,
  ) => Promise<UpdateMeetingMutation['meeting_Update'] | null>;
}

export const useUpdateMeeting = ({ meetingId }: Props): Result => {
  const [updateMeetingMutation, { loading, error, data }] =
    useUpdateMeetingMutation();

  // const handleUpdateCacheAfterAddingMeeting = (
  //   cache: ApolloCache<any>,
  //   { data: { meeting_Update } }: any,
  // ) => {
  //   const data: GetContactTimelineQuery | null = client.readQuery({
  //     query: GetOrganizationTimelineDocument,
  //     variables: {
  //       organizationId,
  //       from: NOW_DATE,
  //       size: 10,
  //     },
  //   });
  //
  //   if (data === null) {
  //     client.writeQuery({
  //       query: GetOrganizationTimelineDocument,
  //       data: {
  //         contact: {
  //           organizationId,
  //           timelineEvents: [meeting_Update],
  //         },
  //         variables: { organizationId, from: NOW_DATE, size: 10 },
  //       },
  //     });
  //     return;
  //   }
  //
  //   const newData = {
  //     contact: {
  //       ...data.contact,
  //       timelineEvents: [
  //         ...(data.contact?.timelineEvents || []),
  //         meeting_Update,
  //       ],
  //     },
  //   };
  //
  //   client.writeQuery({
  //     query: GetOrganizationTimelineDocument,
  //     data: newData,
  //     variables: {
  //       organizationId,
  //       from: NOW_DATE,
  //       size: 10,
  //     },
  //   });
  // };

  const handleUpdateMeetingFromContact: Result['onUpdateMeeting'] = async (
    meeting,
  ) => {
    try {
      const optimisticItem = { id: 'optimistic-id', ...meeting };
      const response = await updateMeetingMutation({
        variables: { meetingId: meetingId, meetingInput: meeting },
        // @ts-expect-error fixme
        // update: handleUpdateCacheAfterAddingMeeting,
      });

      toast.success(`Added draft meeting to the timeline`);
      return response.data?.meeting_Update ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while updating meeting `, {
        toastId: 'id-test',
      });
      return null;
    }
  };

  return {
    onUpdateMeeting: handleUpdateMeetingFromContact,
  };
};
