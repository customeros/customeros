import {
  useUpdateOrganizationOwnerMutation,
  GetOrganizationOwnerDocument,
  UpdateOrganizationOwnerMutation,
  GetOrganizationOwnerQuery,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';

import client from '../../apollo-client';

interface Props {
  organizationId: string;
  userId: string;
}

interface Result {
  saving: boolean;
  onLinkOrganizationOwner: () => Promise<
    UpdateOrganizationOwnerMutation['organization_SetOwner'] | null
  >;
}

export const useLinkOrganizationOwner = ({
  organizationId,
  userId,
}: Props): Result => {
  const [linkOrganizationOwnerMutation, { loading }] =
    useUpdateOrganizationOwnerMutation();

  const handleUpdateCacheAfterAddingLocation = (
    cache: ApolloCache<any>,
    { data: { organization_SetOwner } }: any,
  ) => {
    const data: GetOrganizationOwnerQuery | null = client.readQuery({
      query: GetOrganizationOwnerDocument,
      variables: {
        id: organizationId,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationOwnerDocument,
        variables: {
          id: organizationId,
        },
        data: {
          organization: {
            id: organizationId,
            owner: [organization_SetOwner],
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        owner: { ...organization_SetOwner.owner },
      },
    };

    client.writeQuery({
      query: GetOrganizationOwnerDocument,
      data: newData,
      variables: {
        id: organizationId,
      },
    });
  };

  const handleLinkOrganizationOwner: Result['onLinkOrganizationOwner'] =
    async () => {
      try {
        const response = await linkOrganizationOwnerMutation({
          variables: { organizationId, userId },
          optimisticResponse: {
            organization_SetOwner: {
              id: 'new-organization-owner-id',
            },
          },
          update: handleUpdateCacheAfterAddingLocation,
        });
        if (response.data) {
          toast.success('Owner set!', {
            toastId: `owner-set-${response.data?.organization_SetOwner.id}`,
          });
        }
        return response.data?.organization_SetOwner ?? null;
      } catch (err) {
        toast.error('Something went wrong while setting the owner', {
          toastId: `owner-set-error-${organizationId}`,
        });
        return null;
      }
    };

  return {
    saving: loading,
    onLinkOrganizationOwner: handleLinkOrganizationOwner,
  };
};
