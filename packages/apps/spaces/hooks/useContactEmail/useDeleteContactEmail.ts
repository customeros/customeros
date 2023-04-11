import { useRemoveEmailFromContactMutation } from '../useContact/types';
import { RemoveEmailFromContactMutation } from '../../graphQL/__generated__/generated';

interface Props {
  contactId: string;
}

interface Result {
  onRemoveEmailFromContact: (
    emailId: string,
  ) => Promise<
    RemoveEmailFromContactMutation['emailRemoveFromContactById'] | null
  >;
}
export const useRemoveEmailFromContactEmail = ({
  contactId,
}: Props): Result => {
  const [removeEmailFromContactMutation, { loading, error, data }] =
    useRemoveEmailFromContactMutation();

  const handleRemoveEmailFromContact: Result['onRemoveEmailFromContact'] =
    async (emailId) => {
      try {
        const response = await removeEmailFromContactMutation({
          variables: { contactId, id: emailId },
          refetchQueries: ['GetContactCommunicationChannels'],
          update(cache) {
            const normalizedId = cache.identify({
              id: emailId,
              __typename: 'Email',
            });
            cache.evict({ id: normalizedId });
            cache.gc();
          },
        });
        return response.data?.emailRemoveFromContactById ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveEmailFromContact: handleRemoveEmailFromContact,
  };
};
