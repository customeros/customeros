import {
  HideOrganizationsMutation,
  useHideOrganizationsMutation,
  DashboardView_OrganizationsDocument,
} from './types';
import { useSetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '@spaces/finder/state';
import { toastError } from '@ui/presentation/Toast';

interface Result {
  onHideOrganizations: ({
    ids,
  }: {
    ids: string[];
  }) => Promise<HideOrganizationsMutation['organization_HideAll'] | null>;
}
export const useHideOrganizations = (): Result => {
  const [createOrganizationMutation] = useHideOrganizationsMutation();
  const setSelectedItems = useSetRecoilState(selectedItemsIds);
  const setMode = useSetRecoilState(tableMode);

  const handleHideOrganization: Result['onHideOrganizations'] = async ({
    ids,
  }) => {
    try {
      const response = await createOrganizationMutation({
        variables: { ids },
        awaitRefetchQueries: true,
        refetchQueries: [DashboardView_OrganizationsDocument],
      });
      if (response.data?.organization_HideAll?.result) {
        setSelectedItems([]);
        setMode('PREVIEW');
      }

      return null;
    } catch (err) {
      toastError(
        `We couldn’t hide ${ids.length === 1 ? 'this' : 'these'} organization${
          ids.length === 1 ? '' : 's'
        }. Please try again.`,
        `organzations-hide-error`,
      );
      return null;
    }
  };

  return {
    onHideOrganizations: handleHideOrganization,
  };
};
