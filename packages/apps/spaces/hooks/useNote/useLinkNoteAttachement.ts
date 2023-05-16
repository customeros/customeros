import {
  GetContactTimelineDocument,
  GetOrganizationTimelineDocument,
  useNoteLinkAttachmentMutation,
} from './types';
import { toast } from 'react-toastify';
import { NoteLinkAttachmentMutation } from '../../graphQL/__generated__/generated';
import { ApolloCache } from '@apollo/client/cache';
import { useRouter } from 'next/router';

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
  const { query, pathname } = useRouter();

  const handleUpdateCacheAfterAddingNoteABC = (
    cache: ApolloCache<any>,
    { data: { note_LinkAttachment } }: any,
  ) => {
    cache.modify({
      id: cache.identify({ id: noteId, __typename: 'Note' }),
      broadcast: true,
      optimistic: true,
      fields: {
        includes: () => {
          return [...note_LinkAttachment.includes];
        },
      },
    });
  };

  // Fixme this should not be needed, cache modification should be enough but in this cace cache is always devalidated
  //  this code ensures that we get latest data with properly set "from"
  const getRefetchQueries = () => {
    const isContact = pathname.includes('contact');
    return [
      {
        query: isContact
          ? GetContactTimelineDocument
          : GetOrganizationTimelineDocument,
        variables: {
          [isContact ? 'contactId' : 'organizationId']: query.id,
          from: new Date().toISOString(),
          size: 15,
        },
      },
    ];
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
        refetchQueries: getRefetchQueries(),
        update: handleUpdateCacheAfterAddingNoteABC,
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
