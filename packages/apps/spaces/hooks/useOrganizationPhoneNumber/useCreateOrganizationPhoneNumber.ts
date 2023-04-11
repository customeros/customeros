import {
  AddPhoneToOrganizationMutation,
  useAddPhoneToOrganizationMutation,
} from './types';

interface Props {
  organizationId: string;
}

interface Result {
  onCreateOrganizationPhoneNumber: (
    input: any, //FIXME
  ) => Promise<
    AddPhoneToOrganizationMutation['phoneNumberMergeToOrganization'] | null
  >;
}
export const useCreateOrganizationPhoneNumber = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationPhoneNumberMutation, { loading, error, data }] =
    useAddPhoneToOrganizationMutation();

  const handleCreateOrganizationPhoneNumber: Result['onCreateOrganizationPhoneNumber'] =
    async (input) => {
      try {
        const response = await createOrganizationPhoneNumberMutation({
          variables: { organizationId, input },
          refetchQueries: ['GetOrganizationCommunicationChannels'],
        });
        return response.data?.phoneNumberMergeToOrganization ?? null;
      } catch (err) {
        console.error(err);
        return null;
      }
    };

  return {
    onCreateOrganizationPhoneNumber: handleCreateOrganizationPhoneNumber,
  };
};
