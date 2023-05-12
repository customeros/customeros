import {
  GetContactTimelineDocument,
  GetContactTimelineQuery,
  NOW_DATE,
  useNoteLinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { NoteLinkAttachmentMutation } from '../../graphQL/__generated__/generated';
import { ApolloCache } from '@apollo/client/cache';
import client from '../../apollo-client';
import { gql } from '@apollo/client';

export interface Props {
  noteId: string;
}

export interface Result {
  onLinkNoteAttachment: (
    attachmentId: string,
  ) => Promise<NoteLinkAttachmentMutation['note_LinkAttachment'] | null>;
}

export const useLinkNoteAttachment = ({ noteId }: Props): Result => {
  const [linkNoteAttachmentMutation, { loading, error, data }] =
    useNoteLinkAttachmentMutation();

  const handleUpdateCacheAfterAddingNoteABC = (
    cache: ApolloCache<any>,
    { data: { note_LinkAttachment } }: any,
  ) => {
    const note = cache.identify({ id: noteId, __typename: 'Note' });
    console.log('🏷️ ----- note: ', note);
    cache.modify({
      id: cache.identify({ id: noteId, __typename: 'Note' }),
      fields: {
        includes() {
          console.log('🏷️ ----- note_LinkAttachment: ', note_LinkAttachment);
          return [...note_LinkAttachment.includes];
        },
      },
    });
  };
  const handleLinkNoteAttachment: Result['onLinkNoteAttachment'] = async (
    attachmentId,
  ) => {
    try {
      const response = await linkNoteAttachmentMutation({
        variables: {
          noteId,
          attachmentId,
        },
      });

      toast.success(`Added attachment to note`);
      return response.data?.note_LinkAttachment ?? null;
    } catch (err) {
      console.error(err);
      toast.error(`Something went wrong while attaching file to the note`);
      return null;
    }
  };

  return {
    onLinkNoteAttachment: handleLinkNoteAttachment,
  };
};
