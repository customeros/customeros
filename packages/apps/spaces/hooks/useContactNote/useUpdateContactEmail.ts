import { Note, NoteUpdateInput } from '../../graphQL/generated';
import {
  UpdateContactNoteMutation,
  useUpdateContactNoteMutation,
} from '../../graphQL/generated';

interface Result {
  onUpdateContactNote: (
    input: NoteUpdateInput,
    oldValue: Note,
  ) => Promise<UpdateContactNoteMutation['note_Update'] | null>;
}
export const useUpdateContactNote = (): Result => {
  const [updateContactNoteMutation, { loading, error, data }] =
    useUpdateContactNoteMutation();

  const handleUpdateContactNote: Result['onUpdateContactNote'] = async (
    input,
    oldValue,
  ) => {
    try {
      const response = await updateContactNoteMutation({
        variables: { input },
        optimisticResponse: {
          note_Update: {
            __typename: 'Note',
            ...oldValue,
            ...input,
          },
        },
      });

      return response.data?.note_Update ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateContactNote: handleUpdateContactNote,
  };
};
