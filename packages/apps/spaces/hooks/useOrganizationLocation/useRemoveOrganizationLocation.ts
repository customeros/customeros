import { toast } from 'react-toastify';

import { useRemoveLocationFromOrganizationMutation } from '@spaces/graphql';

interface Props {
  organizationId: string;
}

interface Result {
  saving: boolean;
  onRemoveOrganizationLocation: (locationId:string) => void;
}

export const useRemoveOrganizationLocation = ({
  organizationId,
}: Props): Result => {
  const [removeOrganizationLocationMutation, { loading }] =
    useRemoveLocationFromOrganizationMutation({
      onError: () => {
        toast.error('Something went wrong while removing location', {
          toastId: `Location-remove-error-${organizationId}`,
        });
      },
    });

  return {
    saving: loading,
    onRemoveOrganizationLocation: (locationId: string) =>
      removeOrganizationLocationMutation({
        variables: { organizationId, locationId },
      }),
  };
};
