import { useNoteLinkAttachmentMutation } from './types';
import { toast } from 'react-toastify';
import { NoteLinkAttachmentMutation } from '../../graphQL/__generated__/generated';

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
