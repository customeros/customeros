import {
  OrganizationInput,
  UpdateOrganizationMutation,
  useUpdateOrganizationMutation,
} from './types';

interface Result {
  onUpdateOrganization: UpdateOrganizationMutation['organization_Update'];
}
export const useUpdateOrganization = (): Result => {
  const [updateOrganizationMutation, { loading, error, data }] =
    useUpdateOrganizationMutation();

  const handleUpdateOrganization: Result['onUpdateOrganization'] = async (
    input: OrganizationInput,
  ) => {
    try {
      const response = await updateOrganizationMutation({
        variables: { input },
      });
      return response.data?.organization_Update ?? null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateOrganization: handleUpdateOrganization,
  };
};
