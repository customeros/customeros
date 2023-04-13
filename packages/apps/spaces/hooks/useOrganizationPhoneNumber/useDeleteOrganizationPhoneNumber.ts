import {
  RemovePhoneNumberFromOrganizationMutation,
  useRemovePhoneNumberFromOrganizationMutation,
} from './types';
import { toast } from 'react-toastify';

interface Result {
  onRemovePhoneNumberFromOrganization: (
    phoneNumberId: string,
  ) => Promise<
    | RemovePhoneNumberFromOrganizationMutation['phoneNumberRemoveFromOrganizationById']
    | null
  >;
  loading: boolean;
}
export const useRemovePhoneNumberFromOrganization = ({
  organizationId,
}: {
  organizationId: string;
}): Result => {
  const [removePhoneNumberFromOrganizationMutation, { loading, error, data }] =
    useRemovePhoneNumberFromOrganizationMutation();

  const handleRemovePhoneNumberFromOrganization = async (id: string) => {
    try {
      const response = await removePhoneNumberFromOrganizationMutation({
        variables: { id: id, organizationId },
        refetchQueries: ['GetOrganizationCommunicationChannels'],
        awaitRefetchQueries: true,
        update(cache) {
          const normalizedId = cache.identify({
            id,
            __typename: 'PhoneNumber',
          });
          cache.evict({ id: normalizedId });
          cache.gc();
        },
      });
      return response.data?.phoneNumberRemoveFromOrganizationById ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while deleting phone number! Please contact us or try again later',
      );

      console.error(err);
      return null;
    }
  };

  return {
    onRemovePhoneNumberFromOrganization:
      handleRemovePhoneNumberFromOrganization,
    loading,
  };
};
