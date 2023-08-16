import {
  MergeOrganizationsMutation,
  useMergeOrganizationsMutation,
} from './types';
import { useSetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '@spaces/finder/state';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

interface Result {
  onMergeOrganizations: (input: {
    primaryOrganizationId: string;
    mergedOrganizationIds: Array<string>;
    onSuccess: () => void;
  }) => Promise<MergeOrganizationsMutation['organization_Merge'] | null>;
}
export const useMergeOrganizations = (): Result => {
  const setMode = useSetRecoilState(tableMode);

  const [mergeOrganizationsMutation] = useMergeOrganizationsMutation();
  const setSelectedItems = useSetRecoilState(selectedItemsIds);
  const handleMergeOrganizations: Result['onMergeOrganizations'] = async ({
    onSuccess,
    ...input
  }) => {
    try {
      const response = await mergeOrganizationsMutation({
        variables: input,
        refetchQueries: ['dashboardView_Organizations'],
      });

      if (response.data?.organization_Merge !== null) {
        onSuccess();
        setSelectedItems([]);
        setMode('PREVIEW');
        toastSuccess(
          'Organizations merged',
          `merge-organizations-success-${input.primaryOrganizationId}`,
        );
      }
      return response.data?.organization_Merge ?? null;
    } catch (err) {
      toastError(
        'We couldn’t merge these organizations. Please try again.',
        `merge-organizations-error-${input.primaryOrganizationId}`,
      );

      return null;
    }
  };

  return {
    onMergeOrganizations: handleMergeOrganizations,
  };
};
