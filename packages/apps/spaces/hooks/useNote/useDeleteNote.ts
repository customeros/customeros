import { RemoveNoteMutation, useRemoveNoteMutation } from './types';
import { toast } from 'react-toastify';

interface Result {
  onRemoveNote: (
    id: string,
  ) => Promise<RemoveNoteMutation['note_Delete'] | null>;
}
export const useDeleteNote = (): Result => {
  const [removeNoteMutation, { loading, error, data }] =
    useRemoveNoteMutation();

  const handleRemoveNote: Result['onRemoveNote'] = async (id) => {
    try {
      const response = await removeNoteMutation({
        variables: { id: id },
      });
      toast.success('Note deleted!', {
        toastId: `remove-note-success-${id}`,
      });
      return response.data?.note_Delete ?? null;
    } catch (err) {
      toast.error('Something went wrong while removing the note', {
        toastId: `remove-note-error-${id}`,
      });

      return null;
    }
  };

  return {
    onRemoveNote: handleRemoveNote,
  };
};
