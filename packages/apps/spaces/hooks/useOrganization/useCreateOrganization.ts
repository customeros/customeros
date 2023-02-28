import {
  CreateOrganizationMutation,
  OrganizationInput,
  useCreateOrganizationMutation,
} from './types';

interface Result {
  onCreateOrganization: CreateOrganizationMutation['organization_Create'];
}
export const useCreateOrganization = (): Result => {
  const [createOrganizationMutation, { loading, error, data }] =
    useCreateOrganizationMutation();

  const handleCreateOrganization: Result['onCreateOrganization'] = async (
    input: OrganizationInput,
  ) => {
    try {
      const response = await createOrganizationMutation({
        variables: { input },
      });
      return response.data?.organization_Create ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onCreateOrganization: handleCreateOrganization,
  };
};
