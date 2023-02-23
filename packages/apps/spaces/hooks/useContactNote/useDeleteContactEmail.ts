import {
  RemoveContactNoteMutation,
  useRemoveContactNoteMutation,
} from '../../graphQL/generated';

interface Result {
  onRemoveContactNote: (
    emailId: string,
  ) => Promise<RemoveContactNoteMutation['note_Delete'] | null>;
}
export const useRemoveEmailFromContactEmail = (): Result => {
  const [removeEmailFromContactMutation, { loading, error, data }] =
    useRemoveContactNoteMutation();

  const handleRemoveEmailFromContact: Result['onRemoveContactNote'] = async (
    id,
  ) => {
    try {
      const response = await removeEmailFromContactMutation({
        variables: { id: id },
      });
      return response.data?.note_Delete ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onRemoveContactNote: handleRemoveEmailFromContact,
  };
};
