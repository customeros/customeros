import {
  AddLocationToOrganizationMutation,
  GetOrganizationLocationsQuery,
  useAddLocationToOrganizationMutation,
  GetOrganizationLocationsDocument,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';

import client from '../../apollo-client';

interface Props {
  organizationId: string;
}

interface Result {
  saving: boolean;
  onCreateOrganizationLocation: () => Promise<
    AddLocationToOrganizationMutation['organization_AddNewLocation'] | null
  >;
}

export const useCreateOrganizationLocation = ({
  organizationId,
}: Props): Result => {
  const [createOrganizationLocationMutation, { loading }] =
    useAddLocationToOrganizationMutation();

  const handleUpdateCacheAfterAddingLocation = (
    cache: ApolloCache<any>,
    { data: { organization_AddNewLocation } }: any,
  ) => {
    const data: GetOrganizationLocationsQuery | null = client.readQuery({
      query: GetOrganizationLocationsDocument,
      variables: {
        id: organizationId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationLocationsDocument,
        variables: {
          id: organizationId,
        },
        data: {
          organization: {
            id: organizationId,
            locations: [organization_AddNewLocation],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        locations: [
          ...(data.organization?.locations || []),
          { ...organization_AddNewLocation },
        ],
      },
    };

    client.writeQuery({
      query: GetOrganizationLocationsDocument,
      data: newData,
      variables: {
        id: organizationId,
      },
    });
  };

  const handleCreateOrganizationLocation: Result['onCreateOrganizationLocation'] =
    async () => {
      try {
        const response = await createOrganizationLocationMutation({
          variables: { organzationId: organizationId },
          optimisticResponse: {
            organization_AddNewLocation: {
              id: 'new-organization-location-id',
            },
          },
          update: handleUpdateCacheAfterAddingLocation,
        });
        return response.data?.organization_AddNewLocation ?? null;
      } catch (err) {
        toast.error('Something went wrong while adding location', {
          toastId: `Location-add-error-${organizationId}`,
        });
        return null;
      }
    };

  return {
    saving: loading,
    onCreateOrganizationLocation: handleCreateOrganizationLocation,
  };
};
