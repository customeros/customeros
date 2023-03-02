import {
  MergeOrganizationsMutation,
  useMergeOrganizationsMutation,
} from './types';
import { toast } from 'react-toastify';
import { useRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../../components/finder/state';

interface Result {
  onMergeOrganizations: (input: {
    primaryOrganizationId: string;
    mergedOrganizationIds: Array<string>;
  }) => Promise<MergeOrganizationsMutation['organization_Merge'] | null>;
}
export const useMergeOrganizations = (): Result => {
  const [mergeOrganizationsMutation, { loading, error, data }] =
    useMergeOrganizationsMutation();
  const [mode, setMode] = useRecoilState(tableMode);
  const [selectedItems, setSelectedItems] = useRecoilState(selectedItemsIds);
  const handleMergeOrganizations: Result['onMergeOrganizations'] = async (
    input,
  ) => {
    try {
      const response = await mergeOrganizationsMutation({
        variables: input,
        refetchQueries: ['GetDashboardData'],
      });

      if (response.data?.organization_Merge !== null) {
        setMode('PREVIEW');
        setSelectedItems([]);
        toast.success('Organizations were successfully merged!', {
          toastId: `merge-organizations-success-${input.primaryOrganizationId}`,
        });
      }
      return response.data?.organization_Merge ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong and selected organizations could not me merged!',
        {
          toastId: `merge-organizations-error-${input.primaryOrganizationId}`,
        },
      );

      return null;
    }
  };

  return {
    onMergeOrganizations: handleMergeOrganizations,
  };
};
