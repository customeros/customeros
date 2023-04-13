import {
  AddTagToContactMutation,
  ContactTagInput,
  useAddTagToContactMutation,
} from './types';
import { GetContactTagsDocument } from '../../graphQL/__generated__/generated';
import { gql, useApolloClient } from '@apollo/client';
import { toast } from 'react-toastify';

interface Result {
  onAddTagToContact: (
    input: ContactTagInput,
  ) => Promise<AddTagToContactMutation['contact_AddTagById'] | null>;
}
export const useAddTagToContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const client = useApolloClient();
  const [addTagToContactMutation, { loading, error, data }] =
    useAddTagToContactMutation();

  const handleAddTagToContact: Result['onAddTagToContact'] = async (
    contactTagInput,
  ) => {
    try {
      const response = await addTagToContactMutation({
        variables: { input: contactTagInput },
      });
      const data = client.cache.readQuery({
        query: GetContactTagsDocument,
        variables: { id: contactTagInput.contactId },
      });

      client.cache.writeFragment({
        id: `Contact:${contactId}`,
        fragment: gql`
          fragment Tags on Contact {
            id
            tags
          }
        `,
        data: {
          // @ts-expect-error revisit
          ...data.contact,
          // @ts-expect-error revisit
          tags: [...data.tags, response.data?.contact_AddTagById.tags],
        },
      });
      // Update the cache with the new object
      return response.data?.contact_AddTagById ?? null;
    } catch (err) {
      toast.error('Something went wrong while adding tag', {
        toastId: `contact-add-tag-error`,
      });
      console.error(err);
      return null;
    }
  };

  return {
    onAddTagToContact: handleAddTagToContact,
  };
};
