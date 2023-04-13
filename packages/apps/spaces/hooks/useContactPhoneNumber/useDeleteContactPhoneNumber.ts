import {
  RemovePhoneNumberFromContactMutation,
  useRemovePhoneNumberFromContactMutation,
} from '../../graphQL/__generated__/generated';
import { toast } from 'react-toastify';

interface Result {
  onRemovePhoneNumberFromContact: (
    phoneNumberId: string,
  ) => Promise<
    | RemovePhoneNumberFromContactMutation['phoneNumberRemoveFromContactById']
    | null
  >;
  loading: boolean;
}
export const useRemovePhoneNumberFromContact = ({
  contactId,
}: {
  contactId: string;
}): Result => {
  const [removePhoneNumberFromContactMutation, { loading, error, data }] =
    useRemovePhoneNumberFromContactMutation();

  const handleRemovePhoneNumberFromContact = async (id: string) => {
    try {
      const response = await removePhoneNumberFromContactMutation({
        variables: { id: id, contactId },
        refetchQueries: ['GetContactCommunicationChannels'],
        update(cache) {
          const normalizedId = cache.identify({
            id,
            __typename: 'PhoneNumber',
          });
          cache.evict({ id: normalizedId });
          cache.gc();
        },
      });
      return response.data?.phoneNumberRemoveFromContactById ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while deleting phone number. Please contact us or try again later',
        {
          toastId: `phone-number-${contactId}-${id}-delete-error`,
        },
      );
      console.error(err);
      return null;
    }
  };

  return {
    onRemovePhoneNumberFromContact: handleRemovePhoneNumberFromContact,
    loading,
  };
};
