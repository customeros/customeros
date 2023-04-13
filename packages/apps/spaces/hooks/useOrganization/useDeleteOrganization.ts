import {
  DeleteOrganizationMutation,
  useDeleteOrganizationMutation,
} from './types';
import { useRouter } from 'next/router';
import { toast } from 'react-toastify';

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
        update(cache) {
          const normalizedId = cache.identify({
            id: id,
            __typename: 'Organization',
          });
          cache.evict({ id: normalizedId });
          cache.gc();
        },
      });

      if (response.data?.organization_Delete?.result) {
        push('/').then(() =>
          toast.success('Organization successfully deleted!', {
            toastId: `organzation-${id}-delete-success`,
          }),
        );
      }

      return null;
    } catch (err) {
      toast.error(
        'Something went wrong while deleting organization. Please contact us or try again later',
        {
          toastId: `organzation-${id}-delete-error`,
        },
      );
      return null;
    }
  };

  return {
    onDeleteOrganization: handleDeleteOrganization,
  };
};
