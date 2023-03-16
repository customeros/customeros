import {
  NoteUpdateInput,
  UpdateNoteMutation,
  useUpdateNoteMutation,
} from './types';
import { toast } from 'react-toastify';

interface Result {
  onUpdateNote: (
    input: NoteUpdateInput,
  ) => Promise<UpdateNoteMutation['note_Update'] | null>;
}
export const useUpdateNote = (): Result => {
  const [updateNoteMutation, { loading, error, data }] =
    useUpdateNoteMutation();

  const handleUpdateNote: Result['onUpdateNote'] = async (note) => {
    try {
      const response = await updateNoteMutation({
        variables: { input: note },
      });
      if (response.data) {
        toast.success('Note updated!', {
          toastId: `note-update-success-${note.id}`,
        });
      }
      return response.data?.note_Update ?? null;
    } catch (err) {
      toast.error('Something went wrong while updating the note', {
        toastId: `note-update-error-${note.id}`,
      });
      return null;
    }
  };

  return {
    onUpdateNote: handleUpdateNote,
  };
};
