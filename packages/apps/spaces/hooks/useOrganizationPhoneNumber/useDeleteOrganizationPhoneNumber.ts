import {
  RemovePhoneNumberFromOrganizationMutation,
  useRemovePhoneNumberFromOrganizationMutation,
} from './types';

interface Result {
  onRemovePhoneNumberFromOrganization: (
    phoneNumberId: string,
  ) => Promise<
    | RemovePhoneNumberFromOrganizationMutation['phoneNumberRemoveFromOrganizationById']
    | null
  >;
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
      });
      return response.data?.phoneNumberRemoveFromOrganizationById ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onRemovePhoneNumberFromOrganization:
      handleRemovePhoneNumberFromOrganization,
  };
};
