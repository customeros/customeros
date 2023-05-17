import {
  RemoveOrganizationSubsidiaryMutation,
  useRemoveOrganizationSubsidiaryMutation,
} from './types';
import { toast } from 'react-toastify';

interface Result {
  onRemoveOrganizationSubsidiary: ({
    subsidiaryId,
  }: {
    subsidiaryId: string;
  }) => Promise<
    RemoveOrganizationSubsidiaryMutation['organization_RemoveSubsidiary'] | null
  >;
}
export const useRemoveOrganizationSubsidiary = ({
  organizationId,
}: {
  organizationId: string;
}): Result => {
  const [removeOrganizationSubsidiaryMutation] =
    useRemoveOrganizationSubsidiaryMutation();

  const handleDeleteOrganization: Result['onRemoveOrganizationSubsidiary'] =
    async ({ subsidiaryId }: { subsidiaryId: string }) => {
      try {
        await removeOrganizationSubsidiaryMutation({
          variables: { organizationId, subsidiaryId },
          update(cache) {
            const normalizedId = cache.identify({
              id: subsidiaryId,
              __typename: 'Organization',
            });

            cache.evict({ id: normalizedId });
            cache.gc();
          },
        });

        return null;
      } catch (err) {
        toast.error(
          'Something went wrong while deleting organization branch. Please contact us or try again later',
          {
            toastId: `organzation--subsidiary${subsidiaryId}-delete-error`,
          },
        );
        return null;
      }
    };

  return {
    onRemoveOrganizationSubsidiary: handleDeleteOrganization,
  };
};
