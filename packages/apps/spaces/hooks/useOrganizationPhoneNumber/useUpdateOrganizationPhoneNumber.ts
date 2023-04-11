import {
  PhoneNumber,
  PhoneNumberUpdateInput,
  UpdateOrganizationPhoneNumberMutation,
  useUpdateOrganizationPhoneNumberMutation,
} from './types';

interface Result {
  onUpdateOrganizationPhoneNumber: (
    input: PhoneNumberUpdateInput,
  ) => Promise<
    | UpdateOrganizationPhoneNumberMutation['phoneNumberUpdateInOrganization']
    | null
  >;
}
export const useUpdateOrganizationPhoneNumber = ({
  organizationId,
}: {
  organizationId: string;
}): Result => {
  const [updateOrganizationNoteMutation, { loading, error, data }] =
    useUpdateOrganizationPhoneNumberMutation();

  const handleUpdateOrganizationPhoneNumber: Result['onUpdateOrganizationPhoneNumber'] =
    async (input) => {
      const payload = {
        ...input,
      };
      try {
        const response = await updateOrganizationNoteMutation({
          variables: { input: payload, organizationId },
          optimisticResponse: {
            phoneNumberUpdateInOrganization: {
              __typename: 'PhoneNumber',
              ...payload,
              e164: payload.phoneNumber,
              rawPhoneNumber: payload.phoneNumber,
              primary: input.primary || false,
            },
          },
        });

        return response.data?.phoneNumberUpdateInOrganization ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onUpdateOrganizationPhoneNumber: handleUpdateOrganizationPhoneNumber,
  };
};
