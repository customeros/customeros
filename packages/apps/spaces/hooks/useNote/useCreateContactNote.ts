import {
  NoteInput,
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
} from './types';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
  saving: boolean;
  onCreateContactNote: (
    input: NoteInput,
  ) => Promise<CreateContactNoteMutation['note_CreateForContact'] | null>;
}
export const useCreateContactNote = ({ contactId }: Props): Result => {
  const [createContactNoteMutation, { loading, error, data }] =
    useCreateContactNoteMutation();

  const handleCreateContactNote: Result['onCreateContactNote'] = async (
    note,
  ) => {
    try {
      const response = await createContactNoteMutation({
        variables: { contactId, input: note },
        refetchQueries: ['GetContactTimeline'],
      });
      if (response.data) {
        // toast.success('Note added!', {
        //   toastId: `note-added-${response.data?.note_CreateForContact.id}`,
        // });
      }
      return response.data?.note_CreateForContact ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding a note', {
        toastId: `note-add-error-${contactId}`,
      });
      return null;
    }
  };

  return {
    saving: loading,
    onCreateContactNote: handleCreateContactNote,
  };
};
