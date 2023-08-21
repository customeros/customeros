import {
  GetOrganizationOwnerDocument,
  GetOrganizationOwnerQuery,
  RemoveOrganizationOwnerMutation,
  useRemoveOrganizationOwnerMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';

import client from '../../apollo-client';

interface Props {
  organizationId: string;
  onCompleted?: () => void;
}

interface Result {
  saving: boolean;
  onUnlinkOrganizationOwner: () => Promise<
    RemoveOrganizationOwnerMutation['organization_UnsetOwner'] | null
  >;
}

export const useUnlinkOrganizationOwner = ({
  organizationId,
  onCompleted,
}: Props): Result => {
  const [unlinkOrganizationOwnerMutation, { loading }] =
    useRemoveOrganizationOwnerMutation({
      onCompleted,
    });

  const handleUpdateCacheAfterUnlinkingOwner = (
    cache: ApolloCache<any>,
    { data: { organization_UnsetOwner } }: any,
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
            owner: null,
          },
        },
      });
      return;
    }

    const newData = {
      organization: {
        ...data.organization,
        owner: null,
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

  const handleUnlinkOrganizationOwner: Result['onUnlinkOrganizationOwner'] =
    async () => {
      try {
        const response = await unlinkOrganizationOwnerMutation({
          variables: { organizationId },
          optimisticResponse: {
            organization_UnsetOwner: {
              id: organizationId,
            },
          },
          update: handleUpdateCacheAfterUnlinkingOwner,
        });

        return response.data?.organization_UnsetOwner ?? null;
      } catch (err) {
        toast.error('Something went wrong while removing the owner', {
          toastId: `owner-unset-error-${organizationId}`,
        });
        return null;
      }
    };

  return {
    saving: loading,
    onUnlinkOrganizationOwner: handleUnlinkOrganizationOwner,
  };
};
