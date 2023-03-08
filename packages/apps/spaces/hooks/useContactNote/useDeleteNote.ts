import {
  RemoveContactNoteMutation,
  useRemoveContactNoteMutation,
} from '../../graphQL/__generated__/generated';

interface Result {
  onRemoveNote: (
    id: string,
  ) => Promise<RemoveContactNoteMutation['note_Delete'] | null>;
}
export const useDeleteNote = (): Result => {
  const [removeEmailFromContactMutation, { loading, error, data }] =
    useRemoveContactNoteMutation();

  const handleRemoveNote: Result['onRemoveNote'] = async (id) => {
    try {
      const response = await removeEmailFromContactMutation({
        variables: { id: id },
        refetchQueries: ['GetContactNotes'],
      });
      return response.data?.note_Delete ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onRemoveNote: handleRemoveNote,
  };
};
