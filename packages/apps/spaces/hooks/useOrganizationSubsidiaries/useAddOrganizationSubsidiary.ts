import { useAddOrganizationSubsidiaryMutation } from './types';
import { toast } from 'react-toastify';
import { ApolloCache } from '@apollo/client/cache';
import {
  GetOrganizationSubsidiariesDocument,
  LinkOrganizationsInput,
} from '@spaces/graphql';
import client from '../../apollo-client';

interface Result {
  saving: boolean;
  onAddOrganizationSubsidiary: (input: LinkOrganizationsInput) => void;
}
export const useAddOrganizationSubsidiary = ({
  id,
}: {
  id: string;
}): Result => {
  const handleUpdateCacheAfterAddingOrgSubsidiary = (
    cache: ApolloCache<any>,
    { data: { organization_AddSubsidiary } }: any,
  ) => {
    const data: any | null = client.readQuery({
      query: GetOrganizationSubsidiariesDocument,
      variables: {
        id,
      },
    });

    if (data === null) {
      client.writeQuery({
        query: GetOrganizationSubsidiariesDocument,
        variables: {
          id,
        },
        data: {
          organization: {
            id,
            subsidiaries: [organization_AddSubsidiary],
          },
        },
      });
      return;
    }

    client.writeQuery({
      query: GetOrganizationSubsidiariesDocument,
      variables: {
        id,
      },
      data: {
        organization: {
          id,
          subsidiaries: [organization_AddSubsidiary],
        },
      },
    });
  };

  const [addOrganizationMutation, { loading, error, data }] =
    useAddOrganizationSubsidiaryMutation({
      update: handleUpdateCacheAfterAddingOrgSubsidiary,
      onError: () =>
        toast.error(
          'Something went wrong while adding organization subsidiary',
        ),
    });

  const handleAddOrganizationSubsidiary: Result['onAddOrganizationSubsidiary'] =
    async (input: LinkOrganizationsInput) => {
      return addOrganizationMutation({
        variables: { input },
      });
    };

  return {
    onAddOrganizationSubsidiary: handleAddOrganizationSubsidiary,
    saving: loading,
  };
};
