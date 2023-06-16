import { ContactTagInput, useAddTagToContactMutation } from './types';
import { GetContactTagsDocument } from '@spaces/graphql';
import { gql } from '@apollo/client';
import { toast } from 'react-toastify';

interface Result {
  onAddTagToContact: (input: ContactTagInput) => void;
}
export const useAddTagToContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  // const client = useApolloClient();
  const [addTagToContactMutation] = useAddTagToContactMutation();

  const onAddTagToContact: Result['onAddTagToContact'] = (contactTagInput) => {
    addTagToContactMutation({
      variables: { input: contactTagInput },
      onError: () => {
        toast.error('Something went wrong while adding tag', {
          toastId: `contact-add-tag-error`,
        });
      },
      update: (cache, response) => {
        const data = cache.readQuery({
          query: GetContactTagsDocument,
          variables: { id: contactTagInput.contactId },
        });

        cache.writeFragment({
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
            tags: [
              // @ts-expect-error revisit
              ...(data.tags ?? []),
              // @ts-expect-error revisit
              ...(response.data ? response.data.contact_AddTagById.tags : []),
            ],
          },
        });
      },
    });
  };

  return {
    onAddTagToContact,
  };
};
