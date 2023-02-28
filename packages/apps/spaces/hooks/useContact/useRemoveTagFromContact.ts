import { AddTagToContactMutation, ContactTagInput } from './types';
import {
  RemoveTagFromContactMutation,
  useRemoveTagFromContactMutation,
} from '../../graphQL/__generated__/generated';

interface Result {
  onRemoveTagFromContact: (
    input: Omit<ContactTagInput, 'contactId'>,
  ) => Promise<RemoveTagFromContactMutation['contact_RemoveTagById'] | null>;
}
export const useRemoveTagFromContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [removeTagFromContactMutation, { loading, error, data }] =
    useRemoveTagFromContactMutation();

  const handleRemoveTagFromContact: Result['onRemoveTagFromContact'] = async (
    contactTagInput,
  ) => {
    try {
      const response = await removeTagFromContactMutation({
        variables: { input: { ...contactTagInput, contactId } },
        refetchQueries: ['GetContactTags'],
      });
      return response.data?.contact_RemoveTagById ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onRemoveTagFromContact: handleRemoveTagFromContact,
  };
};
