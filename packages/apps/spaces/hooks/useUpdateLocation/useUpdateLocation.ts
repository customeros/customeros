import {
  LocationUpdateInput,
  UpdateLocationMutation,
  useUpdateLocationMutation,
} from './types';
import { toast } from 'react-toastify';

interface Result {
  saving: boolean;
  onUpdateLocation: (
    input: LocationUpdateInput,
  ) => Promise<UpdateLocationMutation['location_Update'] | null>;
}
export const useUpdateLocation = (): Result => {
  const [updateLocationMutation, { loading }] = useUpdateLocationMutation();

  const handleUpdateLocation: Result['onUpdateLocation'] = async (input) => {
    try {
      const response = await updateLocationMutation({
        variables: { input: { ...input } },
      });

      return response.data?.location_Update ?? null;
    } catch (err) {
      toast.error(
        'Something went wrong while updating location! Please contact us or try again later',
        {
          toastId: `update-location-error-${input.id}`,
        },
      );
      console.error(err);
      return null;
    }
  };

  return {
    onUpdateLocation: handleUpdateLocation,
    saving: loading,
  };
};
