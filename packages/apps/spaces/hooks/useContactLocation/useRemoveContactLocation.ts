import { toast } from 'react-toastify';
import { useRemoveLocationFromContactMutation } from '@spaces/graphql';

interface Props {
  contactId: string;
}

interface Result {
  saving: boolean;
  onRemoveContactLocation: (locationId: string) => void;
}

export const useRemoveContactLocation = ({ contactId }: Props): Result => {
  const [removeContactLocationMutation, { loading }] =
    useRemoveLocationFromContactMutation({
      onError: () => {
        toast.error('Something went wrong while removing location', {
          toastId: `Location-remove-error-${contactId}`,
        });
      },
    });

  return {
    saving: loading,
    onRemoveContactLocation: (locationId: string) =>
      removeContactLocationMutation({
        variables: { contactId, locationId },
      }),
  };
};
