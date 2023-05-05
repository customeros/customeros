import {
  GetContactTimelineDocument,
  NOW_DATE,
  useMeetingLinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import {
  GetContactTimelineQuery,
  LinkMeetingAttachmentMutation,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  meetingId: string;
  contactId?: string;
}

export interface Result {
  onLinkMeetingAttachement: (
    attachmentId: string,
  ) => Promise<LinkMeetingAttachmentMutation['meeting_LinkAttachment'] | null>;
}

export const useLinkMeetingAttachement = ({
  meetingId,
  contactId,
}: Props): Result => {
  const [linkMeetingAttachementMutation, { loading, error, data }] =
    useMeetingLinkAttachmentMutation();

  const handleLinkMeetingAttachement: Result['onLinkMeetingAttachement'] =
    async (attachmentId) => {
      try {
        const response = await linkMeetingAttachementMutation({
          variables: {
            meetingId,
            attachmentId,
          },
        });

        toast.success(`Added attachment to meeting`);
        return response.data?.meeting_LinkAttachment ?? null;
      } catch (err) {
        console.error(err);
        toast.error(`Something went wrong while attaching file to the meeting`);
        return null;
      }
    };

  return {
    onLinkMeetingAttachement: handleLinkMeetingAttachement,
  };
};
