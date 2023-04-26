import { useLinkMeetingAttendeeMutation } from './types';
import { toast } from 'react-toastify';
import { LinkMeetingAttachmentMutation } from '../../graphQL/__generated__/generated';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  meetingId: string;
}

export interface Result {
  onLinkMeetingAttendee: (
    participantId: string,
  ) => Promise<LinkMeetingAttachmentMutation['meeting_LinkAttachment'] | null>;
}

export const useLinkMeetingAttendee = ({ meetingId }: Props): Result => {
  const [linkMeetingAttendeeMutation, { loading, error, data }] =
    useLinkMeetingAttendeeMutation();
  const loggedInUserData = useRecoilValue(userData);

  // const handleUpdateCacheAfterAddingMeeting = (
  //   cache: ApolloCache<any>,
  //   { data: { meeting_Create } }: any,
  // ) => {
  //   const data: GetContactTimelineQuery | null = client.readQuery({
  //     query: GetContactTimelineDocument,
  //     variables: {
  //       contactId,
  //       from: NOW_DATE,
  //       size: 10,
  //     },
  //   });
  //
  //   if (data === null) {
  //     client.writeQuery({
  //       query: GetContactTimelineDocument,
  //       data: {
  //         contact: {
  //           contactId,
  //           timelineEvents: [meeting_Create],
  //         },
  //         variables: { contactId, from: NOW_DATE, size: 10 },
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
  //         meeting_Create,
  //       ],
  //     },
  //   };
  //
  //   client.writeQuery({
  //     query: GetContactTimelineDocument,
  //     data: newData,
  //     variables: {
  //       contactId,
  //       from: NOW_DATE,
  //       size: 10,
  //     },
  //   });
  // };

  const handleLinkMeetingAttendee: Result['onLinkMeetingAttendee'] = async (
    participantId,
  ) => {
    try {
      const response = await linkMeetingAttendeeMutation({
        variables: {
          meetingId,
          meetingParticipant: participantId,
        },

        // //@ts-expect-error fixme
        // update: handleUpdateCacheAfterAddingMeeting,
      });
      console.log('🏷️ ----- response: link attendee ', response);

      // toast.success(`Added draft meeting to the timeline`);
      return response.data?.meeting_LinkAttendedBy ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while adding attendee`);
      return null;
    }
  };

  return {
    onLinkMeetingAttendee: handleLinkMeetingAttendee,
  };
};
