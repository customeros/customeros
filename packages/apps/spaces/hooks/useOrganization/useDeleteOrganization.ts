import {
  DeleteOrganizationMutation,
  useDeleteOrganizationMutation,
} from './types';
import { useRouter } from 'next/router';

interface Result {
  onDeleteOrganization: () => Promise<
    DeleteOrganizationMutation['organization_Delete'] | null
  >;
}
export const useDeleteOrganization = ({ id }: { id: string }): Result => {
  const [createOrganizationMutation] = useDeleteOrganizationMutation();
  const { push } = useRouter();

  const handleDeleteOrganization: Result['onDeleteOrganization'] = async () => {
    try {
      const response = await createOrganizationMutation({
        variables: { id },
        refetchQueries: ['GetDashboardData'],
      });

      if (response.data?.organization_Delete) {
        push('/');
      }

      return null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  return {
    onDeleteOrganization: handleDeleteOrganization,
  };
};
