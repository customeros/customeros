import { CreateTagMutation, TagInput, useCreateTagMutation } from './types';
import { toast } from 'react-toastify';

interface Result {
  onCreateTag: (
    input: TagInput,
  ) => Promise<CreateTagMutation['tag_Create'] | null>;
}
export const useCreateTag = (): Result => {
  const [createContactMutation, { loading, error, data }] =
    useCreateTagMutation();

  const handleCreateTag: Result['onCreateTag'] = async (tag) => {
    try {
      // const optimisticItem = { id: 'optimistic-id', ...tag };
      const response = await createContactMutation({
        variables: { input: tag },
      });
      return response.data?.tag_Create ?? null;
    } catch (err) {
      toast.error('Something went wrong while creating tag', {
        toastId: `tag-create-error`,
      });
      console.error(err);
      return null;
    }
  };

  return {
    onCreateTag: handleCreateTag,
  };
};
