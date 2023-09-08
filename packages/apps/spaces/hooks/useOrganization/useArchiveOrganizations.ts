import {
  ArchiveOrganizationsMutation,
  useArchiveOrganizationsMutation,
  DashboardView_OrganizationsDocument,
} from './types';
import { toast } from 'react-toastify';
import { useSetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '@spaces/finder/state';
import { toastError } from '@ui/presentation/Toast';

interface Result {
  onArchiveOrganization: ({
    ids,
  }: {
    ids: string[];
  }) => Promise<ArchiveOrganizationsMutation['organization_ArchiveAll'] | null>;
}
export const useArchiveOrganizations = (): Result => {
  const [createOrganizationMutation] = useArchiveOrganizationsMutation();
  const setSelectedItems = useSetRecoilState(selectedItemsIds);
  const setMode = useSetRecoilState(tableMode);

  const handleArchiveOrganization: Result['onArchiveOrganization'] = async ({
    ids,
  }) => {
    try {
      const response = await createOrganizationMutation({
        variables: { ids },
        awaitRefetchQueries: true,
        refetchQueries: [DashboardView_OrganizationsDocument],
      });
      if (response.data?.organization_ArchiveAll?.result) {
        setSelectedItems([]);
        setMode('PREVIEW');
      }

      return null;
    } catch (err) {
      toastError(
        `We couldn’t archive ${
          ids.length === 1 ? 'this' : 'these'
        } organization${ids.length === 1 ? '' : 's'}. Please try again.`,
        `organzations-archive-error`,
      );
      return null;
    }
  };

  return {
    onArchiveOrganization: handleArchiveOrganization,
  };
};
