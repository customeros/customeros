import {
  OrganizationInput,
  UpdateOrganizationMutation,
  useUpdateOrganizationMutation,
} from './types';
import { OrganizationUpdateInput } from '../../graphQL/__generated__/generated';

interface Props {
  organizationId: string;
}

interface Result {
  onUpdateOrganization: (
    input: Omit<OrganizationUpdateInput, 'id'>,
  ) => Promise<UpdateOrganizationMutation['organization_Update'] | null>;
}
export const useUpdateOrganization = ({ organizationId }: Props): Result => {
  const [updateOrganizationMutation, { loading, error, data }] =
    useUpdateOrganizationMutation();

  const handleUpdateOrganization: Result['onUpdateOrganization'] = async (
    input,
  ) => {
    try {
      const response = await updateOrganizationMutation({
        variables: { input: { ...input, id: organizationId } },
        refetchQueries: ['GetOrganizationDetails'],
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
