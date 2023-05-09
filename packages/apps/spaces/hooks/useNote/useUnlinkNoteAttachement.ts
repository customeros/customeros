import {
  GetContactTimelineDocument,
  NOW_DATE,
  useNoteUnlinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client';
import {
  GetContactTimelineQuery,
  NoteUnlinkAttachmentMutation,
} from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import { useRecoilValue } from 'recoil';
import { userData } from '../../state';

export interface Props {
  noteId: string;
}

export interface Result {
  onUnlinkNoteAttachment: (
    fileId: string,
  ) => Promise<NoteUnlinkAttachmentMutation['note_UnlinkAttachment'] | null>;
}

export const useUnlinkNoteAttachment = ({
  noteId,
}: Props): Result => {
  const [unlinkNoteAttachmentMutation, { loading, error, data }] =
    useNoteUnlinkAttachmentMutation();

  const handleUnlinkNoteAttachment: Result['onUnlinkNoteAttachment'] = async (
    attachmentId,
  ) => {
    try {
      const response = await unlinkNoteAttachmentMutation({
        variables: {
          noteId,
          attachmentId,
        },
      });

      return response.data?.note_UnlinkAttachment ?? null;
    } catch (err) {
      console.error(err);
      toast.error(
        `Something went wrong while adding draft note to the timeline`,
      );
      return null;
    }
  };

  return {
    onUnlinkNoteAttachment: handleUnlinkNoteAttachment,
  };
};
