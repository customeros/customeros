import { MergeContactsMutation, useMergeContactsMutation } from './types';
import { toast } from 'react-toastify';
import { useSetRecoilState } from 'recoil';
import { selectedItemsIds, tableMode } from '../../components/finder/state';

interface Result {
  onMergeContacts: (input: {
    primaryContactId: string;
    mergedContactIds: Array<string>;
  }) => Promise<MergeContactsMutation['contact_Merge'] | null>;
}
export const useMergeContacts = (): Result => {
  const [mergeContactsMutation, { loading, error, data }] =
    useMergeContactsMutation();
  const setSelectedItems = useSetRecoilState(selectedItemsIds);
  const handleMergeContacts: Result['onMergeContacts'] = async (input) => {
    try {
      const response = await mergeContactsMutation({
        variables: input,
        refetchQueries: ['GetDashboardData'],
      });

      if (response.data?.contact_Merge !== null) {
        setSelectedItems([]);
        toast.success('Contacts were successfully merged!', {
          toastId: `merge-Contacts-success-${input.primaryContactId}`,
        });
      }
      return response.data?.contact_Merge ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong and selected contacts could not be merged!',
        {
          toastId: `merge-Contacts-error-${input.primaryContactId}`,
        },
      );

      return null;
    }
  };

  return {
    onMergeContacts: handleMergeContacts,
  };
};
