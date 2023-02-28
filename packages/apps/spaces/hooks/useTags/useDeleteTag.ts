import { DeleteTagMutation, useDeleteTagMutation } from './types';

interface Result {
  onDeleteTag: (id: string) => Promise<DeleteTagMutation['tag_Delete'] | null>;
}
export const useDeleteTag = (): Result => {
  const [deleteTagMutation, { loading, error, data }] = useDeleteTagMutation();

  const handleDeleteTag = async (id: string) => {
    try {
      const response = await deleteTagMutation({
        variables: { id: id },
      });
      return response.data?.tag_Delete ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onDeleteTag: handleDeleteTag,
  };
};
