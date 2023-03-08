import {
  NoteUpdateInput,
  UpdateContactNoteMutation,
  useUpdateContactNoteMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Result {
  onUpdateContactNote: (
    input: NoteUpdateInput,
  ) => Promise<UpdateContactNoteMutation['note_Update'] | null>;
}
export const useUpdateContactNote = (): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactNoteMutation();

  const handleUpdateContactNote: Result['onUpdateContactNote'] = async (
    note,
  ) => {
    try {
      const response = await updateContactNoteMutation({
        variables: { input: note },
        refetchQueries: ['GetContactNotes'],
      });
      if (response.data) {
        toast.success('Note added!');
      }
      return response.data?.note_Update ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding a note');
      return null;
    }
  };

  return {
    onUpdateContactNote: handleUpdateContactNote,
  };
};
