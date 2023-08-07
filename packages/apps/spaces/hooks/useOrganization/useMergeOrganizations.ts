import {
  MergeOrganizationsMutation,
  useMergeOrganizationsMutation,
} from './types';
import { useRecoilState } from 'recoil';
import { selectedItemsIds } from '../../components/finder/state';
import { toastError, toastSuccess } from '@ui/presentation/Toast';

interface Result {
  onMergeOrganizations: (input: {
    primaryOrganizationId: string;
    mergedOrganizationIds: Array<string>;
  }) => Promise<MergeOrganizationsMutation['organization_Merge'] | null>;
}
export const useMergeOrganizations = (): Result => {
  const [mergeOrganizationsMutation, { loading, error, data }] =
    useMergeOrganizationsMutation();
  const [selectedItems, setSelectedItems] = useRecoilState(selectedItemsIds);
  const handleMergeOrganizations: Result['onMergeOrganizations'] = async (
    input,
  ) => {
    try {
      const response = await mergeOrganizationsMutation({
        variables: input,
        refetchQueries: ['dashboardView_Organizations'],
      });

      if (response.data?.organization_Merge !== null) {
        setSelectedItems([]);
        toastSuccess(
          'Organizations merged',
          `merge-organizations-success-${input.primaryOrganizationId}`,
        );
      }
      return response.data?.organization_Merge ?? null;
    } catch (err) {
      toastError(
        'Something went wrong and selected organizations could not be merged!',
        `merge-organizations-error-${input.primaryOrganizationId}`,
      );

      return null;
    }
  };

  return {
    onMergeOrganizations: handleMergeOrganizations,
  };
};
