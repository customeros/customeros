import { useRemoveEmailFromContactMutation } from '../useContact/types';
import { RemoveEmailFromContactMutation } from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

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
        toast.error(
          'Something went wrong while deleting email. Please contact us or try again later',
          {
            toastId: `email-${emailId}-delete-error`,
          },
        );
        console.error(err);
        return null;
      }
    };

  return {
    onRemoveEmailFromContact: handleRemoveEmailFromContact,
  };
};
