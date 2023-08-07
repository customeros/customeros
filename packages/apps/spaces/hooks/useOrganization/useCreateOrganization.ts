import {
  CreateOrganizationMutation,
  OrganizationInput,
  useCreateOrganizationMutation,
} from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import { GetOrganizationsOptionsDocument } from '../../graphQL/__generated__/generated';
import client from '../../apollo-client';
import CheckWaves from '@spaces/atoms/icons/CheckWaves';
import Times from '@spaces/atoms/icons/Times';
import ExclamationWaves from '@spaces/atoms/icons/ExclamationWaves';

interface Result {
  saving: boolean;
  createdId?: string;
  onCreateOrganization: (
    input: OrganizationInput,
  ) => Promise<CreateOrganizationMutation['organization_Create'] | null>;
}
export const useCreateOrganization = (): Result => {
  const [createOrganizationMutation, { loading, data }] =
    useCreateOrganizationMutation();
  const handleUpdateCacheAfterAddingOrg = (
    cache: ApolloCache<any>,
    { data: { organization_Create } }: any,
  ) => {
    const data: any | null = client.readQuery({
      query: GetOrganizationsOptionsDocument,
      variables: {
        pagination: {
          limit: 9999,
          page: 1,
        },
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationsOptionsDocument,
        variables: {
          pagination: {
            limit: 9999,
            page: 1,
          },
        },
        data: {
          organizations: {
            content: [organization_Create],
          },
        },
      });
      return;
    }

    client.writeQuery({
      query: GetOrganizationsOptionsDocument,
      variables: {
        pagination: {
          limit: 9999,
          page: 1,
        },
      },
      data: {
        organizations: {
          content: {
            organization_Create,
            ...data.organizations?.content,
          },
        },
      },
    });
  };

  const handleCreateOrganization: Result['onCreateOrganization'] = async (
    input: OrganizationInput,
  ) => {
    const createOrganizationToast = toast.loading('Creating organization');
    try {
      const response = await createOrganizationMutation({
        variables: { input },
        update: handleUpdateCacheAfterAddingOrg,
      });
      if (response.data?.organization_Create) {
        toast.update(createOrganizationToast, {
          render: 'New organization created',
          type: 'success',
          isLoading: false,
          autoClose: 2000,
          icon: CheckWaves,
        });
      }
      return response.data?.organization_Create ?? null;
    } catch (err) {
      toast.update(createOrganizationToast, {
        render: 'Something went wrong while creating organization',
        type: 'error',
        isLoading: false,
        autoClose: 2000,
        icon: ExclamationWaves,
      });
      return null;
    }
  };

  return {
    onCreateOrganization: handleCreateOrganization,
    saving: loading,
    createdId: data?.organization_Create?.id,
  };
};
