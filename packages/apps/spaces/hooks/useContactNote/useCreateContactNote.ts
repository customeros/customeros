import {
  NoteInput,
  CreateContactNoteMutation,
  useCreateContactNoteMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Props {
  contactId: string;
}

interface Result {
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
      });
      if (response.data) {
        toast.success('Note added!');
      }
      return response.data?.note_CreateForContact ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding a note');
      return null;
    }
  };

  return {
    onCreateContactNote: handleCreateContactNote,
  };
};
